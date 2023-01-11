# This is an annotated log of tanzu installation

IT takes all of the logs from a typical TKG CI job for

- Installing a Kind cluster / Bootstrapping TKG Client (https://github.com/jayunit100/k8sprototypes/tree/master/tkg-annotated-install/glasgow/0-bootstrap-cluster.md).
- Installing a Mgmt cluster on that kind cluster (https://github.com/jayunit100/k8sprototypes/tree/master/tkg-annotated-install/glasgow/1-mgmt-cluster.log.md).
- Installing a workload cluster (https://github.com/jayunit100/k8sprototypes/tree/master/tkg-annotated-install/glasgow/2-workload-cluster.md).

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








