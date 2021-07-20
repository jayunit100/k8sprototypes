# Interpretting netpol matrix tests 

Taken from 
https://storage.googleapis.com/kubernetes-jenkins/logs/ci-kubernetes-e2e-ubuntu-gce-network-policies/1417411869290795008/build-log.txt

Creating pods... 

```
I0720 10:17:14.944] Jul 20 10:00:56.142: INFO: Using cluster.local as the default dns domain for this cluster... 
I0720 10:17:14.944] Jul 20 10:00:56.142: INFO: DnsDomain cluster.local
I0720 10:17:14.944] Jul 20 10:00:56.142: INFO: initializing cluster: ensuring namespaces, deployments, and pods exist and are ready
I0720 10:17:14.944] Jul 20 10:00:56.192: INFO: creating/updating pod netpol-8043-x/a
I0720 10:17:14.944] Jul 20 10:00:56.192: INFO: creating pod netpol-8043-x/a
I0720 10:17:14.944] Jul 20 10:00:56.296: INFO: creating/updating pod netpol-8043-x/b
I0720 10:17:14.945] Jul 20 10:00:56.296: INFO: creating pod netpol-8043-x/b
I0720 10:17:14.945] Jul 20 10:00:56.383: INFO: creating/updating pod netpol-8043-x/c
I0720 10:17:14.945] Jul 20 10:00:56.383: INFO: creating pod netpol-8043-x/c
I0720 10:17:14.945] Jul 20 10:00:56.520: INFO: creating/updating pod netpol-8043-y/a
I0720 10:17:14.945] Jul 20 10:00:56.520: INFO: creating pod netpol-8043-y/a
I0720 10:17:14.945] Jul 20 10:00:56.614: INFO: creating/updating pod netpol-8043-y/b
I0720 10:17:14.945] Jul 20 10:00:56.614: INFO: creating pod netpol-8043-y/b
I0720 10:17:14.946] Jul 20 10:00:56.705: INFO: creating/updating pod netpol-8043-y/c
I0720 10:17:14.946] Jul 20 10:00:56.705: INFO: creating pod netpol-8043-y/c
I0720 10:17:14.946] Jul 20 10:00:56.832: INFO: creating/updating pod netpol-8043-z/a
I0720 10:17:14.946] Jul 20 10:00:56.832: INFO: creating pod netpol-8043-z/a
I0720 10:17:14.946] Jul 20 10:00:56.920: INFO: creating/updating pod netpol-8043-z/b
I0720 10:17:14.946] Jul 20 10:00:56.920: INFO: creating pod netpol-8043-z/b
I0720 10:17:14.946] Jul 20 10:00:57.007: INFO: creating/updating pod netpol-8043-z/c
I0720 10:17:14.946] Jul 20 10:00:57.007: INFO: creating pod netpol-8043-z/c
I0720 10:17:14.946] Jul 20 10:02:42.259: INFO: finished initializing cluster state
I0720 10:17:14.946] Jul 20 10:02:42.259: INFO: waiting for HTTP servers (ports 80 and 81) to become ready
I0720 10:17:14.947] Jul 20 10:02:42.259: INFO: Starting probe from pod a to 10.0.195.148
I0720 10:17:14.947] Jul 20 10:02:42.259: INFO: Starting probe from pod a to 10.0.98.56
I0720 10:17:14.947] Jul 20 10:02:42.259: INFO: Starting probe from pod a to 10.0.107.214
```
then 

```
I0720 10:17:14.946] Jul 20 10:02:42.259: INFO: waiting for HTTP servers (ports 80 and 81) to become ready
...
I0720 10:17:14.986] Jul 20 10:05:43.987: INFO: server 81->80,TCP is ready
I0720 10:17:15.025] Jul 20 10:08:58.371: INFO: server 81->80,UDP is ready
I0720 10:17:15.120] Jul 20 10:15:06.773: INFO: server 81->81,UDP is ready
I0720 10:19:32.751] Jul 20 10:17:40.855: INFO: server 81->80,TCP is ready
```

then 

```
I0720 10:19:32.752] Jul 20 10:17:40.855: INFO: Network Policy creating netpol-3608-x/allow-ns-y-pod-a-via-namespace-pod-selector 
I0720 10:19:32.752] metadata:
I0720 10:19:32.752]   creationTimestamp: null
I0720 10:19:32.752]   name: allow-ns-y-pod-a-via-namespace-pod-selector
I0720 10:19:32.752] spec:
I0720 10:19:32.752]   ingress:
I0720 10:19:32.752]   - from:
I0720 10:19:32.752]     - namespaceSelector:
I0720 10:19:32.752]         matchLabels:
I0720 10:19:32.753]           ns: netpol-3608-y
I0720 10:19:32.753]       podSelector:
I0720 10:19:32.753]         matchLabels:
I0720 10:19:32.753]           pod: a
I0720 10:19:32.753]   podSelector:
I0720 10:19:32.753]     matchLabels:
I0720 10:19:32.753]       pod: a
I0720 10:19:32.753]   policyTypes:
I0720 10:19:32.753]   - Ingress
``` 

