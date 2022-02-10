The Kubernetes ecosystem currently supports 100s of different components.

This project attempts to generate a combinatorial expansion of how end users
might upgrade these components, over time, and quantify the amount of configuration
drift that might need to be tested and verified by a vendor.

Our end goal is to suggest what constraints on versioning and upgrading are
reasonable as suggestions for K8s vendors wanting to support the penumbra of
expanding CNI, CSI, Authentication, Backup, Loadbalancing, and Security solutions
which often come bundled in a K8s distribution.

# Assume
- Testing an upgrade of a single Kubernetes configuration to any other, takes 4 hours.
... 
