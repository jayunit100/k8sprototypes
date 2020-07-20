### Spin up 3 workload clusters with different CNIs using clusterctl and CAPV

- Set variables under ./scripts/vsphere_vars.sh
- Run ./install.sh

### To reach clusters:

#### Management cluster:

kubectl get pods -A

#### Workload clusters:

kubectl get nodes -A --kubeconfig=antrea-cluster_kubeconfig
kubectl get nodes -A --kubeconfig=calico-cluster_kubeconfig
kubectl get nodes -A --kubeconfig=cilium-cluster_kubeconfig

