#/bin/bash

# PURPOSE: Install tanzu cli so its locally callable, and upload all OVAs to vsphere
# INPUT: tanzu-framework tar gz, all OVAs
# OUTPUT: None

# initialize Tanzu

pushd payload

#tar -xvf tanzu-cli-bundle-darwin-amd64.tar.gz

# same for linux or os x...
tar -xvf tanzu-cli-bundle-*.tar.gz
cp cli/core/v0.25.0/tanzu-core-* ./tanzu

./tanzu init
./tanzu plugin sync
./tanzu plugin list

popd