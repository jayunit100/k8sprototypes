# What logs does CAP produce ? 

During the bootstrapping of a CAP* cluster, you might get containers such as the following...
   
   
- capi-controller-manager-_capi-webhook-system_kube-rbac-proxy-.log
- capi-controller-manager--lcx64_capi-webhook-system_manager-.log
- capi-controller-manager--l9vgn_capi-system_kube-rbac-proxy-.log
- capi-controller-manager--l9vgn_capi-system_manager-.log
- capi-kubeadm-bootstrap-controller-manager--5pstb_capi-kubeadm-bootstrap-system_kube-rbac-proxy-.log
- capi-kubeadm-bootstrap-controller-manager--5pstb_capi-kubeadm-bootstrap-system_manager-.log
- capi-kubeadm-bootstrap-controller-manager--4c6z9_capi-webhook-system_kube-rbac-proxy-.log
- capi-kubeadm-bootstrap-controller-manager--4c6z9_capi-webhook-system_manager-.log
- capi-kubeadm-control-plane-controller-manager--mdfvm_capi-webhook-system_kube-rbac-proxy-.log
- capi-kubeadm-control-plane-controller-manager--mdfvm_capi-webhook-system_manager-.log
- capi-kubeadm-control-plane-controller-manager--phl77_capi-kubeadm-control-plane-system_kube-rbac-proxy-.log
- capi-kubeadm-control-plane-controller-manager--phl77_capi-kubeadm-control-plane-system_manager-.log
- capz-controller-manager--kng7m_capz-system_kube-rbac-proxy-.log
- capz-controller-manager--kng7m_capz-system_manager-.log
- capz-controller-manager--9fj76_capi-webhook-system_kube-rbac-proxy-.log
- capz-controller-manager--9fj76_capi-webhook-system_manager-.log
- cert-manager-5c7d4596b-_cert-manager_cert-manager-.log
- cert-manager-cainjector--928rw_cert-manager_cainjector-.log
- cert-manager-webhook--stz8p_cert-manager_cert-manager-.log
- coredns--w9pqf_kube-system_coredns-.log
- coredns--wm5b6_kube-system_coredns-.log
- etcd-tkg-kind--control-plane_kube-system_etcd-.log
- kindnet-9ncml_kube-system_kindnet-cni-.log
- kube-apiserver-tkg-kind--control-plane_kube-system_kube-apiserver-.log
- kube-controller-manager-tkg-kind--control-plane_kube-system_kube-controller-manager-.log
- kube-proxy-9rb8h_kube-system_kube-proxy-.log
- kube-scheduler-tkg-kind-bt775dc6n3gdf1rfcupg-control-plane_kube-system_kube-scheduler-.log
- local-path-provisioner--ds69h_local-path-storage_local-path-provisioner-.log

# Which logs are worth looking at ? 

Lets  break down the information present in all of the containers for a clsuter bootstrap, starting with capi-controller-manager.

## capi controller manager
  - capi-webhook-system
    - kube-rbac-proxy
This container is straightforward.  It listens on port 8443, and creates a cert.
        ```
            2020-09-01T16:13:53.991639179Z stderr F I0901 16:13:53.991280       1 main.go:213] Generating self signed cert as no cert is provided
            2020-09-01T16:14:10.506527477Z stderr F I0901 16:14:10.500954       1 main.go:243] Starting TCP socket on 0.0.0.0:8443
            2020-09-01T16:14:10.508482444Z stderr F I0901 16:14:10.508385       1 main.go:250] Listening securely on 0.0.0.0:8443
        ```
            - webhook-system_manager
        This container takes a minute to startup, and then ends with...
      ```
            2020-09-01T16:14:02.269887863Z stderr F I0901 16:14:02.269776       1 certwatcher.go:83] controller-runtime/certwatcher "msg"="Starting certificate watcher"
      ```
    - capi-system_kube-rbac-proxy
      similar to kube-rbac-proxy, this container just starts up with 3 lines of logging for port 8443
    - capi-system_manager

    Now, in the capi-system-manager logs , we see some informatino which actually is the  ROOT CAUSE of our failure in this example, related to a 1.19 kubernetes installation issue.
