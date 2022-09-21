# How to Run these ! 


```
CLUSTER=calico ./kind-local-up.sh
CLUSTER=antrea ./kind-local-up.sh
CLUSTER=cillium ./kind-local-up.sh

```
Run the kind-local-up.sh script.  It will create a kind cluster w/ a real CNI (antrea/calico/cillium) for you.  

##  antrea: 
Modify the `antrea` parameters (turn them on) and comment out the calico params

```
  # todo this should be cni=antrea
  cluster=antrea
  conf=xyz.yaml
  ./kind-local-up.sh
```

## Other CNIs
Just change the `cluster` variable to...
- cluster=calico
- cluster=cillium
..
