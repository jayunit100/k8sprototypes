# End to end k8s windows setup guide (local VMs)

This is based on the guide at https://github.com/vmware-tanzu/antrea/blob/master/docs/windows.md. so read that first and make sure you understand it.

This is done in VMWare workstation.  It can be generically followed in other hypervisors as well.     This is an opinionated guide:

- Uses antrea for networking, which relies on Wins
- 2 node cluster
- Only tested in VMWare workstation
- Tested with free trial license for microsoft windows

# step 1 : Create two VMs

- Create VM1, ubuntu
  - Setup sshd
- Create VM2, windows server 2019, apply updates
  - Setup sshd service on windows.  You can do this from powershell via `Restart-service sshd`
  - Enable virtualize intel vt-x/EPT and Virtualize CPU performance counters under "virtual machine settings" so HyperV can work. 
  - Go to "control panne" -> Administrative tools and start openssh
	  - You can add a new user with a password, and microsoft will allow that to ssh in automatically
    - Also, install hyperv on windows so that the hcsschim and so on work properly like so
    ```
      Install-WindowsFeature -Name Hyper-V -IncludeManagementTools -Restart
      # Now install docker enterprise
      Install-Package -Name docker -ProviderName DocekrMsftProvider
    ```
    
# step 2 : Setup a single node kubeadm on ubuntu

- download kubeadm by installing it per the directions on website
    ```
    sudo apt-get update && sudo apt-get install -y apt-transport-https curl
      curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
      cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
        deb https://apt.kubernetes.io/ kubernetes-xenial main
      EOF
    sudo apt-get update
    sudo apt-get install -y kubelet kubeadm kubectl
    sudo apt-mark hold kubelet kubeadm kubectl
     modprobe br_netfilter
     echo '1' > /proc/sys/net/bridge/bridge-nf-call-iptables
     history
     kubeadm init --pod-network-cidr=100.1.1.0/24 --service-cidr=100.2.2.0/24
    ```
  
 
     And now put the basic utils :
     ```
      mkdir -p $HOME/.kube
   sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
     sudo chown $(id -u):$(id -g) $HOME/.kube/config    
     ```

And now install your CNI provider:
```
 kubectl apply -f https://github.com/vmware-tanzu/antrea/releases/download/v0.10.1/antrea.yml
```

Now, we apply a antrea kube-proxy that'll schedule to our windows nodes in the next step
- hostNetwork: true
- uses userspace `win cli process run` arguments.
```
kubectl apply -f https://raw.githubusercontent.com/jayunit100/k8sprototypes/master/windows/kube-proxy-antrea.yaml
kubectl apply -f https://github.com/vmware-tanzu/antrea/releases/download/v0.10.1/antrea-windows.yml
```

 Now... on your windows node ! 

# step 3 : Setup your dang windows node 

Start the OVS installer from antrea which is unsigened.  This requires hacking bcedit... 
```
Bcdedit.exe -set TESTSIGNING ON
Restart-Computer
curl.exe -LO https://raw.githubusercontent.com/vmware-tanzu/antrea/master/hack/windows/Install-OVS.ps1
.\Install-OVS.ps1
```

Now we prepare the nodes k8s stuff: 
```
curl.exe -LO https://github.com/kubernetes-sigs/sig-windows-tools/releases/latest/download/PrepareNode.ps1
.\PrepareNode.ps1 -KubernetesVersion v1.18.0
```

And then try to run kubeadm join, first we also cleanup a possible symlink issue which results in kubelet failing to startup bc of missing pki folder.

## the tricky part 
Now if your env isnt setup right, this is where things might fall down
.. You may run into issues here, the first TWO commands are only need to be run if it fails, and your retrying

- On your linux node, et the join token using `kubeadm token create --print-join-token`.

- Now on your windows node, clean up the kubelet if needed:
```
kubeadm reset # <-- dont run this unless it failed the last time
 New-Item -path $env:SystemDrive\var\lib\kubelet\etc\kubernetes\pki -type SymbolicLink -value  $env:SystemDrive
\etc\kubernetes\pki\ -Force
```
- And finally, join the cluster by copying the above output from the `kubeadm token create... print-join-token` step.  i.e. run something like this:

