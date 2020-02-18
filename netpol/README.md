# Test netpolicies with truth tables!

Here are a few ways to try this out.

## Create 

Create the policy probe tests... 

```
kubectl create clusterrolebinding netpol --clusterrole=admin --serviceaccount=kube-system:netpol
kubectl create sa netpol -n kube-system
kubectl create -f https://raw.githubusercontent.com/jayunit100/k8sprototypes/master/netpol/install.yml
```

Now, look at tthe results of the network policy probe... 

```
 kubectl logs `kubectl get pods -n kube-system | grep netpol | cut -d' ' -f 1` -n kube-system  
```
 
## Create a cluster if you dont have one  and run from source....
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