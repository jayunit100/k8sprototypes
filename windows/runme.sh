source ~/env.sh
cat peri-min.yaml.sh | sed s/\$VSPHERE_USERNAME/administrator@vsphere.local/g > peri-min.yaml
echo "x"

sed -i s/\$VSPHERE_PASSWORD/$VSPHERE_PASSWORD/  peri-min.yaml
sed -i s/\$VSPHERE_DATACENTER/$VSPHERE_DATACENTER/  peri-min.yaml
sed -i s/\$VSPHERE_NETWORK/$VSPHERE_NETWORK/  peri-min.yaml
sed -i s/\$VSPHERE_SERVER/$VSPHERE_SERVER/  peri-min.yaml 
sed -i s/\$VSPHERE_FOLDER/$VSPHERE_FOLDER/  peri-min.yaml 
sed -i s/\$VSPHERE_DATASTORE/$VSPHERE_DATASTORE/  peri-min.yaml 

echo "done substituting !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"

#peri-min.yaml

