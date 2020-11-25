apiVersion: cluster.x-k8s.io/v1alpha3
kind: Cluster
metadata:
  labels:
    cluster.x-k8s.io/cluster-name: clusterapi-peri
  name: clusterapi-peri
  namespace: default
spec:
  clusterNetwork:
    pods:
      cidrBlocks:
      - 192.168.0.0/16
  controlPlaneRef:
    apiVersion: controlplane.cluster.x-k8s.io/v1alpha3
    kind: KubeadmControlPlane
    name: clusterapi-peri
  infrastructureRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
    kind: VSphereCluster
    name: clusterapi-peri
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: VSphereCluster
metadata:
  name: clusterapi-peri
  namespace: default
spec:
  cloudProviderConfiguration:
    global:
      insecure: true
      secretName: cloud-provider-vsphere-credentials
      secretNamespace: kube-system
    network:
      name: VM Network
    providerConfig:
      cloud:
        controllerImage: gcr.io/cloud-provider-vsphere/cpi/release/manager:v1.2.1
    virtualCenter:
      $VSPHERE_SERVER:
        datacenters: $VSPHERE_DATACENTER
    workspace:
      datacenter: $VSPHERE_DATACENTER
      datastore: $VSPHERE_DATASTORE
      folder: $VSPHERE_FOLDER
      server: $VSPHERE_SERVER
  controlPlaneEndpoint:
    host: $VSPHERE_CP_IP
    port: 6443
  server: $VSPHERE_SERVER
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: VSphereMachineTemplate
metadata:
  name: clusterapi-peri
  namespace: default
spec:
  template:
    spec:
      cloneMode: linkedClone
      datacenter: $VSPHERE_DATACENTER
      datastore: $VSPHERE_DATASTORE
      diskGiB: 25
      folder: $VSPHERE_FOLDER
      memoryMiB: 8192
      network:
        devices:
        - dhcp4: true
          networkName: $VSPHERE_NETWORK
      numCPUs: 2
      os: Linux
      server: $VSPHERE_SERVER
      template: /kubo-dc/vm/photon-3-kube-v1.19.3
---
apiVersion: controlplane.cluster.x-k8s.io/v1alpha3
kind: KubeadmControlPlane
metadata:
  name: clusterapi-peri
  namespace: default
spec:
  infrastructureTemplate:
    apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
    kind: VSphereMachineTemplate
    name: clusterapi-peri
  kubeadmConfigSpec:
    clusterConfiguration:
      apiServer:
        extraArgs:
          cloud-provider: external
      controllerManager:
        extraArgs:
          cloud-provider: external
    files:
    - content: |
        apiVersion: v1
        kind: Pod
        metadata:
          creationTimestamp: null
          name: kube-vip
          namespace: kube-system
        spec:
          containers:
          - args:
            - start
            env:
            - name: vip_arp
              value: "true"
            - name: vip_leaderelection
              value: "true"
            - name: vip_address
              value: $VSPHERE_CP_IP
            - name: vip_interface
              value: eth0
            - name: vip_leaseduration
              value: "15"
            - name: vip_renewdeadline
              value: "10"
            - name: vip_retryperiod
              value: "2"
            image: registry.tkg.vmware.run/kube-vip:v0.1.8_vmware.1
            imagePullPolicy: IfNotPresent
            name: kube-vip
            resources: {}
            securityContext:
              capabilities:
                add:
                - NET_ADMIN
                - SYS_TIME
            volumeMounts:
            - mountPath: /etc/kubernetes/admin.conf
              name: kubeconfig
          hostNetwork: true
          volumes:
          - hostPath:
              path: /etc/kubernetes/admin.conf
              type: FileOrCreate
            name: kubeconfig
        status: {}
      owner: root:root
      path: /etc/kubernetes/manifests/kube-vip.yaml
    initConfiguration:
      nodeRegistration:
        criSocket: /var/run/containerd/containerd.sock
        kubeletExtraArgs:
          cloud-provider: external
        name: '{{ ds.meta_data.hostname }}'
    joinConfiguration:
      nodeRegistration:
        criSocket: /var/run/containerd/containerd.sock
        kubeletExtraArgs:
          cloud-provider: external
        name: '{{ ds.meta_data.hostname }}'
    preKubeadmCommands:
    - hostname "{{ ds.meta_data.hostname }}"
    - echo "::1         ipv6-localhost ipv6-loopback" >/etc/hosts
    - echo "127.0.0.1   localhost" >>/etc/hosts
    - echo "127.0.0.1   {{ ds.meta_data.hostname }}" >>/etc/hosts
    - echo "{{ ds.meta_data.hostname }}" >/etc/hostname
    useExperimentalRetryJoin: true
    users:
    - name: capv
      sshAuthorizedKeys:
      - $VSPHERE_SSH_AUTHORIZED_KEY
      sudo: ALL=(ALL) NOPASSWD:ALL
  replicas: 1
  version: v1.19.1
---
apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
kind: KubeadmConfigTemplate
metadata:
  name: clusterapi-peri-md-0
  namespace: default
spec:
  template:
    spec:
      joinConfiguration:
        nodeRegistration:
          criSocket: /var/run/containerd/containerd.sock
          kubeletExtraArgs:
            cloud-provider: external
          name: '{{ ds.meta_data.hostname }}'
      preKubeadmCommands:
      - hostname "{{ ds.meta_data.hostname }}"
      - echo "::1         ipv6-localhost ipv6-loopback" >/etc/hosts
      - echo "127.0.0.1   localhost" >>/etc/hosts
      - echo "127.0.0.1   {{ ds.meta_data.hostname }}" >>/etc/hosts
      - echo "{{ ds.meta_data.hostname }}" >/etc/hostname
      users:
      - name: capv
        sshAuthorizedKeys:
        - '$VSPHERE_SSH_AUTHORIZED_KEY'
        sudo: ALL=(ALL) NOPASSWD:ALL
---
apiVersion: cluster.x-k8s.io/v1alpha3
kind: MachineDeployment
metadata:
  labels:
    cluster.x-k8s.io/cluster-name: clusterapi-peri
  name: clusterapi-peri-md-0
  namespace: default
spec:
  clusterName: clusterapi-peri
  replicas: 2
  selector:
    matchLabels: {}
  template:
    metadata:
      labels:
        cluster.x-k8s.io/cluster-name: clusterapi-peri
    spec:
      bootstrap:
        configRef:
          apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
          kind: KubeadmConfigTemplate
          name: clusterapi-peri-md-0
      clusterName: clusterapi-peri
      infrastructureRef:
        apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
        kind: VSphereMachineTemplate
        name: clusterapi-peri
      version: v1.19.1
---
apiVersion: cluster.x-k8s.io/v1alpha3
kind: MachineDeployment
metadata:
  labels:
    cluster.x-k8s.io/cluster-name: clusterapi-peri
  name: clusterapi-peri-md-0-windows-containerd
  namespace: default
spec:
  clusterName: clusterapi-peri
  replicas: 1
  selector:
    matchLabels: {}
  template:
    metadata:
      labels:
        cluster.x-k8s.io/cluster-name: clusterapi-peri
    spec:
      bootstrap:
        configRef:
          apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
          kind: KubeadmConfigTemplate
          name: clusterapi-peri-md-0-windows-containerd
      clusterName: clusterapi-peri
      infrastructureRef:
        apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
        kind: VSphereMachineTemplate
        name: clusterapi-peri-windows-containerd
      version: v1.19.1
---
apiVersion: cluster.x-k8s.io/v1alpha3
kind: MachineDeployment
metadata:
  labels:
    cluster.x-k8s.io/cluster-name: clusterapi-peri
  name: clusterapi-peri-md-0-windows-docker
  namespace: default
spec:
  clusterName: clusterapi-peri
  replicas: 1
  selector:
    matchLabels: {}
  template:
    metadata:
      labels:
        cluster.x-k8s.io/cluster-name: clusterapi-peri
    spec:
      bootstrap:
        configRef:
          apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
          kind: KubeadmConfigTemplate
          name: clusterapi-peri-md-0-windows-docker
      clusterName: clusterapi-peri
      infrastructureRef:
        apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
        kind: VSphereMachineTemplate
        name: clusterapi-peri-windows-docker
      version: v1.19.1
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: VSphereMachineTemplate
metadata:
  name: clusterapi-peri-windows-containerd
  namespace: default
spec:
  template:
    spec:
      cloneMode: linkedClone
      datacenter: $VSPHERE_DATACENTER
      datastore: $VSPHERE_DATASTORE
      diskGiB: 79
      folder: $VSPHERE_FOLDER
      memoryMiB: 4192
      network:
        devices:
        - dhcp4: true
          networkName: $VSPHERE_NETWORK
      numCPUs: 2
      os: Windows
      server: $VSPHERE_SERVER
      template: windows-2019-kube-v1.19.1-containerd
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: VSphereMachineTemplate
metadata:
  name: clusterapi-peri-windows-docker
  namespace: default
spec:
  template:
    spec:
      cloneMode: linkedClone
      datacenter: $VSPHERE_DATACENTER
      datastore: $VSPHERE_DATASTORE
      diskGiB: 79
      folder: $VSPHERE_FOLDER
      memoryMiB: 4192
      network:
        devices:
        - dhcp4: true
          networkName: VM Network
      numCPUs: 2
      os: Windows
      server: $VSPHERE_SERVER
      template: windows-2019-kube-v1.19.1-docker
---
apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
kind: KubeadmConfigTemplate
metadata:
  name: clusterapi-peri-md-0-windows-docker
  namespace: default
