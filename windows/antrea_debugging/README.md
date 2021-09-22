```
PS C:\Users\capv> cat C:\var\log\antrea\antrea-agent.exe.INFO
```

Will give you the logs...

```
PS C:\Users\capv> cat C:\var\log\antrea\antrea-agent.exe.ERROR
Log file created at: 2021/09/22 14:57:17
Running on machine: tkg-vc-antrea-md-0-windows-containerd-57756fcb9-j9fp7
Binary: Built with gc go1.15.3 for windows/amd64
Log line format: [IWEF]mmdd hh:mm:ss.uuuuuu threadid file:line] msg
E0922 14:57:17.932499    2076 server.go:491] Failed to remove interfaces for container db072a92158874a6df9dfc4ee12f95dda0d82e3009877c4e7d070ed0ecb5332c: timeout when deleting HNSEndpoint pod-f72e-69e818
E0922 17:02:59.036499    2076 agent.go:58] Failed to partially update agent monitoring CRD: etcdserver: request timed out
PS C:\Users\capv> cat C:\var\log\antrea\antrea-agent.exe.INFO
Log file created at: 2021/09/22 14:38:59
Running on machine: tkg-vc-antrea-md-0-windows-containerd-57756fcb9-j9fp7
Binary: Built with gc go1.15.3 for windows/amd64
Log line format: [IWEF]mmdd hh:mm:ss.uuuuuu threadid file:line] msg
```

Now the agent starts

