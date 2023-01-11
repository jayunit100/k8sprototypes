
# TKG Management Cluster Creation (post-kind bootstrapping)

Now, we have some CAPI controllers, running in a `kind` cluster.  Let's create our first real CAPI cluster.  This cluster will eventually become our persistant Mangement Cluster for our TKG distribution, but it has a long way to go before it grows into that.   We will

- Create a CAPI object representing the mgmt cluster
- Wait for the capv controllers to reconcile this, and make VMs in vsphere running k8s
- Install the `tkg-pkg` package onto that cluster
    - Create `PackageInstall`s for all of the packages in `tkg-pkg`
    - Wait for `kapp controller` to install all of the `PackageInstall`
    - Wait for TKG Packages to come alive on the cluster and start running
- Then we'll graduate the cluster: We'll clusterctl move the KIND objects INTO THE management cluster

At that point, our ephemeral kind cluster, which we used as a bootloader for the mgmt cluster, is no longer needed. 

Ok, lets see how this all works:

## Create the CAPI object

Now, we're going to start creating what will be the *real* Management Cluster..... 

```
    [2023-01-10T19:50:48.605Z] Start creating management cluster...
    [2023-01-10T19:50:48.605Z] patch cluster object with operation status: 
    [2023-01-10T19:50:48.606Z] 	{
    [2023-01-10T19:50:48.606Z] 		"metadata": {
    [2023-01-10T19:50:48.606Z] 			"annotations": {
    [2023-01-10T19:50:48.606Z] 				"TKGOperationInfo" : "{\"Operation\":\"Create\",\"OperationStartTimestamp\":\"2023-01-10 19:50:42.118795146 +0000 UTC\",\"OperationTimeout\":1800}",
    [2023-01-10T19:50:48.606Z] 				"TKGOperationLastObservedTimestamp" : "2023-01-10 19:50:42.118795146 +0000 UTC"
    [2023-01-10T19:50:48.606Z] 			}
    [2023-01-10T19:50:48.606Z] 		}
    [2023-01-10T19:50:48.606Z] 	}
    [2023-01-10T19:50:48.606Z] Applying patch to resource tkg-mgmt-vc of type *v1beta1.Cluster ...
    [2023-01-10T19:50:48.606Z] zero or multiple KCP objects found for the given cluster, 0 tkg-mgmt-vc tkg-system, retrying
    [2023-01-10T19:50:52.925Z] zero or multiple KCP objects found for the given cluster, 0 tkg-mgmt-vc tkg-system, retrying
    [2023-01-10T19:51:03.091Z] zero or multiple KCP objects found for the given cluster, 0 tkg-mgmt-vc tkg-system, retrying
    [2023-01-10T19:51:13.277Z] control plane is not available yet, retrying
...

[2023-01-10T19:54:23.036Z] control plane is not available yet, retrying
...

    [2023-01-10T19:54:33.277Z] Management cluster control plane is available, means API server is ready to receive requests
    [2023-01-10T19:54:33.277Z] getting secret for cluster
    [2023-01-10T19:54:33.277Z] waiting for resource tkg-mgmt-vc-kubeconfig of type *v1.Secret to be up and running
    [2023-01-10T19:54:33.277Z] Saving management cluster kubeconfig into /home/kubo/.kube/config
    [2023-01-10T19:54:33.277Z] Installing kapp-controller on management cluster...
    [2023-01-10T19:54:33.277Z] User ConfigValues File: /tmp/2323403141.yaml
    [2023-01-10T19:54:33.277Z] Kapp-controller values-file: /tmp/4151289121.yaml
    [2023-01-10T19:54:33.554Z] Kapp-controller configuration file: /tmp/543269407
    [2023-01-10T19:54:35.004Z] waiting for resource kapp-controller of type *v1.Deployment to be up and running
    [2023-01-10T19:54:36.455Z] pods are not yet running for deployment 'kapp-controller' in namespace 'tkg-system', retrying
    ...
    [2023-01-10T19:55:07.427Z] pods are not yet running for deployment 'kapp-controller' in namespace 'tkg-system', retrying
    [2023-01-10T19:55:12.878Z] Installing providers on management cluster...
    [2023-01-10T19:55:12.878Z] Installing the clusterctl inventory CRD
    ...
    [2023-01-10T19:55:12.878Z] Creating CustomResourceDefinition="providers.clusterctl.cluster.x-k8s.io"
```


We now have some code from cluster API that we call out to : cluster-api/cmd/clusterctl/client/init.go: 
```
	// checks if the cluster already contains a Core provider.
	// if not we consider this the first time init is executed, and thus we enforce the installation of a core provider,
	// a bootstrap provider and a control-plane provider (if not already explicitly requested by the user)
	log.Info("Fetching providers")
```
We can see the result of this command below: 
```
    [2023-01-10T19:55:14.880Z] Fetching providers
```


