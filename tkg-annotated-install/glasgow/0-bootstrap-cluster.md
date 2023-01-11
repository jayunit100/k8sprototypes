# TKG Client and Kind Bootstrapping

```
  [2023-01-10T19:36:11.674Z] + curl https://build-artifactory.eng.vmware.com/kscom-generic-local/TKG/channels/442519250544895703/_boltArtifacts/tkg-v2.1.0-rc.2.buildinfo.yaml
```

```
  sudo -S -p '[sudo] password: ' ln -sf $(realpath /home/kubo/tanzu_tools/cli/core/v0.28.0-dev/ INFO:root:====== 99   CMD: tanzu config get | yq eval '.clientOptions.features.global.context-aware-cli-for-plugins' -
  true
  INFO:root:context-aware-cli-for-plugins is ON

  tanzu plugin clean
  ✔  successfully cleaned up all plugins 


  tanzu plugin repo update -b tanzu-cli-framework core
  ✔  successfully updated repository configuration for core
```

##  Making sure the client plugins are there... 
```
  ℹ  Installing plugin 'pinniped-auth:v0.28.0-dev'
  ℹ  Installing plugin 'secret:v0.28.0-dev' with target 'kubernetes'
  ℹ  Installing plugin 'telemetry:v0.28.0-dev' with target 'kubernetes'
  ℹ  Installing plugin 'isolated-cluster:v0.28.0-dev'
  ℹ  Installing plugin 'login:v0.28.0-dev'
  ℹ  Installing plugin 'management-cluster:v0.28.0-dev' with target 'kubernetes'
  ℹ  Installing plugin 'package:v0.28.0-dev' with target 'kubernetes'
  ℹ  Successfully installed all required plugins 
  ✔  Done 
  INFO:root:====== 103  CMD: tanzu plugin list
  Standalone Plugins
    NAME                DESCRIPTION                                                        TARGET      DISCOVERY  VERSION      STATUS  
    isolated-cluster    isolated-cluster operations                                                    default    v0.28.0-dev  installe
    login               Login to the platform                                                          default    v0.28.0-dev  installe
    pinniped-auth       Pinniped authentication operations (usually not directly invoked)              default    v0.28.0-dev  installe
    management-cluster  Kubernetes management-cluster operations                           kubernetes  default    v0.28.0-dev  installe
    package             Tanzu package management                                           kubernetes  default    v0.28.0-dev  installe
    secret              Tanzu secret management                                            kubernetes  default    v0.28.0-dev  installe
    telemetry           Configure cluster-wide telemetry settings                          kubernetes  default    v0.28.0-dev  installe
  INFO:root:====== 104  CMD: tanzu version
  version: v0.28.0-dev
  buildDate: 2023-01-09
  sha: d293a0881-dirty
```

# Uploading the OVAs to Vsphere


- INFO:root:Downloading OVA from http://build-squid.eng.vmware.com/build/mts/release/bora-21093855/publish/lin64/tkg_release/node/ova-ubuntu-2004-v1.24.9+vmware.1-tkg.1-b030088fe71fea7ff1ecb87a4d425c93/ubuntu-2004-kube-v1.24.9+vmware.1-tkg.1-b030088fe71fea7ff1ecb87a4d425c93.ova to 
temp_ova_dir-NVZOBCZG

-  INFO:root:Importing OVA from http://build-squid.eng.vmware.com/build/mts/release/bora-21093855/publish/lin64/tkg_release/node/ova-ubuntu-2004-v1.24.9+vmware.1-tkg.1-b030088fe71fea7ff1ecb87a4d425c93/ubuntu-2004-kube-v1.24.9+vmware.1-tkg.1-b030088fe71fea7ff1ecb87a4d425c93.ova to /dc0/vm/.  This might take a while...

# TKG Kind Cluster: Create the bootstrap (i.e. kind + "Management Cluster") 

The `InitRegion` function lives inside Tanzu CLI.  IT has ALL the logic for kind and mgmt cluster bootstrapping.  

- You can't have a workload cluster w/o a persistent management cluster.
- You can't have a persistent management cluster without a kind bootstrap cluster.
- You can't have a kind bootstrap cluster that works w/ TKG without installing tkg-pkg and other pre-requisites

Keep in mind that MANY OF THE THINGS the kind cluster does has to happen AGAIN when we make the persistent mangaement cluster.  So, the purpose of the `kind` cluster is to

