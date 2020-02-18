# Test netpolicies with truth tables!

```
kind create cluster
kind get kubeconfig > ~/.kube/config
go build ./main
./main
```

# Alternate setup

```
kind create cluster --name netpols
kubectl cluster-info --context kind-netpols

git clone git@github.com:jayunit100/k8sprototypes.git 
cd k8sprototypes/netpol
go run pkg/main/main.go
```