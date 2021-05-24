package objectstore

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
)

type informerKey struct {
	namespace string
	gvr       schema.GroupVersionResource
}

type interuptibleInformer struct {
	stopCh   chan struct{}
	informer informers.GenericInformer
	gvr      schema.GroupVersionResource
}

func (i interuptibleInformer) Stop() {
	close(i.stopCh)
}
