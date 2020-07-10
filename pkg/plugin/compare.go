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
}
