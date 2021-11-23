# Automating Antrea CNI lifecycle with Carvel

- https://github.com/vmware-tanzu/tanzu-framework/pull/959/, addon reconcilers
- How tanzu addons for CNI, etc, work
    - kubetail on https://github.com/vmware-tanzu/tanzu-framework/tree/main/addons#workflow-of-tanzu-addons-manager
    - https://carvel.dev/kapp-controller/docs/latest/packaging/
        - PackageRepository
        - PackageMetadata
        - Package
        - PackageInstall 
    - https://carvel.dev/kapp-controller/docs/latest/app-spec/
        - note the `cluster` field... its multicluster by default