- *define* the management cluster as a set of CRDs that run a particular K8s version, CNI, and so on.
- *run* a few `capv-controller` objects that can make VMs in the vsphere cloud, where those contorllers read the CRDs, and act on them (i.e. by making VMs on vsphere that run k8s )
- *move* those `capv-controller` and other `capi` objects into the vsphere cluster, once it is up
- *self-destruct* once the persistant management cluster is up and running.

This page ONLY defines the creation of the `kind` cluster.  IT doesn't ACTUALLY create any CAPI clusters, that is on the *next* markdown page, where we define the `Management Cluster` creation.

Run:
```
tanzu management-cluster create  --yes -v 9 --deploy-tkg-on-vSphere7 -f tkg-mgmt-vc.yaml
```
Result: 
```
 compatibility file (/home/kubo/.config/tanzu/tkg/compatibility/tkg-compatibility.yaml) already exists, skipping download
 BOM files inside /home/kubo/.config/tanzu/tkg/bom already exists, skipping download
 CEIP Opt-in status: true
 cluster log directory does not exist. Creating new one at "/home/kubo/.config/tanzu/tkg/logs"
 
 Validating the pre-requisites...
 
 vSphere 7.0 Environment Detected.
 
 You have connected to a vSphere 7.0 environment which does not have vSphere with Tanzu enabled. vSphere with Tanzu includes
 an integrated Tanzu Kubernetes Grid Service which turns a vSphere cluster into a platform for running Kubernetes workloads in dedicated
 resource pools. Configuring Tanzu Kubernetes Grid Service is done through vSphere HTML5 client.
 
 Tanzu Kubernetes Grid Service is the preferred way to consume Tanzu Kubernetes Grid in vSphere 7.0 environments. Alternatively you may
 deploy a non-integrated Tanzu Kubernetes Grid instance on vSphere 7.0.
 Deploying TKG management cluster on vSphere 7.0 ...
 Identity Provider not configured. Some authentication features won't work.
 Using default value for CONTROL_PLANE_MACHINE_COUNT = 1. Reason: CONTROL_PLANE_MACHINE_COUNT variable is not set
 Using default value for WORKER_MACHINE_COUNT = 1. Reason: WORKER_MACHINE_COUNT variable is not set
 Setting config variable "VSPHERE_DATACENTER" to value "/dc0"
 Setting config variable "VSPHERE_NETWORK" to value "/dc0/network/VM Network"
 Setting config variable "VSPHERE_RESOURCE_POOL" to value "/dc0/host/cluster0/Resources/rp0"
 Setting config variable "VSPHERE_DATASTORE" to value "/dc0/datastore/sharedVmfs-0"
 Setting config variable "VSPHERE_FOLDER" to value "/dc0/vm/folder0"
 Checking if VSPHERE_CONTROL_PLANE_ENDPOINT  is already in use
```

At this point, 

```
func (c *TkgClient) InitRegion(options *InitRegionOptions) error { //nolint:funlen,gocyclo
```

This function lives in tkg/client/init.go.    Once it takes over, and it does quite a bit... 

- Setting up kind
- Setting up a management cluster on kind
- Calling the `InstallOrUpgradeManagementComponents` method, which will wait for all the management components and packages to come up.

```
 Setting up management cluster...
 Validating configuration...
 Setting CLUSTER_TOPOLOGY to "true"
 Using infrastructure provider vsphere:v1.5.1
 Generating cluster configuration...
 Setting up bootstrapper...
 Fetching configuration for kind node image...
 kindConfig: 
  &{{Cluster kind.x-k8s.io/v1alpha4}  [{  map[] [{/var/run/docker.sock /var/run/docker.sock false false }] [] [] []}] { 0  100.96.0.0/11 100.64.0.0/13 false } map[] map[] [apiVersion: kubeadm.k8s.io/v1beta3
```

