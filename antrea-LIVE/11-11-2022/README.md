# Antrea LIVE: Episode 3 (Multus, Whereabouts, and host-local IPAM)

# Show Details

Show link: https://youtu.be/Q1CBFoMAG2g

Live notes: https://hackmd.io/wxNOmhZdRNm_hJzZNFKwRg 

https://github.com/antrea-io/antrea/tree/main/docs/cookbooks/multus


- host ports on windows, do they work? or do you need nodeportlocal
- IF -> elastic IP (aws CNI) -> limited
    - encap 
- `service-proxy-name` <-- annotation on k8s services that disables proxying
- why multiple nics
    - bypass CNI network for perf, minimize hops
    - sr-iov/dpdk/edp
- multus quickstart: 
  - hostlocal allocation per node
  - skip dchp, allocate from a pool
  - fork of `static` plugin
- shim + thick component (k8s controller)
- sriov has a device plugin.

## notes from jiunjen:
```
We are developing native CNF/secondary network support in Antrea in collaboration with Intel, but it is in progress. The first PR for SR-IOV is here: https://github.com/antrea-io/antrea/pull/2651
#2651 Antrea secondary network implementation for SRIOV
Implement secondary network for Antrea CNI configured pods using SRIOV interface. (CNF use case)
Native support means it does not require Multus, Whereabouts, etc.
```
