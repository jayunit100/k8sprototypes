# Install prometheus metrics service in antrea 

Some advanced usage of antrea projects

## Antctl

Antctl gives you a fine grained view of some of the finer antrea processes...

To build it:

```
make antctl
```

Then

```

```


## Prometheus

```
	enable prometheus metrics in the configmap:		
	this must be done in BOTH places where you see enablePrometheusMetrics
```

Now

```
	kubectl apply -f build/yamls/antrea-prometheus.yml
```

then expose them:

```
	kubectl port-forward pod/prometheus-deployment-68b648df9c-28q6d --address 0.0.0.0 9090:9090 -n monitoring
```

now go to 

```
	http://localhost:9090/targets
```

You should see the targets running green.


Now, you can see how many flow table rules are in occurance... 

![Image description](flowtables.png)