```
kubeadm join 10.0.0.42:6443 --token kcwnj0.n93z9n1mwddicrr5     --discovery-token-ca-cert-hash sha256:326c18f332ad604a3174ecce0a97eb26fc5b1ba69709cebec499ee2e93e54b31  
```

Finally, you should see the kube services running.  Note that this can still result in the antrea-agent-windows service being broken , i.e. because it isn't able to access the apiserver.

``` 
root@ubuntuk8s:/home/jayunit100# kubectl get pods --all-namespaces
NAMESPACE     NAME                                READY   STATUS             RESTARTS   AGE
kube-system   antrea-agent-vl55p                  2/2     Running            0          165m
kube-system   antrea-agent-windows-hfzd5          0/1     CrashLoopBackOff   4          5m13s <-- interesting issue... investigation below
kube-system   antrea-controller-c59795d4d-sdk9p   1/1     Running            0          165m
kube-system   coredns-f9fd979d6-7zrq7             1/1     Running            0          3h7m
kube-system   coredns-f9fd979d6-zt76c             1/1     Running            0          3h7m
kube-system   etcd-ubuntuk8s                      1/1     Running            0          3h7m
kube-system   kube-apiserver-ubuntuk8s            1/1     Running            0          3h7m
kube-system   kube-controller-manager-ubuntuk8s   1/1     Running            0          3h7m
kube-system   kube-proxy-pvxcr                    1/1     Running            0          3h7m
kube-system   kube-proxy-windows-l9r6d            1/1     Running            0          5m13s
kube-system   kube-scheduler-ubuntuk8s            1/1     Running            0          3h7m     
 ``` 

In the above scenario , it appears that the antrea agent might fail to come up upon accessing the apiserver service:

```
I1027 11:53:23.896868    7564 log_file.go:99] Set log file max size to 104857600                                                                                 I1027 11:53:23.899290    7564 agent.go:63] Starting Antrea agent (version v0.10.1)                                                                               I1027 11:53:23.900424    7564 client.go:34] No kubeconfig file was specified. Falling back to in-cluster config                                                   W1027 11:53:23.904702    7564 env.go:64] Environment variable POD_NAMESPACE not found                                                                             1027 11:53:23.907120    7564 cacert_controller.go:79] Failed to get Pod Namespace from environment. Using "kube-system" as the CA ConfigMap Namespace             I1027 11:53:23.907120    7564 ovs_client.go:67] Connecting to OVSDB at address \\.\pipe\C:openvswitchvarrunopenvswitchdb.sock                                     I1027 11:53:23.910886    7564 agent.go:197] Setting up node network                                                                                               E1027 11:53:44.938614    7564 agent.go:567] Failed to get node from K8s with name win-h0c364gqvjh: Get https://100.2.2.1:443/api/v1/nodes/win-h0c364gqvjh: dial tcp 100.2.2.1:443: connectex: A connection attempt failed because the connected party did not properly resp
F1027 11:53:44.941660    7564 main.go:58] Error running agent: error initializing agent: Get https://100.2.2.1:443/api/v1/nodes/win-h0c364gqvjh: dial tcp 100.2.2.1:443: connectex: A connection attempt failed because the connected party did not properly respond after
```

# Troubleshooting !

## Kubeadm reset / init cycle results in broken cluster

After kubeadm reset you need to reset a symlink due to a cleanup error, so if you run kubeadm reset on windows, also do:

```
 New-Item -path $env:SystemDrive\var\lib\kubelet\etc\kubernetes\pki -type SymbolicLink -value  $env:SystemDrive\etc\kubernetes\pki\ -Force
```


## Cant start windows containers bc of cni-install issues


