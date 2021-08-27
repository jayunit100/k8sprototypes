These files can be used to unit test tanzu framework YTT without having tanzu cli.

Just clone down tanzu-framework, and run 

```
cat ../infrastructure-windows-vsphere/v0.7.8/ytt/base-template.yaml | \

ytt --data-values-file /tmp/bom.yaml --data-values-file /tmp/config.yaml -f- | \

ytt --data-values-file /tmp/bom.yaml --data-values-file /tmp/config.yaml -f ../infrastructure-windows-vsphere/v0.7.8/ytt/overlay.yaml -f ./ -f-

```
