cluster=management
calico_cluster=calico-cluster
antrea_cluster=antrea-cluster
cilium_cluster=cilium-cluster

function create_kind_cluster() {
    kind create cluster --name=${cluster} --config=${conf}
    export KUBECONFIG="$(kind get kubeconfig-path --name=kind-${cluster})"
    chmod 755 ~/.kube/kind-config-kind
    export KUBECONFIG="$(kind get kubeconfig-path --name=kind-${cluster})"
    until kubectl cluster-info;  do
        echo "`date` waiting for cluster..."
        sleep 2
    done
}

function initialize_mgmt_with_clusterctl() {
    if ! clusterctl version; then
        curl -L https://github.com/kubernetes-sigs/cluster-api/releases/download/v0.3.7/clusterctl-darwin-amd64 -o clusterctl
        chmod +x ./clusterctl
        sudo mv ./clusterctl /usr/local/bin/clusterctl
    fi
    source vsphere_vars.sh
    clusterctl init --core cluster-api --bootstrap kubeadm --control-plane kubeadm --infrastructure vsphere

    until [[ `kubectl get pods --field-selector=status.phase=Running -n capi-webhook-system |  wc -l` == "       5" ]] &&\
    [[ `kubectl get pods --field-selector=status.phase=Running -n capv-system |  wc -l` == "       2" ]] &&\
    [[ `kubectl get pods --field-selector=status.phase=Running -n capi-system |  wc -l` == "       2" ]];  do
        echo "`date` waiting for management cluster to initialize..."
        sleep 2
    done
}

function cleanup(){
    kubectl delete -f ${calico_cluster}.yaml
    kubectl delete -f ${antrea_cluster}.yaml
    kubectl delete -f ${cilium_cluster}.yaml
     until [[ `kubectl get clusters -A | wc -l` == "       0" ]];  do
        echo "`date` waiting for cleaning up all workload clusters..."
        sleep 2
    done

    kind delete cluster --name ${cluster}
}

function create_workload_cluster_with_cni(){
    clusterctl config cluster $1 --kubernetes-version v1.18.2 --control-plane-machine-count=1 --worker-machine-count=1 > $1.yaml
    kubectl apply -f $1.yaml

    until [[ `kubectl get cluster $1 -o jsonpath='{.status.phase}'` == "Provisioned" ]]; do
        echo "`date` waiting for workload cluster $1 to be provisioned..."
        sleep 5
    done

    data="false"
    until [[ ! ${data} == "false" ]]; do
        echo "`date` waiting for workload cluster $1 to be reachable..."
        kubectl get secret  $1-kubeconfig -o=jsonpath='{.data.value}'  | { base64 -d 2>/dev/null || base64 -D; } >$1_kubeconfig
        if test -f "$1_kubeconfig"; then
            data="true"
        fi
        sleep 2
    done

    until kubectl cluster-info --kubeconfig=$1_kubeconfig;  do
        echo "`date` waiting for $1 workload cluster to start running now"
        sleep 2
    done

    kubectl apply -f cni/$2.yaml --kubeconfig=$1_kubeconfig
}

function create_all_cni_clusters() {
    create_workload_cluster_with_cni $calico_cluster calico
    create_workload_cluster_with_cni $antrea_cluster antrea
    create_workload_cluster_with_cni $cilium_cluster cilium

    kubectl -n kube-system set env daemonset/calico-node FELIX_IGNORELOOSERPF=true --kubeconfig=${calico_cluster}_kubeconfig
    kubectl -n kube-system set env daemonset/calico-node FELIX_XDPENABLED=false --kubeconfig=${calico_cluster}_kubeconfig
}

cleanup
create_kind_cluster
initialize_mgmt_with_clusterctl
create_all_cni_clusters