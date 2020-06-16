/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package access

import "k8s.io/apimachinery/pkg/runtime/schema"

const (
	List = "list"
)

type Access interface {
	Allowed(namespace, verb string, kind schema.GroupVersionKind) bool
}