The install-cni command below needs to talk to the hcsship CreateComputeSystem.  If windows is buggy or doesnt have up to date OS libraries, you can get this.  Make sure you run a cumulative update such as https://www.catalog.update.microsoft.com/Search.aspx?q=KB4577668 (which worked for me), but other updates may work. 
```
  Normal   Scheduled  2m49s               default-scheduler  Successfully assigned kube-system/antrea-agent-windows-rtfct to win-h0c364gqvjh
  Normal   Pulling    2m48s               kubelet            Pulling image "antrea/antrea-windows:v0.10.1"
  Normal   Pulled     2m2s                kubelet            Successfully pulled image "antrea/antrea-windows:v0.10.1"
  Normal   Created    31s (x5 over 2m2s)  kubelet            Created container install-cni
  Warning  Failed     31s (x5 over 2m)    kubelet            Error: failed to start container "install-cni": Error response from daemon: hcsshim::CreateComputeSystem install-cni: The parameter is incorrect.
(extra info: {"SystemType":"Container","Name":"install-cni","Owner":"docker","VolumePath":"\\\\?\\Volume{dabbf3fa-a41a-4388-aa50-36d529403c0b}","IgnoreFlushesDuringBoot":true,"LayerFolderPath":"C:\\ProgramData\\docker\\windowsfilter\\install-cni","Layers":[{"ID":"9c2cb086-af18-56d1-9024-e59158a354a1","Path":"C:\\ProgramData\\docker\\windowsfilter\\6e2f677d79dd67832b11ab81068fe66b0b121d13bc55d54eb90ab5e2c0c333b2"},{"ID":"745864ff-1dd8-5b50-a989-
... 
d4623fb7affe","Path":"C:\\ProgramData\\docker\\windowsfilter\\87831842756b6c9519993c33e7352491f42a9c7997f0ce239b5ee96df9d06022"}],"HostName":"6543ae24fe43","MappedDirectories":[{"HostPath":"c:\\","ContainerPath":"c:\\host","ReadOnly":false,"BandwidthMaximum":0,"IOPSMaximum":0,"CreateInUtilityVM":false},{"HostPath":"c:\\var\\lib\\kubelet\\pods\\cc149550-96d2-4b87-a158-d9ab6bae8f01\\volumes\\kubernetes.io~configmap\\antrea-windows-config","ContainerPath":"c:\\etc\\antrea","ReadOnly":true,"BandwidthMaximum":0,"IOPSMaximum":0,"CreateInUtilityVM":false},{"HostPath":"c:\\k\\antrea","ContainerPath":"c:\\host\\k\\antrea","ReadOnly":false,"BandwidthMaximum":0,"IOPSMaximum":0,"CreateInUtilityVM":false},{"HostPath":"c:\\etc\\cni\\net.d","ContainerPath":"c:\\host\\etc\\cni\\net.d","ReadOnly":false,"BandwidthMaximum":0,"IOPSMaximum":0,"CreateInUtilityVM":false},{"HostPath":"c:\\opt\\cni\\bin","ContainerPath":"c:\\host\\opt\\cni\\bin","ReadOnly":false,"BandwidthMaximum":0,"IOPSMaximum":0,"CreateInUtilityVM":false},{"HostPath":"c:\\var\\lib\\kubelet\\pods\\cc149550-96d2-4b87-a158-d9ab6bae8f01\\volumes\\kubernetes.io~secret\\antrea-agent-token-6vtpp","ContainerPath":"c:\\var\\run\\secrets\\kubernetes.io\\serviceaccount","ReadOnly":true,"BandwidthMaximum":0,"IOPSMaximum":0,"CreateInUtilityVM":false}],"HvPartition":false,"NetworkSharedContainerName":"6543ae24fe43b912e5f035c81825bc7c504919ee7e4665b89f969080fe8afc7d","EndpointList":["A546008E-0FF2-41A3-881E-E6ABBF636B15"]})
```
## KCM controller crashing

This can happen if you dont set the pod-cidr correctly for your windows node.  This is a known issue with KCM, its very fragile to misconfigured podCidrs.

## Agent not initialized due to OVS version

In this case youll see  `Error running agent: error initializing agent:` with a

```
instances of the ROOT/StandardCimv2/MSFT_NetAdapterAdvancedPropertySettingData class on the  CIM server: SELECT * FROM
MSFT_NetAdapterAdvancedPropertySettingData  WHERE ((Name LIKE 'br-int')) AND ((RegistryKeyword = 'NetworkAddress')).
```

Error.  This means that you need to set your OVS_VERSION, you can do that with 

