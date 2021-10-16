#!/bin/bash
 
rounds=$1
DEPLOYMENT_YAML=deployment.yaml
 
start_date=`date  +%s`
 
echo -----------------------------------------------------
echo ------------ Start `date --date=@$start_date`
echo ----------------------------------------------------
 
for k in $(seq 1 $rounds)
do
  echo "===========> Round $k/$rounds"
  podName=()
 
  # Create service
  kubectl apply -f ${DEPLOYMENT_YAML}
 
  # Wait for all pods created
  while true
  do
    pod=`kubectl get po -o name| grep service-agn`
    if [ "$pod" == "" ];then
      echo "Not found pod service-agn"
      sleep 1
      continue
    fi
    podName=($pod)
    break
  done
 
  # Wait for all pods ready
  #kubectl wait --for=condition=available --timeout=900s deployment/service-agn
  for pn in ${podName[*]}
  do
    kubectl wait --for=condition=Ready $pn
  done
 
  # Wait for servcie available
  while true
  do
    svc=`kubectl get svc -A -o name | grep service-agn`
    if [ "$svc" != "" ];then
      echo "Service is running: $svc"
      break
    fi
    echo "Not found service service-agn"
    sleep 1
  done
 
  # Delete servcie
  kubectl delete -f ${DEPLOYMENT_YAML}
 
  # Wait for all pods deleted
  for pn in ${podName[*]}
  do
    kubectl wait --for=delete $pn
  done
done
 
end_date=`date  +%s`
 
echo -----------------------------------------------------
echo ------------ End `date --date=@$end_date`
echo ----------------------------------------------------
 
let dura=($end_date-$start_date)
 
echo "Total duration for $rounds rounds is $dura seconds!!!"
