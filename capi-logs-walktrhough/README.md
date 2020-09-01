# What logs does CAP produce ? 

During the bootstrapping of a CAP* cluster, you might get containers such as the following...
   
   
- capi-controller-manager-5f7cd97899-lcx64_capi-webhook-system_kube-rbac-proxy-d0dcbb52abc9ceabe19cd7e8bd328614354abe442de105f70821fd052db5d3ae.log
- capi-controller-manager-5f7cd97899-lcx64_capi-webhook-system_manager-49ef807d3b1ece77b84c509df4d7884b2dfb92b39b71c4f176729a6522042211.log
- capi-controller-manager-69c99b87fd-l9vgn_capi-system_kube-rbac-proxy-13fa8267c06ee226d514ac9b4d6aa21f72836d3c09ffcfc780652c16e6923478.log
- capi-controller-manager-69c99b87fd-l9vgn_capi-system_manager-0b5b9d42745c79c2700b52890bc067e5c92e485c5a179aeb9aca65934ee790bb.log
- capi-kubeadm-bootstrap-controller-manager-5f4656cd7c-5pstb_capi-kubeadm-bootstrap-system_kube-rbac-proxy-f44c8dc917a0851a7718a4ffccf158e501b1da38e77725d6a57483909aedeeec.log
- capi-kubeadm-bootstrap-controller-manager-5f4656cd7c-5pstb_capi-kubeadm-bootstrap-system_manager-fd49d391e4451f54326f3b3b798ca8e954df62e398ba8b6041105ae4f7f33183.log
- capi-kubeadm-bootstrap-controller-manager-85b4884469-4c6z9_capi-webhook-system_kube-rbac-proxy-10e19162a6cffecdb372afb4bfade5755f2e26e4c1e8082b0eff065b04726040.log
- capi-kubeadm-bootstrap-controller-manager-85b4884469-4c6z9_capi-webhook-system_manager-98272c370d20d828ae4950132f81cc988cce2871ec428a6ff1e75117d30a4aec.log
- capi-kubeadm-control-plane-controller-manager-5c5484d68d-mdfvm_capi-webhook-system_kube-rbac-proxy-7ef94213690ced0eb7e7ebbbd0c4ad5455315f586b0cf31c53c991dc5e4fab87.log
- capi-kubeadm-control-plane-controller-manager-5c5484d68d-mdfvm_capi-webhook-system_manager-fd1b64c97a1271f3acc38a594a27724c8103bd6dbaebe075a3927f5e61cd6769.log
- capi-kubeadm-control-plane-controller-manager-5f87b79f4f-phl77_capi-kubeadm-control-plane-system_kube-rbac-proxy-a457c32a7b97f247f61a8c48a3a6f7d44bd0f833323792bc1614f7ea1638ec3c.log
- capi-kubeadm-control-plane-controller-manager-5f87b79f4f-phl77_capi-kubeadm-control-plane-system_manager-429070fc17a39f7c76321712945c8dba0fb046291d149d4f82ef3fbf44edb010.log
- capz-controller-manager-558b498bf7-kng7m_capz-system_kube-rbac-proxy-6abe90a9ecc64c9eec1d187ec5074b1a396844ec8d7875e22c1a2a3cd2421717.log
- capz-controller-manager-558b498bf7-kng7m_capz-system_manager-3d2152ce423f95b8fbc0966e59bc2034efe5d531e4314eac31285a3b3f676351.log
- capz-controller-manager-596f66cccc-9fj76_capi-webhook-system_kube-rbac-proxy-a96d925863e0b392e05f4b0e6762b562b4590a45c0f0d969d060daba8ccd74af.log
- capz-controller-manager-596f66cccc-9fj76_capi-webhook-system_manager-e30562f0d6ff412b076bbdb63ee931d73a97c330570c6480482d6068ce4e9810.log
- cert-manager-5c7d4596b-4s7zn_cert-manager_cert-manager-596d9611df9d9d55924f4d6bec2ec633a46375cc715303c91424f443c3091d2e.log
- cert-manager-cainjector-664f57c845-928rw_cert-manager_cainjector-b91b8cd8612c32b137ec92a30fbb26547dd682f36839b6658a6f235ccd2f4e07.log
- cert-manager-webhook-66759c9f87-stz8p_cert-manager_cert-manager-55e2a33a9a98f9b0f2d5d7da85fc33ca897127e0dc4acea9b661a341aef74b38.log
- coredns-774fbc4754-w9pqf_kube-system_coredns-03c2aed6a65afd243f5833e01546d4c229dd2254a1493e4db63bac58b0f107aa.log
- coredns-774fbc4754-wm5b6_kube-system_coredns-45582a4882c1339b87f7ad99e3ff4beb51af5ae563581abe5f2ecc54c3659ef9.log
- etcd-tkg-kind-bt775dc6n3gdf1rfcupg-control-plane_kube-system_etcd-4ffbe01cd538827a02aa22241c107894232cd7dd7bfd02282eff7ab9571ed9b9.log
- kindnet-9ncml_kube-system_kindnet-cni-48c952d3ca153e1cdad50b32c05f40c351325a9004751fe175af853ccac7f6ed.log
- kube-apiserver-tkg-kind-bt775dc6n3gdf1rfcupg-control-plane_kube-system_kube-apiserver-46148ac8dde1c7c614c636c721d6a7412129a95903de4308a1c6935056d207c7.log
- kube-controller-manager-tkg-kind-bt775dc6n3gdf1rfcupg-control-plane_kube-system_kube-controller-manager-52e0a0c5efaa14c0f3f014d18043098a80bdd50c1d62d92244404cf4e13f6021.log
- kube-proxy-9rb8h_kube-system_kube-proxy-f5d6f7118ec50bf1b6eaf1ea057286b50160246cca3eac6468bbcdfe1b661b9a.log
- kube-scheduler-tkg-kind-bt775dc6n3gdf1rfcupg-control-plane_kube-system_kube-scheduler-baff068d5fb2271534d5a2ec25c257e573406784eb2768ab7e3e42e6df0fc5ff.log
- local-path-provisioner-8b46957d4-ds69h_local-path-storage_local-path-provisioner-947e10680813178b7498dd4df0ab40579641998ad8c4401d869a7c02289eb634.log