```
ovs-vsctl.exe --no-wait set Open_vSwitch . ovs_version=2.14.0

# or if you think you have a different ovs_version, you can do this dynamically... 
# ovs-vsctl.exe --no-wait set Open_vSwitch . ovs_version=$(Get-Item c:\openvswitch\driver\ovsext.sys).VersionInfo.ProductVersion
```
from a powershell terminal

## Concerned that ovs isnt making the right bridges 

Run the `ovs-vsctl.exe` command in your windows powershell administrative terminal , to see what the routes look like
they should look something like this , if not, then you havent setup openvswitc (or antrea) properly

```
PS C:\Windows\system32> ovs-vsctl.exe show
880e9275-12bc-44f7-9b98-a4652a5ba0cb
    Bridge br-int
        datapath_type: system
        Port Ethernet0
            Interface Ethernet0
        Port sonobuoy-c7cd88
            Interface sonobuoy-c7cd88
                type: internal
        Port br-int
            Interface br-int
                type: internal
        Port antrea-gw0
            Interface antrea-gw0
                type: internal
        Port antrea-tun0
            Interface antrea-tun0
                type: geneve
                options: {csum="true", key=flow, local_ip="10.0.0.44", remote_ip=flow}
    ovs_version: "2.14.0"
PS C:\Windows\system32>
```

## IIS or other apps arent reachable bc kube-proxy on windows isnt proxying traffic

In this case, sometimes the HpyerV VMM service must ALWAYS be started "before" the ovs-vswitchdb service.
Thus, restarting the ovs-vswitchdb service should be restarted periodically if your seeint this issue.
FINALLY after both hyperv-vmm, AND ovs-switchdb have come up ... THEN restart the kubelet service.

## unable to read client-cert c:\var\lib\kubelet\pki\kubelet-client-current.pem 

Sometimes the kubelet cant read directories propery.  running `kubeadm.exe reset` followed by `New-Item -path $env:SystemDrive\var\lib\kubelet\etc\kubernetes\pki -type 3SymbolicLink -value  $env:SystemDrive\etc\kubernetes\pki\ -Force`, which cleans up the PKI symlinks , can fix this issue (possibly)

## Kube proxy windows... is it working ?

You can run `netsh dump` in powershell to see wheter or not the kube proxy is working to setup pod routes from services.

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
```
## Concerned that pods arent getting attached to hns

HNS will be able to give you information about how pods are routed. You can view wether pods have valid devices attached by running:

```
PS C:\Users\jayunit100> hnsdiag.exe list all
Networks:
Name             ID
fe83824bf7b9e3989f9abcc8306dde3926a6130c2ede016ffd0f09836f8f8867 74B299E4-84A8-41EE-8EA7-C894522C7027
nat              E0AE102B-F2F7-48CE-A966-3E3DA1455BCF
antrea-hnsnetwork 85B3F02B-1BCC-48D1-92DE-AAC7E9775F33

Endpoints:
Name             ID                                   Virtual Network Name
Ethernet         0c7119e3-3eb4-4c72-a22f-c47450a9340c fe83824bf7b9e3989f9abcc8306dde3926a6130c2ede016ffd0f09836f8f8867  
Ethernet         5ad857af-355f-4c9f-89e8-4b68ee68cfee fe83824bf7b9e3989f9abcc8306dde3926a6130c2ede016ffd0f09836f8f8867  
sonobuoy-8b946a  9610568d-0778-413d-9075-c67709812bfc antrea-hnsnetwork

Namespaces:
ID                                   | Endpoint IDs

LoadBalancers:
ID                                   | Virtual IPs      | Direct IP IDs
PS C:\Users\jayunit100> hnsdiag.exe list all
Networks:
Name             ID
fe83824bf7b9e3989f9abcc8306dde3926a6130c2ede016ffd0f09836f8f8867 74B299E4-84A8-41EE-8EA7-C894522C7027
nat              E0AE102B-F2F7-48CE-A966-3E3DA1455BCF
antrea-hnsnetwork 85B3F02B-1BCC-48D1-92DE-AAC7E9775F33