Followed by https://github.com/kubernetes-sigs/cluster-api/blob/main/cmd/clusterctl/client/repository/metadata_client.go#L76, which goes and fetches many things
```
    [2023-01-10T19:55:14.880Z] Fetching File="core-components.yaml" Provider="cluster-api" Type="CoreProvider" Version="v1.2.8"
    [2023-01-10T19:55:14.880Z] Fetching File="bootstrap-components.yaml" Provider="kubeadm" Type="BootstrapProvider" Version="v1.2.8"
    [2023-01-10T19:55:14.880Z] Fetching File="control-plane-components.yaml" Provider="kubeadm" Type="ControlPlaneProvider" Version="v1.2.8"
    ...
    [2023-01-10T19:55:29.851Z] Creating Deployment="cert-manager-webhook" Namespace="cert-manager"
    [2023-01-10T19:55:30.129Z] Creating MutatingWebhookConfiguration="cert-manager-webhook"
    [2023-01-10T19:55:30.129Z] Creating ValidatingWebhookConfiguration="cert-manager-webhook"
    [2023-01-10T19:55:30.403Z] Waiting for cert-manager to be available...
    [2023-01-10T19:55:30.403Z] Updating Namespace="cert-manager-test"
    [2023-01-10T19:55:30.403Z] Creating Issuer="test-selfsigned" Namespace="cert-manager-test"\
    ...
    23-01-10T20:06:32.719Z] Creating Issuer="test-selfsigned" Namespace="cert-manager-test"
    [2023-01-10T20:06:32.994Z] Creating Certificate="selfsigned-cert" Namespace="cert-manager-test"
    [2023-01-10T20:06:32.994Z] Deleting Namespace="cert-manager-test"
    [2023-01-10T20:06:33.269Z] Deleting Issuer="test-selfsigned" Namespace="cert-manager-test"
    ...
    [2023-01-10T20:07:00.515Z] Creating Issuer="caip-in-cluster-selfsigned-issuer" Namespace="caip-in-cluster-system"
    [2023-01-10T20:07:00.515Z] Creating MutatingWebhookConfiguration="caip-in-cluster-mutating-webhook-configuration"
    [2023-01-10T20:07:00.515Z] Creating ValidatingWebhookConfiguration="caip-in-cluster-validating-webhook-configuration"
```

Still in clusterctl, now https://github.com/kubernetes-sigs/cluster-api/blob/main/cmd/clusterctl/client/cluster/installer.go 
```    
    [2023-01-10T20:07:00.515Z] Creating inventory entry Provider="infrastructure-ipam-in-cluster" Version="v0.1.0" 
    TargetNamespace="caip-in-cluster-system"
```

Now, back in TKG Client.  We have tkg/client/init.go -> We are in the `WaitForProviders` method..... 

```
func (c *TkgClient) WaitForProviders(clusterClient clusterclient.Client, options waitForProvidersOptions) error {
```

##### These logs from from tanzu cli, 

Again , like in the kind cluster, we're going to iterate through all the "installed" components that clusterctl setup....

```
    [2023-01-10T20:07:00.515Z] installed  Component=="cluster-api"  Type=="CoreProvider"  Version=="v1.2.8"
    [2023-01-10T20:07:00.515Z] installed  Component=="kubeadm"  Type=="BootstrapProvider"  Version=="v1.2.8"
    [2023-01-10T20:07:00.515Z] installed  Component=="kubeadm"  Type=="ControlPlaneProvider"  Version=="v1.2.8"
    [2023-01-10T20:07:00.515Z] installed  Component=="vsphere"  Type=="InfrastructureProvider"  Version=="v1.5.1"
    [2023-01-10T20:07:00.515Z] installed  Component=="ipam-in-cluster"  Type=="InfrastructureProvider"  Version=="v0.1.0"
    [2023-01-10T20:07:00.515Z] Waiting for provider bootstrap-kubeadm
    [2023-01-10T20:07:00.515Z] Waiting for provider infrastructure-ipam-in-cluster
    [2023-01-10T20:07:00.515Z] Fetching File="ipam-components.yaml" Provider="ipam-in-cluster" Type="InfrastructureProvider" Version="v0.1.0"
    [2023-01-10T20:07:00.515Z] Fetching File="bootstrap-components.yaml" Provider="kubeadm" Type="BootstrapProvider" Version="v1.2.8"
    [2023-01-10T20:07:00.515Z] Waiting for provider infrastructure-vsphere
    [2023-01-10T20:07:00.515Z] Fetching File="infrastructure-components.yaml" Provider="vsphere" Type="InfrastructureProvider" Version="v1.5.1"
    [2023-01-10T20:07:00.515Z] Waiting for provider cluster-api
    [2023-01-10T20:07:00.515Z] Fetching File="core-components.yaml" Provider="cluster-api" Type="CoreProvider" Version="v1.2.8"
    [2023-01-10T20:07:00.515Z] Waiting for provider control-plane-kubeadm
    [2023-01-10T20:07:00.515Z] Fetching File="control-plane-components.yaml" Provider="kubeadm" Type="ControlPlaneProvider" Version="v1.2.8"
    [2023-01-10T20:07:00.787Z] waiting for resource caip-in-cluster-controller-manager of type *v1.Deployment to be up and running
    [2023-01-10T20:07:00.787Z] pods are not yet running for deployment 'caip-in-cluster-controller-manager' in namespace 'caip-in-cluster-system', retrying
    [2023-01-10T20:07:01.058Z] waiting for resource capi-kubeadm-bootstrap-controller-manager of type *v1.Deployment to be up and running
    [2023-01-10T20:07:01.058Z] waiting for resource capi-kubeadm-control-plane-controller-manager of type *v1.Deployment to be up and running
    [2023-01-10T20:07:01.058Z] pods are not yet running for deployment 'capi-kubeadm-control-plane-controller-manager' in namespace 'capi-kubeadm-control-plane-system', retrying
    [2023-01-10T20:07:01.058Z] pods are not yet running for deployment 'capi-kubeadm-bootstrap-controller-manager' in namespace 'capi-kubeadm-bootstrap-system', retrying
    [2023-01-10T20:07:01.058Z] waiting for resource capv-controller-manager of type *v1.Deployment to be up and running
    [2023-01-10T20:07:01.058Z] pods are not yet running for deployment 'capv-controller-manager' in namespace 'capv-system', retrying
    [2023-01-10T20:07:01.058Z] waiting for resource capi-controller-manager of type *v1.Deployment to be up and running
```

 # Wait for CAPI Providers on the management cluster 
 Remember this section in the bootstrap cluster?  Same thing here.  Just that it takes a little longer... 