# Which logs are worth looking at ? 

Lets  break down the information present in all of the containers for a clsuter bootstrap, starting with capi-controller-manager.

- capi controller manager
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
- capi-kubeadm-bootstrap-controller-manager-
  - kubeadm-bootstrap-system_kube-rbac-proxy
    This is similar to kube-rbac-proxy above, just a 3 liner to serve 8443.
  - capi-kubeadm-bootstrap-system_kube-rbac-proxy
    Ditto
  - capi-kubeadm-bootstrap-system_manager-
        ```
        2020-09-01T16:14:03.894049966Z stderr F I0901 16:14:03.893879       1 listener.go:44] controller-runtime/metrics "msg"="metrics server is starting to listen"  "addr"="127.0.0.1:8080"
        2020-09-01T16:14:04.07577407Z stderr F I0901 16:14:04.075641       1 main.go:151] setup "msg"="starting manager"  "version"=""
        2020-09-01T16:14:04.076001479Z stderr F I0901 16:14:04.075932       1 leaderelection.go:242] attempting to acquire leader lease  capi-kubeadm-bootstrap-system/kubeadm-bootstrap-manager-leader-election-capi...
        2020-09-01T16:14:04.076454987Z stderr F I0901 16:14:04.076389       1 internal.go:356] controller-runtime/manager "msg"="starting metrics server"  "path"="/metrics"
        2020-09-01T16:14:06.501861635Z stderr F I0901 16:14:06.501734       1 leaderelection.go:252] successfully acquired lease capi-kubeadm-bootstrap-system/kubeadm-bootstrap-manager-leader-election-capi
        2020-09-01T16:14:06.511239559Z stderr F I0901 16:14:06.511120       1 controller.go:152] controller-runtime/controller "msg"="Starting EventSource" "controller"="kubeadmconfig" "source"={"Type":{"metadata":{"creationTimestamp":null},"spec":{},"status":{}}}
        2020-09-01T16:14:06.631750618Z stderr F I0901 16:14:06.631625       1 controller.go:152] controller-runtime/controller "msg"="Starting EventSource" "controller"="kubeadmconfig" "source"={"Type":{"metadata":{"creationTimestamp":null},"spec":{"clusterName":"","bootstrap":{},"infrastructureRef":{}},"status":{"bootstrapReady":false,"infrastructureReady":false}}}
        2020-09-01T16:14:06.739950693Z stderr F I0901 16:14:06.739827       1 controller.go:152] controller-runtime/controller "msg"="Starting EventSource" "controller"="kubeadmconfig" "source"={"Type":{"metadata":{"creationTimestamp":null},"spec":{"controlPlaneEndpoint":{"host":"","port":0}},"status":{"infrastructureReady":false,"controlPlaneInitialized":false}}}
        2020-09-01T16:14:06.840415965Z stderr F I0901 16:14:06.840262       1 controller.go:159] controller-runtime/controller "msg"="Starting Controller" "controller"="kubeadmconfig"
        2020-09-01T16:14:06.840447644Z stderr F I0901 16:14:06.840307       1 controller.go:180] controller-runtime/controller "msg"="Starting workers" "controller"="kubeadmconfig" "worker count"=10
        2020-09-01T16:14:22.644156969Z stderr F I0901 16:14:22.644021       1 kubeadmconfig_controller.go:225] controllers/KubeadmConfig "msg"="Cluster infrastructure is not ready, waiting" "kind"="Machine" "kubeadmconfig"={"Namespace":"tkg-system","Name":"tkg-mgmt-azure-md-0-9djj6"} "name"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "version"="1480"
        2020-09-01T16:14:22.866892667Z stderr F I0901 16:14:22.866708       1 kubeadmconfig_controller.go:225] controllers/KubeadmConfig "msg"="Cluster infrastructure is not ready, waiting" "kind"="Machine" "kubeadmconfig"={"Namespace":"tkg-system","Name":"tkg-mgmt-azure-md-0-9djj6"} "name"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "version"="1491"
        2020-09-01T16:14:22.89235363Z stderr F I0901 16:14:22.892242       1 kubeadmconfig_controller.go:225] controllers/KubeadmConfig "msg"="Cluster infrastructure is not ready, waiting" "kind"="Machine" "kubeadmconfig"={"Namespace":"tkg-system","Name":"tkg-mgmt-azure-md-0-9djj6"} "name"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "version"="1491"
        root@tkg-kind-bt775dc6n3gdf1rfcupg-control-plane:/var/log/containers# cat capi-kubeadm-bootstrap-controller-manager-85b4884469-4c6z9_capi-webhook-system_kube-rbac-proxy-10e19162a6cffecdb372afb4bfade5755f2e26e4c1e8082b0eff065b04726040.log
        2020-09-01T16:13:54.116394862Z stderr F I0901 16:13:54.116154       1 main.go:213] Generating self signed cert as no cert is provided
        2020-09-01T16:14:03.338580618Z stderr F I0901 16:14:03.142016       1 main.go:243] Starting TCP socket on 0.0.0.0:8443
        2020-09-01T16:14:03.34070673Z stderr F I0901 16:14:03.340536       1 main.go:250] Listening securely on 0.0.0.0:8443
        ```