```
          2020-09-01T21:58:18.111047133Z stderr F E0901 21:58:18.110931       1 controller.go:248] controller-runtime/controller "msg"="Reconciler error" "error"="admission webhook \"validation.kubeadmcontrolplane.controlplane.cluster.x-k8s.io\" denied the request: KubeadmControlPlane.controlplane.cluster.x-k8s.io \"tkg-mgmt-azure-control-plane\" is invalid: spec.kubeadmConfigSpec.clusterConfiguration.dns.imageTag: Forbidden: cannot migrate CoreDNS up to '1.7.0' from '1.7.0'" "controller"="cluster" "name"="tkg-mgmt-azure" "namespace"="tkg-system"

          2020-09-01T21:58:30.541098179Z stderr F I0901 21:58:30.540871       1 machine_controller_noderef.go:52] controllers/Machine "msg"="Machine doesn't have a valid ProviderID yet" "cluster"="tkg-mgmt-azure" "machine"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "namespace"="tkg-system"

          2020-09-01T21:58:30.541158512Z stderr F E0901 21:58:30.540923       1 machine_controller.go:249] controllers/Machine "msg"="Reconciliation for Machine asked to requeue" "error"="Bootstrap provider for Machine \"tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c\" in namespace \"tkg-system\" is not ready, requeuing: requeue in 30s" "cluster"="tkg-mgmt-azure" "machine"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "namespace"="tkg-system"

          ...
          2020-09-01T22:01:00.607203511Z stderr F I0901 22:01:00.607003       1 machine_controller_noderef.go:52] controllers/Machine "msg"="Machine doesn't have a valid ProviderID yet" "cluster"="tkg-mgmt-azure" "machine"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "namespace"="tkg-system"
          2020-09-01T22:01:00.607320883Z stderr F E0901 22:01:00.607139       1 machine_controller.go:249] controllers/Machine "msg"="Reconciliation for Machine asked to requeue" "error"="Bootstrap provider for Machine \"tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c\" in namespace \"tkg-system\" is not ready, requeuing: requeue in 30s" "cluster"="tkg-mgmt-azure" "machine"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "namespace"="tkg-system"
```