```
I0922 14:38:59.839925    2076 log_file.go:99] Set log file max size to 104857600
I0922 14:38:59.842928    2076 agent.go:66] Starting Antrea agent (version v0.13.3-unknown)
W0922 14:38:59.846933    2076 env.go:67] Environment variable POD_NAMESPACE not found
W0922 14:38:59.849020    2076 env.go:105] Failed to get Pod Namespace from environment. Using "kube-system" as the Antrea Service Namespace
I0922 14:38:59.849563    2076 prometheus.go:151] Initializing prometheus metrics
I0922 14:38:59.849563    2076 ovs_client.go:67] Connecting to OVSDB at address \\.\pipe\C:openvswitchvarrunopenvswitchdb.sock
I0922 14:38:59.855664    2076 agent.go:205] Setting up node network
I0922 14:38:59.856628    2076 env.go:44] Environment variable NODE_NAME not found, using hostname instead
I0922 14:38:59.886639    2076 agent.go:656] Setting Node MTU=1450
I0922 14:39:21.653145    2076 net_windows.go:370] Created HNSNetwork with name antrea-hnsnetwork id 3F2BBA3C-ADF0-4414-8C82-36521B908902
I0922 14:39:21.656809    2076 ovs_client.go:118] Created bridge: f8ca2cd5-60e4-477e-bdc5-50a0ae356614
I0922 14:39:45.265385    2076 agent.go:787] No round number found in OVSDB, using 1
I0922 14:39:45.265385    2076 agent.go:799] Using round number 1
I0922 14:39:46.277039    2076 ofctrl_bridge.go:220] OFSwitch is connected: 00:00:00:50:56:a7:dc:59
I0922 14:39:46.292764    2076 agent.go:242] Agent initialized NodeConfig=NodeName: tkg-vc-antrea-md-0-windows-containerd-57756fcb9-j9fp7, OVSBridge: br-int, PodIPv4CIDR: 100.96.1.0/24, PodIPv6CIDR: <nil>, NodeIP: 10.216.113.157/20, Gateway: Name antrea-gw0: IPv4 100.96.1.1, IPv6 <nil>, MAC 00:15:5d:71:9d:03, NetworkConfig=&{Encap geneve false }
I0922 14:39:46.296731    2076 metrics.go:124] Registering Antrea Proxy prometheus metrics
I0922 14:39:46.306739    2076 server.go:641] Reconciliation for CNI server
W0922 14:39:46.329820    2076 agent.go:425] flow-restore-wait was not true before the delete call was made, will retry
W0922 14:39:46.539206    2076 agent.go:425] flow-restore-wait was not true before the delete call was made, will retry
W0922 14:39:46.736011    2076 agent.go:425] flow-restore-wait was not true before the delete call was made, will retry
W0922 14:39:46.938303    2076 agent.go:425] flow-restore-wait was not true before the delete call was made, will retry
W0922 14:39:47.144897    2076 agent.go:425] flow-restore-wait was not true before the delete call was made, will retry
W0922 14:39:47.330823    2076 agent.go:425] flow-restore-wait was not true before the delete call was made, will retry
W0922 14:39:47.531084    2076 agent.go:425] flow-restore-wait was not true before the delete call was made, will retry
W0922 14:39:47.731178    2076 agent.go:425] flow-restore-wait was not true before the delete call was made, will retry
W0922 14:39:47.945635    2076 agent.go:425] flow-restore-wait was not true before the delete call was made, will retry
W0922 14:39:48.130987    2076 agent.go:425] flow-restore-wait was not true before the delete call was made, will retry
W0922 14:39:48.330493    2076 agent.go:425] flow-restore-wait was not true before the delete call was made, will retry
W0922 14:39:48.331501    2076 agent.go:425] flow-restore-wait was not true before the delete call was made, will retry
I0922 14:39:48.332511    2076 agent.go:433] flow-restore-wait was not true, skip cleaning it up
I0922 14:39:48.345497    2076 client.go:130] Updating Antrea client with the new CA bundle
I0922 14:39:48.346490    2076 log_file.go:127] Starting log file monitoring. Maximum log file number is 4
I0922 14:39:48.346490    2076 server.go:575] Starting CNI server
I0922 14:39:48.347478    2076 traceflow_controller.go:155] Starting AntreaAgentTraceflowController
I0922 14:39:48.348110    2076 configmap_cafile_content.go:202] Starting antrea-ca::kube-system::antrea-ca::ca.crt
I0922 14:39:48.348394    2076 shared_informer.go:223] Waiting for caches to sync for antrea-ca::kube-system::antrea-ca::ca.crt
I0922 14:39:48.347478    2076 node_route_controller.go:283] Starting AntreaAgentNodeRouteController
I0922 14:39:48.348936    2076 shared_informer.go:223] Waiting for caches to sync for AntreaAgentNodeRouteController
I0922 14:39:48.348110    2076 networkpolicy_controller.go:372] Waiting for Antrea client to be ready
I0922 14:39:48.348936    2076 networkpolicy_controller.go:383] Antrea client is ready
I0922 14:39:48.348110    2076 agent.go:46] Starting Antrea Agent Monitor
I0922 14:39:48.348110    2076 server.go:585] CNI server is listening ...
I0922 14:39:48.348285    2076 config.go:223] Starting service config controller
I0922 14:39:48.350009    2076 shared_informer.go:223] Waiting for caches to sync for service config
I0922 14:39:48.348285    2076 config.go:132] Starting endpoints config controller
I0922 14:39:48.348394    2076 shared_informer.go:223] Waiting for caches to sync for AntreaAgentTraceflowController
I0922 14:39:48.348936    2076 networkpolicy_controller.go:392] Waiting for all watchers to complete full sync
I0922 14:39:48.348936    2076 networkpolicy_controller.go:577] Starting watch for NetworkPolicy
I0922 14:39:48.348936    2076 networkpolicy_controller.go:577] Starting watch for AppliedToGroup
I0922 14:39:48.348936    2076 networkpolicy_controller.go:577] Starting watch for AddressGroup
I0922 14:39:48.350597    2076 shared_informer.go:223] Waiting for caches to sync for endpoints config
W0922 14:39:48.365860    2076 env.go:58] Environment variable POD_NAME not found
W0922 14:39:48.366164    2076 env.go:67] Environment variable POD_NAMESPACE not found
I0922 14:39:48.374667    2076 networkpolicy_controller.go:584] Started watch for AddressGroup
I0922 14:39:48.374667    2076 networkpolicy_controller.go:584] Started watch for AppliedToGroup
I0922 14:39:48.375215    2076 networkpolicy_controller.go:584] Started watch for NetworkPolicy
```
After it starts youll see cmdAdd cmdDel and so on, if its healthy...
```
0922 14:39:48.376271    2076 networkpolicy_controller.go:394] All watchers have completed full sync, installing flows for init events
I0922 14:39:48.376271    2076 networkpolicy_controller.go:398] Starting NetworkPolicy workers now
I0922 14:39:48.376271    2076 networkpolicy_controller.go:404] Starting IDAllocator worker to maintain the async rule cache
I0922 14:39:48.448748    2076 shared_informer.go:230] Caches are synced for antrea-ca::kube-system::antrea-ca::ca.crt
I0922 14:39:48.449754    2076 shared_informer.go:230] Caches are synced for AntreaAgentNodeRouteController
I0922 14:39:48.449754    2076 node_route_controller.go:257] Reconciliation for AntreaAgentNodeRouteController
I0922 14:39:48.450754    2076 shared_informer.go:230] Caches are synced for service config
I0922 14:39:48.450754    2076 client.go:130] Updating Antrea client with the new CA bundle
I0922 14:39:48.451793    2076 shared_informer.go:230] Caches are synced for AntreaAgentTraceflowController
I0922 14:39:48.453776    2076 shared_informer.go:230] Caches are synced for endpoints config
I0922 14:39:49.087427    2076 serving.go:313] Generated self-signed cert in-memory
I0922 14:39:49.492110    2076 configmap_cafile_content.go:202] Starting client-ca::kube-system::extension-apiserver-authentication::requestheader-client-ca-file
I0922 14:39:49.493113    2076 shared_informer.go:223] Waiting for caches to sync for client-ca::kube-system::extension-apiserver-authentication::requestheader-client-ca-file
I0922 14:39:49.493113    2076 secure_serving.go:178] Serving securely on [::]:10350
I0922 14:39:49.493113    2076 tlsconfig.go:240] Starting DynamicServingCertificateController
I0922 14:39:49.492110    2076 configmap_cafile_content.go:202] Starting client-ca::kube-system::extension-apiserver-authentication::client-ca-file
I0922 14:39:49.494470    2076 shared_informer.go:223] Waiting for caches to sync for client-ca::kube-system::extension-apiserver-authentication::client-ca-file
I0922 14:39:49.504202    2076 node_route_controller.go:417] Adding routes and flows to Node tkg-vc-antrea-control-plane-j9hms, podCIDRs: [100.96.0.0/24], addresses: [{Hostname tkg-vc-antrea-control-plane-j9hms} {InternalIP 10.216.126.9} {ExternalIP 10.216.126.9}]
I0922 14:39:49.505222    2076 node_route_controller.go:423] Adding routes and flows to Node tkg-vc-antrea-control-plane-j9hms, podCIDR: 100.96.0.0/24, addresses: [{Hostname tkg-vc-antrea-control-plane-j9hms} {InternalIP 10.216.126.9} {ExternalIP 10.216.126.9}]
I0922 14:39:49.594362    2076 shared_informer.go:230] Caches are synced for client-ca::kube-system::extension-apiserver-authentication::requestheader-client-ca-file
I0922 14:39:49.595363    2076 shared_informer.go:230] Caches are synced for client-ca::kube-system::extension-apiserver-authentication::client-ca-file
I0922 14:39:56.293117    2076 agent.go:358] Deleting stale flows from previous round if any
I0922 14:39:56.293117    2076 agent.go:249] Persisting round number 1 to OVSDB
I0922 14:39:56.294114    2076 agent.go:254] Round number 1 was persisted to OVSDB
I0922 14:55:17.504452    2076 server.go:375] Received CmdAdd request cni_args:<container_id:"db072a92158874a6df9dfc4ee12f95dda0d82e3009877c4e7d070ed0ecb5332c" netns:"12d7094d-ef08-4358-9e45-cc007b65954a" ifname:"eth0" args:"IgnoreUnknown=1;K8S_POD_NAMESPACE=hybrid-network-4867;K8S_POD_NAME=pod-f72ed608-89fd-466b-8c9a-06ec231c66b8;K8S_POD_INFRA_CONTAINER_ID=db072a92158874a6df9dfc4ee12f95dda0d82e3009877c4e7d070ed0ecb5332c
" path:"C:/opt/cni/bin" network_configuration:"{\"capabilities\":{\"dns\":true},\"cniVersion\":\"0.3.0\",\"ipam\":{\"type\":\"host-local\"},\"name\":\"antrea\",\"runtimeConfig\":{\"dns\":{\"Servers\":[\"100.64.0.10\"],\"Searches\":[\"hybrid-network-4867.svc.cluster.local\",\"svc.cluster.local\",\"cluster.local\"],\"Options\":[\"ndots:5\"]}},\"type\":\"antrea\"}" >
I0922 14:55:19.866638    2076 server.go:436] Requested ip addresses for container db072a92158874a6df9dfc4ee12f95dda0d82e3009877c4e7d070ed0ecb5332c: &{0.4.0 [] [{Version:4 Interface:<nil> Address:{IP:100.96.1.2 Mask:ffffff00} Gateway:100.96.1.1}] [] {[]  [] []}}
I0922 14:55:19.868648    2076 server_windows.go:41] Got runtime DNS configuration: {[100.64.0.10]  [hybrid-network-4867.svc.cluster.local svc.cluster.local cluster.local] []}
```
We can see here that a vEther device is being created 
```
I0922 14:55:23.548371    2076 pod_configuration_windows.go:47] Waiting for interface vEthernet (pod-f72e-69e818) to be created
I0922 14:55:23.548371    2076 pod_configuration.go:259] Configured interfaces for container db072a92158874a6df9dfc4ee12f95dda0d82e3009877c4e7d070ed0ecb5332c
I0922 14:55:23.549486    2076 server.go:461] CmdAdd for container db072a92158874a6df9dfc4ee12f95dda0d82e3009877c4e7d070ed0ecb5332c succeeded
I0922 14:55:25.936446    2076 pod_configuration_windows.go:60] Waiting for interface vEthernet (pod-f72e-69e818) to be created
I0922 14:56:40.398500    2076 server.go:375] Received CmdAdd request cni_args:<container_id:"f2743784db2a3e3e31fcaa7b35e46318f3026ef6be79b84aa4ddc75c0fea0021" netns:"5948c427-32c3-45a1-bdbc-cc184d8bd918" ifname:"eth0" args:"IgnoreUnknown=1;K8S_POD_NAMESPACE=windows-run-as-username-7991;K8S_POD_NAME=run-as-username-c4b738d7-3537-426f-8009-ffebb8a60b58;K8S_POD_INFRA_CONTAINER_ID=f2743784db2a3e3e31fcaa7b35e46318f3026ef6be7
9b84aa4ddc75c0fea0021" path:"C:/opt/cni/bin" network_configuration:"{\"capabilities\":{\"dns\":true},\"cniVersion\":\"0.3.0\",\"ipam\":{\"type\":\"host-local\"},\"name\":\"antrea\",\"runtimeConfig\":{\"dns\":{\"Servers\":[\"100.64.0.10\"],\"Searches\":[\"windows-run-as-username-7991.svc.cluster.local\",\"svc.cluster.local\",\"cluster.local\"],\"Options\":[\"ndots:5\"]}},\"type\":\"antrea\"}" >
I0922 14:56:40.435161    2076 server.go:436] Requested ip addresses for container f2743784db2a3e3e31fcaa7b35e46318f3026ef6be79b84aa4ddc75c0fea0021: &{0.4.0 [] [{Version:4 Interface:<nil> Address:{IP:100.96.1.3 Mask:ffffff00} Gateway:100.96.1.1}] [] {[]  [] []}}
I0922 14:56:40.436119    2076 server_windows.go:41] Got runtime DNS configuration: {[100.64.0.10]  [windows-run-as-username-7991.svc.cluster.local svc.cluster.local cluster.local] []}
I0922 14:56:42.529140    2076 pod_configuration.go:259] Configured interfaces for container f2743784db2a3e3e31fcaa7b35e46318f3026ef6be79b84aa4ddc75c0fea0021
I0922 14:56:42.529140    2076 server.go:461] CmdAdd for container f2743784db2a3e3e31fcaa7b35e46318f3026ef6be79b84aa4ddc75c0fea0021 succeeded
I0922 14:56:42.529140    2076 pod_configuration_windows.go:47] Waiting for interface vEthernet (run-as-u-8c9665) to be created
I0922 14:57:03.815487    2076 server.go:469] Received CmdDel request cni_args:<container_id:"f2743784db2a3e3e31fcaa7b35e46318f3026ef6be79b84aa4ddc75c0fea0021" netns:"5948c427-32c3-45a1-bdbc-cc184d8bd918" ifname:"eth0" args:"K8S_POD_NAMESPACE=windows-run-as-username-7991;K8S_POD_NAME=run-as-username-c4b738d7-3537-426f-8009-ffebb8a60b58;K8S_POD_INFRA_CONTAINER_ID=f2743784db2a3e3e31fcaa7b35e46318f3026ef6be79b84aa4ddc75c0fe
a0021;IgnoreUnknown=1" path:"C:/opt/cni/bin" network_configuration:"{\"capabilities\":{\"dns\":true},\"cniVersion\":\"0.3.0\",\"ipam\":{\"type\":\"host-local\"},\"name\":\"antrea\",\"runtimeConfig\":{\"dns\":{\"Servers\":[\"100.64.0.10\"],\"Searches\":[\"windows-run-as-username-7991.svc.cluster.local\",\"svc.cluster.local\",\"cluster.local\"],\"Options\":[\"ndots:5\"]}},\"type\":\"antrea\"}" >
I0922 14:57:03.822487    2076 server.go:469] Received CmdDel request cni_args:<container_id:"f2743784db2a3e3e31fcaa7b35e46318f3026ef6be79b84aa4ddc75c0fea0021" netns:"5948c427-32c3-45a1-bdbc-cc184d8bd918" ifname:"eth0" args:"IgnoreUnknown=1;K8S_POD_NAMESPACE=windows-run-as-username-7991;K8S_POD_NAME=run-as-username-c4b738d7-3537-426f-8009-ffebb8a60b58;K8S_POD_INFRA_CONTAINER_ID=f2743784db2a3e3e31fcaa7b35e46318f3026ef6be7
9b84aa4ddc75c0fea0021" path:"C:/opt/cni/bin" network_configuration:"{\"capabilities\":{\"dns\":true},\"cniVersion\":\"0.3.0\",\"ipam\":{\"type\":\"host-local\"},\"name\":\"antrea\",\"runtimeConfig\":{\"dns\":{\"Servers\":[\"100.64.0.10\"],\"Searches\":[\"windows-run-as-username-7991.svc.cluster.local\",\"svc.cluster.local\",\"cluster.local\"],\"Options\":[\"ndots:5\"]}},\"type\":\"antrea\"}" >
I0922 14:57:03.840493    2076 server.go:488] Deleted IP addresses for container f2743784db2a3e3e31fcaa7b35e46318f3026ef6be79b84aa4ddc75c0fea0021
```

