package plugin

import xpresource "github.com/crossplane/crossplane-runtime/pkg/resource"

type ResourceMerger interface {
	MergeResources(xpresource.Managed, xpresource.Managed) MergeDescription
}

type MergeDescription struct {
	LateInitializedSpec bool
	StatusUpdated       bool
	AnnotationsUpdated  bool
	NeedsProviderUpdate bool
	AnyFieldUpdated     bool
}

func CompareInt64Slices(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}

	lookup := make(map[int64]struct{})
	for _, x := range a {
		lookup[x] = struct{}{}
	}
	for _, x := range b {
		if _, ok := lookup[x]; !ok {
			return false
		}
	}
	return true
}

func CompareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	lookup := make(map[string]struct{})
	for _, x := range a {
		lookup[x] = struct{}{}
	}
	for _, x := range b {
		if _, ok := lookup[x]; !ok {
			return false
		}
	}
	return true
}

func CompareMapString(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for key, val := range a {
		bv, ok := b[key]
		if !ok || bv != val {
			return false
		}
	}

	return true
}

func CompareMapInt64(a, b map[string]int64) bool {
	if len(a) != len(b) {
		return false
	}

	for key, val := range a {
		bv, ok := b[key]
		if !ok || bv != val {
			return false
		}
	}

	return true
}

func CompareMapBool(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}

	for key, val := range a {
		bv, ok := b[key]
		if !ok || bv != val {
			return false
		}
	}

	return true
}