Endpoints:
Name             ID                                   Virtual Network Name
Ethernet         0c7119e3-3eb4-4c72-a22f-c47450a9340c fe83824bf7b9e3989f9abcc8306dde3926a6130c2ede016ffd0f09836f8f8867  
Ethernet         5ad857af-355f-4c9f-89e8-4b68ee68cfee fe83824bf7b9e3989f9abcc8306dde3926a6130c2ede016ffd0f09836f8f8867  
sonobuoy-8b946a  9610568d-0778-413d-9075-c67709812bfc antrea-hnsnetwork
iis-site-8cd659  c9b218fb-58f5-4fe7-b9ab-72ccc8ef421f antrea-hnsnetwork
iis-site-9582b0  41acce28-1b23-4a4e-ba22-698060f9b417 antrea-hnsnetwork
```


# Miscellaneous notes about windows development

# PowerShell

To hack around on nodes with linux like fluency, you can use powershell.

- `ip a | grep 192` can be replaced with `Get-NetIPAddress | Select-String "192"`

# Windows networking

Things like hostNetwork and hostPort aren't really options.  If you need to access services such
as ingress controllers, you do so by: 

1) Implementing a NodePort service
2) DaemonSet for pods
3) Setting the NodePort Service to `externalTrafficPolicy=local`

# MSFT_NetAdapter issues
 
Sometimes, you might also get another ovs querying issue,  related to MSFT_NetAdapter creation/removal in antrea agent_windows.go.  not sure yet what the workaround is, but this is an example.

```
    Directory: C:\host\k\antrea
Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
d----           12/8/2020  7:03 AM                bin
I1208 07:06:29.129060    6152 log_file.go:99] Set log file max size to 104857600
I1208 07:06:29.134058    6152 agent.go:63] Starting Antrea agent (version v0.10.1)
I1208 07:06:29.134058    6152 client.go:34] No kubeconfig file was specified. Falling back to in-cluster config
W1208 07:06:29.145063    6152 env.go:64] Environment variable POD_NAMESPACE not found
W1208 07:06:29.147063    6152 cacert_controller.go:79] Failed to get Pod Namespace from environment. Using "kube-system" as the CA ConfigMap Namespace
I1208 07:06:29.147063    6152 ovs_client.go:67] Connecting to OVSDB at address \\.\pipe\C:openvswitchvarrunopenvswitchdb.sock
I1208 07:06:29.153063    6152 agent.go:197] Setting up node network
I1208 07:06:29.177073    6152 agent.go:584] Setting Node MTU=1450
I1208 07:06:44.141160    6152 net_windows.go:265] Created HNSNetwork with name antrea-hnsnetwork id BDCA11C4-C0C9-44A6-868A-9082209B14F1
I1208 07:06:44.142155    6152 ovs_client.go:118] Created bridge: 4d2064d3-fb2a-44d9-8088-09b6d36f4e36
F1208 07:06:50.489808    6152 main.go:58] Error running agent: error initializing agent: Remove-NetIPAddress : No matching MSFT_NetIPAddress objects found by CIM query for instances of the
ROOT/StandardCimv2/MSFT_NetIPAddress class on the  CIM server: SELECT * FROM MSFT_NetIPAddress  WHERE ((InterfaceAlias
LIKE 'br-int')) AND ((AddressFamily = 2)). Verify query parameters and retry.
At line:1 char:1
+ Remove-NetIPAddress  -Confirm:$false -AddressFamily IPv4 -InterfaceAl ...
+ ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    + CategoryInfo          : ObjectNotFound: (MSFT_NetIPAddress:String) [Remove-NetIPAddress], CimJobException
    + FullyQualifiedErrorId : CmdletizationQuery_NotFound,Remove-NetIPAddress
: Remove-NetIPAddress  -Confirm:$false -AddressFamily IPv4 -InterfaceAlias br-int
goroutine 1 [running]:
k8s.io/klog.stacks(0xc00047f700, 0xc00028fb00, 0x344, 0x45d)
	C:/gopath/pkg/mod/k8s.io/klog@v1.0.0/klog.go:875 +0xbf
k8s.io/klog.(*loggingT).output(0x3061d60, 0xc000000003, 0xc0000feaf0, 0x2f9dd8d, 0x7, 0x3a, 0x0)
	C:/gopath/pkg/mod/k8s.io/klog@v1.0.0/klog.go:826 +0x337