spec:
  template:
    spec:
      joinConfiguration:
        nodeRegistration:
          kubeletExtraArgs:
            cloud-provider: external
            register-with-taints: os=windows:NoSchedule
          name: '{{ ds.meta_data.hostname }}'
      files:
      - content: |
          $service = Get-Service -Name ovs-vswitchd -ErrorAction SilentlyContinue
          if($service -ne $null) {
            exit
          }
          invoke-expression "bcdedit /set TESTSIGNING ON"
          New-Item -ItemType Directory -Force -Path C:\k\antrea
          $trigger = New-JobTrigger -AtStartup 
          $options = New-ScheduledJobOption -RunElevated
          $startupScript = "
            `$service = Get-Service -Name ovs-vswitchd -ErrorAction SilentlyContinue
            if(`$service -eq `$null) {  
              Push-Location C:\k\antrea\
              curl.exe -LO `"https://raw.githubusercontent.com/vmware-tanzu/antrea/master/hack/windows/Install-OVS.ps1`"
              curl.exe -LO `"https://raw.githubusercontent.com/vmware-tanzu/antrea/master/hack/windows/Prepare-ServiceInterface.ps1`"
              curl.exe -LO `"https://raw.githubusercontent.com/vmware-tanzu/antrea/master/hack/windows/Clean-AntreaNetwork.ps1`"
              & ./Install-OVS.ps1
              Set-NetFirewallProfile -Profile Domain,Public,Private -Enabled False
              & ./Prepare-ServiceInterface.ps1
            } else {
              Invoke-expression `"C:\k\antrea\Clean-AntreaNetwork.ps1`" 
              Invoke-expression `"C:\k\antrea\Prepare-ServiceInterface.ps1`"
              Restart-service vmms
              Stop-Service ovs-vswitchd
              Stop-Service ovsdb-server
              start-service ovsdb-server
              Start-Sleep -Seconds 5
              Start-Service ovs-vswitchd
            }
            Restart-service kubelet  
          "
          $env:HostIP = (
              Get-NetIPConfiguration |
              Where-Object {
                  $_.IPv4DefaultGateway -ne $null -and $_.NetAdapter.Status -ne "Disconnected"
              }
          ).IPv4Address.IPAddress
          $file = 'C:\var\lib\kubelet\kubeadm-flags.env'
          $newstr="--node-ip=" + $env:HostIP
          $raw = Get-Content -Path $file -TotalCount 1
          $raw = $raw -replace ".$"
          $new = "$($raw) $($newstr)`""
          Set-Content $file $new
          $startupScript | out-file c:\k\antrea\antrea-bootprep.ps1
          Register-ScheduledJob -Name PrepareAntrea -Trigger $trigger -FilePath C:\k\antrea\antrea-bootprep.ps1 -ScheduledJobOption $options
          invoke-expression "docker pull harbor-repo.vmware.com/dockerhub-proxy-cache/sigwindowstools/kube-proxy:v1.19.1"
          invoke-expression "docker pull harbor-repo.vmware.com/dockerhub-proxy-cache/antrea/antrea-windows:v0.10.1"
          start-sleep -seconds 5
          Restart-Computer -Force
          
        path: 'C:\Temp\antrea.ps1'
      postKubeadmCommands:
        - powershell C:/Temp/antrea.ps1 -ExecutionPolicy Bypass

      users:
      - name: capv
        passwd: capv!!
        groups: Administrators
        sshAuthorizedKeys:
        - $VSPHERE_SSH_AUTHORIZED_KEY
        sudo: ALL=(ALL) NOPASSWD:ALL
---
apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
kind: KubeadmConfigTemplate
metadata:
  name: clusterapi-peri-md-0-windows-containerd
  namespace: default
spec:
  template:
    spec:
      joinConfiguration:
        nodeRegistration:
          criSocket: npipe:////./pipe/containerd-containerd
          kubeletExtraArgs:
            cloud-provider: external
            register-with-taints: os=windows:NoSchedule
          name: '{{ ds.meta_data.hostname }}'
      files:
      - content: |
          $service = Get-Service -Name ovs-vswitchd -ErrorAction SilentlyContinue
          if($service -ne $null) {
            exit
          }
          curl -LO https://github.com/Microsoft/SDN/raw/master/Kubernetes/windows/hns.psm1 "c:\k\hns.psm1"
          Import-Module "c:\k\hns.psm1" 
          New-HnsNetwork -Type NAT -Name nat
          invoke-expression "bcdedit /set TESTSIGNING ON"
          New-Item -ItemType Directory -Force -Path C:\k\antrea
          Remove-item "c:\program files\containerd\config.toml"
          $env:path += "c:\program files\containerd\;"
          containerd.exe config default | Out-File "c:\program files\containerd\config.toml" -Encoding ascii
          #config file fixups
          $config = Get-Content "c:\program files\containerd\config.toml"
          $config = $config -replace "bin_dir = (.)*$", "bin_dir = `"c:/opt/cni/bin`""
          $config = $config -replace "conf_dir = (.)*$", "conf_dir = `"c:/etc/cni/net.d`""
          $config | Set-Content "c:\program files\containerd\config.toml" -Force 

          mkdir -Force c:\opt\cni\bin | Out-Null
          mkdir -Force c:\etc\cni\net.d | Out-Null

          $trigger = New-JobTrigger -AtStartup 
          $options = New-ScheduledJobOption -RunElevated
          $startupScript = "
            `$service = Get-Service -Name ovs-vswitchd -ErrorAction SilentlyContinue
            if(`$service -eq `$null) {  
              Push-Location C:\k\antrea\
              curl.exe -LO `"https://raw.githubusercontent.com/vmware-tanzu/antrea/master/hack/windows/Install-OVS.ps1`"
              curl.exe -LO `"https://raw.githubusercontent.com/vmware-tanzu/antrea/master/hack/windows/Prepare-ServiceInterface.ps1`"
              curl.exe -LO `"https://raw.githubusercontent.com/vmware-tanzu/antrea/master/hack/windows/Clean-AntreaNetwork.ps1`"
              & ./Install-OVS.ps1
              Set-NetFirewallProfile -Profile Domain,Public,Private -Enabled False
            } else {
              Invoke-expression `"C:\k\antrea\Clean-AntreaNetwork.ps1`" 
              Restart-service vmms
              Stop-Service ovs-vswitchd
              Stop-Service ovsdb-server
              start-service ovsdb-server
              Start-Sleep -Seconds 5
              Start-Service ovs-vswitchd
            }
            Restart-service kubelet  
          "
          $env:HostIP = (
              Get-NetIPConfiguration |
              Where-Object {
                  $_.IPv4DefaultGateway -ne $null -and $_.NetAdapter.Status -ne "Disconnected"
              }
          ).IPv4Address.IPAddress
          $file = 'C:\var\lib\kubelet\kubeadm-flags.env'
          $newstr="--node-ip=" + $env:HostIP
          $raw = Get-Content -Path $file -TotalCount 1
          $raw = $raw -replace ".$"
          $new = "$($raw) $($newstr)`""
          Set-Content $file $new
          $startupScript | out-file c:\k\antrea\antrea-bootprep.ps1
          Register-ScheduledJob -Name PrepareAntrea -Trigger $trigger -FilePath C:\k\antrea\antrea-bootprep.ps1 -ScheduledJobOption $options
          invoke-expression "ctr.exe images pull -n k8s.io harbor-repo.vmware.com/dockerhub-proxy-cache/sigwindowstools/kube-proxy:v1.19.1"
          invoke-expression "ctr.exe images pull -n k8s.io projects.registry.vmware.com/antrea/antrea-windows:v0.10.1"
          start-sleep -seconds 5
          Restart-Computer -Force
          
        path: 'C:\Temp\antrea.ps1'
      postKubeadmCommands:
        - powershell C:/Temp/antrea.ps1 -ExecutionPolicy Bypass

      users:
      - name: capv
        passwd: capv!!
        groups: Administrators
        sshAuthorizedKeys:
        - 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCeWhZAGsKC7VezEIpVhzX9fDfxZ80/oIbO03DRVrLBlhgGdewRiHcZQDvOsHaEReKsM2HOFFH108s1e2bRbf+ZMyAlSLgq+Q3gTWzUD8Lu3UuW7ANk2CAaiAI4+3cPG+eCUsvKngWliZ4dDaTFv+kJ8mgiAZmfIJK7KRzVUM52V48pWp3+ve+IEG5NkHVvIwbOYtKj5UiQ+/MPKCh7WOyFuvT6bMSRHgfreCOeuu2AydiqkhQIKFBjtqrwutaKs8wXOlK1RzSXJttv31hQGN9+mFanDJNIziqZe9BJztE0p2hFdNt4WqkGsaC/Yqt8eDGW0UqnWn+MuGl7awTQDBoP4JV5tjZ5N9AF00Jgp/zYSL4o7m1jrIAN0Nx0/2GT/UmhdMiCgRL7vxKVNrKFEbkcMxmrgsE/MhMYLnKB9L5amZ3uNvZpNZRyCRjjiXqQAzOm0OPS8MS0G9xBmrNx+joNYK+z5q+1PMmKd9/8h2t96/yTJLBWynKWzXzyj4ra3gscGpmqf+fzkVixziW7pbyi2+e/B+idAS/d3oGq5xixHYn+/PcbUSnMitOmJ1rmQf4HfcHYokxgNnJ5Zeikn7hYjy+tjEKtk8eX2WNtLcuOmJ0IVo5fUWfeph6nWSPwNUCRElVSR26j7xW0+bm0swmZkhnUzl76PTrmfXuyXg+AtQ== perit@vmware.com '
        sudo: ALL=(ALL) NOPASSWD:ALL
