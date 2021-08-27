These files can be used to unit test tanzu framework YTT without having tanzu cli.

Just clone down tanzu-framework, and run 

```
cat ../infrastructure-windows-vsphere/v0.7.8/ytt/base-template.yaml | \

ytt --data-values-file /tmp/bom.yaml --data-values-file /tmp/config.yaml -f- | \

ytt --data-values-file /tmp/bom.yaml --data-values-file /tmp/config.yaml -f ../infrastructure-windows-vsphere/v0.7.8/ytt/overlay.yaml -f ./ -f-

```

# Dedication

This file is dedicated to the original TKG-Windows team, Gab Satchi, PeriThompson, and Jay Vyas, who created 
the original CAPV windows YTT installation templates, named "peri.min.sh".  

Please honor them with a moment of silence every time that you run this script.
Or, every time that you ever run any kind of Kubernetes container on any cluster API Installation.

