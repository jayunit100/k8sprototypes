

# Enabling antrea in Tanzu on 1.2.3

Edit the configMap to use NodePortLocal... 
```
  apiVersion: v1
  data:
    antrea-agent.conf: "featureGates:\n  AntreaProxy: true\n  EndpointSlice: false\n
      \ Traceflow: true\n  NodePortLocal: true\n  AntreaPolicy: true\n  FlowExporter:
      false\n  NetworkPolicyStats: false\n  Egress: false\n  \ntrafficEncapMode: encap\nnoSNAT:
      false\nserviceCIDR: 100.64.0.0/13\ntlsCipherSuites: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384\n"
    antrea-cni.conflist: |
      {
          "cniVersion":"0.3.0",
          "name": "antrea",
          "plugins": [
              {
                  "type": "antrea",
                  "ipam": {
                      "type": "host-local"
                  }
              },
              {
                  "type": "portmap",
                  "capabilities": {"portMappings": true}
              },
              {
                  "type": "bandwidth",
                  "capabilities": {"bandwidth": true}
              }
          ]
      }
    antrea-controller.conf: |
      featureGates:
        Traceflow: true
        AntreaPolicy: true
        NetworkPolicyStats: false
        Egress: true
        NodePortLocal: true
```

# Now delete all your antrea pods

```
kubectl delete pod -n kube-ststem `kubectl get pods -n kube-system | grep antrea | cut -d' ' -f 1` -n kube-system
```

# Now make a nodeport local service to test:



```
    apiVersion: v1
    kind: Service
    metadata:
      name: nginx
      annotations:
        nodeportlocal.antrea.io/enabled: "true"
    spec:
      ports:
      - name: web
        port: 80
        protocol: TCP
        targetPort: 80
      selector:
        app: nginx
    ---
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: nginx
    spec:
      selector:
        matchLabels:
          app: nginx
      replicas: 3
      template:
        metadata:
          labels:
            app: nginx
          annotations:
            warning.omg.antrea.io.omg : "this nodeport below gets overrwitten... just FYI :) :) :) :)"
            nodeportlocal.antrea.io: '[{"podPort":8080,"nodeIP":"10.186.58.200","nodePort":61234}]'
        spec:
          containers:
          - name: nginx
            image: gcr.io/google-containers/nginx
```

# Note: Antrea rewrites the nodeportlocal, but whatever...
```
apiVersion: v1
kind: Pod
metadata:
  annotations:
    ### Reminder that this nodeport was overwritten from the original value
    nodeportlocal.antrea.io: '[{"podPort":80,"nodeIP":"10.186.56.87","nodePort":61000}]'
  creationTimestamp: "2022-01-29T15:20:27Z"
```

# Now look at the pods

```
kubo@jumper:~$ kubectl get pods -o yaml | grep nodeport
```
Reveals that antrea autoassigned monotonically increasing nodeport addresses... 
```
      nodeportlocal.antrea.io: '[{"podPort":80,"nodeIP":"10.186.56.87","nodePort":61003}]'
            f:nodeportlocal.antrea.io: {}
      nodeportlocal.antrea.io: '[{"podPort":80,"nodeIP":"10.186.56.87","nodePort":61000}]'
            f:nodeportlocal.antrea.io: {}
      nodeportlocal.antrea.io: '[{"podPort":80,"nodeIP":"61002"
```

Either way it works though:

```
kubo@jumper:~$ curl 10.186.56.87:61000 | head
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx on Debian!</title>
<style>
    body {
```