---
apiVersion: addons.cluster.x-k8s.io/v1alpha3
kind: ClusterResourceSet
metadata:
  labels:
    cluster.x-k8s.io/cluster-name: clusterapi-peri
  name: clusterapi-peri-crs-0
  namespace: default
spec:
  clusterSelector:
    matchLabels:
      cluster.x-k8s.io/cluster-name: clusterapi-peri
  resources:
  - kind: Secret
    name: vsphere-csi-controller
  - kind: ConfigMap
    name: vsphere-csi-controller-role
  - kind: ConfigMap
    name: vsphere-csi-controller-binding
  - kind: Secret
    name: csi-vsphere-config
  - kind: ConfigMap
    name: csi.vsphere.vmware.com
  - kind: ConfigMap
    name: vsphere-csi-node
  - kind: ConfigMap
    name: vsphere-csi-controller
  - kind: ConfigMap
    name: antrea
  - kind: ConfigMap
    name: antrea-windows-kubeproxy
  - kind: ConfigMap
    name: antrea-windows
  - kind: ConfigMap
    name: windowsruntimeclass
---
apiVersion: v1
kind: Secret
metadata:
  name: vsphere-csi-controller
  namespace: default
stringData:
  data: |
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: vsphere-csi-controller
      namespace: kube-system
type: addons.cluster.x-k8s.io/resource-set
---
apiVersion: v1
data:
  data: |
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
      name: vsphere-csi-controller-role
    rules:
    - apiGroups:
      - storage.k8s.io
      resources:
      - csidrivers
      verbs:
      - create
      - delete
    - apiGroups:
      - ""
      resources:
      - nodes
      - pods
      - secrets
      verbs:
      - get
      - list
      - watch
    - apiGroups:
      - ""
      resources:
      - persistentvolumes
      verbs:
      - get
      - list
      - watch
      - update
      - create
      - delete
      - patch
    - apiGroups:
      - storage.k8s.io
      resources:
      - volumeattachments
      verbs:
      - get
      - list
      - watch
      - update
      - patch
    - apiGroups:
      - ""
      resources:
      - persistentvolumeclaims
      verbs:
      - get
      - list
      - watch
      - update
    - apiGroups:
      - storage.k8s.io
      resources:
      - storageclasses
      - csinodes
      verbs:
      - get
      - list
      - watch
    - apiGroups:
      - ""
      resources:
      - events
      verbs:
      - list
      - watch
      - create
      - update
      - patch
    - apiGroups:
      - coordination.k8s.io
      resources:
      - leases
      verbs:
      - get
      - watch
      - list
      - delete
      - update
      - create
    - apiGroups:
      - snapshot.storage.k8s.io
      resources:
      - volumesnapshots
      verbs:
      - get
      - list
    - apiGroups:
      - snapshot.storage.k8s.io
      resources:
      - volumesnapshotcontents
      verbs:
      - get
      - list
kind: ConfigMap
metadata:
  name: vsphere-csi-controller-role
  namespace: default
---
apiVersion: v1
data:
  data: |
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
      name: vsphere-csi-controller-binding
    roleRef:
      apiGroup: rbac.authorization.k8s.io
      kind: ClusterRole
      name: vsphere-csi-controller-role
    subjects:
    - kind: ServiceAccount
      name: vsphere-csi-controller
      namespace: kube-system
kind: ConfigMap
metadata:
  name: vsphere-csi-controller-binding
  namespace: default
---
apiVersion: v1
kind: Secret
metadata:
  name: csi-vsphere-config
  namespace: default
stringData:
  data: |
    apiVersion: v1
    kind: Secret
    metadata:
      name: csi-vsphere-config
      namespace: kube-system
    stringData:
      csi-vsphere.conf: |+
        [Global]
        insecure-flag = true
        cluster-id = "default/clusterapi-peri"

        [VirtualCenter "$VSPHERE_SERVER"]
        user = "$VSPHERE_USERNAME"
        password = "$VSPHERE_PASSWORD"
        datacenters = "$VSPHERE_DATACENTER"

        [Network]
        public-network = "$VSPHERE_NETWORK"

    type: Opaque
type: addons.cluster.x-k8s.io/resource-set
---
apiVersion: v1
data:
  data: |
    apiVersion: storage.k8s.io/v1
    kind: CSIDriver
    metadata:
      name: csi.vsphere.vmware.com
    spec:
      attachRequired: true
kind: ConfigMap
metadata:
  name: csi.vsphere.vmware.com
  namespace: default
---
apiVersion: v1
data:
  data: |
    apiVersion: apps/v1
    kind: DaemonSet
    metadata:
      name: vsphere-csi-node
      namespace: kube-system
    spec:
      selector:
        matchLabels:
          app: vsphere-csi-node
      template:
        metadata:
          labels:
            app: vsphere-csi-node
            role: vsphere-csi
        spec:
          containers:
          - args:
            - --v=5
            - --csi-address=$(ADDRESS)
            - --kubelet-registration-path=$(DRIVER_REG_SOCK_PATH)
            env:
            - name: ADDRESS
              value: /csi/csi.sock
            - name: DRIVER_REG_SOCK_PATH
              value: /var/lib/kubelet/plugins/csi.vsphere.vmware.com/csi.sock
            image: quay.io/k8scsi/csi-node-driver-registrar:v1.2.0
            lifecycle:
              preStop:
                exec:
                  command:
                  - /bin/sh
                  - -c
                  - rm -rf /registration/csi.vsphere.vmware.com-reg.sock /csi/csi.sock
            name: node-driver-registrar
            resources: {}
            securityContext:
              privileged: true
            volumeMounts:
            - mountPath: /csi
              name: plugin-dir
            - mountPath: /registration
              name: registration-dir
          - env:
            - name: CSI_ENDPOINT
              value: unix:///csi/csi.sock
            - name: X_CSI_MODE
              value: node
            - name: X_CSI_SPEC_REQ_VALIDATION
              value: "false"
            - name: VSPHERE_CSI_CONFIG
              value: /etc/cloud/csi-vsphere.conf
            - name: LOGGER_LEVEL
              value: PRODUCTION
            - name: X_CSI_LOG_LEVEL
              value: INFO
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            image: gcr.io/cloud-provider-vsphere/csi/release/driver:v2.0.0
            livenessProbe:
              failureThreshold: 3
              httpGet:
                path: /healthz
                port: healthz
              initialDelaySeconds: 10
              periodSeconds: 5
              timeoutSeconds: 3
            name: vsphere-csi-node
            ports:
            - containerPort: 9808
              name: healthz
              protocol: TCP
            resources: {}
            securityContext:
              allowPrivilegeEscalation: true
              capabilities:
                add:
                - SYS_ADMIN
              privileged: true
            volumeMounts:
            - mountPath: /etc/cloud
              name: vsphere-config-volume
            - mountPath: /csi
              name: plugin-dir
            - mountPath: /var/lib/kubelet
              mountPropagation: Bidirectional
              name: pods-mount-dir
            - mountPath: /dev
              name: device-dir
          - args:
            - --csi-address=/csi/csi.sock
            image: quay.io/k8scsi/livenessprobe:v1.1.0
            name: liveness-probe
            resources: {}
            volumeMounts:
            - mountPath: /csi
              name: plugin-dir
          dnsPolicy: Default
          nodeSelector:
            kubernetes.io/os: linux
          tolerations:
          - effect: NoSchedule
            operator: Exists
          - effect: NoExecute
            operator: Exists
          volumes:
          - name: vsphere-config-volume
            secret:
              secretName: csi-vsphere-config
          - hostPath:
              path: /var/lib/kubelet/plugins_registry
              type: Directory
            name: registration-dir
          - hostPath:
              path: /var/lib/kubelet/plugins/csi.vsphere.vmware.com/
              type: DirectoryOrCreate
            name: plugin-dir
          - hostPath:
              path: /var/lib/kubelet
              type: Directory
            name: pods-mount-dir
          - hostPath:
              path: /dev
            name: device-dir
      updateStrategy:
        type: RollingUpdate
kind: ConfigMap
metadata:
  name: vsphere-csi-node
  namespace: default
