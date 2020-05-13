/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/store"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

// Drain drains a node and prepares it for maintenance.
type Drain struct {
	store         store.Store
	clusterClient cluster.ClientInterface
}

var _ action.Dispatcher = (*Drain)(nil)

// NewDrain creates an instances of drain
func NewDrain(objectStore store.Store, clusterClient cluster.ClientInterface) *Drain {
	drain := &Drain{
		store:         objectStore,
		clusterClient: clusterClient,
	}

	return drain
}

// ActionName returns the name of this action
func (u *Drain) ActionName() string {
	return "action.octant.dev/drain"
}

// Handle executing drain
func (u *Drain) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	logger := log.From(ctx).With("actionName", u.ActionName())
	logger.With("payload", payload).Infof("received action payload")

	key, err := store.KeyFromPayload(payload)
	if err != nil {
		return err
	}


	object, err := u.store.Get(ctx, key)
	if err != nil {
		return err
	}

	if object == nil {
		return errors.New("object store cannot get node")
	}

	var node *corev1.Node
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, &node); err != nil {
		return err
	}

	message := fmt.Sprintf("Node %q marked as schedulable", key.Name)
	alertType := action.AlertTypeInfo
	if err := u.Drain(node); err != nil {
		message = fmt.Sprintf("Unable to drain node %q: %s", key.Name, err)
		alertType = action.AlertTypeWarning
		logger := log.From(ctx)
		logger.WithErr(err).Errorf("drain node")
	}
	alert := action.CreateAlert(alertType, message, action.DefaultAlertExpiration)
	alerter.SendAlert(alert)
	return nil
}

// Drain deletes resources from the node to prepare it for maintenance.
func (u *Drain) Drain(node *corev1.Node) error {
	if node == nil {
		return errors.New("nil node")
	}

	client, err := u.clusterClient.KubernetesClient()
	if err != nil {
		return err
	}

	currentNode, err := client.CoreV1().Nodes().Get(node.Name, metav1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "unable to find node %q", node.Name)
	}

	originalNode, err := json.Marshal(currentNode)
	if err != nil {
		return err
	}

	if !currentNode.Spec.Unschedulable {
		message := fmt.Sprintf("node %q already unmarked", node.Name)
		return errors.New(message)
	}
	currentNode.Spec.Unschedulable = false

	modifiedNode, err := json.Marshal(currentNode)
	if err != nil {
		return err
	}

	patchBytes, patchErr := strategicpatch.CreateTwoWayMergePatch(originalNode, modifiedNode, node)
	if patchErr != nil {
		_, err = client.CoreV1().Nodes().Patch(node.Name, types.StrategicMergePatchType, patchBytes)
	} else {
		_, err = client.CoreV1().Nodes().Update(currentNode)
		return errors.Wrapf(err, "failed to drain %q", node.Name)
	}

	return err
}
