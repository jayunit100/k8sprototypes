package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/jayunit100/k8sprototypes/netpol/pkg/utils"
)

func bootstrap(k8s *utils.Kubernetes) {
	pods := []string{"a", "b", "c"}
	namespaces := []string{"x", "y", "z"}
 	//p81 := 81

	for _, ns := range namespaces {
		k8s.CreateNamespace(ns, map[string]string{"ns": ns})
		for _, pod := range pods {
			fmt.Println(ns)
			k8s.CreateDeployment(ns, ns+pod, 1,
				map[string]string{
					"pod": pod,
				}, "nginx:1.8-alpine") // old nginx cause it was before people deleted everything useful from containers
		}
	}
}

func validate(k8s *utils.Kubernetes, m *utils.ReachableMatrix) {
	pods := []string{"a", "b", "c"}
	namespaces := []string{"x", "y", "z"}
	// better as metrics, obviously, this is only for POC.
	for _, n1 := range namespaces {
		for _, p1 := range pods {
			for _, n2 := range namespaces {
				for _, p2 := range pods {
					fmt.Println("main observ:", n1, p1, n2, p2)
					connected, _ := k8s.Probe(n1, p1, n2, p2, 80)
					m.Observe(n1,p1,n2,p2,connected)
					if !connected {
						if m.Expected[n1+"_"+p1][n2+"_"+p2] {
							log.Warnf("FAILED CONNECTION FOR WHITELISTED PODS %v %v -> %v %v !!!! ", n1, p1, n2, p2)
						}
					}
				}
			}
		}
	}
}

func main(){
	k8s := utils.Kubernetes{}
	bootstrap(&k8s)
	matrix := TestPodLabelWhitelistingFromBToA(&k8s)
	validate(&k8s, matrix)
	summary, pass := matrix.Summary()
	fmt.Println(summary, pass)
}

func TestPodLabelWhitelistingFromBToA(k8s *utils.Kubernetes) *utils.ReachableMatrix{
	pods := []string{"a", "b", "c"}
	namespaces := []string{"x", "y", "z"}
	p80 := 80
	m := &utils.ReachableMatrix{
		DefaultExpect: true,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("allow-client-a-via-pod-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress()
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, map[string]string{"ns": "x"}, nil, nil)
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, map[string]string{"ns": "y"}, nil, nil)
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, map[string]string{"ns": "z"}, nil, nil)
	k8s.CreateNetworkPolicy("x", builder.Get())
	m.ExpectAllIngress("x","a",false)
	m.Expect("x", "b", "x", "a", true)
	m.Expect("y", "b", "x", "a", true)
	m.Expect("z", "b", "x", "a", true)
	m.Expect("x", "a", "x", "a", true)
	// TODO move this to a unit test !
	if m.Expected["z_c"]["x_a"] == true  {
		panic("expectatilns are wrongg")
	}
	return m
}
