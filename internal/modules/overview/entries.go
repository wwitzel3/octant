package overview

import (
	"context"
	"fmt"
	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/loading"
	"github.com/vmware-tanzu/octant/pkg/access"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/store"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sync"
)

var _ SectionGenerator = (*sectionGenerator)(nil)

type Entries map[string]EntryValue

type EntryValue struct {
	title      string
	gvk        schema.GroupVersionKind
	objectType interface{}
	listType   interface{}
}

type sectionGenerator struct {
	objectAccess access.Access
	objectStore  store.Store

	entries   Entries
	namespace string

	navHelper navigation.EntriesHelper
	sections  []describer.Describer

	ctx      context.Context
	cancelFn context.CancelFunc

	mu sync.Mutex
}

type SectionGenerator interface {
	Entries(prefix string) ([]navigation.Navigation, bool, error)
	Describer() *describer.Section
}

func NewSectionGenerator(ctx context.Context, objectStore store.Store, objectAccess access.Access, entries Entries, namespace string) *moduleContent {
	c, f := context.WithCancel(ctx)
	return &sectionGenerator{
		entries:      entries,
		namespace:    namespace,
		objectAccess: objectAccess,
		objectStore:  objectStore,
		ctx:          c,
		cancelFn:     f,
	}
}

func (o *sectionGenerator) generate() {
	o.mu.Lock()
	defer o.mu.Unlock()

	neh := navigation.EntriesHelper{}
	neh.Add("Overview", "", false)

	for k, v := range o.entries {
		if o.objectAccess.Allowed(o.namespace, access.List, gvk.CronJob) {
			neh.Add(v.title, k, loading.IsObjectLoading(o.ctx, o.namespace, store.KeyFromGroupVersionKind(v.gvk), o.objectStore))
			o.sections = append(o.sections, describer.NewResource(describer.ResourceOptions{
				Path:           fmt.Sprintf("/workloads/%s", k),
				ObjectStoreKey: store.KeyFromGroupVersionKind(v.gvk),
				ListType:       v.listType,
				ObjectType:     v.objectType,
				Titles:         describer.ResourceTitle{List: v.title, Object: v.title},
				RootPath:       describer.ResourceLink{Title: "Workloads", Url: "/overview/namespace/($NAMESPACE)/workloads"},
			}))
		}
	}

	o.navHelper = neh
}

func (o *sectionGenerator) Entries(prefix string) ([]navigation.Navigation, bool, error) {
	o.generate()

	children, err := o.navHelper.Generate(prefix, o.namespace, "")
	if err != nil {
		return nil, false, err
	}
	return children, false, nil
}

func (o *sectionGenerator) Describer() *describer.Section {
	o.generate()

	return describer.NewSection(
		"/",
		"Overview",
		o.sections...,
	)
}
