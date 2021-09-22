# Problem

What id like to do is... generate some kind of 
CAPI skeleton, and match YTT output to it from tanzu framework

```
targetCluster := Cluster{
..}
targetVsphereCluster := VsphereCluster{
..}
targetMD := MachineDeployment{
..}

ginkgo.Describe("should match targets")...
   ....
```


# Non Solution

Take tanzu framework and generate this and read it by eye :) 
These files can be used to unit test tanzu framework YTT without having tanzu cli.

Just clone down tanzu-framework, and run 

```
cat ../infrastructure-windows-vsphere/v0.7.8/ytt/base-template.yaml | \

ytt --data-values-file /tmp/bom.yaml --data-values-file /tmp/config.yaml -f- | \

ytt --data-values-file /tmp/bom.yaml --data-values-file /tmp/config.yaml -f ../infrastructure-windows-vsphere/v0.7.8/ytt/overlay.yaml -f ./ -f-

```

# Dedication
 Hey! 
This file is dedicated to the original TKG-Windows team, Gab Satchi, PeriThompson, and Jay V, 
And the script that started it all: peri.min.sh !!!