now validate:

```
I0720 10:19:32.753] Jul 20 10:17:40.855: INFO: creating network policy netpol-3608-x/allow-ns-y-pod-a-via-namespace-pod-selector
I0720 10:19:32.754] Jul 20 10:17:40.896: INFO: Denying all traffic *to* netpol-3608-x/a
I0720 10:19:32.754] [1mSTEP[0m: Validating reachability matrix...
I0720 10:19:32.754] [1mSTEP[0m: Validating reachability matrix... (FIRST TRY)
...
I0720 10:19:32.761] Jul 20 10:17:59.600: INFO: ExecWithOptions {Command:[/agnhost connect 10.0.91.168:80 --timeout=1s --protocol=tcp] Namespace:netpol-3608-x PodName:b ContainerName:cont-80-tcp Stdin:<nil> CaptureStdout:true CaptureStderr:true PreserveWhitespace:false Quiet:false}
```

final, print matrix: 

```
cases ignored
I0720 10:19:32.794] Jul 20 10:19:32.264: INFO: reachability: correct:72, incorrect:0, result=true
I0720 10:19:32.794] 
I0720 10:19:32.794] 
I0720 10:19:32.795] Jul 20 10:19:32.264: INFO: expected:
I0720 10:19:32.795] 
I0720 10:19:32.795] -		netpol-3608-x/a	netpol-3608-x/b	netpol-3608-x/c	netpol-3608-y/a	netpol-3608-y/b	netpol-3608-y/c	netpol-3608-z/a	netpol-3608-z/b	netpol-3608-z/c
I0720 10:19:32.795] netpol-3608-x/a	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.795] netpol-3608-x/b	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.795] netpol-3608-x/c	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.795] netpol-3608-y/a	.		.		.		.		.		.		.		.		.	
I0720 10:19:32.795] netpol-3608-y/b	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.795] netpol-3608-y/c	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.795] netpol-3608-z/a	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.796] netpol-3608-z/b	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.796] netpol-3608-z/c	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.796] 
I0720 10:19:32.796] 
I0720 10:19:32.796] 
I0720 10:19:32.796] Jul 20 10:19:32.264: INFO: observed:
I0720 10:19:32.796] 
I0720 10:19:32.796] -		netpol-3608-x/a	netpol-3608-x/b	netpol-3608-x/c	netpol-3608-y/a	netpol-3608-y/b	netpol-3608-y/c	netpol-3608-z/a	netpol-3608-z/b	netpol-3608-z/c
I0720 10:19:32.796] netpol-3608-x/a	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.796] netpol-3608-x/b	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.796] netpol-3608-x/c	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.797] netpol-3608-y/a	.		.		.		.		.		.		.		.		.	
I0720 10:19:32.797] netpol-3608-y/b	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.797] netpol-3608-y/c	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.797] netpol-3608-z/a	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.797] netpol-3608-z/b	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.797] netpol-3608-z/c	X		.		.		.		.		.		.		.		.	
I0720 10:19:32.797] 
I0720 10:19:32.797] 
I0720 10:19:32.797] 
I0720 10:19:32.797] Jul 20 10:19:32.264: INFO: comparison:
I0720 10:19:32.797] 
I0720 10:19:32.797] -		netpol-3608-x/a	netpol-3608-x/b	netpol-3608-x/c	netpol-3608-y/a	netpol-3608-y/b	netpol-3608-y/c	netpol-3608-z/a	netpol-3608-z/b	netpol-3608-z/c
I0720 10:19:32.798] netpol-3608-x/a	.		.		.		.		.		.		.		.		.	
I0720 10:19:32.798] netpol-3608-x/b	.		.		.		.		.		.		.		.		.	
I0720 10:19:32.798] netpol-3608-x/c	.		.		.		.		.		.		.		.		.	
I0720 10:19:32.798] netpol-3608-y/a	.		.		.		.		.		.		.		.		.	
I0720 10:19:32.798] netpol-3608-y/b	.		.		.		.		.		.		.		.		.	
I0720 10:19:32.798] netpol-3608-y/c	.		.		.		.		.		.		.		.		.	
I0720 10:19:32.798] netpol-3608-z/a	.		.		.		.		.		.		.		.		.	
I0720 10:19:32.798] netpol-3608-z/b	.		.		.		.		.		.		.		.		.	
I0720 10:19:32.798] netpol-3608-z/c	.		.		.		.		.		.		.		.		.	
I0720 10:19:32.798] 
```