k8s.io/klog.(*loggingT).printf(0x3061d60, 0x3, 0x1e1dbc8, 0x17, 0xc000741d08, 0x1, 0x1)
	C:/gopath/pkg/mod/k8s.io/klog@v1.0.0/klog.go:707 +0x152
k8s.io/klog.Fatalf(...)
	C:/gopath/pkg/mod/k8s.io/klog@v1.0.0/klog.go:1276
main.newAgentCommand.func1(0xc0001c2500, 0xc0000ef260, 0x0, 0x6)
	C:/antrea/cmd/antrea-agent/main.go:58 +0x218
github.com/spf13/cobra.(*Command).execute(0xc0001c2500, 0xc0000be010, 0x6, 0x7, 0xc0001c2500, 0xc0000be010)
	C:/gopath/pkg/mod/github.com/spf13/cobra@v0.0.5/command.go:830 +0x2b1
github.com/spf13/cobra.(*Command).ExecuteC(0xc0001c2500, 0x2fdb140, 0xc000069f50, 0xc0005d9f50)
	C:/gopath/pkg/mod/github.com/spf13/cobra@v0.0.5/command.go:914 +0x302
github.com/spf13/cobra.(*Command).Execute(...)
	C:/gopath/pkg/mod/github.com/spf13/cobra@v0.0.5/command.go:864
main.main()
	C:/antrea/cmd/antrea-agent/main.go:37 +0x5d
```

### Service proxy up or down  ? 

Sometimes when running antrea, startup will fail bc your kube proxy is misconfigured and antrea agent cant find the apiserver.... 

```W1208 11:14:50.614465     336 env.go:64] Environment variable POD_NAMESPACE not found
W1208 11:14:50.614465     336 cacert_controller.go:79] Failed to get Pod Namespace from environment. Using "kube-system" as the CA ConfigMap Namespace
I1208 11:14:50.614465     336 ovs_client.go:67] Connecting to OVSDB at address \\.\pipe\C:openvswitchvarrunopenvswitchdb.sock
I1208 11:14:50.621349     336 agent.go:197] Setting up node network                            
E1208 11:15:11.638773     336 agent.go:567] Failed to get node from K8s with name win-h0c364gqvjh: Get https://100.2.2.1:443/api/v1/nodes/win-h0c364gqvjh: dial tcp 100.2.2.1:443: connectex: A connection attempt failed because the connected party did not properly respond after a
period of time, or established connection failed because connected host has failed to respond.                                          
F1208 11:15:11.639936     336 main.go:58] Error running agent: error initializing agent: Get https://100.2.2.1:443/api/v1/nodes/win-h0c364gqvjh: dial tcp 100.2.2.1:443: connectex: A connection attempt failed because the connected party did not properly respond after a period of
time, or established connection failed because connected host has failed to respond.
goroutine 1 [running]:
k8s.io/klog.stacks(0xc000569400, 0xc00055c780, 0x16c, 0x1c1)
        C:/gopath/pkg/mod/k8s.io/klog@v1.0.0/klog.go:875 +0xbf
k8s.io/klog.(*loggingT).output(0x3061d60, 0xc000000003, 0xc00011e0e0, 0x2f9dd8d, 0x7, 0x3a, 0x0)

```
To fix this, you need to doublecheck that your kube proxy is running !
Antrea startup: Make sure your kube internal service rules are routing properly !

Your service cidr can be obtained by just doing `kubectl get services -A`.  if its in the, say, `100.1` space, you can look for corresponding kube proxy rules for the default service.

- When kube proxy is down, if you look for service ip addresses, youll see something like this:

```
PS C:\Users\jay> netsh dump | Select-String "100"

set user name = jayunit100 dialin = policy cbpolicy = none
```

i.e., youll see nothing :) 

- Meanwhie, if the kube proxy is up, youll see these rules
```
PS C:\Users\jay> netsh dump | Select-String "100"

add address name="vEthernet (HNS Internal NIC)" address=100.2.2.1 mask=255.0.0.0
add address name="vEthernet (HNS Internal NIC)" address=100.2.2.246 mask=255.0.0.0
add address name="vEthernet (HNS Internal NIC)" address=100.2.2.10 mask=255.0.0.0
set user name = jayunit100 dialin = policy cbpolicy = none
```