InitRegion: This is the kind configuration that we'll use when bootstrapping TKG. .... 
```
 kind: ClusterConfiguration
 imageRepository: projects.registry.vmware.com/tkg
 etcd:
   local:
     imageRepository: projects.registry.vmware.com/tkg
     imageTag: v3.5.6_vmware.3
 dns:
   type: CoreDNS
   imageRepository: projects.registry.vmware.com/tkg
   imageTag: v1.8.6_vmware.15] [] [[plugins]
 [plugins.'io.containerd.grpc.v1.cri']
 [plugins.'io.containerd.grpc.v1.cri'.registry]
 [plugins.'io.containerd.grpc.v1.cri'.registry.configs]
 [plugins.'io.containerd.grpc.v1.cri'.registry.configs.'projects-stg.registry.vmware.com']
 [plugins.'io.containerd.grpc.v1.cri'.registry.configs.'projects-stg.registry.vmware.com'.tls] 
 insecure_skip_verify = false
 ca_file = ''
 ] []}
 ```

InitRegion: Now, the logs continue rolling out for kind cluster bootstrapping  
```

 Creating kind cluster: tkg-kind-ceus13jb4t5tab9k5jk0
 Creating cluster "tkg-kind-ceus13jb4t5tab9k5jk0" ...
 Ensuring node image (projects-stg.registry.vmware.com/tkg/kind/node:v1.24.9_vmware.1-tkg.1_v0.17.0) ...
 Pulling image: projects-stg.registry.vmware.com/tkg/kind/node:v1.24.9_vmware.1-tkg.1_v0.17.0 ...
 Preparing nodes ...
 Writing configuration ...
 U 
[2023-01-10T19:46:42.469Z] Starting control-plane ...

 GET https://tkg-kind-ceus13jb4t5tab9k5jk0-control-plane:6443/healthz?timeout=10s in 0 millisecondsI0110 19:46:46.954014 280 round_trippers.go:553] GET https://tkg-kind-ceus13jb4t5tab9k5jk0-control-plane:6443/healthz?timeout=10s in 0 millisecondsI0110 19:46:47.454402 280 round_trippers.go:553] GET https://
 ...
 tkg-kind-ceus13jb4t5tab9k5jk0-control-plane:6443/healthz?timeout=10s in 0 millisecondsI0110 19:46:54.201890 280 round_trippers.go:553] GET https://
 ...
 tkg-kind-ceus13jb4t5tab9k5jk0-control-plane:6443/healthz?timeout=10s 500 Internal Server Error in 5 millisecondsI0110 19:46:55.956712 280 round_trippers.go:553] GET https://
 ...
 tkg-kind-ceus13jb4t5tab9k5jk0-control-plane:6443/healthz?timeout=10s 200 OK in 2 milliseconds[apiclient] All control plane components are healthy after 10.006544 secondsI0110 
```

InitRegion: Next, the CNI Installation (kind-net, and so on happens).  Note were still just making the bootstrap cluster.  Nothing will go wrong here, ever... 

```
 Installing CNI ...
 Installing StorageClass ...
 Waiting 2m0s for control-plane = Ready ...
 Ready after 28s 
 Bootstrapper created. Kubeconfig: /home/kubo/.kube-tkg/tmp/config_Z1ics3TW
```

## Kind Cluster: TKG And ClusterClass customizations

Now that the kind cluster is up, tanzu cli will start customizing it. 
```
 Warning: unable to find component 'kube_rbac_proxy' under BoM
 Installing kapp-controller on bootstrap cluster...
 User ConfigValues File: /tmp/2991662634.yaml
 Kapp-controller values-file: /tmp/330923605.yaml
 Kapp-controller configuration file: /tmp/3117802898
 
 waiting for resource kapp-controller of type *v1.Deployment to be up and running
 pods are not yet running for deployment 'kapp-controller' in namespace 'tkg-system', retrying
 ...
 pods are not yet running for deployment 'kapp-controller' in namespace 'tkg-system', retrying
 Installing providers on bootstrapper...
 Installing the clusterctl inventory CRD
```

Now youll get ALOT of logs for cluster CTL inventory setup.  Close your eyes, we'll wake you up when its over.

## InitRegion: Cluster API installation onto the Kind cluster

At this point, on our kind cluster, we'll start adding basic primitives necessary to make our bootstrap cluster.
That means, teaching it CAPI.
- Teach the kind cluster about CAPI  <-- you are here
- Tell the kind cluster to make a workload cluster on vsphere (later)
- Tell the kind cluster to nominate the workload cluster to a management cluster (later)
- Tell the kind cluster to go away... forever (later)

### Tell kind about CAPI 

