/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectaccess

import (
	"context"
	"github.com/vmware-tanzu/octant/internal/cluster"
	"sync"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/pkg/access"
)

var _ access.Access = (*accessCache)(nil)

type accessCache struct {
	ctx context.Context
	cancelFn context.CancelFunc

	client cluster.ClientInterface
	nsCache sync.Map
}

func (c *accessCache) Allowed(namespace, verb string, kind schema.GroupVersionKind) bool {
	return true
}

func NewObjectAccessCache(ctx context.Context, client cluster.ClientInterface) (*accessCache, error) {
	c, f := context.WithCancel(ctx)
	return &accessCache{
		ctx: c,
		cancelFn: f,
		client: client,
	}, nil
}
