# envtest-start

`envtest-start` is a cli written in Go which starts an "envtest" environment (an API server + etcd) and creates a kubeconfig file
named `envtest-kubeconfig` in the TMP directory of the system.


## Requisites

Install `setup-envtest`, install a version of envtest and set in `KUBEBUILDER_ASSETS` the path in which it is installed:

```
go install sigs.k8s.io/controller-runtime/tools/setup-envtest@release-0.22
export KUBEBUILDER_ASSETS=$(setup-envtest use -p path)
```