```
    Passed waiting on provider cluster-api after 424.811557ms

    [2023-01-10T20:07:06.446Z] pods are not yet running for deployment 'caip-in-cluster-controller-manager' in namespace 'caip-in-cluster-system', retrying
    [2023-01-10T20:07:06.446Z] pods are not yet running for deployment 'capi-kubeadm-control-plane-controller-manager' in namespace 'capi-kubeadm-control-plane-system', retrying
    
    Passed waiting on provider bootstrap-kubeadm after 5.342992042s
    
    [2023-01-10T20:07:06.446Z] pods are not yet running for deployment 'capv-controller-manager' in namespace 'capv-system', retrying
    [2023-01-10T20:07:10.765Z] pods are not yet running for deployment 'caip-in-cluster-controller-manager' in namespace 'caip-in-cluster-system', retrying
    
    Passed waiting on provider control-plane-kubeadm after 10.343552056s
    Passed waiting on provider infrastructure-vsphere after 10.386474481s
    Passed waiting on provider infrastructure-ipam-in-cluster after 15.119210611s
```

And eventually, we get the same success message" 

```
    [2023-01-10T20:07:16.515Z] Success waiting on all providers.
    [2023-01-10T20:07:20.828Z] ℹ  Updated package repository 'tanzu-management' in namespace 'tkg-system'
    [2023-01-10T20:09:27.780Z] ℹ  
    [2023-01-10T20:09:27.780Z]  Added installed package 'tkg-pkg'waiting for package: tkg-pkg
```

# MAGIC PART: tkg-pkg and "Waiting for package" hot loop 

Just like the kind cluster.  We let `tkg-pkg` pollute our cluster with all of our required management packages.   Interesting to note here that
packages we install are specific to the infrastructure provider.  So, we must either have `ako-operator` on all clouds (probably not), or, we have a different version of `tkg-pkg` meta-package repo that is used depending on if your azure or AWS.

