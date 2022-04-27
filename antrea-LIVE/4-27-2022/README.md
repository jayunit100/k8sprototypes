
# Antrea LIVE Episode 26 etcd and CRDs and stuff


## THE FIRST ANT LIVE CODING CHALLENGE EVER

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