## capi-kubeadm-bootstrap-controller-manager-
  - kubeadm-bootstrap-system_kube-rbac-proxy (This is similar to kube-rbac-proxy above, just a 3 liner to serve 8443)
  - capi-kubeadm-bootstrap-system_kube-rbac-proxy (Ditto)
  - capi-kubeadm-bootstrap-system_manager-
        
  ```
        2020-09-01T16:14:03.      1 listener.go:44] controller-runtime/metrics "msg"="metrics server is starting to listen"  "addr"="127.0.0.1:8080"
  
        2020-09-01T16:14:04.     1 main.go:151] setup "msg"="starting manager"  "version"=""
        
        2020-09-01T16:14:04.      1 leaderelection.go:242] attempting to acquire leader lease  capi-kubeadm-bootstrap-system/kubeadm-bootstrap-manager-leader-election-capi...
        
        2020-09-01T16:14:04.      1 internal.go:356] controller-runtime/manager "msg"="starting metrics server"  "path"="/metrics"
        
        2020-09-01T16:14:06.      1 leaderelection.go:252] successfully acquired lease capi-kubeadm-bootstrap-system/kubeadm-bootstrap-manager-leader-election-capi
        
        2020-09-01T16:14:06.      1 controller.go:152] controller-runtime/controller "msg"="Starting EventSource" "controller"="kubeadmconfig" "source"={"Type":{"metadata":{"creationTimestamp":null},"spec":{},"status":{}}}
        2020-09-01T16:14:06.      1 controller.go:152] controller-runtime/controller "msg"="Starting EventSource" "controller"="kubeadmconfig" "source"={"Type":{"metadata":{"creationTimestamp":null},"spec":{"clusterName":"","bootstrap":{},"infrastructureRef":{}},"status":{"
                      bootstrapReady":false,
                      "infrastructureReady":false
        }}}
        2020-09-01T16:14:06.      1 controller.go:152] controller-runtime/controller "msg"="Starting EventSource" "controller"="kubeadmconfig" "source"={"Type":{"metadata":{"creationTimestamp":null},"spec":{"controlPlaneEndpoint":{"host":"","port":0}},"status":{
                      "infrastructureReady":false,
                      "controlPlaneInitialized":false
        }}}
        
        2020-09-01T16:14:06.      1 controller.go:159] controller-runtime/controller "msg"="Starting Controller" "controller"="kubeadmconfig"
        
        2020-09-01T16:14:06.      1 controller.go:180] controller-runtime/controller "msg"="Starting workers" "controller"="kubeadmconfig" "worker count"=10
        
        2020-09-01T16:14:22.      1 kubeadmconfig_controller.go:225] controllers/KubeadmConfig "msg"="Cluster infrastructure is not ready, waiting" "kind"="
        Machine" "kubeadmconfig"={"Namespace":"tkg-system","Name":"tkg-mgmt-azure-md-0-9djj6"} "name"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "version"="1480"
        
        2020-09-01T16:14:22.      1 kubeadmconfig_controller.go:225] controllers/KubeadmConfig "msg"="Cluster infrastructure is not ready, waiting" "kind"="
        Machine" "kubeadmconfig"={"Namespace":"tkg-system","Name":"tkg-mgmt-azure-md-0-9djj6"} "name"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "version"="1491"
        
        2020-09-01T16:14:22.     1 kubeadmconfig_controller.go:225] controllers/KubeadmConfig "msg"="Cluster infrastructure is not ready, waiting" "kind"="
        Machine" "kubeadmconfig"={"Namespace":"tkg-system","Name":"tkg-mgmt-azure-md-0-9djj6"} "name"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "version"="1491"
        
        root@tkg-kind-bt775dainers# cat capi-kubeadm-bootstrap-controller-manager-85b4884469-4c6z9_capi-webhook-system_kube-rbac-proxy-10e19162
        a6cffecdb372afb4bfade5755f2e26e4c1e8082b0eff065b04726040.log
        
        2020-09-01T16:13:54.      1 main.go:213] Generating self signed cert as no cert is provided
        
        2020-09-01T16:14:03.      1 main.go:243] Starting TCP socket on 0.0.0.0:8443
        
        2020-09-01T16:14:03.     1 main.go:250] Listening securely on 0.0.0.0:8443
        ```

Now for the kubeadm control plane: 

## capi-kubeadm-control-plane
   - capi-kubeadm-system_manager
    ```
        2020-09-01T22:42:32.949745883Z stderr F I0901 22:42:32.949552       1 controller.go:133] controllers/KubeadmControlPlane "msg"="Cluster Controller has not yet set OwnerRef" "kubeadmControlPlane"="tkg-mgmt-azure-control-plane" "namespace"="tkg-system"
    ```
And finally for the capz-contrller-manager logs.  In these youll usually see a good hint of the top level errors if you arent seeing VMs being created.  Interestingly the "controllers/AzureJSONTemplate" is appended to the same OwnerRef message below.  Note to self - need to see why that same error message is printed in both:
        - kubeadmcontroplane
        - capz-system-manager

