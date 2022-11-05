#/bin/bash

# PURPOSE: Install tanzu cli so its locally callable, and upload all OVAs to vsphere
# INPUT: tanzu-framework tar gz, all OVAs
# OUTPUT: None

# initialize Tanzu

pushd payload

export GOVC_URL="10.1.2.3"
export GOVC_USERNAME="administrator@vsphere.local"
export GOVC_PASSWORD="Admin\!23"
export GOVC_DATACENTER="dc0"
export GOVC_INSECURE=true
export GOVC_NETWORK="VM Network"

#govc import.ova ./windows-2019-kube-v1.21.2+vmware.1-tkg.1.ova
#govc vm.markastemplate windows-2019-kube-v1.21.2

# substituted by tkgprotoform...
govc import.ova ./my+vmware.ova
govc vm.markastemplate my.ova

popd