```
    [2023-01-10T20:09:27.780Z] waiting for package: ako-operator
    [2023-01-10T20:09:27.780Z] waiting for package: tanzu-addons-manager
    [2023-01-10T20:09:27.780Z] waiting for package: tanzu-auth
    [2023-01-10T20:09:27.780Z] waiting for package: tanzu-cliplugins
    [2023-01-10T20:09:27.780Z] waiting for package: tanzu-core-management-plugins
    [2023-01-10T20:09:27.780Z] waiting for package: tanzu-featuregates
    [2023-01-10T20:09:27.780Z] waiting for package: tanzu-framework
    [2023-01-10T20:09:27.780Z] waiting for package: tkg-clusterclass
    [2023-01-10T20:09:27.780Z] waiting for package: tkg-clusterclass-vsphere
    [2023-01-10T20:09:27.780Z] waiting for package: tkr-service
    [2023-01-10T20:09:27.780Z] waiting for package: tkr-source-controller
    [2023-01-10T20:09:27.780Z] waiting for package: tkr-vsphere-resolver
    [2023-01-10T20:09:27.780Z] waiting for resource tkr-vsphere-resolver of type *v1alpha1.PackageInstall to be up and running
    [2023-01-10T20:09:27.780Z] waiting for resource tkg-pkg of type *v1alpha1.PackageInstall to be up and running
    [2023-01-10T20:09:27.780Z] waiting for resource ako-operator of type *v1alpha1.PackageInstall to be up and running
    [2023-01-10T20:09:27.780Z] waiting for resource tanzu-addons-manager of type *v1alpha1.PackageInstall to be up and running
    [2023-01-10T20:09:27.780Z] waiting for resource tanzu-auth of type *v1alpha1.PackageInstall to be up and running
    [2023-01-10T20:09:27.780Z] waiting for resource tanzu-cliplugins of type *v1alpha1.PackageInstall to be up and running
    [2023-01-10T20:09:27.780Z] waiting for resource tanzu-core-management-plugins of type *v1alpha1.PackageInstall to be up and running
    [2023-01-10T20:09:27.780Z] waiting for resource tkg-clusterclass of type *v1alpha1.PackageInstall to be up and running
    [2023-01-10T20:09:27.780Z] waiting for resource tkg-clusterclass-vsphere of type *v1alpha1.PackageInstall to be up and running
    [2023-01-10T20:09:27.780Z] waiting for resource tanzu-featuregates of type *v1alpha1.PackageInstall to be up and running
    [2023-01-10T20:09:27.780Z] waiting for resource tkr-service of type *v1alpha1.PackageInstall to be up and running
    [2023-01-10T20:09:27.780Z] waiting for resource tanzu-framework of type *v1alpha1.PackageInstall to be up and running
    [2023-01-10T20:09:27.780Z] waiting for resource tkr-source-controller of type *v1alpha1.PackageInstall to be up and running
```

And again, we now have the main management cluster packages, this time running in the VSphere cluster.

```
    [2023-01-10T20:09:27.780Z] successfully reconciled package: tkr-service
    [2023-01-10T20:09:27.780Z] successfully reconciled package: tkg-clusterclass-vsphere
    [2023-01-10T20:09:27.780Z] successfully reconciled package: tkr-vsphere-resolver
    [2023-01-10T20:09:27.780Z] successfully reconciled package: tanzu-addons-manager
    [2023-01-10T20:09:27.780Z] successfully reconciled package: tkg-pkg
    [2023-01-10T20:09:27.780Z] successfully reconciled package: tanzu-cliplugins
    [2023-01-10T20:09:27.781Z] successfully reconciled package: tanzu-framework
    [2023-01-10T20:09:27.781Z] successfully reconciled package: ako-operator
    [2023-01-10T20:09:27.781Z] successfully reconciled package: tkr-source-controller
    [2023-01-10T20:09:27.781Z] successfully reconciled package: tkg-clusterclass
    [2023-01-10T20:09:27.781Z] successfully reconciled package: tanzu-auth
    [2023-01-10T20:09:27.781Z] successfully reconciled package: tanzu-core-management-plugins
    [2023-01-10T20:09:27.781Z] successfully reconciled package: tanzu-featuregates
```

Now, we get ready to migrate the management cluster to VSphere.

