# Netsh

Netsh is the equialent to iptables for linux clusters.  Ive been using it to compare/understand how the kube proxy works.  It turns out that in the codebase
for the Kubernetes scheduler itself, there are calls to netsh that happen.  The overall codepath looks something like this:

- the Provider interface is implemented by any kube proxy (linux, windows, linux userspace, windows userspace, ipvs, ...).
- the provider provides handlers for services, endpoints, and endpoint slices
- ultimately the windows kube proxy service handler calls various netsh commands

```
type Provider interface {
    config.EndpointsHandler

    config.EndpointSliceHandler
        
    config.ServiceHandler
        -> OnServiceAdd(service *v1.Service)
           -> mergeService
            *** pkg/proxy/winuserspace/proxier.go *** 
            -> addServiceOnPortPortal
               -> netsh add address name,  address, mask
```

You can see the manifestation of these netsh commands in the output for netsh dump, in a windows terminal:

```
  # ----------------------------------
  # IPv4 Configuration
  # ----------------------------------
  pushd interface ipv4

  reset
  set global
  add route prefix=0.0.0.0/0 interface="vEthernet (Ethernet) 2" nexthop=172.19.48.1 publish=Yes
  add route prefix=0.0.0.0/0 interface="br-int" nexthop=10.0.0.1 publish=Yes
  add route prefix=100.1.1.0/24 interface="antrea-gw0" nexthop=100.1.1.1 publish=Yes
  add route prefix=0.0.0.0/0 interface="vEthernet (Ethernet)" nexthop=172.19.48.1 publish=Yes
  set interface interface="Ethernet (Kernel Debugger)" forwarding=enabled advertise=enabled nud=enabled ignoredefaultroutes=disabled
  set interface interface="Ethernet0" forwarding=enabled advertise=enabled nud=enabled ignoredefaultroutes=disabled
  set interface interface="vEthernet (nat)" forwarding=enabled advertise=enabled nud=enabled ignoredefaultroutes=disabled
  set interface interface="vEthernet (fe83824bf7b9e39)" forwarding=enabled advertise=enabled nud=enabled ignoredefaultroutes=disabled
  set interface interface="vEthernet (HNS Internal NIC)" forwarding=enabled advertise=enabled nud=enabled ignoredefaultroutes=disabled
  set interface interface="vEthernet (Ethernet) 2" forwarding=enabled advertise=enabled nud=enabled ignoredefaultroutes=disabled
  set interface interface="br-int" forwarding=enabled advertise=enabled nud=enabled ignoredefaultroutes=disabled
  set interface interface="antrea-gw0" forwarding=enabled advertise=enabled nud=enabled ignoredefaultroutes=disabled
  set interface interface="vEthernet (Ethernet)" forwarding=enabled advertise=enabled nud=enabled ignoredefaultroutes=disabled
  add address name="vEthernet (nat)" address=172.25.208.1 mask=255.255.240.0
  add address name="vEthernet (fe83824bf7b9e39)" address=172.19.48.1 mask=255.255.240.0
  add address name="vEthernet (HNS Internal NIC)" address=100.2.2.246 mask=255.0.0.0 
  add address name="vEthernet (HNS Internal NIC)" address=100.2.2.1 mask=255.0.0.0
  add address name="vEthernet (HNS Internal NIC)" address=100.2.2.10 mask=255.0.0.0
  add address name="vEthernet (Ethernet) 2" address=172.19.56.17 mask=255.255.240.0
  add address name="br-int" address=10.0.0.44 mask=255.255.255.0
  add address name="antrea-gw0" address=100.1.2.1 mask=255.255.255.0
  add address name="vEthernet (Ethernet)" address=172.19.48.93 mask=255.255.240.0


  popd
```
