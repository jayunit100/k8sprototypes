# Phases of TKG Installation (Glasgow)

The below chronicles all of the common logs that youll see during TKG installation, and what phases of the overall
Tanzu Framework and TKG installation process they correspond to.  If you have a broken test environment for TKG, 
scanning through these might help you to identify the specific reason/component/codebase associated with your failure.

This is a work in progress, PRs welcome !

## Part 0: Tanzu cli installation and Kind Cluster

- [TKG Client and Kind Bootstrapping](0-bootstrap-cluster.md#tkg-client-and-kind-bootstrapping)
  * [Making sure the client plugins are there...](0-bootstrap-cluster.md#making-sure-the-client-plugins-are-there)
- [Uploading the OVAs to Vsphere](0-bootstrap-cluster.md#uploading-the-ovas-to-vsphere)
- [Create the bootstrap cluster](0-bootstrap-cluster.md#tkg-kind-cluster--create-the-bootstrap--ie-kind----management-cluster--)
  * [Kind ClusterTKG And ClusterClass customizations](0-bootstrap-cluster.md#kind-cluster--tkg-and-clusterclass-customizations)
  * [InitRegion: Cluster API installation onto the Kind cluster](0-bootstrap-cluster.md#initregion--cluster-api-installation-onto-the-kind-cluster)
    + [Telling bootstrap about CAPI](0-bootstrap-cluster.md#tell-kind-about-capi)

- [Kind Cert Manager setup](0-bootstrap-cluster.md#kind-cert-manager-setup)

### Part 0: Tanzu Customization ~ Setting up the kind cluster w/ TKG and CAPI capabilities
- [Installing CAPI objects onto the Kind cluster](0-bootstrap-cluster.md#installing-capi-objects-onto-the-kind-cluster)
- [Wait for CAPI PRoviders on the kind cluster](0-bootstrap-cluster.md#wait-for-capi-providers-on-the-kind-cluster)
- [MAGIC PART: tkg-pkg and "Waiting for package" hot loop](0-bootstrap-cluster.md#magic-part--tkg-pkg-and--waiting-for-package--hot-loop)
  * [Now, we wait for PackageInstall's to complete on the kind cluster...](0-bootstrap-cluster.md#now--we-wait-for-packageinstall-s-to-complete-on-the-kind-cluster)

## Part 1: Migrating Kind to the Management Cluster 

- [TKG Management Cluster Creation (post-kind bootstrapping)](1-mgmt-cluster.log.md#tkg-management-cluster-creation--post-kind-bootstrapping-)
  * [Create the CAPI object](1-mgmt-cluster.log.md#create-the-capi-object)
        * [These logs from from tanzu cli,](1-mgmt-cluster.log.md#these-logs-from-from-tanzu-cli-)
- [Wait for CAPI Providers on the management cluster](1-mgmt-cluster.log.md#wait-for-capi-providers-on-the-management-cluster)
- [MAGIC PART: tkg-pkg and "Waiting for package" hot loop](1-mgmt-cluster.log.md#magic-part--tkg-pkg-and--waiting-for-package--hot-loop)
  * [Performing the move](1-mgmt-cluster.log.md#performing-the-move)
    + [Move: part 0,  Discovery](1-mgmt-cluster.log.md#move--part-0---discovery)
    + [Move: part 1, Now we know WHAT to move to the mgmt cluster](1-mgmt-cluster.log.md#move--part-1--now-we-know-what-to-move-to-the-mgmt-cluster)
    + [Move: part 2,Patch the management cluster after the fact](1-mgmt-cluster.log.md#move--part-2-patch-the-management-cluster-after-the-fact)

## Part 2: Workload Cluster Creation

- [TKG Workload Cluster Creation](2-workload-cluster.log.md#tkg-workload-cluster-creation)
  * [Now we create our first workload cluster:](2-workload-cluster.log.md#now-we-create-our-first-workload-cluster-)
  * [This is normal, b/c we dont usually install pinniped in testbeds.](2-workload-cluster.log.md#this-is-normal--b-c-we-dont-usually-install-pinniped-in-testbeds)
  * [We dont make old clusters, instead we use cluster class for everything...](2-workload-cluster.log.md#we-dont-make-old-clusters--instead-we-use-cluster-class-for-everything)
  * [Telling Tanzu cli to create a WL cluster](2-workload-cluster.log.md#Telling-Tanzu-cli-to-create-a-WL-cluster)
  * [Why do we patch the workload cluster object ?](2-workload-cluster.log.md#why-do-we-patch-the-workload-cluster-object--)
  * [Wait a while...](2-workload-cluster.log.md#wait-a-while)
  * [Wait for controlplane...](2-workload-cluster.log.md#wait-for-controlplane)
  * [Control plane up now.... now we can poll the workload cluster till its up:](2-workload-cluster.log.md#control-plane-up-now-now-we-can-poll-the-workload-cluster-till-its-up-)
  * [Waiting for addons packages (like antrea) to come online:](2-workload-cluster.log.md#waiting-for-addons-packages--like-antrea--to-come-online-)
  * [Heres where we loop through all the packages:](2-workload-cluster.log.md#heres-where-we-loop-through-all-the-packages-)
  * [Now some of the packages start appearing as installed.](2-workload-cluster.log.md#now-some-of-the-packages-start-appearing-as-installed)

# Some things worth noting

I havent vetted these statements 100% yet, but it appears to be that they are true, or close to true:

- There is a mega function, called `InitRegion` in tanzu-framework that has all the logic for bootstrapping kind, making a mgmt cluster, and migrating to mgmt cluster.
- Reading Tanzu CLI logs means being able to interpolate back and forth between the `clusterctl` logs, and the `tanzu cli` logs
- We install `tkg-pkg` on both the kind cluster and the mgmt cluster.  It is a super-package - it has other packages in it.
  - Tanzu CLI after installing tkg-pkg, makes sure there's a PackageInstall object for everything in `tkg-pkg`.
  - This makes Glasgow different from Fuji: In Fuji, the entire kind cluster wasnt able to bootstrap all of its own packages, and neither was the ManagementCluster.
  - Thus, in Glasgow, ClusterResourceSets arent needed anymore
- Question: When Tanzu CLI installs `tkg-pkg` are there any other customizations needed ? 
- You'll see alot of `succesfully reconciled package` statements in both kind and mgmt bootstrapping.  Thats because both phases have the same logical flow:
  - Wait for CAPI to come up
  - Create PackageInstalls on the cluster
  - Wait for all packages to come up

# Open to PRs

This could be published somewhere else.  Theres probably some things that arent well annotated and need more description.  Feel free to PR if theres something
you understand that isnt described properly here!