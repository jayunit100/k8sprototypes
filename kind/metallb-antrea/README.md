# Does antrea work with Metallb

Yes...  the above script can be used with the `cluster=ANTREA kind-local-up.sh` script...
```
Every 2.0s: kubectl get pods -A                                                                                                                 jay-build-box-6: Sun Jun 26 19:32:29 2022

NAMESPACE            NAME                                           READY   STATUS    RESTARTS   AGE
default              bar-app                                        1/1     Running   0          3m38s
default              foo-app                                        1/1     Running   0          3m38s
kube-system          antrea-agent-4679h                             2/2     Running   0          26m
kube-system          antrea-agent-7d5zg                             2/2     Running   0          26m
kube-system          antrea-agent-856m2                             2/2     Running   0          26m
kube-system          antrea-controller-5fdb4d7c8d-j8rbq             1/1     Running   0          26m
kube-system          coredns-b4b5969d4-krmgv                        1/1     Running   0          26m
kube-system          coredns-b4b5969d4-r82j6                        1/1     Running   0          26m
kube-system          etcd-antrea-control-plane                      1/1     Running   0          27m
kube-system          kube-apiserver-antrea-control-plane            1/1     Running   0          27m
```

and metallb version 12.1...

```
-> % kubectl get pods -A | grep metal
metallb-system       controller-7476b58756-j6fxw                    1/1     Running   0          27m
metallb-system       speaker-2hkj2                                  1/1     Running   0          26m
metallb-system       speaker-dlnh2                                  1/1     Running   0          26m
```
As we can see here

```
-> % curl 172.18.255.200:5678
foo
ubuntu@jay-build-box-6 [19:34:43] [~/SOURCE/k8sprototypes/kind/metallb] [2d472d6 *]
-> % kubectl exec -t -i antrea-agent-4679h -n kube-system  antctl version
kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl kubectl exec [POD] -- [COMMAND] instead.
Defaulting container name to antrea-agent.
Use 'kubectl describe pod/antrea-agent-4679h -n kube-system' to see all of the containers in this pod.
agentVersion: 1.4.0-bd88118.dirty
antctlVersion: v1.4.0
```


# What about on antrea 1.5? BGP? 

Yes, this works as well contact scott@terasky.com for details !







