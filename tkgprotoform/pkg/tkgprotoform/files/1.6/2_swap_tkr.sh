#! /bin/bash

# Input: arg $1 -> a file with a TKR definition in it 
#    (i.e. w/ OVAs similar to https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.6/vmware-tanzu-kubernetes-grid-16/GUID-build-images-linux.html )
# Output: A Configmap w/ a TKR in it that references 
bomTag=v1.23.8---vmware.2-tkg.1
name=v1.23.8---vmware.2-tkg.1-jay2

cat <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
 namespace: tkr-system
 annotations:
   bomImageTag: ${bomTag}
 name: ${name}
 labels:
   tanzuKubernetesRelease: ${name}
binaryData:
 # multiline binary content below...
 bomContent: |
EOF

base64 <$1 | sed 's,^,  ,'