We now see in the logs, something like this.  At this point, tanzu cli is now `clusterctl`, internally, and clusterctl is doing a bunch
of magic to turn our kind cluster into a cluster API Management cluster... this means... just adding a bunch of containers, provider CRDs, and so on...

```
 Creating CustomResourceDefinition="providers.clusterctl.cluster.x-k8s.io"
```

Now we "Fetch" providers.   This happens in clusterctl https://github.com/kubernetes-sigs/cluster-api/blob/main/cmd/clusterctl/client/init.go
- ipam provider
- bootstrap provider
- control plane provider
- infrastructure provider

```
    Fetching providers
    Fetching File="core-components.yaml" Provider="cluster-api" Type="CoreProvider" Version="v1.2.8"
    Fetching File="bootstrap-components.yaml" Provider="kubeadm" Type="BootstrapProvider" Version="v1.2.8"
    Fetching File="control-plane-components.yaml" Provider="kubeadm" Type="ControlPlaneProvider" Version="v1.2.8"
    Fetching File="infrastructure-components.yaml" Provider="vsphere" Type="InfrastructureProvider" Version="v1.5.1"
    Fetching File="ipam-components.yaml" Provider="ipam-in-cluster" Type="InfrastructureProvider" Version="v0.1.0"
    Fetching File="metadata.yaml" Provider="cluster-api" Type="CoreProvider" Version="v1.2.8"
    Fetching File="metadata.yaml" Provider="kubeadm" Type="BootstrapProvider" Version="v1.2.8"
    Fetching File="metadata.yaml" Provider="kubeadm" Type="ControlPlaneProvider" Version="v1.2.8"
    Fetching File="metadata.yaml" Provider="vsphere" Type="InfrastructureProvider" Version="v1.5.1"
    Fetching File="metadata.yaml" Provider="ipam-in-cluster" Type="InfrastructureProvider" Version="v0.1.0"
```

At this poing we see the cert-manager components being created, https://github.com/kubernetes-sigs/cluster-api/blob/main/cmd/clusterctl/client/cluster/cert_manager.go#L495. 
Again, remember, this function is being called by tanzu cli, but under the hood tanzu cli is calling `clusterctl`, which knows natively how to set up certmanager.

# Kind Cert Manager setup 

Cert manager is how CAPI webhooks authenticate to each other.  It's an implementation detail of how CAPI communicates within itself. 

```
      Creating Namespace="cert-manager-test"
      I0110 19:47:49.143398    5793 request.go:601] Waited for 1.047538531s due to client-side throttling, not priority and fairness, request: GET:https://127...
      Installing cert-manager Version="v1.9.1"
      Fetching File="cert-manager.yaml" Provider="cert-manager" Type="" Version="v1.9.1"
      Creating Namespace="cert-manager"
      Creating CustomResourceDefinition="certificaterequests.cert-manager.io"
      Creating CustomResourceDefinition="certificates.cert-manager.io"
      ...
      Creating Deployment="cert-manager-webhook" Namespace="cert-manager"
      Creating MutatingWebhookConfiguration="cert-manager-webhook"
      Creating ValidatingWebhookConfiguration="cert-manager-webhook"
      Waiting for cert-manager to be available...
      Updating Namespace="cert-manager-test"
      Creating Issuer="test-selfsigned" Namespace="cert-manager-test"
      ...
      Deleting Certificate="selfsigned-cert" Namespace="cert-manager-test"
```

# Installing CAPI objects onto the Kind cluster

Now, still seting up CAPI on kind , we are now installing CAPI objects: These "Creating" logs come from 
https://github.com/kubernetes-sigs/cluster-api/blob/main/cmd/clusterctl/client/cluster/components.go the internal `createObj` tools in cluster api, which know how to install arbitrary kubernetes objects.... 

```
      Installing Provider="cluster-api" Version="v1.2.8" TargetNamespace="capi-system"
      Creating objects Provider="cluster-api" Version="v1.2.8" TargetNamespace="capi-system"
      Creating Namespace="capi-system"
      Creating CustomResourceDefinition="clusterclasses.cluster.x-k8s.io"
        
        ... (you'll see hundreds of other "Creating" logs here... )
      
      pods are not yet running for deployment 'caip-in-cluster-controller-manager' in namespace 'caip-in-cluster-system', retrying
```

InitRegion.... still going ! We're still in the bootstrap cluster: Next, we'll wait for "packages" and "providers" to come online.  
Still in the init.go function of tkg/client... 

