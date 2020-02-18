package main

import (
	"fmt"

	"github.com/jayunit100/k8sprototypes/netpol/pkg/utils"
)

func main() {
	pods := []string{"a", "b", "c"}
	namespaces := []string{"x", "y", "z"}
	k8s := utils.Kubernetes{}

	p80 := 80
	//p81 := 81

	for _, ns := range namespaces {
		k8s.CreateNamespace(ns, nil)
		for _, pod := range pods {
			fmt.Println(ns)
			k8s.CreateDeployment(ns, ns+pod, 1,
				map[string]string{
					"pod": pod,
				}, "nginx:1.8-alpine") // old nginx cause it was before people deleted everything useful from containers
		}
	}

	// An example test:
	m := utils.ReachableMatrix{
		DefaultExpect: true,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	// Verify all connectivity
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("allow-client-a-via-pod-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress()
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, nil, nil, nil)
	k8s.CreateNetworkPolicy("x", builder.Get())
	for _, n1 := range namespaces {
		for _, p1 := range pods {
			for _, n2 := range namespaces {
				for _, p2 := range pods {
					if p2 == "a" && n1 == "x" {
						m.Expect(n1, p1, n2, p2, p1 == "b")
					}
				}
			}
		}
	}

	// better as metrics, obviously, this is only for POC.
	for _, n1 := range namespaces {
		for _, p1 := range pods {
			for _, n2 := range namespaces {
				for _, p2 := range pods {
					p1pod := k8s.GetPods(n1, "pod", p1)[0].GetName()
					p2pod := k8s.GetPods(n2, "pod", p2)[0].GetName()
					connected := k8s.Probe(n1, p1pod, n2, p2pod, 80)
					m.Observe(n1, p1, n2, p2, connected)
				}
			}
		}
	}
	summary, pass := m.Summary()
	fmt.Println(summary, pass)
}
