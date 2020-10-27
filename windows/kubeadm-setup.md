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
I1027 11:53:23.896868    7564 log_file.go:99] Set log file max size to 104857600                                                                                                I1027 11:53:23.899290    7564 agent.go:63] Starting Antrea agent (version v0.10.1)                                                                                               I1027 11:53:23.900424    7564 client.go:34] No kubeconfig file was specified. Falling back to in-cluster config                                                                 W1027 11:53:23.904702    7564 env.go:64] Environment variable POD_NAMESPACE not found                                                                                           W1027 11:53:23.907120    7564 cacert_controller.go:79] Failed to get Pod Namespace from environment. Using "kube-system" as the CA ConfigMap Namespace                            I1027 11:53:23.907120    7564 ovs_client.go:67] Connecting to OVSDB at address \\.\pipe\C:openvswitchvarrunopenvswitchdb.sock                                                    I1027 11:53:23.910886    7564 agent.go:197] Setting up node network                                                                                                              E1027 11:53:44.938614    7564 agent.go:567] Failed to get node from K8s with name win-h0c364gqvjh: Get https://100.2.2.1:443/api/v1/nodes/win-h0c364gqvjh: dial tcp 100.2.2.1:443: connectex: A connection attempt failed because the connected party did not properly resp
F1027 11:53:44.941660    7564 main.go:58] Error running agent: error initializing agent: Get https://100.2.2.1:443/api/v1/nodes/win-h0c364gqvjh: dial tcp 100.2.2.1:443: connectex: A connection attempt failed because the connected party did not properly respond after
```