```
      func (c *TkgClient) WaitForProviders(clusterClient clusterclient.Client, options waitForProvidersOptions) error {
```
# Wait for CAPI PRoviders on the kind cluster
```

      Passed waiting on provider bootstrap-kubeadm after 15.51390451s
      Passed waiting on provider control-plane-kubeadm after 15.566739105s
      Passed waiting on provider infrastructure-vsphere after 15.602801637s
      Passed waiting on provider cluster-api after 15.694037955s
      Passed waiting on provider infrastructure-ipam-in-cluster after 20.223833161s
      Success waiting on all providers.
      [ ℹ  Updated package repository 'tanzu-management' in namespace 'tkg-system'
```

Still in the bootstrap cluster: Next, we'll wait for packages. 

# MAGIC PART: tkg-pkg and "Waiting for package" hot loop 


Below we'll "wait for package", i.e. we'll wait for individual packages to come online:
```
  [2023-01-10T19:50:48.603Z]  Added installed package 'tkg-pkg'waiting for package: tkg-pkg
  [2023-01-10T19:50:48.604Z] waiting for package: tanzu-addons-manager

  [2023-01-10T19:50:48.604Z] waiting for package: tanzu-auth <--- not surey we need THIS PACKAEG on the bootstrap cluster though???

  waiting for package: tanzu-cliplugins
  waiting for package: tanzu-core-management-plugins
  waiting for resource tanzu-addons-manager of type *v1alpha1.PackageInstall to be up and running
  waiting for package: tanzu-featuregates
  waiting for package: tanzu-framework
  waiting for package: tkg-clusterclass
  waiting for resource tanzu-cliplugins of type *v1alpha1.PackageInstall to be up and running
  waiting for package: tkg-clusterclass-vsphere
  waiting for package: tkr-service
  waiting for resource tanzu-framework of type *v1alpha1.PackageInstall to be up and running
  waiting for package: tkr-source-controller
  waiting for package: tkr-vsphere-resolver

```

## Now, we wait for PackageInstall's to complete on the kind cluster...

Were now inside of `WaitForManagementPackages` in tanzu cli.  We're going to wait for all the children of `tkg-pkg` to be
installed.... once these are installed, `kind` can do everything it needs to do, in order to bootstrap a management cluster.

Why does our "bootstrap" cluster need: tanzu-auth ?  (which comes with `tkg-pkg`) to be installed as a `PackageInstall` ? 


```
 waiting for resource tkr-service of type *v1alpha1.PackageInstall to be up and running
 successfully reconciled package: tkr-source-controller
 successfully reconciled package: tanzu-auth
 successfully reconciled package: tanzu-featuregates
 successfully reconciled package: tkg-clusterclass-vsphere
 successfully reconciled package: tanzu-framework
 successfully reconciled package: tkg-pkg
 successfully reconciled package: tanzu-core-management-plugins
 successfully reconciled package: tanzu-cliplugins
 successfully reconciled package: ako-operator
 successfully reconciled package: tkr-service
 successfully reconciled package: tkg-clusterclass
 successfully reconciled package: tanzu-addons-manager
 successfully reconciled package: tkr-vsphere-resolver
```
Carrying on 
```
[2023-01-10T19:50:48.605Z] Get AVI_CONTROL_PLANE_HA_PROVIDER from user config 
[2023-01-10T19:50:48.605Z] Installing AKO on bootstrapper...
[2023-01-10T19:50:48.605Z] Using default value for CONTROL_PLANE_MACHINE_COUNT = 1. Reason: CONTROL_PLANE_MACHINE_COUNT variable is not set
[2023-01-10T19:50:48.605Z] Fetching File="cluster-template-definition-devcc.yaml" Provider="vsphere" Type="InfrastructureProvider" Version="v1.5.1"
[2023-01-10T19:50:48.605Z] Management cluster config file has been generated and stored at: '/home/kubo/.config/tanzu/tkg/clusterconfigs/tkg-mgmt-vc.yaml'
[2023-01-10T19:50:48.605Z] Checking Tkr v1.24.9---vmware.1-tkg.1-rc.2 in bootstrap cluster...
[2023-01-10T19:50:48.605Z] waiting for resource v1.24.9---vmware.1-tkg.1-rc.2 of type *v1alpha3.TanzuKubernetesRelease to be up and running
```
 