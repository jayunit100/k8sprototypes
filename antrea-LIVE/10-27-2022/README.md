# Antrea LIVE: Episode 1 (Antrea 1.3.0, FQDN, and K8sNetlook)

# Show Details

https://www.youtube.com/watch?v=aWUwxQ58bEQ

## Antrea 1.3.0

- kubectl apply -f https://github.com/antrea-io/antrea/releases/download/v1.3.0/antrea.yml
- install ~ jayunit100/k8sprototypes/kind/ kind-local-up.sh

## FQDN Policies


```
apiVersion: crd.antrea.io/v1alpha1
kind: ClusterNetworkPolicy
metadata:
  name: acnp-fqdn-all-foobar
spec:
  priority: 1
  appliedTo:
  - podSelector:
      matchLabels:
        app: client
  egress:
  - action: Drop
    to:
      - fqdn: "*foobar.com"
```

## K8sNetLook

- https://github.com/sarun87/k8snetlook

## Notes

- Release Notes: https://github.com/antrea-io/antrea/releases/tag/v1.3.0
- https://github.com/kubernetes/enhancements/pull/2522
- `kubectl get acnp`
- `CLUSTER=antrea CONFIG=kind-conf.yaml ./kind-local-up.sh`
- https://kubernetes.io/docs/concepts/services-networking/network-policies/
- https://github.com/sarun87/k8snetlook/blob/master/examples/run-k8s.yaml
- bitnami/nginx image to avoid docker limits
- Arun's awesome kubecon talk slides: https://kccncna2021.sched.com/event/lV0b
