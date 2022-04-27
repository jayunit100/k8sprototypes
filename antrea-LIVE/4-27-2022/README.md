
# Antrea LIVE Episode 26 etcd and CRDs and stuff

getting started on adv. windowsy networking policy test improvements 
```https://github.com/kubernetes/kubernetes/issues/107853

..........
ubuntu@jay-build-box-6 [21:09:28] [~/SOURCE/kubernetes/test/e2e/network/netpol] [31dba0a435e *]
-> % ### /home/ubuntu/SOURCE/kubernetes/e2e.test --provider=local --kubeconfig=/home/ubuntu/.kube/config --dump-logs-on-failure=false --ginkgo.focus="sig-network" --ginkgo.skip="Driver|Slow|Driver" --ginkgo.dryRun=true
ubuntu@jay-build-box-6 [21:10:11] [~/SOURCE/kubernetes/test/e2e/network/netpol] [31dba0a435e *]
-> % make WHAT=test/e2e/e2e.test
make: *** No targets specified and no makefile found.  Stop.
```

## THE SECOND ANT LIVE CODING CHALLENGE EVER

```
kubectl patch deployment patch-demo --patch "{\"spec\":{\"template\":{\"spec\":{\"containers\":[{\"name\":\"patch-demo-ctr-${i}\",\"image\":\"redis\"}]}}}}"
``` 
- randomly create 100 deployments
- put this into a pod w/ a kubectl, wrap it in bacsh and see if you can use it to 
contintuously patch a those deployment objects
- figure out a way to compare the object as seen in the script over time, depending on what node in the
daemonset was running a get on that resource that is constantly being mutated
- see if things get out of sync....


```
--------------> wether things get out of sync -------------> 
etcdctl --endpoints="https://localhost:2379" --cacert="/etc/kubernetes/pki/etcd/ca.crt" --cert="/etc/kubernetes/pki/etcd/server.crt" --key=/etc/kubernetes/pki/etcd/server.key endpoint status --cluster
``` 


## THE FIRST ANT LIVE CODING CHALLENGE EVER
test snippet

```
-> % for i in {1..10} ; do NAME=$i ./crd-race.sh  ; done
```


MAKE THIS REPRODOCE https://github.com/kubernetes/kubernetes/issues/65517 
```
#!/bin/bash

#NAME=yyy

kubectl apply -f <(cat << EOF
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foo${NAME}s.stable.example.com
spec:
  group: stable.example.com
  versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
  scope: Namespaced
  names:
    plural: foo${NAME}s
    singular: foo${NAME}
    # kind is normally the CamelCased singular type. Your resource manifests use this.
    kind: Foo${NAME}
    shortNames:
    - foo${NAME}s
EOF
)


kubectl get foo${NAME}s

echo 2

kubectl get foo${NAME}s

echo 3
kubectl get foo${NAME}s

kubectl delete CustomResourceDefinition foo${NAME}s.stable.example.com
```



# Antrea LIVE Episode 24 etcd , ricardo, ...

https://github.com/ahrtr/etcd-issues/tree/master/issues/13766

```
while true ; do etcdctl --endpoints="https://localhost:2379" --cacert="/etc/kubernetes/pki/etcd/ca.crt" --cert="/etc/kubernetes/pki/etcd/server.crt" --key=/etc/kubernetes/pki/etcd/server.key endpoint status --cluster ; sleep 1 ; done
```

perf

```
etcdctl --endpoints="https://localhost:2379" --cacert="/etc/kubernetes/pki/etcd/ca.crt" --cert="/etc/kubernetes/pki/etcd/server.crt" --key=/etc/kubernetes/pki/etcd/server.key check perf --load=xl
```
