package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/jayunit100/k8sprototypes/netpol/pkg/utils"
)

type NetPolConfig struct {
	pods []string
	namespaces []string
	k8s *utils.Kubernetes
}
var p80 int = 80
var p81 int = 81

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

func validate(k8s *utils.Kubernetes, m *utils.ReachableMatrix, reachability *utils.Reachability) {
	pods := []string{"a", "b", "c"}
	namespaces := []string{"x", "y", "z"}
	// better as metrics, obviously, this is only for POC.
	for _, n1 := range namespaces {
		for _, p1 := range pods {
			for _, n2 := range namespaces {
				for _, p2 := range pods {
					log.Infof("main observe: %s-%s, %s-%s", n1, p1, n2, p2)
					connected, err := k8s.Probe(n1, p1, n2, p2, 80)
					if err != nil {
						log.Errorf("unable to make main observation on %s-%s -> %s-%s: %s", n1, p1, n2, p2, err)
					}
					m.Observe(n1,p1,n2,p2,connected)
					reachability.Observe(utils.NewPod(n1, p1), utils.NewPod(n2, p2), connected)
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
	k8s, err := utils.NewKubernetes()
	if err != nil {
		panic(err)
	}
	bootstrap(k8s)
	matrix, reachability := TestPodLabelWhitelistingFromBToA(k8s)
	validate(k8s, matrix, reachability)

	summary, pass := matrix.Summary()
	fmt.Println(summary, pass)

	right, wrong, comparison := reachability.Summary()
	fmt.Printf("reachability: correct:%v, incorrect:%v, result=%t\n\n", right, wrong, wrong == 0)
	fmt.Printf("expected:\n\n%s\n\n\n", reachability.Expected.PrettyPrint())
	fmt.Printf("observed:\n\n%s\n\n\n", reachability.Observed.PrettyPrint())
	fmt.Printf("comparion:\n\n%s\n\n\n", comparison.PrettyPrint())
}

func TestPodLabelWhitelistingFromBToA(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability) {
	pods := []string{"a", "b", "c"}
	namespaces := []string{"x", "y", "z"}

	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("allow-client-a-via-pod-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress()
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, map[string]string{"ns": "x"}, nil, nil)
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, map[string]string{"ns": "y"}, nil, nil)
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, map[string]string{"ns": "z"}, nil, nil)
	k8s.CreateNetworkPolicy("x", builder.Get())

	m := &utils.ReachableMatrix{
		DefaultExpect: true,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	m.ExpectAllIngress("x","a",false)
	m.Expect("x", "b", "x", "a", true)
	m.Expect("y", "b", "x", "a", true)
	m.Expect("z", "b", "x", "a", true)
	m.Expect("x", "a", "x", "a", true)

	allPods := []utils.Pod{}
	for _, podName := range pods {
		for _, ns := range namespaces {
			allPods = append(allPods, utils.NewPod(ns, podName))
		}
	}
	reachability := utils.NewReachability(allPods)
	reachability.ExpectAllIngress(utils.NewPod("x", "a"), false)
	reachability.Expect(utils.NewPod("x", "b"), utils.NewPod("x", "a"), true)
	reachability.Expect(utils.NewPod("y", "b"), utils.NewPod("x", "a"), true)
	reachability.Expect(utils.NewPod("z", "b"), utils.NewPod("x", "a"), true)
	reachability.Expect(utils.NewPod("x", "a"), utils.NewPod("x", "a"), true)

	// TODO move this to a unit test !
	if m.Expected["z_c"]["x_a"] == true  {
		panic("expectations are wrong")
	}
	return m, reachability
}
