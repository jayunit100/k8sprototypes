## How to loop through a map in powershell

```
$antreaInstallationFiles = @{
      "https://raw.githubusercontent.com/antrea-io/antrea/main/build/yamls/base/conf/antrea-cni.conflist" = "C:\etc\cni\net.d\10-antrea.conflist"
      "https://raw.githubusercontent.com/antrea-io/antrea/main/hack/windows/Install-OVS.ps1" =  "C:\k\antrea\Install-OVS.ps1"
      "https://raw.githubusercontent.com/antrea-io/antrea/main/hack/windows/helper.psm1" = "C:\k\antrea\helper.psm1"
      "https://github.com/antrea-io/antrea/releases/download/v1.1.0/antrea-agent-windows-x86_64.exe" = "C:\opt\cni\bin\antrea.exe"
      "https://github.com/containernetworking/plugins/releases/download/v0.9.1/cni-plugins-windows-amd64-v0.9.1.tgz" = "C:\k\antrea\bin"
      "https://dl.k8s.io/release/v1.21.0/bin/windows/amd64/kubectl.exe" = "C:/k/bin/kubectl.exe"
}

foreach ($theFile in $antreaInstallationFiles.keys) {
  Write-Output "Downloading $theFile if not available..."
  $outPath = $antreaInstallationFiles[$theFile]
  Write-Output("$theFile")
}
```

## One liners

- Tail a file: `Get-Content myfile.txt -Wait`
- Find a file: `Get-Childitem –Path C:\ -Recurse -Name *ctr*`
- List processes:`Get-Process *ovs* | Format-Table -Property Name`
- Chocolatey: `Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))`
- use Vim: `choco install vim -y`
- You can run `ctr --namespace=k8s.io containers list`
- Look at VNIC Events `Get-WinEvent Microsoft-Windows-Hyper-V-VmSwitch-Operational and Microsoft-Windows-Hyper-V-Compute-Operational`
- sort cmd line output:`| Sort-Object`
- look at the containers `hcsdiag.exe list all`
- look at container networks `hnsdiag.exe list all`
- look at all containers `ctr.exe --namespace=k8s.io c ls`
