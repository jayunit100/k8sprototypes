#!/bin/bash

source ~/env.sh
if [[ "$VSPHERE_CP_IP" == "" ]]; then
	echo "missinc VSPHERE_CP_IP"
        exit 1
fi
#wget https://build-artifactory.eng.vmware.com/artifactory/webapp/#/artifacts/browse/tree/General/k8simages-windows-local/windows-2019-kube-v1.19.1-containerd-nohyper.ova
# govc import.ova https://build-artifactory.eng.vmware.com/k8simages-windows-local/windows-2019-kube-v1.19.1-containerd-nohyper.ova
govc object.rename /SDDC-Datacenter/vm/Templates/windows-2019-kube-v1.19.1 /SDDC-Datacenter/vm/Templates/windows-2019-kube-v1.19.1-containerd 

### make the vm hyper0v independent...
govc vm.change -vm /SDDC-Datacenter/vm/Templates/windows-2019-kube-v1.19.1-containerd/windows-2019-kube-v1.19.1-containerd -nested-hv-enabled=false 
govc vm.change -vm /SDDC-Datacenter/vm/Workloads/tkg-infra/users/jvyas/photon-3-kube-v1.19.1+vmware.2-take3_WINDOWS -nested-hv-enabled=false

kubectl delete -f vspheremachinetemplates.crd.yaml
kubectl create -f vspheremachinetemplates.crd.yaml

## wget peri-min.yaml.sh

cat peri-min.yaml.sh | sed s/xVSPHERE_USERNAME/jvyas@vmware.ci/g > peri-min.yaml

# Set image to run gabi's capv controller
# 1/28, totally optional since we dont really do any special logic in it anymore :)
#kubectl set image deployment/capv-controller-manager manager=harbor-repo.vmware.com/tkgwindows/vsphere-manager:dev -n capv-system
#kubectl set image deployment/capv-controller-manager manager=harbor-repo.vmware.com/tkgwindows/vsphere-manager:dev -n capi-webhook-system


sed -i s/xVSPHERE_PASSWORD/$VSPHERE_PASSWORD/  peri-min.yaml
sed -i s/xVSPHERE_DATACENTER/$VSPHERE_DATACENTER/  peri-min.yaml
sed -i s/xVSPHERE_MACHINE_TEMPLATE/$VSPHERE_MACHINE_TEMPLATE/  peri-min.yaml
sed -i s/xVSPHERE_MACHINE_TEMPLATE_WINDOWS/$VSPHERE_MACHINE_TEMPLATE_WINDOWS/ peri-min.yaml
sed -i s/xVSPHERE_RESOURCE_POOL/$VSPHERE_RESOURCE_POOL/ peri-min.yaml
sed -i s/xVSPHERE_NETWORK/$VSPHERE_NETWORK/ peri-min.yaml

# total hack
if [[ "$VSPHERE_NETWORK" == "VM Network" ]] ;  then
        VSPHERE_NETWORK='VM\ Network'
fi
echo $VSPHERE_NETWORK
sed -i s/xVSPHERE_NETWORK/"$VSPHERE_NETWORK"/  peri-min.yaml
sed -i s/xVSPHERE_SERVER/$VSPHERE_SERVER/  peri-min.yaml

sed -i s/xVSPHERE_FOLDER/$VSPHERE_FOLDER/  peri-min.yaml

sed -i s/xVSPHERE_CP_IP/$VSPHERE_CP_IP/  peri-min.yaml

KEYYYY=$(printf %q "$VSPHERE_SSH_AUTHORIZED_KEY")
echo $KEYYYY

sed -i 's/\xVSPHERE_SSH_AUTHORIZED_KEY/\"ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA6VfBKd6hqd5h7k5f+AtjJSV1hdW5u9\/3uAolK3SD2\/5GD9+rn+FMSdbtkeaKuuVJPi2HjnsVMO+r8WcuyN5ZSYHywiSoh4S7PamAxra1CLISsFHPYFlGrtdHC70wnoT7+\/wAJk2D3CYkCNMWIxs5eR0cefDOytipBfDplhkJByyrcnXuhI8St3XJzpjlXu454diJOxfsk6axanWLOr\/WZFmUi1U6V4gRE7XtKG9WFUm1bmNgkgd7lehKzi+isTjnI+b4tnD0yIzKFcsgIvLdGJTI6Lluj33CeBHIocwu0LbvowTyYSqhP6DzGhGuKfK9rMnJh\/ll0Bnu1xf\/ok0NSQ== tkg\"/' peri-min.yaml

sed -i s/xVSPHERE_DATASTORE/$VSPHERE_DATASTORE/  peri-min.yaml

echo "done substituting !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"

#peri-min.yaml
# FOR DEV UNCOMMENT NEXT 2 OLINES
kubectl delete -f peri-min.yaml
sleep 60
kubectl create -f peri-min.yaml
