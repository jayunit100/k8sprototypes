# Work in progress

Just noting things here as I confirm that this workflow is actually valid. YMMV 

## Start csi-proxy 

## SMB

From https://github.com/kubernetes-csi/csi-driver-smb/tree/master/deploy/example/smb-provisioner... Create a smb server in a pod: 

```
  kubectl create secret generic smbcreds --from-literal username=USERNAME --from-literal password="PASSWORD"
  kubectl create -f https://raw.githubusercontent.com/kubernetes-csi/csi-driver-smb/master/deploy/example/smb-provisioner/smb-server.yaml 
``` 

Note csi-proy might not be able to access that via clusterIP ^ .  Need to make sure of this.

## CSI Proxy 

After that, https://github.com/kubernetes-csi/csi-proxy#setup-for-csi-driver-deployment shows how to setup CSI-proxy

## CSI driver SMB

Now that the CSI proxy is up and running, the csi driver for SMB can be installed:

https://github.com/kubernetes-csi/csi-driver-smb

Which should use the CSI wins proxy to do the various privileged mounting operations