```
    [2023-01-10T20:09:27.781Z] Waiting for the management cluster to get ready for move...
    [2023-01-10T20:09:27.781Z] waiting for resource tkg-mgmt-vc of type *v1beta1.Cluster to be up and running
    [2023-01-10T20:09:27.781Z] waiting for resources type *v1beta1.KubeadmControlPlaneList to be up and running
    [2023-01-10T20:09:27.781Z] waiting for resources type *v1beta1.MachineDeploymentList to be up and running
    [2023-01-10T20:09:27.781Z] waiting for resources type *v1beta1.MachineList to be up and running
    [2023-01-10T20:09:27.781Z] Waiting for addons installation...
    [2023-01-10T20:09:27.781Z] waiting for resources type *v1beta1.ClusterResourceSetList to be up and running
    [2023-01-10T20:09:27.781Z] waiting for resource antrea-controller of type *v1.Deployment to be up and running
    [2023-01-10T20:09:27.781Z] Applying ClusterBootstrap and its associated resources on management cluster
    [2023-01-10T20:09:27.781Z] User ConfigValues File: /tmp/1826823029.yaml
    [2023-01-10T20:09:27.781Z] Checking if TKr v1.24.9---vmware.1-tkg.1-rc.2 is created on management cluster
    [2023-01-10T20:09:27.781Z] waiting for resource v1.24.9---vmware.1-tkg.1-rc.2 of type *v1alpha3.TanzuKubernetesRelease to be up and running
    [2023-01-10T20:09:27.781Z] Applying ClusterBootstrap: apiVersion: v1
    [2023-01-10T20:09:27.781Z] kind: Secret
    [2023-01-10T20:09:27.781Z] metadata:
    [2023-01-10T20:09:27.781Z]   name: tkg-mgmt-vc-pinniped-package
    [2023-01-10T20:09:27.781Z]   namespace: tkg-system
    [2023-01-10T20:09:27.781Z]   labels:
    [2023-01-10T20:09:27.781Z]     tkg.tanzu.vmware.com/addon-name: pinniped
    [2023-01-10T20:09:27.781Z]     tkg.tanzu.vmware.com/cluster-name: tkg-mgmt-vc
    [2023-01-10T20:09:27.781Z]     clusterctl.cluster.x-k8s.io/move: ""
    [2023-01-10T20:09:27.781Z]   annotations:
    [2023-01-10T20:09:27.781Z]     tkg.tanzu.vmware.com/addon-type: authentication/pinniped
    [2023-01-10T20:09:27.781Z] type: clusterbootstrap-secret
    [2023-01-10T20:09:27.781Z] stringData:
    [2023-01-10T20:09:27.781Z]   values.yaml: |
    [2023-01-10T20:09:27.781Z]     infrastructure_provider: vsphere
    [2023-01-10T20:09:27.781Z]     tkg_cluster_role: workload
    [2023-01-10T20:09:27.781Z]     identity_management_type: none
    [2023-01-10T20:09:27.781Z] ---
    [2023-01-10T20:09:27.781Z] apiVersion: run.tanzu.vmware.com/v1alpha3
    [2023-01-10T20:09:27.781Z] kind: ClusterBootstrap
    [2023-01-10T20:09:27.781Z] metadata:
    [2023-01-10T20:09:27.781Z]   name: tkg-mgmt-vc
    [2023-01-10T20:09:27.781Z]   namespace: tkg-system
    [2023-01-10T20:09:27.781Z]   annotations:
    [2023-01-10T20:09:27.781Z]     tkg.tanzu.vmware.com/add-missing-fields-from-tkr: v1.24.9---vmware.1-tkg.1-rc.2
    [2023-01-10T20:09:27.781Z] spec:
    [2023-01-10T20:09:27.781Z]   kapp:
    [2023-01-10T20:09:27.781Z]     refName: kapp-controller*
    [2023-01-10T20:09:27.781Z]   additionalPackages:
    [2023-01-10T20:09:27.781Z]   - refName: metrics-server*
    [2023-01-10T20:09:27.781Z]   - refName: secretgen-controller*
    [2023-01-10T20:09:27.781Z]   - refName: pinniped*
    [2023-01-10T20:09:27.781Z]   - refName: tkg-storageclass*
    [2023-01-10T20:09:27.781Z]     valuesFrom:
    [2023-01-10T20:09:27.781Z]       inline:
    [2023-01-10T20:09:27.781Z]         metadata:
    [2023-01-10T20:09:27.781Z]           infraProvider: vsphere
```

## Performing the move

At this point we have:
- kind running
  - Cluster API object for a VSPhere cluster
  - Vsphere VMs that are a K8s cluster
    - They are running Cluster API controllers
    - They are running Cluster API VSphere controllers
    - They are running AVI/AKO/Antrea/etc

HOWEVER:

- The VSphere K8s cluster is not managing itself...
- The Kind cluster is still capable of manageming the management cluster !
- That means kind needs to go away... otherwise it will be a security issue, and it will be confusing for people.

```
    [2023-01-10T20:09:27.781Z] Moving all Cluster API objects from bootstrap cluster to management cluster...
    [2023-01-10T20:09:27.781Z] Performing move...
```
### Move: part 0,  Discovery 

Moving CAPI resources requires listing them all, first...  Theres lots of stuff (Even cloud provider credentials, for example), that need
to get migrated over.  For details on this, check https://cluster-api.sigs.k8s.io/clusterctl/provider-contract.html#move.  This describes yhe
discovery process...  Let's look at the logs from TKG for the discovery.  Remember here, 

- we're reading Cluster API objects that are living in our Kind cluster  
- with the intent of migrating them to our PERMANANT managemnt cluster.

