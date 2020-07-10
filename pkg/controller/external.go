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
	"context"
	"fmt"

	kubeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/terraform-provider-runtime/pkg/api"
	"github.com/crossplane/terraform-provider-runtime/pkg/client"
	"github.com/crossplane/terraform-provider-runtime/pkg/plugin"
)

const (
	errNotMyType                  = "managed resource is not a MyType custom resource"
	errProviderNotRetrieved       = "provider could not be retrieved"
	errProviderSecretNil          = "cannot find Secret reference on Provider"
	errProviderSecretNotRetrieved = "secret referred in provider could not be retrieved"

	errNewClient = "cannot create new Service"
)

type External struct {
	KubeClient kubeclient.Client
	Invoker    *plugin.Invoker
	Callbacks  managed.ExternalClientFns
	logger     logging.Logger
	provider   *client.Provider
}

func (c *External) Observe(ctx context.Context, res resource.Managed) (managed.ExternalObservation, error) {
	c.entryLog(res, "Observe")
	gvk := res.GetObjectKind().GroupVersionKind()
	c.logger.Debug(fmt.Sprintf("terraform.External.Observe: %s", gvk.String()))
	if c.Callbacks.ObserveFn != nil {
		return c.Callbacks.Observe(ctx, res)
	}

	ares, err := api.Read(c.provider, c.Invoker, res)
	if err != nil {
		if err == api.ErrNotFound {
			return managed.ExternalObservation{}, nil
		}
		return managed.ExternalObservation{}, err
	}

	description, err := c.Invoker.MergeResources(res, ares)
	if err != nil {
		return managed.ExternalObservation{}, err
	}
	if description.AnnotationsUpdated || description.LateInitializedSpec {
		if err := c.KubeClient.Update(ctx, res); err != nil {
			return managed.ExternalObservation{}, err
		}
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: !description.NeedsProviderUpdate,
		// ConnectionDetails: getConnectionDetails(cr, instance),
	}, nil
}

func (c *External) Create(ctx context.Context, res resource.Managed) (managed.ExternalCreation, error) {
	c.entryLog(res, "Create")
	if c.Callbacks.CreateFn != nil {
		return c.Callbacks.Create(ctx, res)
	}

	created, err := api.Create(c.provider, c.Invoker, res)
	if err != nil {
		return managed.ExternalCreation{}, err
	}

	description, err := c.Invoker.MergeResources(res, created)
	if err != nil {
		return managed.ExternalCreation{}, err
	}
	if description.AnnotationsUpdated {
		if err = c.KubeClient.Update(ctx, res); err != nil {
			return managed.ExternalCreation{}, err
		}
	}
	return managed.ExternalCreation{}, nil
}

func (c *External) Update(ctx context.Context, res resource.Managed) (managed.ExternalUpdate, error) {
	c.entryLog(res, "Update")
	if c.Callbacks.UpdateFn != nil {
		return c.Callbacks.Update(ctx, res)
	}

	updated, err := api.Update(c.provider, c.Invoker, res)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}
	description, err := c.Invoker.MergeResources(res, updated)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}
	if description.AnnotationsUpdated || description.LateInitializedSpec {
		if err := c.KubeClient.Update(ctx, res); err != nil {
			return managed.ExternalUpdate{}, err
		}
	}

	return managed.ExternalUpdate{}, nil
}

func (c *External) Delete(ctx context.Context, res resource.Managed) error {
	c.entryLog(res, "Delete")
	if c.Callbacks.DeleteFn != nil {
		return c.Callbacks.Delete(ctx, res)
	}

	return api.Delete(c.provider, c.Invoker, res)
}

func (c *External) entryLog(res resource.Managed, method string) {
	gvk := res.GetObjectKind().GroupVersionKind()
	c.logger.Debug(fmt.Sprintf("terraform.External.%s: %s", method, gvk.String()))
}
