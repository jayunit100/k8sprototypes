# vCenter config/credentials
export VSPHERE_SERVER=                # (required) The vCenter server IP or FQDN
export VSPHERE_USERNAME=   # (required) The username used to access the remote vSphere endpoint
export VSPHERE_PASSWORD=   # (required) The password used to access the remote vSphere endpoint

# vSphere deployment configs
export VSPHERE_DATACENTER=         # (required) The vSphere datacenter to deploy the management cluster on
export VSPHERE_DATASTORE=         # (required) The vSphere datastore to deploy the management cluster on
export VSPHERE_NETWORK=              # (required) The VM network to deploy the management cluster on
export VSPHERE_RESOURCE_POOL=           # (required) The vSphere resource pool for your VMs
export VSPHERE_FOLDER=                        # (optional) The VM folder for your VMs, defaults to the root vSphere folder if not set.
export VSPHERE_TEMPLATE=  # (required) The VM template to use for your management cluster.
export VSPHERE_MACHINE_TEMPLATE=
export VSPHERE_SSH_AUTHORIZED_KEY=
export KUBERNETES_VERSION=
export VSPHERE_HAPROXY_TEMPLATE=
export VSPHERE_DISK_GIB=                       # (optional) The VM Disk size in GB, defaults to 20 if not set
export VSPHERE_NUM_CPUS=                         # (optional) The # of CPUs for control plane nodes in your management cluster, defaults to 2 if not set
export VSPHERE_MEM_MIB=                       # (optional) The memory (in MiB) for control plane nodes in your management cluster, defaults to 2048 if not set
