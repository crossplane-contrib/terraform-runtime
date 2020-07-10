# terraform-provider-runtime


`terraform-provider-runtime` specifies a set of interfaces and dependency injection points for a code-generated provider to use to offer CRUD
methods to fulfill the crossplane ExternalClient API.

## Install

If you would like to install `terraform-provider-runtime` without modifications create
the following `ClusterPackageInstall` in a Kubernetes cluster where Crossplane is
installed:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: terraform-plugin
---
apiVersion: packages.crossplane.io/v1alpha1
kind: ClusterPackageInstall
metadata:
  name: terraform-provider-runtime
  namespace: terraform-plugin
spec:
  package: "crossplane/terraform-provider-runtime:latest"
```

## Developing

Run against a Kubernetes cluster:
```
make run
```

Install `latest` into Kubernetes cluster where Crossplane is installed:
```
make install
```

Install local build into [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
cluster where Crossplane is installed:
```
make install-local
```

Build, push, and install:
```
make all
```

Build image:
```
make image
```

Push image:
```
make push
```

Build binary:
```
make build
```

Build package:
```
make build-package
```