```
    # introspect the CAPI objects on the kind cluster so that the mgmt cluster can become self-aware and we dont lose any info after kind dies.
    [2023-01-10T20:09:27.781Z] Discovering Cluster API objects
    [2023-01-10T20:09:27.781Z] Certificate Count=4
    [2023-01-10T20:09:27.781Z] KubeadmControlPlane Count=1
    [2023-01-10T20:09:27.781Z] KubeadmControlPlaneTemplate Count=1
    [2023-01-10T20:09:27.781Z] VSphereClusterTemplate Count=1
    [2023-01-10T20:09:27.781Z] CertificateRequest Count=4
    [2023-01-10T20:09:27.781Z] ClusterClass Count=1
    [2023-01-10T20:09:27.781Z] KubeadmConfigTemplate Count=2
    [2023-01-10T20:09:27.781Z] MachineSet Count=1
    [2023-01-10T20:09:27.781Z] VSphereMachineTemplate Count=4
    [2023-01-10T20:09:27.781Z] Issuer Count=3
    [2023-01-10T20:09:27.781Z] Machine Count=2
    [2023-01-10T20:09:27.781Z] Secret Count=51
    [2023-01-10T20:09:27.781Z] ConfigMap Count=42
    [2023-01-10T20:09:27.781Z] KubeadmConfig Count=2
    [2023-01-10T20:09:27.781Z] MachineDeployment Count=1
    [2023-01-10T20:09:27.781Z] VSphereVM Count=2
    [2023-01-10T20:09:27.781Z] Cluster Count=1
    [2023-01-10T20:09:27.781Z] VSphereCluster Count=1
    [2023-01-10T20:09:27.781Z] MachineHealthCheck Count=2
    [2023-01-10T20:09:27.781Z] VSphereMachine Count=2
    [2023-01-10T20:09:27.781Z] Total objects Count=145

    [2023-01-10T20:09:27.781Z] Excluding secret from move (not linked with any Cluster) name="ako-operator-v2-values"
    [2023-01-10T20:09:27.782Z] Excluding secret from move (not linked with any Cluster) name="tanzu-framework-values"
    ...
    [2023-01-10T20:09:27.782Z] Excluding secret from move (not linked with any Cluster) name="tkr-source-controller-values"
    [2023-01-10T20:09:27.782Z] Excluding secret from move (not linked with any Cluster) name="tkr-vsphere-resolver-values"

    ...
    [2023-01-10T20:09:27.782Z] Object won't be moved because it's not included in GVK considered for move kind="PackageRepository" 
    [2023-01-10T20:09:27.782Z] Object won't be moved because it's not included in GVK considered for move kind="PackageInstall" name="tanzu-addons-manager"

    [2023-01-10T20:09:27.782Z] Object won't be moved because it's not included in GVK considered for move kind="ClusterBootstrap" name="tkg-mgmt-vc"
```

### Move: part 1, Now we know WHAT to move to the mgmt cluster

Now finally, we start the "move" .  FIRST WE HAVE TO PAUSE the existing `Cluster` object!.

```
    [2023-01-10T20:09:27.782Z] Moving Cluster API objects Clusters=1
    [2023-01-10T20:09:27.782Z] Moving Cluster API objects ClusterClasses=1

    #### The PAUSE generally is IMPORTANT !!!!!!!!  (in TKG, its not a huge deal bc nobody is using bootstrap cluster right now...) But
    #### In the real world, (i.e. in a backup restore situation) when using cluster API, this is a non-trivial operation.
    #### Hence we ALWAYS pause the source cluster before migraion as a matter of how Clusterctl works.

    [2023-01-10T20:09:27.782Z] Pausing the source cluster

    [2023-01-10T20:09:27.782Z] Set Cluster.Spec.Paused Paused=true Cluster="tkg-mgmt-vc" Namespace="tkg-system"
    [2023-01-10T20:09:27.782Z] Pausing the source cluster classes
    [2023-01-10T20:09:27.782Z] Set Paused annotation ClusterClass="tkg-vsphere-default-v1.0.0" Namespace="tkg-system"
    [2023-01-10T20:09:27.782Z] Creating target namespaces, if missing
    [2023-01-10T20:09:27.782Z] Creating objects in the target cluster
    [2023-01-10T20:09:27.782Z] Creating ClusterClass="tkg-vsphere-default-v1.0.0" Namespace="tkg-system"
    ...
    [2023-01-10T20:09:32.471Z] Deleting VSphereVM="tkg-mgmt-vc-md-0-69pnq-858fcd866b-4h59q" Namespace="tkg-system"
    [2023-01-10T20:09:32.471Z] Deleting Secret="tkg-mgmt-vc-md-0-bootstrap-mzn9v-cgbx2" Namespace="tkg-system"
    [2023-01-10T20:09:32.747Z] Deleting KubeadmConfig="tkg-mgmt-vc-md-0-bootstrap-mzn9v-cgbx2" Namespace="tkg-system"
    [2023-01-10T20:09:32.747Z] Deleting VSphereVM="tkg-mgmt-vc-4b4wz-6vhgl" Namespace="tkg-system"
    
    [2023-01-10T20:09:34.409Z] Deleting Secret="tkg-mgmt-vc-antrea-data-values" Namespace="tkg-system"
    [2023-01-10T20:09:34.409Z] Deleting VSphereMachineTemplate="tkg-mgmt-vc-control-plane-6g4h4" Namespace="tkg-system"
    ...
    [2023-01-10T20:09:36.122Z] Deleting VSphereMachineTemplate="tkg-vsphere-default-v1.0.0-control-plane" Namespace="tkg-system"
    [2023-01-10T20:09:36.396Z] Deleting VSphereClusterTemplate="tkg-vsphere-default-v1.0.0-cluster" Namespace="tkg-system"

    [2023-01-10T20:09:37.271Z] Remove Paused annotation ClusterClass="tkg-vsphere-default-v1.0.0" Namespace="tkg-system"
```

