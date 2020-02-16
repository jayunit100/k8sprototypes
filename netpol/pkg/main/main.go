package main

import (
	"fmt"

	"github.com/jayunit100/k8sprototypes/netpol/pkg/utils"
)

func main() {
	pods := []string{"a", "b", "c"}
	namespaces := []string{"x", "y", "z"}
	k8s := utils.Kubernetes{}

	for _, ns := range namespaces {
		k8s.CreateNamespace(ns, nil)
		for _, pod := range pods {
			fmt.Println(ns)
			k8s.CreateDeployment(ns, ns+pod, 1,
				map[string]string{
					"pod": "a",
				}, "busybox")
		}
	}

	// An example test:
	m := utils.ReachableMatrix{
		DefaultExpect: false,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	m.Expect("A", "a", "B", "b", 5)

	fmt.Println(k8s.Probe("x", "xb-768f8cd4-z8gsh", "y", "ya-69c5d95599-q2kg5", 8080))
}