- capz-controller-manager 
  - capz-system_manager
    ```
         F I0901 22:36:22.787327       1 azurejson_machinetemplate_controller.go:82] controllers/AzureJSONTemplate "msg"="Cluster Controller has not yet set OwnerRef" "AzureMachineTemplate"="tkg-mgmt-azure-control-plane" "namespace"="tkg-system"
        ...
         F I0901 22:49:13.059145       1 azurecluster_controller.go:154] controllers/AzureCluster "msg"="Reconciling AzureCluster" "AzureCluster"="tkg-mgmt-azure" "cluster"="tkg-mgmt-azure" "namespace"="tkg-system"
         F I0901 22:49:13.059448       1 azuremachine_controller.go:238] controllers/AzureMachine "msg"="Bootstrap data secret reference is not yet available" "AzureCluster"="tkg-mgmt-azure" "azureMachine"="tkg-mgmt-azure-md-0-j9c6c" "cluster"="tkg-mgmt-azure" "machine"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "namespace"="tkg-system"
         F I0901 22:49:32.487403       1 azuremachine_controller.go:216] controllers/AzureMachine "msg"="Reconciling AzureMachine" "AzureCluster"="tkg-mgmt-azure" "azureMachine"="tkg-mgmt-azure-md-0-j9c6c" "cluster"="tkg-mgmt-azure" "machine"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "namespace"="tkg-system"
         F I0901 22:49:32.488962       1 azuremachine_controller.go:238] controllers/AzureMachine "msg"="Bootstrap data secret reference is not yet available" "AzureCluster"="tkg-mgmt-azure" "azureMachine"="tkg-mgmt-azure-md-0-j9c6c" "cluster"="tkg-mgmt-azure" "machine"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "namespace"="tkg-system"
    ```

## cert-manager
  - cert-manager_cert-manager
    This container has about a 1 minute startup time ?   And then ends something like this... 
    ```
        2020-09-01T16:14:06.786971901Z stderr F I0901 16:14:06.786910       1 sync.go:303] cert-manager/controller/certificates "level"=0 "msg"="certificate does not require re-issuance. certificate renewal scheduled near expiry time." "related_resource_kind"="CertificateRequest" "related_resource_name"="capz-serving-cert-3774028192" "related_resource_namespace"="capi-webhook-system" "resource_kind"="Certificate" "resource_name"="capz-serving-cert" "resource_namespace"="capi-webhook-system"
2020-09-01T16:14:06.787257362Z stderr F I0901 16:14:06.787211       1 controller.go:135] cert-manager/controller/certificates "level"=0 "msg"="finished processing work item" "key"="capi-webhook-system/capz-serving-cert"
    ```

## core-dns
A healthy coredns container starts up like this... 
```
        2020-09-01T16:13:23.227516361Z stdout F .:53
        2020-09-01T16:13:23.227559282Z stdout F [INFO] plugin/reload: Running configuration MD5 = db32ca3650231d74073ff4cf814959a7
        2020-09-01T16:13:23.227573754Z stdout F CoreDNS-1.7.0
        2020-09-01T16:13:23.227578794Z stdout F linux/amd64, go1.15,
```
## api-server
Its notable that the webhook failures can ALSO be detected in the api-server logs ! this might be the next best place to check for failures during bootstrapping.

```
        2020-09-01T22:59:16.350213378Z stderr F W0901 22:59:16.350058       1 dispatcher.go:142] rejected by webhook "validation.kubeadmcontrolplane.controlplane.cluster.x-k8s.io": &errors.StatusError{ErrStatus:v1.Status{TypeMeta:v1.TypeMeta{Kind:"", APIVersion:""}, ListMeta:v1.ListMeta{SelfLink:"", ResourceVersion:"", Continue:"", RemainingItemCount:(*int64)(nil)}, Status:"Failure", Message:"admission webhook \"validation.kubeadmcontrolplane.controlplane.cluster.x-k8s.io\" denied the request: KubeadmControlPlane.controlplane.cluster.x-k8s.io \"tkg-mgmt-azure-control-plane\" is invalid: spec.kubeadmConfigSpec.clusterConfiguration.dns.imageTag: Forbidden: cannot migrate CoreDNS up to '1.7.0' from '1.7.0'", Reason:"KubeadmControlPlane.controlplane.cluster.x-k8s.io \"tkg-mgmt-azure-control-plane\" is invalid: spec.kubeadmConfigSpec.clusterConfiguration.dns.imageTag: Forbidden: cannot migrate CoreDNS up to '1.7.0' from '1.7.0'", Details:(*v1.StatusDetails)(nil), Code:403}}
```