Now for the kubeadm control plane: 

-  capi-kubeadm-control-plane
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
        2020-09-01T22:36:22.787485237Z stderr F I0901 22:36:22.787327       1 azurejson_machinetemplate_controller.go:82] controllers/AzureJSONTemplate "msg"="Cluster Controller has not yet set OwnerRef" "AzureMachineTemplate"="tkg-mgmt-azure-control-plane" "namespace"="tkg-system"
        ...
        2020-09-01T22:49:13.059254532Z stderr F I0901 22:49:13.059145       1 azurecluster_controller.go:154] controllers/AzureCluster "msg"="Reconciling AzureCluster" "AzureCluster"="tkg-mgmt-azure" "cluster"="tkg-mgmt-azure" "namespace"="tkg-system"
        2020-09-01T22:49:13.059502182Z stderr F I0901 22:49:13.059448       1 azuremachine_controller.go:238] controllers/AzureMachine "msg"="Bootstrap data secret reference is not yet available" "AzureCluster"="tkg-mgmt-azure" "azureMachine"="tkg-mgmt-azure-md-0-j9c6c" "cluster"="tkg-mgmt-azure" "machine"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "namespace"="tkg-system"
        2020-09-01T22:49:32.4876033Z stderr F I0901 22:49:32.487403       1 azuremachine_controller.go:216] controllers/AzureMachine "msg"="Reconciling AzureMachine" "AzureCluster"="tkg-mgmt-azure" "azureMachine"="tkg-mgmt-azure-md-0-j9c6c" "cluster"="tkg-mgmt-azure" "machine"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "namespace"="tkg-system"
        2020-09-01T22:49:32.489048512Z stderr F I0901 22:49:32.488962       1 azuremachine_controller.go:238] controllers/AzureMachine "msg"="Bootstrap data secret reference is not yet available" "AzureCluster"="tkg-mgmt-azure" "azureMachine"="tkg-mgmt-azure-md-0-j9c6c" "cluster"="tkg-mgmt-azure" "machine"="tkg-mgmt-azure-md-0-59d6c89dc4-xvz7c" "namespace"="tkg-system"
    ```

- cert-manager
  - cert-manager_cert-manager
    This container has about a 1 minute startup time ?   And then ends something like this... 
    ```
        2020-09-01T16:14:06.786971901Z stderr F I0901 16:14:06.786910       1 sync.go:303] cert-manager/controller/certificates "level"=0 "msg"="certificate does not require re-issuance. certificate renewal scheduled near expiry time." "related_resource_kind"="CertificateRequest" "related_resource_name"="capz-serving-cert-3774028192" "related_resource_namespace"="capi-webhook-system" "resource_kind"="Certificate" "resource_name"="capz-serving-cert" "resource_namespace"="capi-webhook-system"
