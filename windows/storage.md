# Work in progress

Just noting things here as I confirm that this workflow is actually valid. YMMV 

## Start csi-proxy 


## SMB

From https://github.com/kubernetes-csi/csi-driver-smb/tree/master/deploy/example/smb-provisioner

... Create a smb server in a pod: 

```
  kubectl create secret generic smbcreds --from-literal username=USERNAME --from-literal password="PASSWORD"
  kubectl create -f https://raw.githubusercontent.com/kubernetes-csi/csi-driver-smb/master/deploy/example/smb-provisioner/smb-server.yaml 
``` 

Note csi-proy might not be able to access that via clusterIP ^ .  Need to make sure of this.

## Install the CSI Driver:

https://github.com/kubernetes-csi/csi-driver-smb/blob/master/docs/install-csi-driver-v0.4.0.md

Now youll have a storage class for smb but you need to make it the default... 
```
  kubectl patch storageclass smb -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
```
Then create a PVC to test that it gets bound:
```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: myclaim
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Mi
```
And modify an app like iis to mount that volume:

```
   spec:
      containers:
      - image: gcr.io/dotnet-atamel/iis-site-windows
      ...
      volumes:
      - name: test-data
        persistentVolumeClaim:
          claimName: myclaim
```
## CSI Proxy 

After that, https://github.com/kubernetes-csi/csi-proxy#setup-for-csi-driver-deployment shows how to setup CSI-proxy

## CSI driver SMB

Now that the CSI proxy is up and running, the csi driver for SMB can be installed:

https://github.com/kubernetes-csi/csi-driver-smb

Which should use the CSI wins proxy to do the various privileged mounting operations

# Issues

Currently i seem to be hitting https://github.com/kubernetes-csi/csi-driver-smb/issues/150, wherein the windows mount operation fails 
because it cant find the moutn for the pvc which was created 
