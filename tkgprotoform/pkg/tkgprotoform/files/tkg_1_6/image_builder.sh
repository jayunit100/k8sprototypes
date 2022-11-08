#!/bin/bash
# See https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.6/vmware-tanzu-kubernetes-grid-16/GUID-build-images-linux.html#linux-create
# Download image-builder from https://developer.vmware.com/samples/7984/tkg-image-builder-for-kubernetes-v1.23.8-on-tkg-v1.6.0?h=Image%20Builder

docker pull projects.registry.vmware.com/tkg/linux-resource-bundle:v1.23.8_vmware.2-tkg.1
docker run -d -p 3000:3000 projects.registry.vmware.com/tkg/linux-resource-bundle:v1.23.8_vmware.2-tkg.1

# TODO: Support node hardening / internet-restricted workflows by adding "custom_role_name" and "https_proxy" settings.

mv ~/Downloads/TKG-Image-Builder-for-Kubernetes-v1.23.8-on-TKG-v1.6.0-master.zip .
unzip ./TKG-Image-Builder-for-Kubernetes-v1.23.8-on-TKG-v1.6.0-master.zip
pushd  TKG-Image-Builder-for-Kubernetes-v1.23.8-on-TKG-v1.6.0-master
  pushd TKG-Image-Builder-for-Kubernetes-v1_23_8---vmware_2-tkg_v1_6_0



  # pop out of 1_6_0
  popd
# pop out of TKG-v1.6.0-master
popd