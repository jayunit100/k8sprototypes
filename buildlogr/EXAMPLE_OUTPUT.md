 # POSSIBLE "driver" failures
 
 
``` 
   2 {"msg":"PASSED [k8s.io] Container Lifecycle Hook when create a pod with lifecycle hook should execute poststart http hook properly [NodeConformance] [Conformance]","total":-1,"completed":10,"skipped":95,"failed":0}
   2 {"msg":"PASSED [k8s.io] Kubelet when scheduling a busybox command that always fails in a pod should have an terminated reason [NodeConformance] [Conformance]","total":-1,"completed":12,"skipped":211,"failed":0}
   2 {"msg":"PASSED [k8s.io] Variable Expansion should succeed in writing subpaths in container [sig-storage][Slow] [Conformance]","total":-1,"completed":6,"skipped":99,"failed":0}
   2 {"msg":"PASSED [sig-api-machinery] Garbage collector should keep the rc around until all its pods are deleted if the deleteOptions says so [Conformance]","total":-1,"completed":7,"skipped":45,"failed":0}
   2 {"msg":"PASSED [sig-api-machinery] Secrets should fail to create secret due to empty secret key [Conformance]","total":-1,"completed":3,"skipped":20,"failed":1,"failures":["[sig-network] NetworkPolicy [LinuxOnly] NetworkPolicy between server and client should enforce policy based on NamespaceSelector with MatchExpressions[Feature:NetworkPolicy]"]}
   2 {"msg":"PASSED [sig-api-machinery] Watchers should be able to restart watching from the last resource version observed by the previous watch [Conformance]","total":-1,"completed":9,"skipped":108,"failed":1,"failures":["[sig-network] NoSNAT [Feature:NoSNAT] [Slow] Should be able to send traffic between Pods without SNAT"]}
   2 {"msg":"PASSED [sig-cli] Kubectl client Kubectl run pod should create a pod from an image when restart is Never  [Conformance]","total":-1,"completed":10,"skipped":93,"failed":0}
   2 {"msg":"PASSED [sig-network] Netpol [LinuxOnly] NetworkPolicy between server and client should allow ingress access from updated pod [Feature:NetworkPolicy]","total":-1,"completed":4,"skipped":39,"failed":0}
   2 {"msg":"PASSED [sig-network] Netpol [LinuxOnly] NetworkPolicy between server and client should enforce multiple ingress policies with ingress allow-all policy taking precedence [Feature:NetworkPolicy]","total":-1,"completed":20,"skipped":181,"failed":0}
   2 {"msg":"PASSED [sig-network] Netpol [LinuxOnly] NetworkPolicy between server and client should enforce policy to allow traffic only from a pod in a different namespace based on PodSelector and NamespaceSelector [Feature:NetworkPolicy]","total":-1,"completed":11,"skipped":151,"failed":0}
   2 {"msg":"PASSED [sig-network] Netpol [LinuxOnly] NetworkPolicy between server and client should enforce updated policy [Feature:NetworkPolicy]","total":-1,"completed":5,"skipped":39,"failed":0}
   2 {"msg":"PASSED [sig-network] Netpol [LinuxOnly] NetworkPolicy between server and client should ensure an IP overlapping both IPBlock.CIDR and IPBlock.Except is allowed [Feature:NetworkPolicy]","total":-1,"completed":20,"skipped":256,"failed":0}
   2 {"msg":"PASSED [sig-network] Netpol [LinuxOnly] NetworkPolicy between server and client should work with Ingress, Egress specified together [Feature:NetworkPolicy]","total":-1,"completed":25,"skipped":304,"failed":0}
   2 {"msg":"PASSED [sig-network] NetworkPolicy [LinuxOnly] NetworkPolicy between server and client should deny ingress access to updated pod [Feature:NetworkPolicy]","total":-1,"completed":15,"skipped":177,"failed":0}
   2 {"msg":"PASSED [sig-network] Networking Granular Checks: Services should function for endpoint-Service: http","total":-1,"completed":18,"skipped":363,"failed":0}
   2 {"msg":"PASSED [sig-network] Services should be able to switch session affinity for NodePort service [LinuxOnly] [Conformance]","total":-1,"completed":10,"skipped":168,"failed":0}
   2 {"msg":"PASSED [sig-node] ConfigMap should run through a ConfigMap lifecycle [Conformance]","total":-1,"completed":6,"skipped":45,"failed":0}
   2 {"msg":"PASSED [sig-storage] Projected configMap should be consumable from pods in volume as non-root [NodeConformance] [Conformance]","total":-1,"completed":8,"skipped":82,"failed":0}
  29 
  29 {"msg":"
  33 --
  58 ********
```
