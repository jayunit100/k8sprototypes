Working Example:


Lets assume your kube-proxy can't talk to the APIServer, and your not sure why. 


Create two simple files in your host node, which runs the apiserver.  If your running the apiserver in a pod, then that looks like this:


# Tell your apiserver to audit everything 

Write this file to "/etc/auditperi/audit.yaml"
```
apiVersion: audit.k8s.io/v1
kind: Policy
rules:
- level: Metadata
```

# Launch your apiserver with this argument as input:

(you can edit a CAPI clusters apiserver by going into /etc/kubernetes/manifests on a master node)

```
spec:
  containers:
  - command:
    - kube-apiserver
    - --audit-policy-file=/etc/auditperi/audit.yaml # add this !!!
    - --audit-log-path=/etc/auditperi/audit.log # this is an arbitrary path on your container, make sure its host mounted though !  this should work since we already are mounting our /etc/auditperi directory
```
Now later in your pod spec for the apiserver, 

```
   
  volumes:
  - hostPath:
      path: /etc/auditperi # mount it !
    name: policy
```

Ok... not when your apiserver mounts, it will output all failed api calls to this file.  Then, you can run `audit2rbac` which will reverse engineer the exact objects you need to make to allow your kube proxy to come back online


## Finally you can run:
```
 ./audit2rbac -f /etc/auditperi/audit.log --serviceaccount=kube-system:kube-proxy
```