... WAIT FOR IT ...

```
    #### This is the last step of the clusterctl move !!! 
    [2023-01-10T20:09:37.271Z] Set Cluster.Spec.Paused Paused=false Cluster="tkg-mgmt-vc" Namespace="tkg-system"
    [2023-01-10T20:09:37.271Z] Resuming the target cluster

```

How is this different then Velero migrations ? Clusterctl looks at `ownerRef` fields.  
- CAPI: Preserves exact identities of each object (ownerRef, finalizer, managed fields)
  - This forces CAPI to order `Creation` and `Deletion`
- Velero: DROPS all ownerRefs, finalizers, managed fields.... 

Now WE start back up the management cluster.   We can see the logic for how we patch things in tkg/client/init.go in the `PatchClusterInitOperations` function...

Concretely, one reason to patch mgmt cluster is so that tanzu cli can statelessly determine what verion of TKG is associated w/ the MGMT cluster.
That is super important, for example, when a user wants to upgrade a WL cluster to a New TKR, bc only CERTAIN Mgmt clusters support CERTAIN TKRs.
i.e. you cant run create a WL cluster w k8s 1.25 if you are on a 1.2 TKG management cluster, bc that mgmt cluster only runs k8s 1.20(or something).

### Move: part 2,Patch the management cluster after the fact

```
    [2023-01-10T20:09:37.872Z] Applying patch to resource tkg-mgmt-vc of type *v1beta1.Cluster ...
    [2023-01-10T20:09:37.872Z] Applying patch to resource tkg-mgmt-vc of type *v1beta1.Cluster ...
    [2023-01-10T20:09:38.472Z] Applying patch to resource tkg-vsphere-default-v1.0.0-kcp of type *unstructured.Unstructured ...
    [2023-01-10T20:09:38.743Z] Applying patch to resource tkg-vsphere-default-v1.0.0-cluster of type *unstructured.Unstructured ...
    [2023-01-10T20:09:39.016Z] Applying patch to resource tkg-mgmt-vc-control-plane-6g4h4 of type *unstructured.Unstructured ...
    [2023-01-10T20:09:39.016Z] Applying patch to resource tkg-mgmt-vc-md-0-infra-hjcsr of type *unstructured.Unstructured ...
    [2023-01-10T20:09:39.016Z] Applying patch to resource tkg-vsphere-default-v1.0.0-control-plane of type *unstructured.Unstructured ...
    [2023-01-10T20:09:39.289Z] Applying patch to resource tkg-vsphere-default-v1.0.0-worker of type *unstructured.Unstructured ...
    [2023-01-10T20:09:39.289Z] Applying patch to resource tkg-vsphere-default-v1.0.0 of type *unstructured.Unstructured ...
    [2023-01-10T20:09:39.918Z] Applying patch to resource tkg-mgmt-vc-md-0-bootstrap-mzn9v of type *unstructured.Unstructured ...
    [2023-01-10T20:09:39.918Z] Applying patch to resource tkg-vsphere-default-v1.0.0-md-config of type *unstructured.Unstructured ...
    [2023-01-10T20:09:40.899Z] IsProd: 
    [2023-01-10T20:09:40.899Z] IsOfficialBuild: False
    [2023-01-10T20:09:40.899Z] ---
    [2023-01-10T20:09:40.899Z] apiVersion: v1
    [2023-01-10T20:09:40.899Z] kind: Namespace
    [2023-01-10T20:09:40.899Z] metadata:
    [2023-01-10T20:09:40.899Z]   name: tkg-system-telemetry
    [2023-01-10T20:09:40.899Z] 
    [2023-01-10T20:09:40.899Z] ---
    [2023-01-10T20:09:40.899Z] apiVersion: v1
    [2023-01-10T20:09:40.899Z] kind: ServiceAccount
    [2023-01-10T20:09:40.899Z] metadata:
    [2023-01-10T20:09:40.899Z]   name: tkg-telemetry-sa
    [2023-01-10T20:09:40.899Z]   namespace: tkg-system-telemetry
    [2023-01-10T20:09:40.899Z] 
    [2023-01-10T20:09:40.899Z] ---
    [2023-01-10T20:09:40.899Z] kind: ClusterRole
    [2023-01-10T20:09:40.899Z] apiVersion: rbac.authorization.k8s.io/v1
    [2023-01-10T20:09:40.899Z] metadata:
    [2023-01-10T20:09:40.899Z]   name: tkg-telemetry-cluster-role
    [2023-01-10T20:09:40.899Z] rules:
    [2023-01-10T20:09:40.899Z]   - apiGroups: [""]
    [2023-01-10T20:09:40.900Z]     resources: ["secrets", "namespaces", "configmaps"]
    ...
    [2023-01-10T20:09:40.901Z]           restartPolicy: Never
    [2023-01-10T20:09:41.497Z] Creating tkg-bom versioned ConfigMaps...
```

