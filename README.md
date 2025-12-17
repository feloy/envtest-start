# envtest-start

`envtest-start` is a cli written in Go which starts an "envtest" environment (an API server + etcd) and creates a kubeconfig file
named `envtest-kubeconfig` in the TMP directory of the system, or in the file passed as first argument of the command.

```bash
go build -o envtest-start main.go 

./envtest-start /path/to/kubeconfig
Kubeconfig written to: /path/to/kubeconfig
You can use it with: export KUBECONFIG=/path/to/kubeconfig
```

## Prerequisites

Install `setup-envtest`, install a version of envtest and set in `KUBEBUILDER_ASSETS` the path in which it is installed:

```bash
go install sigs.k8s.io/controller-runtime/tools/setup-envtest@release-0.22
export KUBEBUILDER_ASSETS=$(setup-envtest use -p path)
```
