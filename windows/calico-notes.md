# How to install calico !

0) Install calico on linux using the YAML file in this repo.
1) Make sure to use powershell 5 for the install-calico.ps1 script, bc it uses an api call not compatible w powershell 7.
2) For kubeconfig make sure you copy the kubeconfig from your linux host to c:/k/config
3) Download calicoctl , ` wget https://github.com/projectcalico/calico/releases/download/v3.17.1/release-v3.17.1.tgz` and unzip calico ctl to /usr/local/bin/ , you need it in order to configure `calicoctl ipam configure --strictaffinity=true`
3) Then you can follow: https://docs.projectcalico.org/getting-started/windows-calico/quickstart#install-calico-for-windows ...  which starts off like so:

```
Invoke-WebRequest https://docs.projectcalico.org/scripts/install-calico-windows.ps1 -OutFile c:\install-calico-windows.ps1
c:\install-calico-windows.ps1 -DownloadOnly yes -KubeVersion 1.19.3
```

And then there are other instructions, namely editing the config.ps1 file, and then running the calico installer.  There is a reboot required.

MAKE SURE TO USE POWERSHELL 5 NOT 7 ... to run the ./install-calico.ps1 script....


# Troubleshooting

Seems like the calico namespace isnt created....  solution is to make sure the kubeconfig file is in c:\\k\\config (no kubeconfig will cause those queries to fail).


```
PS C:\Users\jay> c:\install-calico-windows.ps1 -DownloadOnly yes -KubeVersion 1.19.3 #-CalicoNamespace kube-system
WARNING: The names of some imported commands from the module 'helper' include unapproved verbs that might make them less discoverable. To find the commands with unapproved verbs, run the Import-Module command again with the Verbose parameter. For a list of approved verbs, type Get-Verb.
Downloading CNI binaries
[DownloadFile] File c:\k\cni\host-local.exe already exists.
Downloading Windows Kubernetes scripts
[DownloadFile] File c:\k\hns.psm1 already exists.
[DownloadFile] File c:\k\InstallImages.ps1 already exists.
[DownloadFile] File c:\k\Dockerfile already exists.
WARNING: The names of some imported commands from the module 'hns' include unapproved verbs that might make them less discoverable. To find the commands with unapproved verbs, run the Import-Module command again with the Verbose parameter. For a list of approved verbs, type Get-Verb.
Downloaded [https://dl.k8s.io/v1.19.3/kubernetes-node-windows-amd64.tar.gz] => [C:\Users\jay\AppData\Local\Temp\tmp418A.tar.gz]
Download Calico for Windows release...
[DownloadFile] File c:\calico-windows.zip already exists.
Unzip Calico for Windows release...
Setup Calico for Windows...
error: CreateFile c:\\k\\config: The system cannot find the file specified.
Calico running in kube-system namespace
error: CreateFile c:\\k\\config: The system cannot find the file specified.
Exception: C:\install-calico-windows.ps1:177
Line |
 177 |          throw "$SecretName service account does not exist."
     |          ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
     | calico-node service account does not exist.
```

On startup, succesfull launch looks like this (calico-node.log)

```
2020-12-11 16:59:17.190 [INFO][3948] startup/ipam.go 1804: Ensure block for host win-h0c364gqvjh, ipv4 attr &{3 1 windows-reserved-ipam-handle windows host rsvd} ipv6 attr <nil>
2020-12-11 16:59:17.190 [INFO][3948] startup/ipam.go 1875: Looking up existing affinities for host host="win-h0c364gqvjh"
2020-12-11 16:59:17.194 [INFO][3948] startup/ipam.go 346: Looking up existing affinities for host host="win-h0c364gqvjh"
2020-12-11 16:59:17.198 [INFO][3948] startup/ipam.go 429: Trying affinity for 100.1.1.64/26 host="win-h0c364gqvjh"
2020-12-11 16:59:17.200 [INFO][3948] startup/ipam.go 140: Attempting to load block cidr=100.1.1.64/26 host="win-h0c364gqvjh"
2020-12-11 16:59:17.202 [INFO][3948] startup/ipam.go 217: Affinity is confirmed and block has been loaded cidr=100.1.1.64/26 host="win-h0c364gqvjh"
2020-12-11 16:59:17.202 [INFO][3948] startup/ipam.go 1901: Host's block '100.1.1.64/26'  host="win-h0c364gqvjh"
2020-12-11 16:59:17.212 [INFO][3948] startup/dataplane_windows.go 294: Found existing HNS network [&{Id:D9EA6AC0-D7D3-4C7C-8428-DC588A105369 Name:Calico Type:Overlay NetworkAdapterName: SourceMac: Policies:[[123 34 84 121 112 101 34 58 34 72 111 115 116 82 111 117 116 101 34 125] [123 34 68 101 115 116 105 110 97 116 105 111 110 80 114 101 102 105 120 34 58 34 49 48 48 46 49 46 49 46 49 57 50 47 50 54 34 44 34 68 105 115 116 114 105 98 117 116 101 100 82 111 117 116 101 114 77 97 99 65 100 100 114 101 115 115 34 58 34 54 54 45 49 53 45 99 49 45 50 51 45 97 100 45 100 52 34 44 34 73 115 111 108 97 116 105 111 110 73 100 34 58 52 48 57 54 44 34 80 114 111 118 105 100 101 114 65 100 100 114 101 115 115 34 58 34 49 48 48 46 49 46 49 46 49 34 44 34 84 121 112 101 34 58 34 82 101 109 111 116 101 83 117 98 110 101 116 82 111 117 116 101 34 125]] MacPools:[{StartMacAddress:00-15-5D-9D-70-00 EndMacAddress:00-15-5D-9D-7F-FF}] Subnets:[{AddressPrefix:100.1.1.64/26 GatewayAddress:100.1.1.65 Policies:[[123 34 84 121 112 101 34 58 34 86 83 73 68 34 44 34 86 83 73 68 34 58 52 48 57 54 125]]}] DNSSuffix: DNSServerList: DNSServerCompartment:4 ManagementIP:10.0.0.44 AutomaticDNS:false}] subnet="100.1.1.64/26"
2020-12-11 16:59:17.214 [INFO][3948] startup/dataplane_windows.go 498: Found existing remote HNSEndpoint Calico_ep subnet="100.1.1.64/26"
2020-12-11 16:59:17.214 [INFO][3948] startup/startup_windows.go 108: Ensure network is done.
Calico node initialisation succeeded; monitoring kubelet for restarts...
```


VTEP errors, not sure what this means? 
```
2020-12-11 16:59:09.068 [INFO][5516] felix/vxlan_resolver.go 247: Missing vxlan tunnel address for node, cannot send VTEP yet node="win-h0c364gqvjh"
2020-12-11 17:04:09.253 [INFO][5516] felix/vxlan_resolver.go 247: Missing vxlan tunnel address for node, cannot send VTEP yet node="win-h0c364gqvjh"\
```
