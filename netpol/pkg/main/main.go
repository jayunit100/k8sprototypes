package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/jayunit100/k8sprototypes/netpol/pkg/utils"
)

func bootstrap() {
	pods := []string{"a", "b", "c"}
	namespaces := []string{"x", "y", "z"}
	k8s := utils.Kubernetes{}
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

func validate(m *utils.ReachableMatrix) {
	pods := []string{"a", "b", "c"}
	namespaces := []string{"x", "y", "z"}
	k8s := utils.Kubernetes{}

	// better as metrics, obviously, this is only for POC.
	for _, n1 := range namespaces {
		log.Infof("invoked ", n1)
		for _, p1 := range pods {
			for _, n2 := range namespaces {
				for _, p2 := range pods {
					fmt.Println("main observ:", n1, p1, n2, p2)
					connected, err := k8s.Probe(n1, p1, n2, p2, 80)
					m.Observe(n1,p1,n2,p2,connected)
					if err != nil {
						//log.Info("....confirm firewall on ",n1,p1,n2,p2)
					}
				}
			}
		}
	}

}

func main(){
	bootstrap()
	matrix := TestPodLabelWhitelistingFromBToA()
	validate(matrix)
	summary, pass := matrix.Summary()
	fmt.Println(summary, pass)

}

func TestPodLabelWhitelistingFromBToA() *utils.ReachableMatrix{
	k8s := utils.Kubernetes{}
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

	return m
}
