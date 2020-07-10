package plugin

import (
	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// An ImplementationMerger collects a sequence of Implementations
// through the Overlay method, and then can generate a single Implementation
// which is the result of merging all the layers into a single Implementation.
// It does this by picking a non-nil value for each field
// from the highest possible layer. Please see indexer_test.go for clarification.
type ImplementationMerger struct {
	layers []*Implementation
}

// Merge flattens all the layers in the ImplementationMerger, with the invariant:
// a layers' non-empty fields always replace underlying layers' fields.
// So last field Overlayed wins, eg call Overlay w/ the generated base
// Implementation before the user-override Implementations.
// I don't know if we really need an error; if we can assume generated
// Implementations cover everything, we don't.
// if that turns out to be a bad assumption, the error would give us
// an escape hatch to fail on incomplete implementations at runtime.
func (mt *ImplementationMerger) Merge() (Implementation, error) {
	merged := Implementation{}
	for _, ft := range mt.layers {
		if !ft.GVK.Empty() {
			merged.GVK = ft.GVK
		}
		if ft.TerraformResourceName != "" {
			merged.TerraformResourceName = ft.TerraformResourceName
		}
		if ft.CtyEncoder != nil {
			merged.CtyEncoder = ft.CtyEncoder
		}
		if ft.CtyDecoder != nil {
			merged.CtyDecoder = ft.CtyDecoder
		}
		if ft.SchemeBuilder != nil {
			merged.SchemeBuilder = ft.SchemeBuilder
		}
		if ft.ReconcilerConfigurer != nil {
			merged.ReconcilerConfigurer = ft.ReconcilerConfigurer
		}
		if ft.ResourceYAMLUnmarshaller != nil {
			merged.ResourceYAMLUnmarshaller = ft.ResourceYAMLUnmarshaller
		}
		if ft.ResourceYAMLMarshaller != nil {
			merged.ResourceYAMLMarshaller = ft.ResourceYAMLMarshaller
		}
		if ft.ResourceMerger != nil {
			merged.ResourceMerger = ft.ResourceMerger
		}
	}
	return merged, nil
}

func (mt *ImplementationMerger) Overlay(ft *Implementation) {
	mt.layers = append(mt.layers, ft)
}

func NewImplementationMerger() *ImplementationMerger {
	return &ImplementationMerger{
		layers: make([]*Implementation, 0),
	}
}

// Implementation is a collection of callbacks required to implement a resource
// There can be multiple Implementations for a resource, with
type Implementation struct {
	// GVK is used to index other elements of the Entry by GVK
	GVK k8schema.GroupVersionKind
	// TerraformResourceName is needed to map the crossplane type
	// to the Terraform type name. This is needed to find the schema
	// for the type and to identify the type in API calls.
	TerraformResourceName string
	// SchemeBuilder is used to register the controller for this type with the
	// controller runtime. StartTerraformManager (in pkg/controller) iterates
	// through all the registration entries and performs the bindings.
	SchemeBuilder *scheme.Builder
	// ReconcilerConfigurer is the function responsible for calling
	// managed.NewReconciler to bind the reconciler to the managed resource
	// type. It is also invoked in StartTerraformManager.
	ReconcilerConfigurer ReconcilerConfigurer
	// ResourceMerger can update the local kubernetes object with attributes
	// from the cloud provider, late-initializing Spec fields, copying over Status
	// fields and annotations.
	ResourceMerger ResourceMerger
	// CtyEncoder produces the cty.Value (cty-encoded resource for
	// terraform) for a resource.Managed, given the corresponding schema
	// object. Note that we do not try to compile schemas in to the generated
	// code, they are always obtained from the terraform process itself.
	CtyEncoder CtyEncoder
	// CtyDecoder is the complement to EncodeCtyCallback. In addition
	// to the schema and cty.Value, it also requires a resource.Managed, using
	// the deepcopied value from this resource as the base structure (and
	// providing values for .Spec fields and k8s metadata)
	CtyDecoder CtyDecoder
	// ResourceYAMLMarshaller is the complement to UnmarshalResourceCallback, taking
	// a resource.Managed and producing the []byte representation.
	ResourceYAMLMarshaller ResourceYAMLMarshaller
	// ResourceYAMLUnmarshaller is only used for prototyping atm -- it's a
	// function that can parse the []byte representation of a managed resource
	// to a resource.Managed
	ResourceYAMLUnmarshaller ResourceYAMLUnmarshaller
}
