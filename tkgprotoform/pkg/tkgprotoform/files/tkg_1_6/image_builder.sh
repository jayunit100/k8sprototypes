#!/bin/bash
START_DIR=/home/kubo/geetika
mkdir $START_DIR
pushd $START_DIR

mkdir output
# See https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.6/vmware-tanzu-kubernetes-grid-16/GUID-build-images-linux.html#linux-create
# Download image-builder from https://developer.vmware.com/samples/7984/tkg-image-builder-for-kubernetes-v1.23.8-on-tkg-v1.6.0?h=Image%20Builder

docker pull projects-stg.registry.vmware.com/tkg/linux-resource-bundle:v1.23.10_vmware.1-tkg.1-rc.1
docker run -d -p 3000:3000  docker run -d -p 3000:3000 projects-stg.registry.vmware.com/tkg/linux-resource-bundle:v1.23.10_vmware.1-tkg.1-rc.1

# TODO: Support node hardening / internet-restricted workflows by adding "custom_role_name" and "https_proxy" settings.

mv ~/Downloads/TKG-Image-Builder-for-Kubernetes-v1.23.8-on-TKG-v1.6.0-master.zip .
unzip ./TKG-Image-Builder-for-Kubernetes-v1.23.8-on-TKG-v1.6.0-master.zip
pushd  TKG-Image-Builder-for-Kubernetes-v1.23.8-on-TKG-v1.6.0-master
  pushd TKG-Image-Builder-for-Kubernetes-v1_23_8---vmware_2-tkg_v1_6_0

# 1.6.0 another mechanism is used then stig (CIS)
# 1.6.1 need stig
govc vm.destroy /dc0/vm/dc0/vm/ubuntu-2004-efi-kube-v1.23.8 ; rm -rf output/ ; sudo chmod -R 777 ./goss ; mkdir output/ ; sudo chmod -R 777 output ;

IB_ROOT=$START_DIR/TKG-Image-Builder-for-Kubernetes-v1_23_10---vmware_1-tkg_v1_6_1/
DEFAULTS=$START_DIR/tkgprotoform/

cat << EOF > $IB_ROOT/customizations.json
{
        "vmx_version": "17"
}
EOF
cat << EOF > $IB_ROOT/metadata.json
{
        "VERSION": "v1.23.10+vmware.1-capv.1"
}
EOF

### /home/kubo/geetika/TKG-Image-Builder-for-Kubernetes-v1_23_10---vmware_1-tkg_v1_6_1/stig-ubuntu-2004:/home/imagebuilder/stig-ubuntu-2004
docker run -it --rm \
-v $DEFAULTS/image-builder-credentials.json:/home/imagebuilder/vsphere.json \
-v $DEFAULTS/image-builder-tkg.json:/home/imagebuilder/tkg.json \
-v $IB_ROOT/tkg:/home/imagebuilder/tkg \
-v $IB_ROOT/stig_ubuntu_2004/:/home/imagebuilder/stig-ubuntu-2004 \
-v $START_DIR/goss/vsphere-ubuntu-1.23.10+vmware.2-goss-spec.yaml:/home/imagebuilder/goss/goss.yaml \
-v $START_DIR/metadata.json:/home/imagebuilder/metadata.json \
-v $START_DIR/output:/home/imagebuilder/output \
--env PACKER_VAR_FILES="tkg.json vsphere.json" \
--env OVF_CUSTOM_PROPERTIES=/home/imagebuilder/metadata.json \
--env IB_OVFTOOL=1 \
projects-stg.registry.vmware.com/tkg/image-builder:v0.1.13_vmware.1 build-node-ova-vsphere-ubuntu-2004-efi

popd
# pop out of TKG-v1.6.0-master
popd