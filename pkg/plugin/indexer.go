package plugin

import (
	"fmt"

	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

type Indexer struct {
	gvkIndex map[k8schema.GroupVersionKind]*ImplementationMerger
}

func NewIndexer() *Indexer {
	return &Indexer{gvkIndex: make(map[k8schema.GroupVersionKind]*ImplementationMerger)}
}

func (i *Indexer) Overlay(ft *Implementation) error {
	if ft.GVK.Empty() {
		return fmt.Errorf("nil value for functable GVK")
	}
	if _, ok := i.gvkIndex[ft.GVK]; !ok {
		i.gvkIndex[ft.GVK] = NewImplementationMerger()
	}
	i.gvkIndex[ft.GVK].Overlay(ft)
	return nil
}

func (i *Indexer) BuildIndex() (*Index, error) {
	idx := &Index{
		reconcilerConfigurers: make([]ReconcilerConfigurer, 0),
		schemeBuilders:        make([]*scheme.Builder, 0),
		ftMap:                 make(map[k8schema.GroupVersionKind]Implementation),
		gvkToTFName:           make(map[k8schema.GroupVersionKind]string),
		tfNameToGVK:           make(map[string]k8schema.GroupVersionKind),
	}
	for gvk, mt := range i.gvkIndex {
		merged, err := mt.Merge()
		if err != nil {
			return nil, err
		}
		idx.ftMap[gvk] = merged
		idx.reconcilerConfigurers = append(idx.reconcilerConfigurers, merged.ReconcilerConfigurer)
		idx.schemeBuilders = append(idx.schemeBuilders, merged.SchemeBuilder)
		idx.gvkToTFName[gvk] = merged.TerraformResourceName
		idx.tfNameToGVK[merged.TerraformResourceName] = gvk
	}
	return idx, nil
}

type Index struct {
	ftMap                 map[k8schema.GroupVersionKind]Implementation
	reconcilerConfigurers []ReconcilerConfigurer
	schemeBuilders        []*scheme.Builder
	gvkToTFName           map[k8schema.GroupVersionKind]string
	tfNameToGVK           map[string]k8schema.GroupVersionKind
}

func (idx *Index) ReconcilerConfigurers() []ReconcilerConfigurer {
	return idx.reconcilerConfigurers
}

func (idx *Index) SchemeBuilders() []*scheme.Builder {
	return idx.schemeBuilders
}

func (idx *Index) InvokerForGVK(gvk k8schema.GroupVersionKind) (*Invoker, error) {
	ft, ok := idx.ftMap[gvk]
	if !ok {
		return &Invoker{}, fmt.Errorf("Could not look up functable for gvk=%s", gvk.String())
	}
	return &Invoker{ft: ft}, nil
}
