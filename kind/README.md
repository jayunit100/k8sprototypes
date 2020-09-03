# how to run this

run the kind-local-up.sh script.

##  antrea: 
Modify the `antrea` parameters (turn them on) and comment out the calico params

```
  cluster=antrea
  conf=calico-conf.yaml
```

## calico

Same as antrea just uncomment the 

```
  ## calico conf
```

section
