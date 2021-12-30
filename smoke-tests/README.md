To start off, run this 

kubectl create -f https://github.com/jayunit100/k8sprototypes/blob/master/smoke-tests/nginx-pod-svc.yaml


...that will give you the baseline for this experiment... (note, alot of the mods described here you get for free when running that yaml, so read it first before you bother making these YAML snippets, these are really just for walkign you through the important parts of the YAML above)

Now, make a headless Service ... 

```
spec:
  clusterIP: None
  clusterIPs:
  - None
  ipFamilies:
  - IPv4
  ipFamilyPolicy: SingleStack
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: nginx
```


What happens if you curl  a headless service that isnt ready ?  Lets see.  First lets modify our a webserver so that it  has a readinessProbe like this:
```
    spec:
      containers:
      - image: nginx:1.14.2
        imagePullPolicy: IfNotPresent
        name: nginx
        ports:
        - containerPort: 80
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /index.html
            port: 80
            scheme: HTTP
```


 OKAY now ... lets look at coredns!  We can enable logs like so:



Now, with logging enabled ... I can exec into the busybox pod that was created in the original YAML.  Then you can run 
wget headless-svc and if you look in the coredns logs, you'll see:  

```
[INFO] 192.168.88.2:40504 - 2 "AAAA IN headless-svc.default.svc.cluster.local. udp 56 false 512" NOERROR qr,aa,rd 149 0.00022089s

[INFO] 192.168.88.2:40062 - 3 "A IN headless-svc.default.svc.cluster.local. udp 56 false 512" NOERROR qr,aa,rd 164 0.000156625s
```

Okay !  So coredns is very very happy.  Now, what happens if i UN-READINESS the endpoint for the headless service???  i.e. what if i make the httpGet intentionally fail???





NOW make sure and roll out your deployment after modifying the readiness probe ! 
```
-> % kubectl scale deployment nginx-deployment --replicas=0; kubectl scale deployment nginx-deployment --replicas=2
```
OK ! Now lets check coredns again:

its very... not happy.... lots of NXDOMAIN errors...  because theres no endpoint to back the headless services...

In general , you can solve this by serving non ready endpoints, though:

This will basically say "hey, i know the readiness checks fail, thats ok , lets publish these anyways" and then, coredns will not return NXDOMAIN errors anymore.


NOTE

If you want to understand the state of the endpoints, you can just continuously print them out like so:
```
kubectl get -o yaml endpointslices | grep -i read                                                     
      ready: true
      ready: true
      ready: true
      ready: true
      ready: true
      ready: true
      ready: true
```

This is ultimately what your coredns is using to decide wether or not an NXDOMAIN will be returned or not....