---
apiVersion: v1
data:
  data: |
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: vsphere-csi-controller
      namespace: kube-system
    spec:
      replicas: 1
      selector:
        matchLabels:
          app: vsphere-csi-controller
      strategy:
        type: RollingUpdate
      template:
        metadata:
          labels:
            app: vsphere-csi-controller
            role: vsphere-csi
        spec:
          containers:
          - args:
            - --v=4
            - --timeout=300s
            - --csi-address=$(ADDRESS)
            - --leader-election
            env:
            - name: ADDRESS
              value: /csi/csi.sock
            image: quay.io/k8scsi/csi-attacher:v2.0.0
            name: csi-attacher
            resources: {}
            volumeMounts:
            - mountPath: /csi
              name: socket-dir
          - env:
            - name: CSI_ENDPOINT
              value: unix:///var/lib/csi/sockets/pluginproxy/csi.sock
            - name: X_CSI_MODE
              value: controller
            - name: VSPHERE_CSI_CONFIG
              value: /etc/cloud/csi-vsphere.conf
            - name: LOGGER_LEVEL
              value: PRODUCTION
            - name: X_CSI_LOG_LEVEL
              value: INFO
            image: gcr.io/cloud-provider-vsphere/csi/release/driver:v2.0.0
            lifecycle:
              preStop:
                exec:
                  command:
                  - /bin/sh
                  - -c
                  - rm -rf /var/lib/csi/sockets/pluginproxy/csi.vsphere.vmware.com
            livenessProbe:
              failureThreshold: 3
              httpGet:
                path: /healthz
                port: healthz
              initialDelaySeconds: 10
              periodSeconds: 5
              timeoutSeconds: 3
            name: vsphere-csi-controller
            ports:
            - containerPort: 9808
              name: healthz
              protocol: TCP
            resources: {}
            volumeMounts:
            - mountPath: /etc/cloud
              name: vsphere-config-volume
              readOnly: true
            - mountPath: /var/lib/csi/sockets/pluginproxy/
              name: socket-dir
          - args:
            - --csi-address=$(ADDRESS)
            env:
            - name: ADDRESS
              value: /var/lib/csi/sockets/pluginproxy/csi.sock
            image: quay.io/k8scsi/livenessprobe:v1.1.0
            name: liveness-probe
            resources: {}
            volumeMounts:
            - mountPath: /var/lib/csi/sockets/pluginproxy/
              name: socket-dir
          - args:
            - --leader-election
            env:
            - name: X_CSI_FULL_SYNC_INTERVAL_MINUTES
              value: "30"
            - name: LOGGER_LEVEL
              value: PRODUCTION
            - name: VSPHERE_CSI_CONFIG
              value: /etc/cloud/csi-vsphere.conf
            image: gcr.io/cloud-provider-vsphere/csi/release/syncer:v2.0.0
            name: vsphere-syncer
            resources: {}
            volumeMounts:
            - mountPath: /etc/cloud
              name: vsphere-config-volume
              readOnly: true
          - args:
            - --v=4
            - --timeout=300s
            - --csi-address=$(ADDRESS)
            - --feature-gates=Topology=true
            - --strict-topology
            - --enable-leader-election
            - --leader-election-type=leases
            env:
            - name: ADDRESS
              value: /csi/csi.sock
            image: quay.io/k8scsi/csi-provisioner:v1.4.0
            name: csi-provisioner
            resources: {}
            volumeMounts:
            - mountPath: /csi
              name: socket-dir
          dnsPolicy: Default
          serviceAccountName: vsphere-csi-controller
          nodeSelector:
            kubernetes.io/os: linux
          tolerations:
          - effect: NoSchedule
            key: node-role.kubernetes.io/master
            operator: Exists
          volumes:
          - name: vsphere-config-volume
            secret:
              secretName: csi-vsphere-config
          - hostPath:
              path: /var/lib/csi/sockets/pluginproxy/csi.vsphere.vmware.com
              type: DirectoryOrCreate
            name: socket-dir
kind: ConfigMap
metadata:
  name: vsphere-csi-controller
  namespace: default
