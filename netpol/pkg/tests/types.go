package tests

import (
	// "github.com/jayunit100/netpoltests/pkg/utils"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

// A Scenario is a network environment which consists of a bunch of pods
type Scenario interface {
	AddDeployment(MatchExpressions map[string]string, Labels map[string]string, Namespace string, podName string, replicas int)
	MutateDeployment(MatchExpressions map[string]string, Labels map[string]string, Namespace string, podName string, replicas int)
	AddSS(MatchExpressions map[string]string, Labels map[string]string, Namespace string, podName string, replicas int)
	ClearPolicys()
	Probe(namespace string, name string, urls []string) []int
	ApplyPolicys(specs []networkingv1.NetworkPolicySpec)
}

type K8sScenario struct {
	Deployments  []*v1.ReplicationController
	Statefulsets []*v1.ReplicationController
	Pods         []*v1.Pod
	Namespaces   []*v1.Namespace
}

func (m *K8sScenario) WithDeployment(podName string, podLabels map[string]string, podImage string) {

}

func (m *K8sScenario) ClearPolicys() {

}

func (m *K8sScenario) ApplyPolicys(specs []networkingv1.NetworkPolicySpec) {

}
