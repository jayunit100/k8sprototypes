Logs from open vswitch on windows, might look like this (if you look in ovs-vswitchd.log in startup cases where it fails...

```
2020-12-07T22:51:05.001Z|00001|vlog|INFO|opened log file C:/openvswitch/var/log/openvswitch/ovs-vswitchd.log
2020-12-07T22:51:05.008Z|00002|ovs_numa|INFO|Discovered 0 NUMA nodes and 0 CPU cores
2020-12-07T22:51:05.009Z|00003|reconnect|INFO|unix:C:/openvswitch/var/run/openvswitch/db.sock: connecting...
2020-12-07T22:51:05.010Z|00004|reconnect|INFO|unix:C:/openvswitch/var/run/openvswitch/db.sock: connected
2020-12-07T22:51:05.016Z|00005|dpif_netlink|INFO|The kernel module does not support meters.


2020-12-08T07:10:38.645Z|00003|reconnect|INFO|unix:C:/openvswitch/var/run/openvswitch/db.sock: connecting...
2020-12-08T07:10:38.645Z|00004|reconnect|INFO|unix:C:/openvswitch/var/run/openvswitch/db.sock: connected
2020-12-08T13:18:01.428Z|00001|vlog|INFO|opened log file C:/openvswitch/var/log/openvswitch/ovs-vswitchd.log
2020-12-08T13:18:01.430Z|00002|ovs_numa|INFO|Discovered 0 NUMA nodes and 0 CPU cores
2020-12-08T13:18:01.430Z|00003|reconnect|INFO|unix:C:/openvswitch/var/run/openvswitch/db.sock: connecting...
2020-12-08T13:18:01.430Z|00004|reconnect|INFO|unix:C:/openvswitch/var/run/openvswitch/db.sock: connected
2020-12-08T13:18:01.433Z|00005|dpif_netlink|INFO|The kernel module does not support meters.
2020-12-08T13:18:02.372Z|00001|vlog|INFO|opened log file C:/openvswitch/var/log/openvswitch/ovs-vswitchd.log
2020-12-08T13:18:02.374Z|00002|ovs_numa|INFO|Discovered 0 NUMA nodes and 0 CPU cores
2020-12-08T13:18:02.374Z|00003|reconnect|INFO|unix:C:/openvswitch/var/run/openvswitch/db.sock: connecting...


2020-12-08T13:36:48.269Z|00001|vlog|INFO|opened log file C:/openvswitch/var/log/openvswitch/ovs-vswitchd.log
2020-12-08T13:36:48.271Z|00002|ovs_numa|INFO|Discovered 0 NUMA nodes and 0 CPU cores
2020-12-08T13:36:48.271Z|00003|reconnect|INFO|unix:C:/openvswitch/var/run/openvswitch/db.sock: connecting...
2020-12-08T13:36:48.271Z|00004|reconnect|INFO|unix:C:/openvswitch/var/run/openvswitch/db.sock: connected
2020-12-08T13:36:48.273Z|00005|dpif_netlink|INFO|The kernel module does not support meters.
2020-12-08T13:36:48.929Z|00001|vlog|INFO|opened log file C:/openvswitch/var/log/openvswitch/ovs-vswitchd.log
2020-12-08T13:36:48.931Z|00002|ovs_numa|INFO|Discovered 0 NUMA nodes and 0 CPU cores
2020-12-08T13:36:48.931Z|00003|reconnect|INFO|unix:C:/openvswitch/var/run/openvswitch/db.sock: connecting...
2020-12-08T13:36:48.931Z|00004|reconnect|INFO|unix:C:/openvswitch/var/run/openvswitch/db.sock: connected
2020-12-08T13:36:48.934Z|00005|dpif_netlink|INFO|The kernel module does not support meters.
2020-12-08T13:36:50.012Z|00001|vlog|INFO|opened log file C:/openvswitch/var/log/openvswitch/ovs-vswitchd.log
2020-12-08T13:36:50.014Z|00002|ovs_numa|INFO|Discovered 0 NUMA nodes and 0 CPU cores
2020-12-08T13:36:50.014Z|00003|reconnect|INFO|unix:C:/openvswitch/var/run/openvswitch/db.sock: connecting...
2020-12-08T13:36:50.014Z|00004|reconnect|INFO|unix:C:/openvswitch/var/run/openvswitch/db.sock: connected
2020-12-08T13:36:50.018Z|00005|dpif_netlink|INFO|The kernel module does not support meters.
2020-12-08T13:36:51.053Z|00001|vlog|INFO|opened log file C:/openvswitch/var/log/openvswitch/ovs-vswitchd.log
2020-12-08T13:36:51.056Z|00002|ovs_numa|INFO|Discovered 0 NUMA nodes and 0 CPU cores
2020-12-08T13:36:51.056Z|00003|reconnect|INFO|unix:C:/openvswitch/var/run/openvswitch/db.sock: connecting...
2020-12-08T13:36:51.056Z|00004|reconnect|INFO|unix:C:/openvswitch/var/run/openvswitch/db.sock: connected
2020-12-08T13:36:51.059Z|00005|dpif_netlink|INFO|The kernel module does not support meters.
```


The first thing to check is the hns endpoints:


```PS C:\Windows\System32> Get-HnsNetwork

ActivityId             : 33FEAB03-0C2E-4D3A-972E-621A9C8AF6DA
AdditionalParams       :
CurrentEndpointCount   : 2
Extensions             : {@{Id=E7C3B2F0-F3C5-48DF-AF2B-10FED6D72E7A; IsEnabled=False; Name=Microsoft Windows Filtering Platform},
                         @{Id=583CC151-73EC-4A6A-8B47-578297AD7623; IsEnabled=False; Name=Open vSwitch Extension},
                         @{Id=E9B59CFA-2BE1-4B21-828F-B6FBDBDDC017; IsEnabled=False; Name=Microsoft Azure VFP Switch Extension},
                         @{Id=EA24CD6C-D17A-4348-9190-09F0D5BE83DD; IsEnabled=True; Name=Microsoft NDIS Capture}}
Flags                  : 0
Health                 : @{AddressNotificationMissedCount=0; AddressNotificationSequenceNumber=0; InterfaceNotificationMissedCount=0;
                         InterfaceNotificationSequenceNumber=0; LastErrorCode=0; LastUpdateTime=132518531507348144; RouteNotificationMissedCount=0;
                         RouteNotificationSequenceNumber=0}
ID                     : 376B63EF-9FE8-4E51-8937-38F453A8F959
IPv6                   : False
LayeredOn              : B7509E1C-5852-4999-A180-A10C1A724220
MacPools               : {@{EndMacAddress=00-15-5D-E2-2F-FF; StartMacAddress=00-15-5D-E2-20-00}}
MaxConcurrentEndpoints : 2
Name                   : 2bd756dfe237f4b71d924c093ebca042c97e9cad9317c24c4be1f1462a96bac5
NatName                : ICS7376ED07-77F0-47E4-857A-A0BE9D801B2A
Policies               : {}
Resources              : @{AdditionalParams=; AllocationOrder=2; Allocators=System.Collections.ArrayList; Health=;
                         ID=33FEAB03-0C2E-4D3A-972E-621A9C8AF6DA; PortOperationTime=0; State=1; SwitchOperationTime=0; VfpOperationTime=0;
                         parentId=FB23740A-CC2E-494F-ABB8-4A2A5226A551}
State                  : 1
Subnets                : {@{AdditionalParams=; AddressPrefix=172.22.16.0/20; GatewayAddress=172.22.16.1; Health=;
                         ID=4FCF370A-A81A-446E-9505-F5D7EF8363E8; Policies=System.Collections.ArrayList; State=0}}
TotalEndpoints         : 2
Type                   : nat
Version                : 38654705667
RunspaceId             : 45ebf9de-33fc-4e9e-b1d1-18fb1f646b1f

ActivityId             : 590029B7-FC99-4491-B040-715FF4F41B92
AdditionalParams       :
CurrentEndpointCount   : 0
Extensions             : {@{Id=E7C3B2F0-F3C5-48DF-AF2B-10FED6D72E7A; IsEnabled=False; Name=Microsoft Windows Filtering Platform},
                         @{Id=583CC151-73EC-4A6A-8B47-578297AD7623; IsEnabled=False; Name=Open vSwitch Extension},
                         @{Id=E9B59CFA-2BE1-4B21-828F-B6FBDBDDC017; IsEnabled=False; Name=Microsoft Azure VFP Switch Extension},
                         @{Id=EA24CD6C-D17A-4348-9190-09F0D5BE83DD; IsEnabled=True; Name=Microsoft NDIS Capture}}
Flags                  : 0
Health                 : @{AddressNotificationMissedCount=0; AddressNotificationSequenceNumber=0; InterfaceNotificationMissedCount=0;
                         InterfaceNotificationSequenceNumber=0; LastErrorCode=0; LastUpdateTime=132518531500996953; RouteNotificationMissedCount=0;
                         RouteNotificationSequenceNumber=0}
ID                     : EE7D24DB-351C-4211-BB7F-F671A5B5AF37
IPv6                   : False
LayeredOn              : B7509E1C-5852-4999-A180-A10C1A724220
MacPools               : {@{EndMacAddress=00-15-5D-9A-EF-FF; StartMacAddress=00-15-5D-9A-E0-00}}
MaxConcurrentEndpoints : 0
Name                   : nat
NatName                : ICS29524AEC-FE14-49F7-84AA-FF6FF7A7734E
Policies               : {}
Resources              : @{AdditionalParams=; AllocationOrder=2; Allocators=System.Collections.ArrayList; Health=;
                         ID=590029B7-FC99-4491-B040-715FF4F41B92; PortOperationTime=0; State=1; SwitchOperationTime=0; VfpOperationTime=0;
                         parentId=FB23740A-CC2E-494F-ABB8-4A2A5226A551}
State                  : 1
Subnets                : {@{AdditionalParams=; AddressPrefix=172.30.208.0/20; GatewayAddress=172.30.208.1; Health=;
                         ID=F83E48D3-FB43-406A-AFB7-92B981379554; Policies=System.Collections.ArrayList; State=0}}
TotalEndpoints         : 0
Type                   : nat
Version                : 38654705667
RunspaceId             : 45ebf9de-33fc-4e9e-b1d1-18fb1f646b1f
```


