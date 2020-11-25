source ~/env.sh
cat peri-min.yaml.sh | sed s/\$VSPHERE_USERNAME/administrator@vsphere.local/g > peri-min.yaml
echo "x"

sed -i s/\$VSPHERE_PASSWORD/$VSPHERE_PASSWORD/  peri-min.yaml
sed -i s/\$VSPHERE_DATACENTER/$VSPHERE_DATACENTER/  peri-min.yaml
sed -i s/\$VSPHERE_NETWORK/$VSPHERE_NETWORK/  peri-min.yaml
sed -i s/\$VSPHERE_SERVER/$VSPHERE_SERVER/  peri-min.yaml 
sed -i s/\$VSPHERE_FOLDER/$VSPHERE_FOLDER/  peri-min.yaml
KEYYYY=$(printf %q "$VSPHERE_SSH_AUTHORIZED_KEY")
echo $KEYYYY 



sed -i 's/\$VSPHERE_SSH_AUTHORIZED_KEY/\"ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA6VfBKd6hqd5h7k5f+AtjJSV1hdW5u9\/3uAolK3SD2\/5GD9+rn+FMSdbtkeaKuuVJPi2HjnsVMO+r8WcuyN5ZSYHywiSoh4S7PamAxra1CLISsFHPYFlGrtdHC70wnoT7+\/wAJk2D3CYkCNMWIxs5eR0cefDOytipBfDplhkJByyrcnXuhI8St3XJzpjlXu454diJOxfsk6axanWLOr\/WZFmUi1U6V4gRE7XtKG9WFUm1bmNgkgd7lehKzi+isTjnI+b4tnD0yIzKFcsgIvLdGJTI6Lluj33CeBHIocwu0LbvowTyYSqhP6DzGhGuKfK9rMnJh\/ll0Bnu1xf\/ok0NSQ== tkg\"/' peri-min.yaml 

sed -i s/\$VSPHERE_DATASTORE/$VSPHERE_DATASTORE/  peri-min.yaml 



echo "done substituting !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"

#peri-min.yaml

