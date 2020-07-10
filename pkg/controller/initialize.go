/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"

	crossplaneapis "github.com/crossplane/crossplane/apis"
	"github.com/crossplane/terraform-provider-runtime/pkg/client"
	"github.com/crossplane/terraform-provider-runtime/pkg/plugin"
)

//func StartTerraformManager(r *registry.Registry, opts ctrl.Options, ropts *client.RuntimeOptions, log logging.Logger) error {
func StartTerraformManager(idx *plugin.Index, p *plugin.ProviderInit, opts ctrl.Options, ropts *client.RuntimeOptions, log logging.Logger) error {
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return errors.Wrap(err, "Cannot get API server rest config")
	}
	mgr, err := ctrl.NewManager(cfg, opts)
	if err != nil {
		return errors.Wrap(err, "Cannot create controller manager")
	}
	err = crossplaneapis.AddToScheme(mgr.GetScheme())
	if err != nil {
		return errors.Wrap(err, "Cannot add core Crossplane APIs to scheme")
	}
	for _, sb := range idx.SchemeBuilders() {
		if err := sb.AddToScheme(mgr.GetScheme()); err != nil {
			return err
		}
	}
	p.SchemeBuilder.AddToScheme(mgr.GetScheme())
	pool := client.NewProviderPool(p.Initializer, ropts)
	for _, rc := range idx.ReconcilerConfigurers() {
		if err := rc.ConfigureReconciler(mgr, log, idx, pool); err != nil {
			return err
		}
	}
	err = mgr.Start(ctrl.SetupSignalHandler())
	if err != nil {
		return errors.Wrap(err, "Cannot start controller manager")
	}
	return nil
}
