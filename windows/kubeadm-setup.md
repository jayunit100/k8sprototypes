# step 1 : Create two VMs

- Create VM1, ubuntu
  - Setup sshd
- Create VM2, windows server 2019, apply updates
  - Go to "control panne" -> Administrative tools and start openssh
	  - You can add a new user with a password, and microsoft will allow that to ssh in automatically
    - Also, install hyperv on windows so that the hcsschim and so on work properly
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
kubectl apply -f https://github.com/jayunit100/k8sprototypes/blob/master/windows/kube-proxy-antrea.yaml
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
```
 New-Item -path $env:SystemDrive\var\lib\kubelet\etc\kubernetes\pki -type SymbolicLink -value  $env:SystemDrive
\etc\kubernetes\pki\ -Force

kubeadm join 10.0.0.42:6443 --token kcwnj0.n93z9n1mwddicrr5     --discovery-token-ca-cert-hash sha256:326c18f332ad604a3174ecce0a97eb26fc5b1ba69709cebec499ee2e93e54b31  
```


Finally, you should see the kube services running.  Note that this can still result in the antrea-agent-windows service being broken , i.e. because it isn't able to access the apiserver.
Ill investiaget this further as we look for the root cause.
