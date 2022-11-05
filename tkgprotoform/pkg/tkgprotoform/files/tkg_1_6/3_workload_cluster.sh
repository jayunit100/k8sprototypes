#!/bin/sh

# INPUT: MIN_CLUSTER (included), MAX_CLUSTER (excluded)
# i.e. MIN_CLUSTER=0 , MAX_CLUSTER=99 ==> creates 100 clusters.
# OUTPUT: Creates multiple clusters, per input range.

if [ -z "$MIN_CLUSTER" ]; then
  MIN_CLUSTER=0
fi
if [ -z "$MAX_CLUSTER" ]; then
  MAX_CLUSTER=1
fi

i=${MIN_CLUSTER}
# try to create MAX-MIN clusters...
# tanzu cluster create -f cluster.yaml wl-1
# tanzu cluster create -f cluster.yaml wl-2
# tanzu cluster create -f cluster.yaml wl-3
while [ $i -ne ${MAX_CLUSTER} ]
do
        i=$(($i+1))
	tanzu cluster create -f cluster.yaml wl-$i
        echo "$i"
done