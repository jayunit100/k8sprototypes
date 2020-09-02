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
- role: worker
- role: worker
EOF

cat << EOF > kind-conf-ipv6.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  disableDefaultCNI: true
  ipFamily: ipv6
EOF

#cluster=nocni
#conf="calico-conf.yaml"

#cluster=cipv6
#conf=kind-conf-ipv6.yaml

## calico conf == no cni, so use it for antrea/calico/whatever
cluster=antrea
conf=calico-conf.yaml

#cluster=calico
#conf=calico-conf.yaml

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

function install_calico() {
    kubectl get pods
    kubectl apply -f ./calico3.12.3.yaml
    kubectl get pods -n kube-system
    kubectl -n kube-system set env daemonset/calico-node FELIX_IGNORELOOSERPF=true
    kubectl -n kube-system set env daemonset/calico-node FELIX_XDPENABLED=false
}

function install_antrea() {
	if [[ ! -d antrea ]] ; then
	    git clone https://github.com/vmware-tanzu/antrea.git
	fi
	pushd antrea/ci/kind
    	    ./kind-setup.sh create antrea
	popd
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

if [[ $cluster == "antrea" ]] ; then 
	echo "using antrea/master setup script for kind"
	sleep 1
	install_antrea
fi

if [[ $cluster == "" ]]; then
	echo "skipping cni"
fi
if [[ $cluster == "calico" ]]; then
	install_k8s
	install_calico
fi

