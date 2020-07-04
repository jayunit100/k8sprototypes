#!/bin/bash

# Thanks to https://alexbrand.dev/post/creating-a-kind-cluster-with-calico-networking/ for this snippet :)
cat << EOF > kind-conf.yaml
kind: Cluster
apiVersion: kind.sigs.k8s.io/v1alpha3
networking:
  disableDefaultCNI: true # disable kindnet
  podSubnet: 192.168.0.0/16 # set to Calico's default subnet
nodes:
- role: control-plane
- role: worker
EOF

cat << EOF > kind-conf-ipv6.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  disableDefaultCNI: true
  ipFamily: ipv6
EOF

#cluster=cipv6
#conf=kind-conf-ipv6.yaml

#cluster=antrea
#conf=kind-conf

cluster=calico
conf=kind-conf.yaml

function install_k8s() {
    if kind delete cluster --name=${cluster}; then
    	echo "deleted old kind cluster, creating a new one..."
    fi	    
    kind create cluster --name=${cluster} --config=${conf} 
    export KUBECONFIG="$(kind get kubeconfig-path --name=kind-${cluster})"
    chmod 755 ~/.kube/kind-config-kind
    export KUBECONFIG="$(kind get kubeconfig-path --name=kind-${cluster})"
    until kubectl cluster-info;  do
        echo "`date`waiting for cluster..."
        sleep 2
    done
}

function install_antrea() {
   kubectl apply -f kubectl apply -f https://github.com/vmware-tanzu/antrea/releases/download/v0.8.0/antrea.yml -n kube-system  
}

function install_calico() {
    kubectl get pods
    kubectl apply -f ./calico312.yaml
}

function install_calico() {
    kubectl get pods
    kubectl apply -f ./calico.yaml
    kubectl get pods -n kube-system
    
    kubectl -n kube-system set env daemonset/calico-node FELIX_IGNORELOOSERPF=true
    kubectl -n kube-system set env daemonset/calico-node FELIX_XDPENABLED=false
}

function install_antrea() {
	kubectl create ns tkg-system
	kubectl create -f antrea.yml -n tkg-system
}

function wait() {
    sleep 5 ; kubectl -n kube-system get pods 
    echo "will wait for calico/antrea/... to start running now... "
    while true ; do
        kubectl -n kube-system get pods
        sleep 3
    done
}

function testStatefulSets() {
   sonobuoy run --e2e-focus "Basic StatefulSet" --e2e-skip ""   
}

install_k8s

if [[ $cluster == "calico" ]]; then
	install_calico
fi
if [[ $cluster == "antrea" ]]; then
	install_antrea
fi

#testStatefulSets