2020-09-01T16:14:06.787257362Z stderr F I0901 16:14:06.787211       1 controller.go:135] cert-manager/controller/certificates "level"=0 "msg"="finished processing work item" "key"="capi-webhook-system/capz-serving-cert"
    ```

- core-dns
A healthy coredns container starts up like this... 
```
    2020-09-01T16:13:23.227516361Z stdout F .:53
    2020-09-01T16:13:23.227559282Z stdout F [INFO] plugin/reload: Running configuration MD5 = db32ca3650231d74073ff4cf814959a7
    2020-09-01T16:13:23.227573754Z stdout F CoreDNS-1.7.0
    2020-09-01T16:13:23.227578794Z stdout F linux/amd64, go1.15,
```
- api-server
Its notable that the webhook failures can ALSO be detected in the api-server logs ! this might be the next best place to check for failures during bootstrapping.

```
    2020-09-01T22:59:16.350213378Z stderr F W0901 22:59:16.350058       1 dispatcher.go:142] rejected by webhook "validation.kubeadmcontrolplane.controlplane.cluster.x-k8s.io": &errors.StatusError{ErrStatus:v1.Status{TypeMeta:v1.TypeMeta{Kind:"", APIVersion:""}, ListMeta:v1.ListMeta{SelfLink:"", ResourceVersion:"", Continue:"", RemainingItemCount:(*int64)(nil)}, Status:"Failure", Message:"admission webhook \"validation.kubeadmcontrolplane.controlplane.cluster.x-k8s.io\" denied the request: KubeadmControlPlane.controlplane.cluster.x-k8s.io \"tkg-mgmt-azure-control-plane\" is invalid: spec.kubeadmConfigSpec.clusterConfiguration.dns.imageTag: Forbidden: cannot migrate CoreDNS up to '1.7.0' from '1.7.0'", Reason:"KubeadmControlPlane.controlplane.cluster.x-k8s.io \"tkg-mgmt-azure-control-plane\" is invalid: spec.kubeadmConfigSpec.clusterConfiguration.dns.imageTag: Forbidden: cannot migrate CoreDNS up to '1.7.0' from '1.7.0'", Details:(*v1.StatusDetails)(nil), Code:403}}

```
