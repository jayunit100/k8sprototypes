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

mkdir /home/kubo/geetika/
mkdir /home/kubo/geetika/output

docker run -it --rm \
  -v /home/kubo/geetika/tkgprotoform/image-builder-credentials.json:/home/imagebuilder/vsphere.json \
  -v /home/kubo/geetika/tkgprotoform/image-builder-tkg.json:/home/imagebuilder/tkg.json \
  -v /home/kubo/geetika/tkg:/home/imagebuilder/tkg \
  -v /home/kubo/geetika/cis:/home/imagebuilder/cis \
  -v /home/kubo/geetika/goss/vsphere-ubuntu-1.23.10+vmware.2-goss-spec.yaml:/home/imagebuilder/goss/goss.yaml \
  -v /home/kubo/geetika/metadata.json:/home/imagebuilder/metadata.json \
  -v /home/kubo/geetika/output:/home/imagebuilder/output \
  --env PACKER_VAR_FILES="tkg.json vsphere.json" \
  --env OVF_CUSTOM_PROPERTIES=/home/kubo/geetika/metadata.json \
  --env IB_OVFTOOL=1 \
  projects.registry.vmware.com/tkg/image-builder:v0.1.11_vmware.3 \
  build-node-ova-vsphere-ubuntu-2004-efi

  # pop out of 1_6_0
  popd
# pop out of TKG-v1.6.0-master
popd