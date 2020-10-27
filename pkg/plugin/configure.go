package plugin

import (
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane-contrib/terraform-runtime/pkg/client"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ReconcilerConfigurer interface {
	ConfigureReconciler(ctrl.Manager, logging.Logger, *Index, *client.ProviderPool) error
}