To see the names of pod network endpoints, run

```
PS C:\Users\capv> hnsdiag list endpoints
Name             ID                                   Virtual Network Name
pod-f72e-69e818  18ca0b2b-3203-4c49-98f4-a3fcb812e48c antrea-hnsnetwork
busybox-719fb3   ea122d28-0e74-401f-8c18-6e87d423cc52 antrea-hnsnetwork
windows--55662b  dbb6d0dd-ae42-464b-8df4-00813492d83b antrea-hnsnetwork
```

To test wether pods have networks, run 

```
PS C:\Users\capv> Get-HnsEndpoint

ActivityId                : B56C3D4B-7225-423E-818A-ED781692053B
AdditionalParams          : @{attached-mac=; container-id=db072a92158874a6df9dfc4ee12f95dda0d82e3009877c4e7d070ed0ecb5332c; ip-address=100.96.1.2; pod-name=pod-f72ed608-89fd-466b-8c9a-06ec231c66b8; pod-namespace=hybrid-network-4867}
CreateProcessingStartTime : 132768213198814662
DNSServerList             : 100.64.0.10
DNSSuffix                 : hybrid-network-4867.svc.cluster.local,svc.cluster.local,cluster.local
EncapOverhead             : 0
GatewayAddress            : 100.96.1.1
Health                    : @{LastErrorCode=0; LastUpdateTime=132768213198786463}
ID                        : 18CA0B2B-3203-4C49-98F4-A3FCB812E48C
IPAddress                 : 100.96.1.2
MacAddress                : 00-15-5D-25-85-06
Name                      : pod-f72e-69e818
Namespace                 : @{ID=12D7094D-EF08-4358-9E45-CC007B65954A; IsDefault=False}
Policies                  : {}
PrefixLength              : 24
RemoveProcessingStartTime : 132768214068678123
Resources                 : @{AdditionalParams=; AllocationOrder=2; Allocators=System.Object[]; Health=; ID=B56C3D4B-7225-423E-818A-ED781692053B; PortOperationTime=0; State=1; SwitchOperationTime=0; VfpOperationTime=0; parentId=6DC6FD75-C2E0-4BF9-A2BA-E2C0F590290A}
SharedContainers          : {db072a92158874a6df9dfc4ee12f95dda0d82e3009877c4e7d070ed0ecb5332c}
StartTime                 : 132768213261881793
State                     : 3
Type                      : Transparent
Version                   : 38654705667
VirtualNetwork            : 3F2BBA3C-ADF0-4414-8C82-36521B908902
VirtualNetworkName        : antrea-hnsnetwork
```