And finally we have a management cluster........
```
    [2023-01-10T20:09:41.497Z] You can now access the management cluster tkg-mgmt-vc by running 'kubectl config use-context tkg-mgmt-vc-admin@tkg-mgmt-vc'
    [2023-01-10T20:09:41.497Z] Deleting kind cluster: tkg-kind-ceus13jb4t5tab9k5jk0
    [2023-01-10T20:09:45.829Z] 
    [2023-01-10T20:09:45.829Z] Management cluster created!
    [2023-01-10T20:09:45.829Z] 
    [2023-01-10T20:09:45.829Z] 
    [2023-01-10T20:09:45.829Z] You can now create your first workload cluster by running the following:
    [2023-01-10T20:09:45.829Z] 
    [2023-01-10T20:09:45.829Z]   tanzu cluster create [name] -f [file]
    [2023-01-10T20:09:45.829Z] 
    [2023-01-10T20:09:45.829Z] 
    [2023-01-10T20:09:45.829Z] Some addons might be getting installed! Check their status by running the following:
    [2023-01-10T20:09:45.829Z] 
    [2023-01-10T20:09:45.829Z]   kubectl get apps -A
    [2023-01-10T20:09:45.829Z] 
```

Now the managementcluster/create.go checks one last time, after MGMT cluster creation, to confirm that
the plugins  it cares about, are installed.  See cmd/cli/plugin/managementcluster/create.go for details. 

```
// Sync plugins if management-cluster creation is successful and --dry-run was not set
	if config.IsFeatureActivated(cliconfig.FeatureContextAwareCLIForPlugins) && !iro.dryRun {
		err = pluginmanager.SyncPlugins()
```
And of course, this works ....

```
    [2023-01-10T20:09:45.829Z] ℹ  Checking for required plugins... 
    [2023-01-10T20:09:46.826Z] ℹ  Installing plugin 'kubernetes-release:v0.28.0-dev' with target 'kubernetes'
    [2023-01-10T20:09:50.250Z] ℹ  Installing plugin 'cluster:v0.28.0-dev' with target 'kubernetes'
    [2023-01-10T20:09:52.879Z] ℹ  Installing plugin 'feature:v0.28.0-dev' with target 'kubernetes'
    [2023-01-10T20:09:53.865Z] ℹ  Successfully installed all required plugins 
```

Now we list available clusters, and we see the mgmt cluster.  Note the need to use -A... 

```
    [2023-01-10T20:09:54.698Z] INFO:root:====== 198  CMD: tanzu cluster list --include-management-cluster -A --output=json
    [2023-01-10T20:09:55.303Z] [
    [2023-01-10T20:09:55.303Z]   {
    [2023-01-10T20:09:55.303Z]     "name": "tkg-mgmt-vc",
    [2023-01-10T20:09:55.303Z]     "namespace": "tkg-system",
    [2023-01-10T20:09:55.303Z]     "status": "running",
    [2023-01-10T20:09:55.303Z]     "plan": "dev",
    [2023-01-10T20:09:55.303Z]     "controlplane": "1/1",
    [2023-01-10T20:09:55.303Z]     "workers": "1/1",
    [2023-01-10T20:09:55.303Z]     "kubernetes": "v1.24.9+vmware.1",
    [2023-01-10T20:09:55.303Z]     "roles": [
    [2023-01-10T20:09:55.303Z]       "management"
    [2023-01-10T20:09:55.303Z]     ],
    [2023-01-10T20:09:55.303Z]     "tkr": "v1.24.9---vmware.1-tkg.1-rc.2",
    [2023-01-10T20:09:55.303Z]     "labels": {
    [2023-01-10T20:09:55.303Z]       "cluster-role.tkg.tanzu.vmware.com/management": "",
    [2023-01-10T20:09:55.303Z]       "cluster.x-k8s.io/cluster-name": "tkg-mgmt-vc",
    [2023-01-10T20:09:55.303Z]       "networking.tkg.tanzu.vmware.com/avi": "install-ako-for-management-cluster",
    [2023-01-10T20:09:55.303Z]       "run.tanzu.vmware.com/tkr": "v1.24.9---vmware.1-tkg.1-rc.2",
    [2023-01-10T20:09:55.303Z]       "tkg.tanzu.vmware.com/cluster-name": "tkg-mgmt-vc",
    [2023-01-10T20:09:55.303Z]       "topology.cluster.x-k8s.io/owned": ""
    [2023-01-10T20:09:55.303Z]     }
    [2023-01-10T20:09:55.303Z]   }
```

Ok thats it.  We now have a fully functional TKG MAnagement cluster !!!