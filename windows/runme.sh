source ~/env.sh
if [[ ! -v VSPHERE_CP_IP ]]; then
        echo "NO ENV SET FOR CONTROL PLANE MAKING A DEFAULT IN THE 12 SUBNET (this will work on corgi nsx jobs) !!!"
        VSPHERE_CP_IP=12.10.10.205
fi
wget https://build-artifactory.eng.vmware.com/artifactory/webapp/#/artifacts/browse/tree/General/k8simages-windows-local/windows-2019-kube-v1.19.1-containerd-nohyper.ova

govc import.ova https://build-artifactory.eng.vmware.com/k8simages-windows-local/windows-2019-kube-v1.19.1-containerd-nohyper.ova
govc object.rename windows-2019-kube-v1.19.1 windows-2019-kube-v1.19.1-containerd
govc vm.change -vm windows-2019-kube-v1.19.1-containerd -nested-hv-enabled=false 

kubectl delete -f vspheremachinetemplates.crd.yaml
kubectl create -f vspheremachinetemplates.crd.yaml

## wget peri-min.yaml.sh

cat peri-min.yaml.sh | sed s/\$VSPHERE_USERNAME/administrator@vsphere.local/g > peri-min.yaml

# Set image to run gabi's capv controller
# 1/28, totally optional since we dont really do any special logic in it anymore :)
kubectl set image deployment/capv-controller-manager manager=harbor-repo.vmware.com/tkgwindows/vsphere-manager:dev -n capv-system
kubectl set image deployment/capv-controller-manager manager=harbor-repo.vmware.com/tkgwindows/vsphere-manager:dev -n capi-webhook-system


sed -i s/\$VSPHERE_PASSWORD/$VSPHERE_PASSWORD/  peri-min.yaml
sed -i s/\$VSPHERE_DATACENTER/$VSPHERE_DATACENTER/  peri-min.yaml

# total hack
if [[ "$VSPHERE_NETWORK" == "VM Network" ]] ;  then
        VSPHERE_NETWORK='VM\ Network'
fi
echo $VSPHERE_NETWORK
sed -i s/\$VSPHERE_NETWORK/"$VSPHERE_NETWORK"/  peri-min.yaml
sed -i s/\$VSPHERE_SERVER/$VSPHERE_SERVER/  peri-min.yaml

sed -i s/\$VSPHERE_FOLDER/$VSPHERE_FOLDER/  peri-min.yaml

sed -i s/\$VSPHERE_CP_IP/$VSPHERE_CP_IP/  peri-min.yaml

KEYYYY=$(printf %q "$VSPHERE_SSH_AUTHORIZED_KEY")
echo $KEYYYY

sed -i 's/\$VSPHERE_SSH_AUTHORIZED_KEY/\"ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA6VfBKd6hqd5h7k5f+AtjJSV1hdW5u9\/3uAolK3SD2\/5GD9+rn+FMSdbtkeaKuuVJPi2HjnsVMO+r8WcuyN5ZSYHywiSoh4S7PamAxra1CLISsFHPYFlGrtdHC70wnoT7+\/wAJk2D3CYkCNMWIxs5eR0cefDOytipBfDplhkJByyrcnXuhI8St3XJzpjlXu454diJOxfsk6axanWLOr\/WZFmUi1U6V4gRE7XtKG9WFUm1bmNgkgd7lehKzi+isTjnI+b4tnD0yIzKFcsgIvLdGJTI6Lluj33CeBHIocwu0LbvowTyYSqhP6DzGhGuKfK9rMnJh\/ll0Bnu1xf\/ok0NSQ== tkg\"/' peri-min.yaml

sed -i s/\$VSPHERE_DATASTORE/$VSPHERE_DATASTORE/  peri-min.yaml

echo "done substituting !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"

#peri-min.yaml

kubectl delete -f peri-min.yaml
sleep 60
kubectl create -f peri-min.yaml