---
apiVersion: v1
data:
  data: |
    apiVersion: v1
    data:
      run-script.ps1: |-
        $ErrorActionPreference = "Stop";
        mkdir -force /host/var/lib/kube-proxy/var/run/secrets/kubernetes.io/serviceaccount
        mkdir -force /host/k/kube-proxy

        cp -force /k/kube-proxy/* /host/k/kube-proxy
        cp -force /var/lib/kube-proxy/* /host/var/lib/kube-proxy
        cp -force /var/run/secrets/kubernetes.io/serviceaccount/* /host/var/lib/kube-proxy/var/run/secrets/kubernetes.io/serviceaccount

        wins cli process run --path /k/kube-proxy/kube-proxy.exe --args "--v=3 --config=/var/lib/kube-proxy/config.conf --proxy-mode=userspace --hostname-override=$env:NODE_NAME"

    kind: ConfigMap
    metadata:
      labels:
        app: kube-proxy
      name: kube-proxy-windows
      namespace: kube-system
    ---
    apiVersion: apps/v1
    kind: DaemonSet
    metadata:
      labels:
        k8s-app: kube-proxy
      name: kube-proxy-windows
      namespace: kube-system
    spec:
      selector:
        matchLabels:
          k8s-app: kube-proxy-windows
      template:
        metadata:
          labels:
            k8s-app: kube-proxy-windows
        spec:
          containers:
          - args:
            - -file
            - /var/lib/kube-proxy-windows/run-script.ps1
            command:
            - powershell
            env:
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  apiVersion: v1
                  fieldPath: spec.nodeName
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            image: harbor-repo.vmware.com/dockerhub-proxy-cache/sigwindowstools/kube-proxy:v1.19.1
            name: kube-proxy
            volumeMounts:
            - mountPath: /host
              name: host
            - mountPath: \\.\pipe\rancher_wins
              name: wins
            - mountPath: /var/lib/kube-proxy
              name: kube-proxy
            - mountPath: /var/lib/kube-proxy-windows
              name: kube-proxy-windows
          hostNetwork: true
          nodeSelector:
            kubernetes.io/os: windows
          serviceAccountName: kube-proxy
          tolerations:
          - key: CriticalAddonsOnly
            operator: Exists
          - operator: Exists
          volumes:
          - configMap:
              defaultMode: 420
              name: kube-proxy-windows
            name: kube-proxy-windows
          - configMap:
              name: kube-proxy
            name: kube-proxy
          - hostPath:
              path: /
            name: host
          - hostPath:
              path: \\.\pipe\rancher_wins
            name: wins
      updateStrategy:
        type: RollingUpdate
kind: ConfigMap
metadata:
  name: antrea-windows-kubeproxy
  namespace: default
---
apiVersion: v1
data:
  data: |
    apiVersion: v1
    data:
      Run-AntreaAgent.ps1: |
        $ErrorActionPreference = "Stop"
        # wins will rename the binary when executing it. So we need to copy the binary everytime before running it.
        mkdir -force /host/k/antrea/bin
        cp /k/antrea/bin/* /host/k/antrea/bin/
        C:/k/antrea/utils/wins.exe cli process run --path /k/antrea/bin/antrea-agent.exe --args "--config=/k/antrea/etc/antrea-agent.conf --logtostderr=false --log_dir=/k/antrea/logs/ --alsologtostderr --log_file_max_size=100 --log_file_max_num=4" --envs "KUBERNETES_SERVICE_HOST=$env:KUBERNETES_SERVICE_HOST KUBERNETES_SERVICE_PORT=$env:KUBERNETES_SERVICE_PORT ANTREA_SERVICE_HOST=$env:ANTREA_SERVICE_HOST ANTREA_SERVICE_PORT=$env:ANTREA_SERVICE_PORT NODE_NAME=$env:NODE_NAME"
    kind: ConfigMap
    metadata:
      labels:
        app: antrea
      name: antrea-agent-windows-h7td2mh9gm
      namespace: kube-system
    ---
    apiVersion: v1
    data:
      antrea-agent.conf: |
        # FeatureGates is a map of feature names to bools that enable or disable experimental features.
        featureGates:
        # Enable antrea proxy which provides ServiceLB for in-cluster services in antrea agent.
        # It should be enabled on Windows, otherwise NetworkPolicy will not take effect on
        # Service traffic.
          AntreaProxy: true

        # Name of the OpenVSwitch bridge antrea-agent will create and use.
        # Make sure it doesn't conflict with your existing OpenVSwitch bridges.
        #ovsBridge: br-int

        # Name of the interface antrea-agent will create and use for host <--> pod communication.
        # Make sure it doesn't conflict with your existing interfaces.
        #hostGateway: antrea-gw0

        # Encapsulation mode for communication between Pods across Nodes, supported values:
        # - geneve (default)
        # - vxlan
        # - stt
        #tunnelType: geneve

        # Default MTU to use for the host gateway interface and the network interface of each Pod.
        # If omitted, antrea-agent will discover the MTU of the Node's primary interface and
        # also adjust MTU to accommodate for tunnel encapsulation overhead.
        #defaultMTU: 1450

        # ClusterIP CIDR range for Services. It's required when AntreaProxy is not enabled, and should be
        # set to the same value as the one specified by --service-cluster-ip-range for kube-apiserver. When
        # AntreaProxy is enabled, this parameter is not needed and will be ignored if provided.
        #serviceCIDR: 10.96.0.0/12

        # The port for the antrea-agent APIServer to serve on.
        #apiPort: 10350

        # Enable metrics exposure via Prometheus. Initializes Prometheus metrics listener.
        #enablePrometheusMetrics: false
        antrea-cni.conflist: |
        {
            "cniVersion":"0.3.0",
            "name": "antrea",
            "plugins": [
                {
                    "type": "antrea",
                    "ipam": {
                        "type": "host-local"
                    },
                    "capabilities": {"dns": true}
                }
            ]
        }
    kind: ConfigMap
    metadata:
      labels:
        app: antrea
      name: antrea-windows-config-4f6g849tgk
      namespace: kube-system
    ---
    apiVersion: apps/v1
    kind: DaemonSet
    metadata:
      labels:
        app: antrea
        component: antrea-agent
      name: antrea-agent-windows
      namespace: kube-system
    spec:
      selector:
        matchLabels:
          app: antrea
          component: antrea-agent
      template:
        metadata:
          labels:
            app: antrea
            component: antrea-agent
        spec:
          containers:
          - args:
            - -file
            - /var/lib/antrea-windows/Run-AntreaAgent.ps1
            command:
            - pwsh
            env:
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            image: harbor-repo.vmware.com/dockerhub-proxy-cache/antrea/antrea-windows:v0.10.1
            name: antrea-agent
            volumeMounts:
            - mountPath: /host
              name: host
            - mountPath: \\.\pipe\rancher_wins
              name: wins
            - mountPath: /etc/antrea
              name: antrea-windows-config
            - mountPath: /var/lib/antrea-windows
              name: antrea-agent-windows
            - mountPath: /host/k/antrea/
              name: host-antrea-home
          hostNetwork: true
          initContainers:
          - args:
            - -File
            - /k/antrea/Install-WindowsCNI.ps1
            command:
            - pwsh
            image: harbor-repo.vmware.com/dockerhub-proxy-cache/antrea/antrea-windows:v0.10.1
            name: install-cni
            volumeMounts:
            - mountPath: /etc/antrea
              name: antrea-windows-config
              readOnly: true
            - mountPath: /host/etc/cni/net.d
              name: host-cni-conf
            - mountPath: /host/opt/cni/bin
              name: host-cni-bin
            - mountPath: /host/k/antrea/
              name: host-antrea-home
            - mountPath: /host
              name: host
          nodeSelector:
            kubernetes.io/os: windows
          priorityClassName: system-node-critical
          serviceAccountName: antrea-agent
          tolerations:
          - key: CriticalAddonsOnly
            operator: Exists
          - effect: NoSchedule
            operator: Exists
          volumes:
          - configMap:
              name: antrea-windows-config-4f6g849tgk
            name: antrea-windows-config
          - configMap:
              defaultMode: 420
              name: antrea-agent-windows-h7td2mh9gm
            name: antrea-agent-windows
          - hostPath:
              path: /etc/cni/net.d
              type: DirectoryOrCreate
            name: host-cni-conf
          - hostPath:
              path: /opt/cni/bin
              type: DirectoryOrCreate
            name: host-cni-bin
          - hostPath:
              path: /k/antrea
              type: DirectoryOrCreate
            name: host-antrea-home
          - hostPath:
              path: /
            name: host
          - hostPath:
              path: \\.\pipe\rancher_wins
              type: null
            name: wins
      updateStrategy:
        type: RollingUpdate
kind: ConfigMap
metadata:
  name: antrea-windows
  namespace: default
---
apiVersion: v1
data:
  data: |
    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    metadata:
      labels:
        app: antrea
      name: antreaagentinfos.clusterinformation.antrea.tanzu.vmware.com
    spec:
      group: clusterinformation.antrea.tanzu.vmware.com
      names:
        kind: AntreaAgentInfo
        plural: antreaagentinfos
        shortNames:
        - aai
        singular: antreaagentinfo
      scope: Cluster
      versions:
      - name: v1beta1
        schema:
          openAPIV3Schema:
            type: object
            x-kubernetes-preserve-unknown-fields: true
        served: true
        storage: true
    ---
    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    metadata:
      labels:
        app: antrea
      name: antreacontrollerinfos.clusterinformation.antrea.tanzu.vmware.com
    spec:
      group: clusterinformation.antrea.tanzu.vmware.com
      names:
        kind: AntreaControllerInfo
        plural: antreacontrollerinfos
        shortNames:
        - aci
        singular: antreacontrollerinfo
      scope: Cluster
      versions:
      - name: v1beta1
        schema:
          openAPIV3Schema:
            type: object
            x-kubernetes-preserve-unknown-fields: true
        served: true
        storage: true
    ---
    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    metadata:
      labels:
        app: antrea
      name: clusternetworkpolicies.security.antrea.tanzu.vmware.com
    spec:
      group: security.antrea.tanzu.vmware.com
      names:
        kind: ClusterNetworkPolicy
        plural: clusternetworkpolicies
        shortNames:
        - cnp
        - acnp
        singular: clusternetworkpolicy
      scope: Cluster
      versions:
      - additionalPrinterColumns:
        - description: The Tier to which this ClusterNetworkPolicy belongs to.
          jsonPath: .spec.tier
          name: Tier
          type: string
        - description: The Priority of this ClusterNetworkPolicy relative to other policies.
          format: float
          jsonPath: .spec.priority
          name: Priority
          type: number
        - jsonPath: .metadata.creationTimestamp
          name: Age
          type: date
        name: v1alpha1
        schema:
          openAPIV3Schema:
            properties:
              spec:
                properties:
                  appliedTo:
                    items:
                      properties:
                        namespaceSelector:
                          x-kubernetes-preserve-unknown-fields: true
                        podSelector:
                          x-kubernetes-preserve-unknown-fields: true
                      type: object
                    type: array
                  egress:
                    items:
                      properties:
                        action:
                          enum:
                          - Allow
                          - Drop
                          type: string
                        ports:
                          items:
                            properties:
                              port:
                                x-kubernetes-int-or-string: true
                              protocol:
                                type: string
                            type: object
                          type: array
                        to:
                          items:
                            properties:
                              ipBlock:
                                properties:
                                  cidr:
                                    format: cidr
                                    type: string
                                type: object
                              namespaceSelector:
                                x-kubernetes-preserve-unknown-fields: true
                              podSelector:
                                x-kubernetes-preserve-unknown-fields: true
                            type: object
                          type: array
                      required:
                      - action
                      type: object
                    type: array
                  ingress:
                    items:
                      properties:
                        action:
                          enum:
                          - Allow
                          - Drop
                          type: string
                        from:
                          items:
                            properties:
                              ipBlock:
                                properties:
                                  cidr:
                                    format: cidr
                                    type: string
                                type: object
                              namespaceSelector:
                                x-kubernetes-preserve-unknown-fields: true
                              podSelector:
                                x-kubernetes-preserve-unknown-fields: true
                            type: object
                          type: array
                        ports:
                          items:
                            properties:
                              port:
                                x-kubernetes-int-or-string: true
                              protocol:
                                type: string
                            type: object
                          type: array
                      required:
                      - action
                      type: object
                    type: array
                  priority:
                    format: float
                    maximum: 10000
                    minimum: 1
                    type: number
                  tier:
                    type: string
                type: object
            type: object
        served: true
        storage: true
    ---
    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    metadata:
      labels:
        app: antrea
      name: externalentities.core.antrea.tanzu.vmware.com
    spec:
      group: core.antrea.tanzu.vmware.com
      names:
        kind: ExternalEntity
        plural: externalentities
        shortNames:
        - ee
        singular: externalentity
      scope: Namespaced
      versions:
      - name: v1alpha1
        schema:
          openAPIV3Schema:
            properties:
              spec:
                properties:
                  endpoints:
                    items:
                      properties:
                        ip:
                          format: ipv4
                          type: string
                        name:
                          type: string
                        ports:
                          items:
                            properties:
                              name:
                                type: string
                              port:
                                x-kubernetes-int-or-string: true
                              protocol:
                                type: string
                            type: object
                          type: array
                      type: object
                    type: array
                  externalNode:
                    type: string
                type: object
            type: object
        served: true
        storage: true
    ---
    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    metadata:
      labels:
        app: antrea
      name: networkpolicies.security.antrea.tanzu.vmware.com
    spec:
      group: security.antrea.tanzu.vmware.com
      names:
        kind: NetworkPolicy
        plural: networkpolicies
        shortNames:
        - netpol
        - anp
        singular: networkpolicy
      scope: Namespaced
      versions:
      - additionalPrinterColumns:
        - description: The Tier to which this Antrea NetworkPolicy belongs to.
          jsonPath: .spec.tier
          name: Tier
          type: string
        - description: The Priority of this Antrea NetworkPolicy relative to other policies.
          format: float
          jsonPath: .spec.priority
          name: Priority
          type: number
        - jsonPath: .metadata.creationTimestamp
          name: Age
          type: date
        name: v1alpha1
        schema:
          openAPIV3Schema:
            properties:
              spec:
                properties:
                  appliedTo:
                    items:
                      properties:
                        podSelector:
                          type: object
                          x-kubernetes-preserve-unknown-fields: true
                      type: object
                    type: array
                  egress:
                    items:
                      properties:
                        action:
                          enum:
                          - Allow
                          - Drop
                          type: string
                        ports:
                          items:
                            properties:
                              port:
                                x-kubernetes-int-or-string: true
                              protocol:
                                type: string
                            type: object
                          type: array
                        to:
                          items:
                            properties:
                              externalEntitySelector:
                                x-kubernetes-preserve-unknown-fields: true
                              ipBlock:
                                properties:
                                  cidr:
                                    format: cidr
                                    type: string
                                type: object
                              namespaceSelector:
                                x-kubernetes-preserve-unknown-fields: true
                              podSelector:
                                x-kubernetes-preserve-unknown-fields: true
                            type: object
                          type: array
                      required:
                      - action
                      type: object
                    type: array
                  ingress:
                    items:
                      properties:
                        action:
                          enum:
                          - Allow
                          - Drop
                          type: string
                        from:
                          items:
                            properties:
                              externalEntitySelector:
                                x-kubernetes-preserve-unknown-fields: true
                              ipBlock:
                                properties:
                                  cidr:
                                    format: cidr
                                    type: string
                                type: object
                              namespaceSelector:
                                x-kubernetes-preserve-unknown-fields: true
                              podSelector:
                                x-kubernetes-preserve-unknown-fields: true
                            type: object
                          type: array
                        ports:
                          items:
                            properties:
                              port:
                                x-kubernetes-int-or-string: true
                              protocol:
                                type: string
                            type: object
                          type: array
                      required:
                      - action
                      type: object
                    type: array
                  priority:
                    format: float
                    maximum: 10000
                    minimum: 1
                    type: number
                  tier:
                    type: string
                required:
                - appliedTo
                - priority
                type: object
            type: object
        served: true
        storage: true
    ---
    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    metadata:
      labels:
        app: antrea
      name: tiers.security.antrea.tanzu.vmware.com
    spec:
      group: security.antrea.tanzu.vmware.com
      names:
        kind: Tier
        plural: tiers
        shortNames:
        - tr
        singular: tier
      scope: Cluster
      versions:
      - additionalPrinterColumns:
        - description: The Priority of this Tier relative to other Tiers.
          jsonPath: .spec.priority
          name: Priority
          type: integer
        - jsonPath: .metadata.creationTimestamp
          name: Age
          type: date
        name: v1alpha1
        schema:
          openAPIV3Schema:
            properties:
              spec:
                properties:
                  description:
                    type: string
                  priority:
                    maximum: 255
                    minimum: 0
                    type: integer
                required:
                - priority
                type: object
            type: object
        served: true
        storage: true
    ---
    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    metadata:
      labels:
        app: antrea
      name: traceflows.ops.antrea.tanzu.vmware.com
    spec:
      group: ops.antrea.tanzu.vmware.com
      names:
        kind: Traceflow
        plural: traceflows
        shortNames:
        - tf
        singular: traceflow
      scope: Cluster
      versions:
      - additionalPrinterColumns:
        - description: The phase of the Traceflow.
          jsonPath: .status.phase
          name: Phase
          type: string
        - description: The name of the source Pod.
          jsonPath: .spec.source.pod
          name: Source-Pod
          priority: 10
          type: string
        - description: The name of the destination Pod.
          jsonPath: .spec.destination.pod
          name: Destination-Pod
          priority: 10
          type: string
        - description: The IP address of the destination.
          jsonPath: .spec.destination.ip
          name: Destination-IP
          priority: 10
          type: string
        - jsonPath: .metadata.creationTimestamp
          name: Age
          type: date
        name: v1alpha1
        schema:
          openAPIV3Schema:
            properties:
              spec:
                properties:
                  destination:
                    oneOf:
                    - required:
                      - pod
                      - namespace
                    - required:
                      - service
                      - namespace
                    - required:
                      - ip
                    properties:
                      ip:
                        format: ipv4
                        type: string
                      namespace:
                        type: string
                      pod:
                        type: string
                      service:
                        type: string
                    type: object
                  packet:
                    properties:
                      ipHeader:
                        properties:
                          flags:
                            type: integer
                          protocol:
                            type: integer
                          srcIP:
                            format: ipv4
                            type: string
                          ttl:
                            type: integer
                        type: object
                      transportHeader:
                        properties:
                          icmp:
                            properties:
                              id:
                                type: integer
                              sequence:
                                type: integer
                            type: object
                          tcp:
                            properties:
                              dstPort:
                                type: integer
                              flags:
                                type: integer
                              srcPort:
                                type: integer
                            type: object
                          udp:
                            properties:
                              dstPort:
                                type: integer
                              srcPort:
                                type: integer
                            type: object
                        type: object
                    type: object
                  source:
                    properties:
                      namespace:
                        type: string
                      pod:
                        type: string
                    required:
                    - pod
                    - namespace
                    type: object
                required:
                - source
                - destination
                type: object
              status:
                properties:
                  dataplaneTag:
                    type: integer
                  phase:
                    type: string
                  reason:
                    type: string
                  results:
                    items:
                      properties:
                        node:
                          type: string
                        observations:
                          items:
                            properties:
                              action:
                                type: string
                              component:
                                type: string
                              componentInfo:
                                type: string
                              dstMAC:
                                type: string
                              networkPolicy:
                                type: string
                              pod:
                                type: string
                              translatedDstIP:
                                type: string
                              translatedSrcIP:
                                type: string
                              ttl:
                                type: integer
                              tunnelDstIP:
                                type: string
                            type: object
                          type: array
                        role:
                          type: string
                        timestamp:
                          type: integer
                      type: object
                    type: array
                type: object
            required:
            - spec
            type: object
        served: true
        storage: true
        subresources:
          status: {}
    ---
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      labels:
        app: antrea
      name: antctl
      namespace: kube-system
    ---
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      labels:
        app: antrea
      name: antrea-agent
      namespace: kube-system
    ---
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      labels:
        app: antrea
      name: antrea-controller
      namespace: kube-system
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
      labels:
        app: antrea
        rbac.authorization.k8s.io/aggregate-to-admin: "true"
        rbac.authorization.k8s.io/aggregate-to-edit: "true"
      name: aggregate-antrea-policies-edit
    rules:
    - apiGroups:
      - security.antrea.tanzu.vmware.com
      resources:
      - clusternetworkpolicies
      - networkpolicies
      verbs:
      - get
      - list
      - watch
      - create
      - update
      - patch
      - delete
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
      labels:
        app: antrea
        rbac.authorization.k8s.io/aggregate-to-view: "true"
      name: aggregate-antrea-policies-view
    rules:
    - apiGroups:
      - security.antrea.tanzu.vmware.com
      resources:
      - clusternetworkpolicies
      - networkpolicies
      verbs:
      - get
      - list
      - watch
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
      labels:
        app: antrea
        rbac.authorization.k8s.io/aggregate-to-admin: "true"
        rbac.authorization.k8s.io/aggregate-to-edit: "true"
      name: aggregate-traceflows-edit
    rules:
    - apiGroups:
      - ops.antrea.tanzu.vmware.com
      resources:
      - traceflows
      verbs:
      - get
      - list
      - watch
      - create
      - update
      - patch
      - delete
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
      labels:
        app: antrea
        rbac.authorization.k8s.io/aggregate-to-view: "true"
      name: aggregate-traceflows-view
    rules:
    - apiGroups:
      - ops.antrea.tanzu.vmware.com
      resources:
      - traceflows
      verbs:
      - get
      - list
      - watch
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
      labels:
        app: antrea
      name: antctl
    rules:
    - apiGroups:
      - controlplane.antrea.tanzu.vmware.com
      - networking.antrea.tanzu.vmware.com
      resources:
      - networkpolicies
      - appliedtogroups
      - addressgroups
      verbs:
      - get
      - list
    - apiGroups:
      - stats.antrea.tanzu.vmware.com
      resources:
      - networkpolicystats
      - antreaclusternetworkpolicystats
      - antreanetworkpolicystats
      verbs:
      - get
      - list
    - apiGroups:
      - system.antrea.tanzu.vmware.com
      resources:
      - controllerinfos
      - agentinfos
      verbs:
      - get
    - apiGroups:
      - system.antrea.tanzu.vmware.com
      resources:
      - supportbundles
      verbs:
      - get
      - post
    - apiGroups:
      - system.antrea.tanzu.vmware.com
      resources:
      - supportbundles/download
      verbs:
      - get
    - nonResourceURLs:
      - /agentinfo
      - /addressgroups
      - /appliedtogroups
      - /networkpolicies
      - /ovsflows
      - /ovstracing
      - /podinterfaces
      verbs:
      - get
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
      labels:
        app: antrea
      name: antrea-agent
    rules:
    - apiGroups:
      - ""
      resources:
      - nodes
      verbs:
      - get
      - watch
      - list
    - apiGroups:
      - ""
      resources:
      - pods
      - endpoints
      - services
      verbs:
      - get
      - watch
      - list
    - apiGroups:
      - clusterinformation.antrea.tanzu.vmware.com
      resources:
      - antreaagentinfos
      verbs:
      - get
      - create
      - update
      - delete
    - apiGroups:
      - controlplane.antrea.tanzu.vmware.com
      - networking.antrea.tanzu.vmware.com
      resources:
      - networkpolicies
      - appliedtogroups
      - addressgroups
      verbs:
      - get
      - watch
      - list
    - apiGroups:
      - controlplane.antrea.tanzu.vmware.com
      resources:
      - nodestatssummaries
      verbs:
      - create
    - apiGroups:
      - authentication.k8s.io
      resources:
      - tokenreviews
      verbs:
      - create
    - apiGroups:
      - authorization.k8s.io
      resources:
      - subjectaccessreviews
      verbs:
      - create
    - apiGroups:
      - ""
      resourceNames:
      - extension-apiserver-authentication
      resources:
      - configmaps
      verbs:
      - get
      - list
      - watch
    - apiGroups:
      - ""
      resourceNames:
      - antrea-ca
      resources:
      - configmaps
      verbs:
      - get
      - watch
      - list
    - apiGroups:
      - ops.antrea.tanzu.vmware.com
      resources:
      - traceflows
      - traceflows/status
      verbs:
      - get
      - watch
      - list
      - update
      - patch
      - create
      - delete
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
      labels:
        app: antrea
      name: antrea-controller
    rules:
    - apiGroups:
      - ""
      resources:
      - nodes
      - pods
      - namespaces
      verbs:
      - get
      - watch
      - list
    - apiGroups:
      - networking.k8s.io
      resources:
      - networkpolicies
      verbs:
      - get
      - watch
      - list
    - apiGroups:
      - clusterinformation.antrea.tanzu.vmware.com
      resources:
      - antreacontrollerinfos
      verbs:
      - get
      - create
      - update
      - delete
    - apiGroups:
      - clusterinformation.antrea.tanzu.vmware.com
      resources:
      - antreaagentinfos
      verbs:
      - list
      - delete
    - apiGroups:
      - authentication.k8s.io
      resources:
      - tokenreviews
      verbs:
      - create
    - apiGroups:
      - authorization.k8s.io
      resources:
      - subjectaccessreviews
      verbs:
      - create
    - apiGroups:
      - ""
      resourceNames:
      - extension-apiserver-authentication
      resources:
      - configmaps
      verbs:
      - get
      - list
      - watch
    - apiGroups:
      - ""
      resourceNames:
      - antrea-ca
      resources:
      - configmaps
      verbs:
      - get
      - update
    - apiGroups:
      - apiregistration.k8s.io
      resourceNames:
      - v1alpha1.stats.antrea.tanzu.vmware.com
      - v1beta1.system.antrea.tanzu.vmware.com
      - v1beta1.controlplane.antrea.tanzu.vmware.com
      - v1beta1.networking.antrea.tanzu.vmware.com
      resources:
      - apiservices
      verbs:
      - get
      - update
    - apiGroups:
      - admissionregistration.k8s.io
      resourceNames:
      - crdvalidator.antrea.tanzu.vmware.com
      resources:
      - validatingwebhookconfigurations
      verbs:
      - get
      - update
    - apiGroups:
      - security.antrea.tanzu.vmware.com
      resources:
      - clusternetworkpolicies
      - networkpolicies
      verbs:
      - get
      - watch
      - list
    - apiGroups:
      - security.antrea.tanzu.vmware.com
      resources:
      - tiers
      verbs:
      - get
      - watch
      - list
      - create
    - apiGroups:
      - ops.antrea.tanzu.vmware.com
      resources:
      - traceflows
      - traceflows/status
      verbs:
      - get
      - watch
      - list
      - update
      - patch
      - create
      - delete
    - apiGroups:
      - core.antrea.tanzu.vmware.com
      resources:
      - externalentities
      verbs:
      - get
      - watch
      - list
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
      labels:
        app: antrea
      name: antctl
      namespace: kube-system
    roleRef:
      apiGroup: rbac.authorization.k8s.io
      kind: ClusterRole
      name: antctl
    subjects:
    - kind: ServiceAccount
      name: antctl
      namespace: kube-system
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
      labels:
        app: antrea
      name: antrea-agent
    roleRef:
      apiGroup: rbac.authorization.k8s.io
      kind: ClusterRole
      name: antrea-agent
    subjects:
    - kind: ServiceAccount
      name: antrea-agent
      namespace: kube-system
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
      labels:
        app: antrea
      name: antrea-controller
    roleRef:
      apiGroup: rbac.authorization.k8s.io
      kind: ClusterRole
      name: antrea-controller
    subjects:
    - kind: ServiceAccount
      name: antrea-controller
      namespace: kube-system
    ---
    apiVersion: v1
    kind: ConfigMap
    metadata:
      labels:
        app: antrea
      name: antrea-ca
      namespace: kube-system
    ---
    apiVersion: v1
    data:
      antrea-agent.conf: |
        # FeatureGates is a map of feature names to bools that enable or disable experimental features.
        featureGates:
        # Enable AntreaProxy which provides ServiceLB for in-cluster Services in antrea-agent.
        # It should be enabled on Windows, otherwise NetworkPolicy will not take effect on
        # Service traffic.
        #  AntreaProxy: false

        # Enable traceflow which provides packet tracing feature to diagnose network issue.
        #  Traceflow: false

        # Enable Antrea ClusterNetworkPolicy feature to complement K8s NetworkPolicy for cluster admins
        # to define security policies which apply to the entire cluster, and Antrea NetworkPolicy
        # feature that supports priorities, rule actions and externalEntities in the future.
        #  AntreaPolicy: false

        # Enable flowexporter which exports polled conntrack connections as IPFIX flow records from each agent to a configured collector.
        #  FlowExporter: false

        # Enable collecting and exposing NetworkPolicy statistics.
        #  NetworkPolicyStats: false

        # Name of the OpenVSwitch bridge antrea-agent will create and use.
        # Make sure it doesn't conflict with your existing OpenVSwitch bridges.
        #ovsBridge: br-int

        # Datapath type to use for the OpenVSwitch bridge created by Antrea. Supported values are:
        # - system
        # - netdev
        # 'system' is the default value and corresponds to the kernel datapath. Use 'netdev' to run
        # OVS in userspace mode. Userspace mode requires the tun device driver to be available.
        #ovsDatapathType: system

        # Name of the interface antrea-agent will create and use for host <--> pod communication.
        # Make sure it doesn't conflict with your existing interfaces.
        #hostGateway: antrea-gw0

        # Encapsulation mode for communication between Pods across Nodes, supported values:
        # - geneve (default)
        # - vxlan
        # - gre
        # - stt
        #tunnelType: geneve

        # Default MTU to use for the host gateway interface and the network interface of each Pod.
        # If omitted, antrea-agent will discover the MTU of the Node's primary interface and
        # also adjust MTU to accommodate for tunnel encapsulation overhead (if applicable).
        #defaultMTU: 1450

        # Whether or not to enable IPsec encryption of tunnel traffic. IPsec encryption is only supported
        # for the GRE tunnel type.
        #enableIPSecTunnel: false

        # ClusterIP CIDR range for Services. It's required when AntreaProxy is not enabled, and should be
        # set to the same value as the one specified by --service-cluster-ip-range for kube-apiserver. When
        # AntreaProxy is enabled, this parameter is not needed and will be ignored if provided.
        #serviceCIDR: 10.96.0.0/12

        # Determines how traffic is encapsulated. It has the following options
        # encap(default): Inter-node Pod traffic is always encapsulated and Pod to outbound traffic is masqueraded.
        # noEncap: Inter-node Pod traffic is not encapsulated, but Pod to outbound traffic is masqueraded.
        #          Underlying network must be capable of supporting Pod traffic across IP subnet.
        # hybrid: noEncap if worker Nodes on same subnet, otherwise encap.
        # networkPolicyOnly: Antrea enforces NetworkPolicy only, and utilizes CNI chaining and delegates Pod IPAM and connectivity to primary CNI.
        #
        #trafficEncapMode: encap

        # The port for the antrea-agent APIServer to serve on.
        # Note that if it's set to another value, the `containerPort` of the `api` port of the
        # `antrea-agent` container must be set to the same value.
        #apiPort: 10350

        # Enable metrics exposure via Prometheus. Initializes Prometheus metrics listener.
        #enablePrometheusMetrics: false

        # Provide flow collector address as string with format <IP>:<port>[:<proto>], where proto is tcp or udp. This also enables
        # the flow exporter that sends IPFIX flow records of conntrack flows on OVS bridge. If no L4 transport proto is given,
        # we consider tcp as default.
        #flowCollectorAddr: ""

        # Provide flow poll interval as a duration string. This determines how often the flow exporter dumps connections from the conntrack module.
        # Flow poll interval should be greater than or equal to 1s (one second).
        # Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
        #flowPollInterval: "5s"

        # Provide flow export frequency, which is the number of poll cycles elapsed before flow exporter exports flow records to
        # the flow collector.
        # Flow export frequency should be greater than or equal to 1.
        #flowExportFrequency: 12
      antrea-cni.conflist: |
        {
            "cniVersion":"0.3.0",
            "name": "antrea",
            "plugins": [
                {
                    "type": "antrea",
                    "ipam": {
                        "type": "host-local"
                    }
                },
                {
                    "type": "portmap",
                    "capabilities": {"portMappings": true}
                }
            ]
        }
      antrea-controller.conf: |
        # FeatureGates is a map of feature names to bools that enable or disable experimental features.
        featureGates:
        # Enable traceflow which provides packet tracing feature to diagnose network issue.
        #  Traceflow: false

        # Enable Antrea ClusterNetworkPolicy feature to complement K8s NetworkPolicy for cluster admins
        # to define security policies which apply to the entire cluster, and Antrea NetworkPolicy
        # feature that supports priorities, rule actions and externalEntities in the future.
        #  AntreaPolicy: false

        # Enable collecting and exposing NetworkPolicy statistics.
        #  NetworkPolicyStats: false

        # The port for the antrea-controller APIServer to serve on.
        # Note that if it's set to another value, the `containerPort` of the `api` port of the
        # `antrea-controller` container must be set to the same value.
        #apiPort: 10349

        # Enable metrics exposure via Prometheus. Initializes Prometheus metrics listener.
        #enablePrometheusMetrics: false

        # Indicates whether to use auto-generated self-signed TLS certificate.
        # If false, A Secret named "antrea-controller-tls" must be provided with the following keys:
        #   ca.crt: <CA certificate>
        #   tls.crt: <TLS certificate>
        #   tls.key: <TLS private key>
        # And the Secret must be mounted to directory "/var/run/antrea/antrea-controller-tls" of the
        # antrea-controller container.
        #selfSignedCert: true
    kind: ConfigMap
    metadata:
      annotations: {}
      labels:
        app: antrea
      name: antrea-config-f4kt4bdh8t
      namespace: kube-system
    ---
    apiVersion: v1
    kind: Service
    metadata:
      labels:
        app: antrea
      name: antrea
      namespace: kube-system
    spec:
      ports:
      - port: 443
        protocol: TCP
        targetPort: api
      selector:
        app: antrea
        component: antrea-controller
    ---
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      labels:
        app: antrea
        component: antrea-controller
      name: antrea-controller
      namespace: kube-system
    spec:
      replicas: 1
      selector:
        matchLabels:
          app: antrea
          component: antrea-controller
      strategy:
        type: Recreate
      template:
        metadata:
          labels:
            app: antrea
            component: antrea-controller
        spec:
          containers:
          - args:
            - --config
            - /etc/antrea/antrea-controller.conf
            - --logtostderr=false
            - --log_dir=/var/log/antrea
            - --alsologtostderr
            - --log_file_max_size=100
            - --log_file_max_num=4
            - --v=0
            command:
            - antrea-controller
            env:
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            image: harbor-repo.vmware.com/dockerhub-proxy-cache/antrea/antrea-ubuntu:v0.10.1
            name: antrea-controller
            ports:
            - containerPort: 10349
              name: api
              protocol: TCP
            readinessProbe:
              failureThreshold: 5
              httpGet:
                host: 127.0.0.1
                path: /healthz
                port: api
                scheme: HTTPS
              initialDelaySeconds: 5
              periodSeconds: 10
              timeoutSeconds: 5
            resources:
              requests:
                cpu: 200m
            volumeMounts:
            - mountPath: /etc/antrea/antrea-controller.conf
              name: antrea-config
              readOnly: true
              subPath: antrea-controller.conf
            - mountPath: /var/run/antrea/antrea-controller-tls
              name: antrea-controller-tls
            - mountPath: /var/log/antrea
              name: host-var-log-antrea
          hostNetwork: true
          nodeSelector:
            kubernetes.io/os: linux
          priorityClassName: system-cluster-critical
          serviceAccountName: antrea-controller
          tolerations:
          - key: CriticalAddonsOnly
            operator: Exists
          - effect: NoSchedule
            key: node-role.kubernetes.io/master
          volumes:
          - configMap:
              name: antrea-config-f4kt4bdh8t
            name: antrea-config
          - name: antrea-controller-tls
            secret:
              defaultMode: 256
              optional: true
              secretName: antrea-controller-tls
          - hostPath:
              path: /var/log/antrea
              type: DirectoryOrCreate
            name: host-var-log-antrea
    ---
    apiVersion: apiregistration.k8s.io/v1
    kind: APIService
    metadata:
      labels:
        app: antrea
      name: v1alpha1.stats.antrea.tanzu.vmware.com
    spec:
      group: stats.antrea.tanzu.vmware.com
      groupPriorityMinimum: 100
      service:
        name: antrea
        namespace: kube-system
      version: v1alpha1
      versionPriority: 100
    ---
    apiVersion: apiregistration.k8s.io/v1
    kind: APIService
    metadata:
      labels:
        app: antrea
      name: v1beta1.controlplane.antrea.tanzu.vmware.com
    spec:
      group: controlplane.antrea.tanzu.vmware.com
      groupPriorityMinimum: 100
      service:
        name: antrea
        namespace: kube-system
      version: v1beta1
      versionPriority: 100
    ---
    apiVersion: apiregistration.k8s.io/v1
    kind: APIService
    metadata:
      labels:
        app: antrea
      name: v1beta1.networking.antrea.tanzu.vmware.com
    spec:
      group: networking.antrea.tanzu.vmware.com
      groupPriorityMinimum: 100
      service:
        name: antrea
        namespace: kube-system
      version: v1beta1
      versionPriority: 100
    ---
    apiVersion: apiregistration.k8s.io/v1
    kind: APIService
    metadata:
      labels:
        app: antrea
      name: v1beta1.system.antrea.tanzu.vmware.com
    spec:
      group: system.antrea.tanzu.vmware.com
      groupPriorityMinimum: 100
      service:
        name: antrea
        namespace: kube-system
      version: v1beta1
      versionPriority: 100
    ---
    apiVersion: apps/v1
    kind: DaemonSet
    metadata:
      labels:
        app: antrea
        component: antrea-agent
      name: antrea-agent
      namespace: kube-system
    spec:
      selector:
        matchLabels:
          app: antrea
          component: antrea-agent
      template:
        metadata:
          labels:
            app: antrea
            component: antrea-agent
        spec:
          containers:
          - args:
            - --config
            - /etc/antrea/antrea-agent.conf
            - --logtostderr=false
            - --log_dir=/var/log/antrea
            - --alsologtostderr
            - --log_file_max_size=100
            - --log_file_max_num=4
            - --v=0
            command:
            - antrea-agent
            env:
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            image: harbor-repo.vmware.com/dockerhub-proxy-cache/antrea/antrea-ubuntu:v0.10.1
            livenessProbe:
              exec:
                command:
                - /bin/sh
                - -c
                - container_liveness_probe agent
              failureThreshold: 5
              initialDelaySeconds: 5
              periodSeconds: 10
              timeoutSeconds: 5
            name: antrea-agent
            ports:
            - containerPort: 10350
              name: api
              protocol: TCP
            readinessProbe:
              failureThreshold: 5
              httpGet:
                host: 127.0.0.1
                path: /healthz
                port: api
                scheme: HTTPS
              initialDelaySeconds: 5
              periodSeconds: 10
              timeoutSeconds: 5
            resources:
              requests:
                cpu: 200m
            securityContext:
              privileged: true
            volumeMounts:
            - mountPath: /etc/antrea/antrea-agent.conf
              name: antrea-config
              readOnly: true
              subPath: antrea-agent.conf
            - mountPath: /var/run/antrea
              name: host-var-run-antrea
            - mountPath: /var/run/openvswitch
              name: host-var-run-antrea
              subPath: openvswitch
            - mountPath: /var/lib/cni
              name: host-var-run-antrea
              subPath: cni
            - mountPath: /var/log/antrea
              name: host-var-log-antrea
            - mountPath: /host/proc
              name: host-proc
              readOnly: true
            - mountPath: /host/var/run/netns
              mountPropagation: HostToContainer
              name: host-var-run-netns
              readOnly: true
            - mountPath: /run/xtables.lock
              name: xtables-lock
          - command:
            - start_ovs
            image: harbor-repo.vmware.com/dockerhub-proxy-cache/antrea/antrea-ubuntu:v0.10.1
            livenessProbe:
              exec:
                command:
                - /bin/sh
                - -c
                - timeout 10 container_liveness_probe ovs
              failureThreshold: 5
              initialDelaySeconds: 5
              periodSeconds: 10
              timeoutSeconds: 10
            name: antrea-ovs
            resources:
              requests:
                cpu: 200m
            securityContext:
              capabilities:
                add:
                - SYS_NICE
                - NET_ADMIN
                - SYS_ADMIN
                - IPC_LOCK
            volumeMounts:
            - mountPath: /var/run/openvswitch
              name: host-var-run-antrea
              subPath: openvswitch
            - mountPath: /var/log/openvswitch
              name: host-var-log-antrea
              subPath: openvswitch
          hostNetwork: true
          initContainers:
          - command:
            - install_cni
            image: harbor-repo.vmware.com/dockerhub-proxy-cache/antrea/antrea-ubuntu:v0.10.1
            name: install-cni
            resources:
              requests:
                cpu: 100m
            securityContext:
              capabilities:
                add:
                - SYS_MODULE
            volumeMounts:
            - mountPath: /etc/antrea/antrea-cni.conflist
              name: antrea-config
              readOnly: true
              subPath: antrea-cni.conflist
            - mountPath: /host/etc/cni/net.d
              name: host-cni-conf
            - mountPath: /host/opt/cni/bin
              name: host-cni-bin
            - mountPath: /lib/modules
              name: host-lib-modules
              readOnly: true
            - mountPath: /sbin/depmod
              name: host-depmod
              readOnly: true
          nodeSelector:
            kubernetes.io/os: linux
          priorityClassName: system-node-critical
          serviceAccountName: antrea-agent
          tolerations:
          - key: CriticalAddonsOnly
            operator: Exists
          - effect: NoSchedule
            operator: Exists
          - effect: NoExecute
            operator: Exists
          volumes:
          - configMap:
              name: antrea-config-f4kt4bdh8t
            name: antrea-config
          - hostPath:
              path: /etc/cni/net.d
            name: host-cni-conf
          - hostPath:
              path: /opt/cni/bin
            name: host-cni-bin
          - hostPath:
              path: /proc
            name: host-proc
          - hostPath:
              path: /var/run/netns
            name: host-var-run-netns
          - hostPath:
              path: /var/run/antrea
              type: DirectoryOrCreate
            name: host-var-run-antrea
          - hostPath:
              path: /var/log/antrea
              type: DirectoryOrCreate
            name: host-var-log-antrea
          - hostPath:
              path: /lib/modules
            name: host-lib-modules
          - hostPath:
              path: /sbin/depmod
            name: host-depmod
          - hostPath:
              path: /run/xtables.lock
              type: FileOrCreate
            name: xtables-lock
      updateStrategy:
        type: RollingUpdate
    ---
    apiVersion: admissionregistration.k8s.io/v1
    kind: ValidatingWebhookConfiguration
    metadata:
      labels:
        app: antrea
      name: crdvalidator.antrea.tanzu.vmware.com
    webhooks:
    - admissionReviewVersions:
      - v1
      - v1beta1
      clientConfig:
        service:
          name: antrea
          namespace: kube-system
          path: /validate/tier
      name: tiervalidator.antrea.tanzu.vmware.com
      rules:
      - apiGroups:
        - security.antrea.tanzu.vmware.com
        apiVersions:
        - v1alpha1
        operations:
        - CREATE
        - UPDATE
        - DELETE
        resources:
        - tiers
        scope: Cluster
      sideEffects: None
      timeoutSeconds: 5
    - admissionReviewVersions:
      - v1
      - v1beta1
      clientConfig:
        service:
          name: antrea
          namespace: kube-system
          path: /validate/acnp
      name: acnpvalidator.antrea.tanzu.vmware.com
      rules:
      - apiGroups:
        - security.antrea.tanzu.vmware.com
        apiVersions:
        - v1alpha1
        operations:
        - CREATE
        - UPDATE
        resources:
        - clusternetworkpolicies
        scope: Cluster
      sideEffects: None
      timeoutSeconds: 5
    - admissionReviewVersions:
      - v1
      - v1beta1
      clientConfig:
        service:
          name: antrea
          namespace: kube-system
          path: /validate/anp
      name: anpvalidator.antrea.tanzu.vmware.com
      rules:
      - apiGroups:
        - security.antrea.tanzu.vmware.com
        apiVersions:
        - v1alpha1
        operations:
        - CREATE
        - UPDATE
        resources:
        - networkpolicies
        scope: Namespaced
      sideEffects: None
      timeoutSeconds: 5
kind: ConfigMap
metadata:
  name: antrea
  namespace: default  
---
apiVersion: v1
data:
  data: |
    apiVersion: node.k8s.io/v1beta1
    kind: RuntimeClass
    metadata:
      name: windows
    handler: 'docker'
    scheduling:
      nodeSelector:
        kubernetes.io/os: 'windows'
        kubernetes.io/arch: 'amd64'
      tolerations:
      - effect: NoSchedule
        key: os
        operator: Equal
        value: "windows"
kind: ConfigMap
metadata:
  name: windowsruntimeclass
  namespace: default
