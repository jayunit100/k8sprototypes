package main

import (
	"fmt"
	"github.com/jayunit100/k8sprototypes/netpol/pkg/utils"
)

func main() {
	fmt.Println("hi")

	// Scenario creation
	pods := []string{"a","b","c"}
	namespaces := []string{"A","B","C"}

	for _,ns := range namespaces {
		k8s := utils.Kubernetes{}
		k8s.CreateNamespace(ns, nil)
		for _,pod := range pods {
			k8s.CreateDeployment(ns, ns+pod, 1 , map[string]string{"pod":"a"}, "nginx" )
		}
	}

	// An example test:
	m := utils.ReachableMatrix{
		DefaultExpect: false,
		Pods: pods,
		Namespaces: namespaces,
	}
	m.Expect("A","a","B","b",5)
	/**
			whitelist := map[string]bool{}...
			In m:
			 For namespaces (a b c)
				m.WithDeployments (a b c)

			p = Use the builder to make a network policy.
			In m:
				ApplyPolicies(p)

			// validate policies on whitelist
			r := newReachableMatrix()
			for namespaces:
				for pods:
					for namespaces:
						for pods:
							ReachableMatrix.add(n1, pod1, n2, pod2, m.Probe(n1, pod1, n2, pod2)

			testResult := false

	**/

}
