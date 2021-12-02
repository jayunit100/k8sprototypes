```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: iis
  labels:
    app: iis
spec:
  replicas: 1
  selector:
    matchLabels:
      app: iis
  template:
    metadata:
      labels:
        app: iis
    spec:
      nodeSelector:
        kubernetes.io/os: windows
      containers:
      - name: iis-server
        image: mcr.microsoft.com/windows/servercore/iis
        ports:
        - containerPort: 80
        volumeMounts:
        - mountPath: C:/k/
          name: stuey
      volumes:
      - name: stuey
        hostPath:
          path: C:/k/
          type: Directory
```

- Then jump into a shell on the container

kubectl exec -t -i iis-54f48785f7-45szv powershell

- Now , try to run a copy command in side the container you made

- Now test that the write persisted into the node:

scp capv@10.187.151.211:C:/k/stuey.png stuey.png
