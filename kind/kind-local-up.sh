#!/bin/bash

# Thanks to https://alexbrand.dev/post/creating-a-kind-cluster-with-calico-networking/ for this snippet :)
cat << EOF > calico-conf.yaml
kind: Cluster
apiVersion: kind.sigs.k8s.io/v1alpha3
networking:
  disableDefaultCNI: true # disable kindnet
  podSubnet: 192.168.0.0/16 # set to Calico's default subnet
nodes:
- role: control-plane
- role: worker
EOF

function install_k8s() {
    if kind delete cluster --name calico-test; then
    	echo "deleted old kind cluster, creating a new one..."
    fi	    
    kind create cluster --name calico-test --config calico-conf.yaml
    export KUBECONFIG="$(kind get kubeconfig-path --name=calico-test)"
    for i in "cni-plugin" "node" "pod2daemon" "kube-controllers"; do 
        echo "...$i"
    done
    chmod 755 ~/.kube/kind-config-kind
    export KUBECONFIG="$(kind get kubeconfig-path --name=calico-test)"
    until kubectl cluster-info;  do
        echo "`date`waiting for cluster..."
        sleep 2
    done
}

function install_calico() {
    kubectl get pods
    kubectl apply -f ./calico.yaml
    kubectl get pods -n kube-system
    
    kubectl -n kube-system set env daemonset/calico-node FELIX_IGNORELOOSERPF=true
    kubectl -n kube-system set env daemonset/calico-node FELIX_XDPENABLED=false
	
    sleep 5 ; kubectl -n kube-system get pods | grep calico-node
    echo "will wait for calico to start running now... "
    while true ; do
        kubectl -n kube-system get pods
        sleep 3
    done
}

install_k8s
install_calico
