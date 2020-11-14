#!/usr/bin/env bash

set -e

CLUSTER=${CLUSTER:-calico}
CONFIG=${CONFIG:-calico-conf.yaml}

# Usage examples:
#CLUSTER=nocni CONFIG="calico-conf.yaml" ./kind-local-up.sh
#CLUSTER=cipv6 CONFIG=kind-conf-ipv6.yaml ./kind-local-up.sh
# Calico usage - CLUSTER=calico CONFIG=calico-conf.yaml ./kind-local-up.sh
# Antrea usage - CLUSTER=antrea CONFIG=calico-conf.yaml ./kind-local-up.sh
# Cilium usage - CLUSTER=cilium CONFIG=cilium-conf.yaml ./kind-local-up.sh

function check_kind() {
    if [ kind > /dev/null ]; then
        echo "Kind binary not found."
        exit 1
    fi
}

function init_configuration() {
    # Thanks to https://alexbrand.dev/post/creating-a-kind-cluster-with-calico-networking/ for this snippet :)
    cat << EOF > calico-conf.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  disableDefaultCNI: true # disable kindnet
  podSubnet: 192.168.0.0/16 # set to Calico's default subnet
nodes:
- role: control-plane
- role: worker
- role: worker
EOF

    cat << EOF > cilium-conf.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
- role: worker
networking:
  disableDefaultCNI: true
EOF

#cluster=nocni
#conf="calico-conf.yaml"

#cluster=cipv6
#conf=kind-conf-ipv6.yaml

## calico conf == no cni, so use it for antrea/calico/whatever
cluster=antrea
conf=calico-conf.yaml

# cluster=calico
# conf=calico-conf.yaml

cluster=calico
conf=calico-conf.yaml

function install_k8s() {
    if kind delete cluster --name=${CLUSTER}; then
    	echo "deleted old kind cluster, creating a new one..."
    fi
    kind create cluster --name=${CLUSTER} --config=${CONFIG}
    until kubectl cluster-info;  do
        echo "`date` waiting for cluster..."
        sleep 2
    done
}

function install_calico() {
    kubectl get pods
    kubectl apply -f ./calico312.yaml
    kubectl get pods -n kube-system
    kubectl -n kube-system set env daemonset/calico-node FELIX_IGNORELOOSERPF=true
    kubectl -n kube-system set env daemonset/calico-node FELIX_XDPENABLED=false
}

function install_antrea() {
	if [[ ! -d antrea ]] ; then
	    git clone https://github.com/vmware-tanzu/antrea.git	
	fi
	pushd antrea
	     git checkout v0.9.0
	     pushd ci/kind
    	      ./kind-setup.sh create antrea
	     popd
	popd
}

function install_cilium() {
    CILIUM_VERSION="1.9.1"

    # Add Cilium Helm repo
    helm repo add cilium https://helm.cilium.io/

    # Pre-load images
    docker pull cilium/cilium:"v${CILIUM_VERSION}"

    # Install cilium with Helm
    helm install cilium cilium/cilium --version ${CILIUM_VERSION} \
         --namespace kube-system \
         --set nodeinit.enabled=true \
         --set kubeProxyReplacement=partial \
         --set hostServices.enabled=false \
         --set externalIPs.enabled=true \
         --set nodePort.enabled=true \
         --set hostPort.enabled=true \
         --set bpf.masquerade=false \
         --set image.pullPolicy=IfNotPresent \
         --set ipam.mode=kubernetes
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

init_configuration
sleep 1

case "$CLUSTER" in
    "antrea")
        echo "Using Antrea/master setup script for kind"
        install_antrea
        ;;
    "calico")
        echo "Using Calico CNI."
        install_k8s
        install_calico
        ;;
    "cilium")
        echo "Using Cilium CNI."
        install_k8s
        install_cilium
        ;;
    "*")
        echo "Skipping CNI"
        ;;
esac
