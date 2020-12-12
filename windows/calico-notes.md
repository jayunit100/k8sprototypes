
Install calico... https://docs.projectcalico.org/getting-started/windows-calico/quickstart#install-calico-for-windows ... notes
```
Invoke-WebRequest https://docs.projectcalico.org/scripts/install-calico-windows.ps1 -OutFile c:\install-calico-windows.ps1
c:\install-calico-windows.ps1 -DownloadOnly yes -KubeVersion 1.19.3
```

Then install the deps

```
Install-WindowsFeature RemoteAccess
Install-WindowsFeature RSAT-RemoteAccess-PowerShell
Install-WindowsFeature Routing
```

Then use POWERSHELL 5 NOT 7 ... to run the ./install-calico.ps1 script....


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
