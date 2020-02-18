# Test netpolicies with truth tables!

This repo implements https://github.com/vmware-tanzu/antrea/blob/community-network-policy-tests/docs/design/cni-testing-initiative-upstream.md, a fast, comprehensive truth table matrix for network policy's which can be used to ensure that you're CNI provider is fast, reliably, and air-tight.

## How is this different then NetworkPolicy tests in upstream K8s ?

We are working to merge this into upstream Kubernetes, in the meanwhile, here's the differences.

- We define tests as *truth tables, and have a 'builder' library* for building up network policy structs with almost no boilerplat, meaning you can define a very sophisticated network policy test in just a few lines of code.
- *Comprehensive:* All pod-to-pod connectivity is validated for every test run.  In a typical network policy test in current upstream we only validate 2 or 3 scenarios, leaving out intra and inner namespace connections which might be comprimised due to a hard to detect CNI inconsitency.  In these tests, we test all 81 connections for 3 identical pods running in 3 different namespaces (i.e. the 9x9 connectivity matrix).
- *Transparent:* Each test prints out a `kubectl ` command you can run to reprobe a given pods connectivty patterns.
- Its *fast:* Because we use `kubectl exec` to run tests with `wget` between pods, all 81 tests can easily finish with 20 seconds or less, even if pod scheduling is slow.  This is because no polling is done, and there is now down/uptime for pods.
- *Easy to reason about:* The pods in this repo stay up forever, so you can reuse the above kubectl commands outputted by your netpol logs to exec into a pod and reproduce any failures.
- *Scalable:* If you want to test 32 policies, all at once ? Just take a look at the example test (in `main`) and copy paste a few lines, and you'll be testing enterprise CNI application patterns in a heartbeat.

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
