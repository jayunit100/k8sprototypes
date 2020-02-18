package utils

import (
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type NetworkPolicySpecBuilder struct {
	Spec networkingv1.NetworkPolicySpec
	Name string
}

func (n *NetworkPolicySpecBuilder) Get() *networkingv1.NetworkPolicy {
	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: n.Name,
		},
		Spec: n.Spec,
	}
}

func (n *NetworkPolicySpecBuilder) SetPodSelector(labels map[string]string) *NetworkPolicySpecBuilder {
	ps := metav1.LabelSelector{
		MatchLabels: labels,
	}
	n.Spec.PodSelector = ps
	return n
}

func (n *NetworkPolicySpecBuilder) SetName(name string) *NetworkPolicySpecBuilder {
	n.Name = name
	return n
}

func (n *NetworkPolicySpecBuilder) AddIngress(protoc *v1.Protocol, port *int, portName *string, cidr *string, podSelector map[string]string, nsSelector map[string]string, podSelectorMatchExp *[]metav1.LabelSelectorRequirement, nsSelectorMatchExp *[]metav1.LabelSelectorRequirement) *NetworkPolicySpecBuilder {

	var ps *metav1.LabelSelector
	var ns *metav1.LabelSelector

	if podSelector != nil {
		ps = &metav1.LabelSelector{
			MatchLabels: podSelector,
		}
	}
	if nsSelector != nil {
		ns = &metav1.LabelSelector{
			MatchLabels: nsSelector,
		}
	}

	r := networkingv1.NetworkPolicyIngressRule{
		From: []networkingv1.NetworkPolicyPeer{{
			PodSelector:       ps,
			NamespaceSelector: ns,
			IPBlock:           nil,
		}},
	}
	if n.Spec.Ingress == nil {
		n.Spec.Ingress = []networkingv1.NetworkPolicyIngressRule{}
	}
	n.Spec.Ingress = append(n.Spec.Ingress, r)

	return n
}

// AddEgressDNS mutates the nth policy rule to allow DNS, convenience method
func (n *NetworkPolicySpecBuilder) WithEgressDNS() *NetworkPolicySpecBuilder {
	protocolUDP := v1.ProtocolUDP
	route53 := networkingv1.NetworkPolicyPort{
		Protocol: &protocolUDP,
		Port:     &intstr.IntOrString{Type: intstr.Int, IntVal: 53},
	}

	for _, e := range n.Spec.Egress {
		e.Ports = append(e.Ports, route53)
	}
	return n
}

func (n *NetworkPolicySpecBuilder) AddEgress(protoc *v1.Protocol, port *int, portName *string, cidr *string, podSelector *map[string]string, nsSelector *map[string]string, podSelectorMatchExp *[]metav1.LabelSelectorRequirement, nsSelectorMatchExp *[]metav1.LabelSelectorRequirement) *NetworkPolicySpecBuilder {
	r := networkingv1.NetworkPolicyEgressRule{
		To: []networkingv1.NetworkPolicyPeer{{
			PodSelector: &metav1.LabelSelector{
				MatchLabels: *podSelector,
			},
			NamespaceSelector: &metav1.LabelSelector{
				MatchLabels: *nsSelector,
			},
		}},
	}
	if n.Spec.Egress == nil {
		n.Spec.Egress = []networkingv1.NetworkPolicyEgressRule{}
	}
	n.Spec.Egress = append(n.Spec.Egress, r)
	return n
}

func (n *NetworkPolicySpecBuilder) SetTypeIngress() *NetworkPolicySpecBuilder {
	n.Spec.PolicyTypes = []networkingv1.PolicyType{networkingv1.PolicyTypeIngress}
	return n
}
func (n *NetworkPolicySpecBuilder) SetTypeEgress() *NetworkPolicySpecBuilder {
	n.Spec.PolicyTypes = []networkingv1.PolicyType{networkingv1.PolicyTypeEgress}
	return n
}
func (n *NetworkPolicySpecBuilder) SetTypeBoth() *NetworkPolicySpecBuilder {
	n.Spec.PolicyTypes = []networkingv1.PolicyType{networkingv1.PolicyTypeEgress, networkingv1.PolicyTypeIngress}
	return n
}
