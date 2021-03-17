#!/bin/bash

function get_antrea_logs() {
	for i in `kubectl get nodes -o jsonpath='{.items[*].status.addresses[?(@.type=="ExternalIP")].address}'` ; do
		ssh -o StrictHostKeyChecking=no capv@$i "cd C:/k/antrea/logs ; Get-ChildItem -Path C:\k\antrea\logs\ -Name | cat"
	done
}

function update_tomls() {
	for i in `kubectl get nodes -o jsonpath='{.items[*].status.addresses[?(@.type=="ExternalIP")].address}'` ; do 
		echo "START --------------  $i"; 
		ssh -o StrictHostKeyChecking=no capv@$i "Stop-Service Kubelet ;"
		ssh -o StrictHostKeyChecking=no capv@$i "curl.exe https://raw.githubusercontent.com/jayunit100/k8sprototypes/master/kind/capi-containerd-win.toml > a.txt"
		ssh -o StrictHostKeyChecking=no capv@$i "Get-Content a.txt | out-file -encoding ASCII 'C:\Program Files\containerd\config.toml'"
		ssh -o StrictHostKeyChecking=no capv@$i "Get-Content 'C:\Program Files\containerd\config.toml' | Measure-Object -Line "
		ssh -o StrictHostKeyChecking=no capv@$i "Restart-Service containerd"
		ssh -o StrictHostKeyChecking=no capv@$i "Start-Service Kubelet" 
		echo "DONE -------------- $i"  ; 
	done
}

function clear_defender() {
	for i in `kubectl get nodes -o jsonpath='{.items[*].status.addresses[?(@.type=="ExternalIP")].address}'` ; do 
		ssh -o StrictHostKeyChecking=no capv@$i "Uninstall-WindowsFeature Windows-Defender" ; echo $i ; 
	done                                       
}

function reboot() {
	for i in `kubectl get nodes -o jsonpath='{.items[*].status.addresses[?(@.type=="ExternalIP")].address}'` ; do 
		ssh -o StrictHostKeyChecking=no capv@$i "Reboot-Machine" ; echo $i ; 
	done
}

if [[ "$1" == "t" ]] ; then
	update_tomls
fi

if [[ "$1" == "c" ]] ; then
	clear_defender
fi

if [[ "$1" == "r" ]] ; then
	reboot
fi

if [[ "$1" == "a" ]] ; then
	get_antrea_logs	
fi
