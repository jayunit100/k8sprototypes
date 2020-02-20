package main

import (
	"fmt"
	"github.com/jayunit100/k8sprototypes/netpol/pkg/utils"
	log "github.com/sirupsen/logrus"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NetPolConfig struct {
	pods       []string
	namespaces []string
	k8s        *utils.Kubernetes
}
// common for all tests.  these get hardcoded into the Expect() clauses,
// so, we cant easily parameterize them (well, we could, but that would
// make the code harder to interpret.
var pods []string
var namespaces []string
var p80 int
var p81 int

func init() {
	p80 = 80
	p81 = 81
	pods = []string{"a", "b", "c"}
	namespaces = []string{"x", "y", "z"}
}

func bootstrap(k8s *utils.Kubernetes) {

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

func validate(k8s *utils.Kubernetes, m *utils.ReachableMatrix, reachability *utils.Reachability, port int) {
	pods := []string{"a", "b", "c"}
	namespaces := []string{"x", "y", "z"}
	// better as metrics, obviously, this is only for POC.
	for _, n1 := range namespaces {
		for _, p1 := range pods {
			for _, n2 := range namespaces {
				for _, p2 := range pods {
					log.Infof("main observe: %s-%s, %s-%s", n1, p1, n2, p2)
					connected, err := k8s.Probe(n1, p1, n2, p2, port)
					if err != nil {
						log.Errorf("unable to make main observation on %s-%s -> %s-%s: %s", n1, p1, n2, p2, err)
					}
					m.Observe(n1, p1, n2, p2, connected)
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
	/**
	testWrapperPort80(TestDefaultDeny)
	testWrapperPort80(testInnerNamespaceTraffic)
	testWrapperPort80(testEnforcePodAndNSSelector)
	testWrapperPort80(testEnforcePodOrNSSelector)
	testWrapperPort8081(testPortsPolicies)
	// This is a stacked test b/c of the true arg.
	testWrapperStacked(testPortsPoliciesStackedOrUpdated, true)
	testWrapperPort80(testAllowAll)
	testWrapperPort80(testNamedPort)
	// This is an update test b/c of the  false arg.
	testWrapperStacked(testPortsPoliciesStackedOrUpdated, false)
	*/
	testWrapperPort80(TestPodLabelWhitelistingFromBToA)
}


// testWrapperStaged is for tests which involve steps of mutation.
type Stack struct {
	ReachableMatrix *utils.ReachableMatrix
	Reachability *utils.Reachability
	NetworkPolicy *networkingv1.NetworkPolicy
	Port int
 }

// catch all for any type of test, where we use stacks.  these are validated one at a time.
// probably use this for *all* tests when we port to upstream.
func testWrapperStacked(theTest func(k8s *utils.Kubernetes)(stack []*Stack)) {
	k8s, err := utils.NewKubernetes()
	if err != nil {
		panic(err)
	}
	bootstrap(k8s)

	stack := theTest(k8s)
	for _, s := range stack{
		matrix := s.ReachableMatrix
		reachability := s.Reachability
		policy := s.NetworkPolicy
		if policy != nil {
			_,err := k8s.CreateOrUpdateNetworkPolicy(policy.Namespace, policy)
			if err != nil {
				panic(err)
			}
		}
		validate(k8s, matrix, reachability, 80)
		summary1, pass1 := matrix.Summary()
		fmt.Println(summary1, pass1)

		right, wrong, comparison := reachability.Summary()
		fmt.Printf("reachability: correct:%v, incorrect:%v, result=%t\n\n", right, wrong, wrong == 0)
		fmt.Printf("expected:\n\n%s\n\n\n", reachability.Expected.PrettyPrint())
		fmt.Printf("observed:\n\n%s\n\n\n", reachability.Observed.PrettyPrint())
		fmt.Printf("comparison:\n\n%s\n\n\n", comparison.PrettyPrint())
	}
}


// For dual port tests... confirms both ports 80 and 81
func testWrapperPort8081(theTest func(k8s *utils.Kubernetes)(*utils.ReachableMatrix, *utils.Reachability, *utils.ReachableMatrix, *utils.Reachability )) {
	k8s, err := utils.NewKubernetes()
	if err != nil {
		panic(err)
	}
	bootstrap(k8s)
	matrix80, reachability80, matrix81, reachability81 := theTest(k8s)
	validate(k8s, matrix80, reachability80, 80)
	validate(k8s, matrix81, reachability81, 81)

	summary1, pass1 := matrix80.Summary()
	summary2, pass2 := matrix81.Summary()

	fmt.Println(summary1, pass1)
	fmt.Println(summary2, pass2)

	for _, reachability := range []*utils.Reachability{reachability80, reachability81} {
		right, wrong, comparison := reachability.Summary()
		fmt.Printf("reachability: correct:%v, incorrect:%v, result=%t\n\n", right, wrong, wrong == 0)
		fmt.Printf("expected:\n\n%s\n\n\n", reachability.Expected.PrettyPrint())
		fmt.Printf("observed:\n\n%s\n\n\n", reachability.Observed.PrettyPrint())
		fmt.Printf("comparison:\n\n%s\n\n\n", comparison.PrettyPrint())
	}
}

// simple type of test, majority of tests use this, just port 80
func testWrapperPort80(theTest func(k8s *utils.Kubernetes)(*utils.ReachableMatrix, *utils.Reachability )) {
	k8s, err := utils.NewKubernetes()
	if err != nil {
		panic(err)
	}
	bootstrap(k8s)
	matrix, reachability := theTest(k8s)
	validate(k8s, matrix, reachability, 80)

	summary, pass := matrix.Summary()
	fmt.Println(summary, pass)

	right, wrong, comparison := reachability.Summary()
	fmt.Printf("reachability: correct:%v, incorrect:%v, result=%t\n\n", right, wrong, wrong == 0)
	fmt.Printf("expected:\n\n%s\n\n\n", reachability.Expected.PrettyPrint())
	fmt.Printf("observed:\n\n%s\n\n\n", reachability.Observed.PrettyPrint())
	fmt.Printf("comparison:\n\n%s\n\n\n", comparison.PrettyPrint())
}

func listAllPods() []utils.Pod {
	allPods := []utils.Pod{}
	for _, podName := range pods {
		for _, ns := range namespaces {
			allPods = append(allPods, utils.NewPod(ns, podName))
		}
	}
	return allPods
}

/**
TODO: These 3 tests should be implemented using a different strategy, possibly combined.
ginkgo.It("should allow ingress access from updated namespace [Feature:NetworkPolicy]", func() {
ginkgo.It("should allow ingress access from updated pod [Feature:NetworkPolicy]", func() {
ginkgo.It("should deny ingress access to updated pod [Feature:NetworkPolicy]", func() {
ginkgo.It("should stop enforcing policies after they are deleted [Feature:NetworkPolicy]", func() {

TODO: These tests should be easy to add later.
ginkgo.It("should enforce egress policy allowing traffic to a server in a different namespace based on PodSelector and NamespaceSelector [Feature:NetworkPolicy]", func() {
ginkgo.It("should enforce multiple ingress policies with ingress allow-all policy taking precedence [Feature:NetworkPolicy]", func() {
ginkgo.It("should enforce multiple egress policies with egress allow-all policy taking precedence [Feature:NetworkPolicy]", func() {
ginkgo.It("should allow egress access to server in CIDR block [Feature:NetworkPolicy]", func() {
ginkgo.It("should enforce policies to check ingress and egress policies can be controlled independently based on PodSelector [Feature:NetworkPolicy]", func() {
*/

// should allow egress access on one named port [Feature:NetworkPolicy]
func testEgressOnNamedPort(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability ) {
	namedPorts := "serve-80"
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("allow-client-a-via-named-port-egress-rule").SetPodSelector(map[string]string{"pod": "a"})

	// TODO, for this test to properly test egress, we need to modify the probe to support probing
	// via Service endpoints.
	builder.SetTypeEgress().WithEgressDNS().AddEgress(nil, nil, &namedPorts, nil, nil, nil , nil, nil)

	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := &utils.ReachableMatrix{
		DefaultExpect: false,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	// No egress rules because we're deny all !
	reachability := utils.NewReachability(listAllPods())
	m.Expect("x", "a", "x", "a", true)
	m.Expect("x", "b", "x", "a", true)
	m.Expect("x", "c", "x", "a", true)
	m.Expect("x", "a", "x", "a", true)

	// TODO, maybe add validation that 81 doesn't work as well?
	return m, reachability
}


// should allow ingress access from namespace on one named port [Feature:NetworkPolicy]
func testNamedPortWNamespace(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability ) {
	namedPorts := "serve-80"
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("allow-client-a-via-named-port-ingress-rule").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress().AddIngress(nil, nil, &namedPorts, nil, nil, map[string]string{"ns":"x"}, nil, nil)


	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := &utils.ReachableMatrix{
		DefaultExpect: false,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	// No egress rules because we're deny all !
	reachability := utils.NewReachability(listAllPods())
	m.Expect("x", "a", "x", "a", true)
	m.Expect("x", "b", "x", "a", true)
	m.Expect("x", "c", "x", "a", true)
	m.Expect("x", "a", "x", "a", true)

	// TODO, add validation that 81 doesn't work.
	return m, reachability
}

// testNamedPort should allow ingress access on one named port [Feature:NetworkPolicy]
func testNamedPort(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability ) {
	namedPorts := "serve-80"
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("allow-client-a-via-named-port-ingress-rule").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress().AddIngress(nil, nil, &namedPorts, nil, nil, nil, nil, nil)

	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := &utils.ReachableMatrix{
		DefaultExpect: true,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	// No egress rules because we're deny all !
	reachability := utils.NewReachability(listAllPods())

	// TODO, add validation that 81 doesn't work.
	return m, reachability
}


// testAllowAll should support allow-all policy [Feature:NetworkPolicy]
func testAllowAll(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability ) {
		builder := &utils.NetworkPolicySpecBuilder{}
		builder = builder.SetName("default-deny").SetPodSelector(map[string]string{"pod": "a"})
		builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, nil, nil, nil, nil)
		k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
		m := &utils.ReachableMatrix{
			DefaultExpect: true,
			Pods:          pods,
			Namespaces:    namespaces,
		}
		// No egress rules because we're deny all !
		reachability := utils.NewReachability(listAllPods())
		return m, reachability
}

// This covers two test cases: stacked policy's and updated policies.
// 1) should enforce policy based on Ports [Feature:NetworkPolicy] (disallow 80) (stacked == false)
// 2) should enforce updated policy (stacked == true)
func testPortsPoliciesStacked(k8s *utils.Kubernetes, stacked bool) []*Stack {
	policyName := "policy-that-will-update-for-ports"
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName(policyName).SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress()
	builder.AddIngress(nil, &p80, nil, nil, nil, nil, nil, nil)

	policy1 := builder.Get()
	m1 := &utils.ReachableMatrix{
		DefaultExpect: false,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	reachability1 := utils.NewReachability(listAllPods())

	// using false makes this a test for 'updated' policies...
	if stacked == true {
		policyName = "policy-that-will-update-for-ports-2"
	}
	builder2 := &utils.NetworkPolicySpecBuilder{}
	// by preserving the same name, this policy will also serve to test the 'updated policy' scenario.
	builder2 = builder2.SetName(policyName).SetPodSelector(map[string]string{"pod": "a"})
	builder2.SetTypeIngress()
	builder2.AddIngress(nil, &p81, nil, nil, nil, nil, nil, nil)

	policy2 := builder2.Get()
		m2 := &utils.ReachableMatrix{
		DefaultExpect: true,
		Pods:          pods,
		Namespaces:    namespaces,
	}
 	reachability2 := utils.NewReachability(listAllPods())

	return []*Stack{
		&Stack{
			m1,
			reachability1,
			policy1,
			p81,
		},
		&Stack{
			m2,
			reachability2,
			policy2,
			p80,
		},
		&Stack{
			m2,
			reachability2,
			nil, // nil policy wont be created, this is just a 2nd validation, this time, of port 81.
			p81,
		},
	}
}

// "should enforce policy based on Ports [Feature:NetworkPolicy] (disallow 80)
func testPortsPolicies(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability, *utils.ReachableMatrix, *utils.Reachability ) {
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("allow-x-via-pod-and-ns-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress()
	builder.AddIngress(nil, &p81, nil, nil, nil, nil, nil, nil)

	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())

	m80 := &utils.ReachableMatrix{
		DefaultExpect: false,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	reachability80 := utils.NewReachability(listAllPods())

	m81 := &utils.ReachableMatrix{
		DefaultExpect: true,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	reachability81 := utils.NewReachability(listAllPods())

	return m80, reachability80, m81, reachability81
}

// should enforce policy to allow traffic only from a pod in a different namespace based on PodSelector and NamespaceSelector [Feature:NetworkPolicy]
// should enforce policy based on PodSelector and NamespaceSelector [Feature:NetworkPolicy]
func testEnforcePodAndNSSelector(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability ) {
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("allow-x-via-pod-and-ns-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress()
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod":"b"}, map[string]string{"ns":"y"}, nil, nil)

	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := &utils.ReachableMatrix{
		DefaultExpect: true,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	reachability := utils.NewReachability(listAllPods())
	m.ExpectAllIngress("x", "a", false)
	m.Expect("y", "b", "x", "a", true)

	return m, reachability
}

// should enforce policy based on PodSelector or NamespaceSelector [Feature:NetworkPolicy]
func testEnforcePodOrNSSelector(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability ) {
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("allow-x-via-pod-or-ns-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress()
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod":"b"}, nil, nil, nil)
	builder.AddIngress(nil, &p80, nil, nil, nil, map[string]string{"ns":"y"}, nil, nil)

	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := &utils.ReachableMatrix{
		DefaultExpect: true,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	reachability := utils.NewReachability(listAllPods())
	m.ExpectAllIngress("x", "a", false)
	m.Expect("y", "a", "x", "a", true)
	m.Expect("y", "b", "x", "a", true)
	m.Expect("y", "c", "x", "a", true)
	m.Expect("x", "b", "x", "a", true)
	m.Expect("y", "b", "x", "a", true)
	m.Expect("z", "b", "x", "a", true)
	m.Expect("x", "a", "x", "a", true)

	return m, reachability
}

// should enforce policy based on NamespaceSelector with MatchExpressions[Feature:NetworkPolicy]
func testNamespaceSelectorMatchExpressions(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability ) {
	builder := &utils.NetworkPolicySpecBuilder{}
	selector := []metav1.LabelSelectorRequirement{{
		Key:      "ns",
		Operator: metav1.LabelSelectorOpIn,
		Values:   []string{"y"},
	}}
	builder = builder.SetName("allow-a-via-ns-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, nil, nil, &selector, nil)
	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := &utils.ReachableMatrix{
		DefaultExpect: true,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	reachability := utils.NewReachability(listAllPods())
	m.ExpectAllIngress("x", "a", false)
	m.Expect("y", "a", "x", "a", true)
	m.Expect("y", "b", "x", "a", true)
	m.Expect("y", "c", "x", "a", true)
	m.Expect("x", "a", "x", "a", true)

	return m, reachability
}

// testPodSelectorMatchExpressions should enforce policy based on PodSelector with MatchExpressions[Feature:NetworkPolicy]
func testPodSelectorMatchExpressions(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability ) {
	builder := &utils.NetworkPolicySpecBuilder{}
	selector := []metav1.LabelSelectorRequirement{{
		Key:      "pod",
		Operator: metav1.LabelSelectorOpIn,
		Values:   []string{"b"},
	}}
	builder = builder.SetName("allow-client-b-via-pod-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, nil, nil, &selector, nil)
	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := &utils.ReachableMatrix{
		DefaultExpect: true,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	reachability := utils.NewReachability(listAllPods())
	m.ExpectAllIngress("x", "a", false)
	m.Expect("x", "b", "x", "a", true)
	m.Expect("y", "b", "x", "a", true)
	m.Expect("z", "b", "x", "a", true)
	m.Expect("x", "a", "x", "a", true)

	return m, reachability
}

// testInnerNamespaceTraffic should enforce policy to allow traffic from pods within server namespace based on PodSelector [Feature:NetworkPolicy]
func testIntraNamespaceTrafficOnly(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability ) {
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("allow-client-b-via-pod-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, nil, map[string]string{"ns":"y"}, nil, nil)
	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := &utils.ReachableMatrix{
		DefaultExpect: true,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	reachability := utils.NewReachability(listAllPods())
	m.ExpectAllIngress("x", "a", false)
	m.Expect("y", "a", "x", "a", true)
	m.Expect("y", "b", "x", "a", true)
	m.Expect("y", "c", "x", "a", true)

	return m, reachability
}

// testInnerNamespaceTraffic should enforce policy to allow traffic from pods within server namespace based on PodSelector [Feature:NetworkPolicy]
func testInnerNamespaceTraffic(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability ) {
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("allow-client-b-via-pod-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, map[string]string{"pod":"b"}, nil, nil, nil)
	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := &utils.ReachableMatrix{
		DefaultExpect: true,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	reachability := utils.NewReachability(listAllPods())
	m.ExpectAllIngress("x", "a", false)
	m.Expect("x", "b", "x", "a", true)
	return m, reachability
}


func TestDefaultDeny(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability ) {
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("default-deny").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, map[string]string{}, nil, nil, nil)
	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := &utils.ReachableMatrix{
		DefaultExpect: false,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	// No egress rules because we're deny all !
	reachability := utils.NewReachability(listAllPods())
	return m, reachability
}

func TestPodLabelWhitelistingFromBToA(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability ) {
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("allow-client-a-via-pod-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress()
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, map[string]string{"ns": "x"}, nil, nil)
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, map[string]string{"ns": "y"}, nil, nil)
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, map[string]string{"ns": "z"}, nil, nil)
	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())

	m := &utils.ReachableMatrix{
		DefaultExpect: true,
		Pods:          pods,
		Namespaces:    namespaces,
	}
	m.ExpectAllIngress("x", "a", false)
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
	if m.Expected["z_c"]["x_a"] == true {
		panic("expectations are wrong")
	}
	return m, reachability
}
