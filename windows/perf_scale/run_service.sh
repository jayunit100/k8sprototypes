#!/bin/sh

YAML=service_${1}.yaml

kubectl apply -f $YAML
date
kubectl wait --for=condition=available --timeout=900s deployment/service-agn
date

kubectl delete -f $YAML
