#!/bin/bash

# PURPOSE: To detect wether the bits needed to install a standard TKG (ubuntu OVA, tanzu framework zip) are available.
# INPUT: NONE
# OUTPUT:
#
# cli              ./payload/tanzu-cli-bundle-darwin-amd64.tar.gz
# ova-ubuntu:      ./payload/ubuntu-2004-kube-v1.23.8+vmware.2-tkg.1-85a434f93857371fccb566a414462981.ova
#
# FAILURE: Fails (exit 1) if both OVA and tanzu framework.tar.gz arent present.

PAYLOAD="./payload"

# TODO make the path customizable
FILE1="`find $PAYLOAD | grep tanzu-cli-bundle | grep tar | grep gz`"

# TODO make this so it finds either photon or ubuntu OVAs.
FILE2="`find $PAYLOAD | grep ova | grep ubuntu | head -1`"

FILE1=$FILE1
FILE2=$FILE2

if test -f "$FILE1" && test -f $FILE2; then
    echo "cli        \t $FILE1"
    echo "ova-ubuntu:\t $FILE2"
    exit 0
fi

# Hope we dont get here...

echo "Missing files, check that "
echo "tanzu cli (at location $FILE1 )"
echo "ubuntu ova (at location $FILE2 )"
echo "are valid paths in this dir."
echo "you must have files that say 'tanzu-cli-bundle'"
echo "and ubuntu in this dir, or else, cant install tanzu!"
exit 1
