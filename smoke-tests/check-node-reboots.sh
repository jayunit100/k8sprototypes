### Poll nodes for their last reboot: 
for i in $(kubectl get nodes -o custom-columns='NAME:.metadata.name,INTERNAL-IP:.status.addresses[?(@.type=="InternalIP")].address' | awk '{print $2}'); do
  echo -n "Node: $i, Uptime: "
  ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR  capv@$i uptime
done
