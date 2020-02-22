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
var allPods []utils.Pod

func init() {
	p80 = 80
	p81 = 81
	pods = []string{"a", "b", "c"}
	namespaces = []string{"x", "y", "z"}

	for _, podName := range pods {
		for _, ns := range namespaces {
			allPods = append(allPods, utils.NewPod(ns, podName))
		}
	}
}

func bootstrap(k8s *utils.Kubernetes) {
	//p81 := 81
	for _, ns := range namespaces {
		k8s.CreateOrUpdateNamespace(ns, map[string]string{"ns": ns})
		for _, pod := range pods {
			fmt.Println(ns)
			k8s.CreateOrUpdateDeployment(ns, ns+pod, 1,
				map[string]string{
					"pod": pod,
				}, "nginx:1.8-alpine") // old nginx cause it was before people deleted everything useful from containers
		}
	}
}

func validate(k8s *utils.Kubernetes, m *utils.ReachableMatrix, reachability *utils.Reachability, port int) {
	// better as metrics, obviously, this is only for POC.
	for _, n1 := range namespaces {
		for _, p1 := range pods {
			for _, n2 := range namespaces {
				for _, p2 := range pods {
					log.Infof("Probing: %s-%s, %s-%s", n1, p1, n2, p2)
					connected, err := k8s.Probe(n1, p1, n2, p2, port)
					log.Infof("... expected %v , got %v", m.Expected[n1+"_"+p1][n2+"_"+p2], connected)
					if err != nil {
						log.Errorf("unable to make main observation on %s-%s -> %s-%s: %s", n1, p1, n2, p2, err)
					}
					m.Observe(n1, p1, n2, p2, connected)
					if reachability != nil {
						reachability.Observe(utils.NewPod(n1, p1), utils.NewPod(n2, p2), connected)
					}
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

func main() {
	k8s, err := utils.NewKubernetes()
	if err != nil {
		panic(err)
	}
	k8s.CleanNetworkPolicies(namespaces)

	// testWrapperPort80(TestDefaultDeny)
	testWrapperPort80(TestPodLabelWhitelistingFromBToA)

	// testWrapperPort80(testInnerNamespaceTraffic)
	// testWrapperPort80(testEnforcePodAndNSSelector)

	// testWrapperPort80(testEnforcePodOrNSSelector)

	// testPortsPolicies(k8s)

	// stacked port policies
	// testWrapperStacked(testPortsPoliciesStackedOrUpdated, true)
	// updated port policies
	// testWrapperStacked(testPortsPoliciesStackedOrUpdated, false)

	// testWrapperPort80(testAllowAll)

	// testWrapperPort80(testNamedPort)

	// testWrapperPort80(testNamedPortWNamespace)

	// testWrapperPort80(testEgressOnNamedPort)

	//testWrapperStacked(TestAllowAllPrecedenceIngress,true )

	/**
		TestEgressAndIngressIntegration
		TestMultipleUpdates()
	**/
}

// testWrapperStaged is for tests which involve steps of mutation.
type Stack struct {
	ReachableMatrix *utils.ReachableMatrix
	Reachability    *utils.Reachability
	NetworkPolicy   *networkingv1.NetworkPolicy
	Port            int
}

// catch all for any type of test, where we use stacks.  these are validated one at a time.
// probably use this for *all* tests when we port to upstream.
func testWrapperStacked(theTest func(k8s *utils.Kubernetes, isStacked bool) (stack []*Stack), stacked bool) {
	k8s, err := utils.NewKubernetes()
	if err != nil {
		panic(err)
	}
	bootstrap(k8s)

	stack := theTest(k8s, stacked)
	for _, s := range stack {
		matrix := s.ReachableMatrix
		reachability := s.Reachability
		policy := s.NetworkPolicy
		if policy != nil {
			_, err := k8s.CreateOrUpdateNetworkPolicy(policy.Namespace, policy)
			if err != nil {
				panic(err)
			}
		}
		validate(k8s, matrix, reachability, s.Port)
		summary1, pass1 := matrix.Summary()
		fmt.Println(summary1, pass1)

		reachability.PrintSummary(true, true, true)
	}
}

// For dual port tests... confirms both ports 80 and 81
func testWrapperPort8081(theTest func(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability, *utils.ReachableMatrix, *utils.Reachability)) {
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
		reachability.PrintSummary(true, true, true)
	}
}

// simple type of test, majority of tests use this, just port 80
func testWrapperPort80(theTest func(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability)) {
	k8s, err := utils.NewKubernetes()
	if err != nil {
		panic(err)
	}
	bootstrap(k8s)
	matrix, reachability := theTest(k8s)
	validate(k8s, matrix, reachability, 80)

	summary, pass := matrix.Summary()
	fmt.Println(summary, pass)

	reachability.PrintSummary(true, true, true)
}

/**
CIDR tests.... todo
*/

/**
	ginkgo.It("should allow ingress access from updated namespace [Feature:NetworkPolicy]", func() {
	ginkgo.It("should allow ingress access from updated pod [Feature:NetworkPolicy]", func() {

	TODO: These 3 tests should be implemented using a different strategy, possibly combined.

	ginkgo.It("should deny ingress access to updated pod [Feature:NetworkPolicy]", func() {
	ginkgo.It("should stop enforcing policies after they are deleted [Feature:NetworkPolicy]", func() {
**/
func TestMultipleUpdates() {
	k8s, err := utils.NewKubernetes()
	if err != nil {
		panic(err)
	}
	bootstrap(k8s)

	func() {
		builder := &utils.NetworkPolicySpecBuilder{}
		builder = builder.SetName("x", "deny-all").SetPodSelector(map[string]string{"pod": "a"})
		builder.SetTypeIngress()
		builder.AddIngress(nil, &p80, nil, nil, nil, map[string]string{"ns-updated": "true", "ns": "y"}, nil, nil)
		builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod-updated": "true", "pod": "b"}, nil, nil, nil)

		k8s.CreateOrUpdateNetworkPolicy("deny-all-to-x", builder.Get())
		m1 := utils.NewReachableMatrix(true, pods, namespaces)
		reachability1 := utils.NewReachability(allPods)
		m1.ExpectAllIngress("x", "a", false)
		validate(k8s, m1, reachability1, 80)
		summary1, pass1 := m1.Summary()
		fmt.Println(summary1, pass1)

		reachability1.PrintSummary(true, true, true)
	}()

	func() {
		k8s.CreateOrUpdateNamespace("y", map[string]string{"ns-updated": "true", "ns": "y"})
		m1 := utils.NewReachableMatrix(true, pods, namespaces)
		reachability1 := utils.NewReachability(allPods)
		m1.ExpectAllIngress("x", "a", false)
		m1.Expect("y", "a", "x", "a", true)
		m1.Expect("y", "b", "x", "a", true)
		m1.Expect("y", "c", "x", "a", true)
		validate(k8s, m1, reachability1, 80)
		summary1, pass1 := m1.Summary()
		fmt.Println(summary1, pass1)

		reachability1.PrintSummary(true, true, true)
	}()

	func() {
		k8s.CreateOrUpdateNamespace("y", map[string]string{"ns-updated": "true", "ns": "y"})
		m1 := utils.NewReachableMatrix(true, pods, namespaces)
		reachability1 := utils.NewReachability(allPods)
		m1.ExpectAllIngress("x", "a", false)
		m1.Expect("y", "a", "x", "a", true)
		m1.Expect("y", "b", "x", "a", true)
		m1.Expect("y", "c", "x", "a", true)
		validate(k8s, m1, reachability1, 80)
		summary1, pass1 := m1.Summary()
		fmt.Println(summary1, pass1)

		reachability1.PrintSummary(true, true, true)
	}()

	func() {
		k8s.CreateOrUpdateDeployment("z", "zb", 1,
			map[string]string{
				"pod":     "b",
				"updated": "true",
			}, "nginx:1.8-alpine") // old nginx cause it was before people deleted everything useful from containers
		m1 := utils.NewReachableMatrix(true, pods, namespaces)
		// copied from above
		reachability1 := utils.NewReachability(allPods)
		m1.ExpectAllIngress("x", "a", false)
		m1.Expect("y", "a", "x", "a", true)
		m1.Expect("y", "b", "x", "a", true)
		m1.Expect("y", "c", "x", "a", true)

		// delta... pod z in b has 'updated=true' so its whitelisted.
		m1.Expect("z", "b", "x", "a", true)

		validate(k8s, m1, reachability1, 80)
		summary1, pass1 := m1.Summary()
		fmt.Println(summary1, pass1)

		reachability1.PrintSummary(true, true, true)
	}()

	// NOTE THIS TEST IS COPIED FROM THE ABOVE TEST, only delta being that we
	// dont have the udpated=true annotation above.
	func() {
		k8s.CreateOrUpdateDeployment("z", "zb", 1,
			map[string]string{
				"pod": "b",
				// REMOVE UPDATED ANNOTATION, otherwise identical to above function.
			}, "nginx:1.8-alpine") // old nginx cause it was before people deleted everything useful from containers
		m1 := utils.NewReachableMatrix(true, pods, namespaces)
		// copied from above
		reachability1 := utils.NewReachability(allPods)
		m1.ExpectAllIngress("x", "a", false)
		m1.Expect("y", "a", "x", "a", true)
		m1.Expect("y", "b", "x", "a", true)
		m1.Expect("y", "c", "x", "a", true)

		// REMOVED DELTA, otherwise identical... this confirms that access is blocked again.
		validate(k8s, m1, reachability1, 80)
		summary1, pass1 := m1.Summary()
		fmt.Println(summary1, pass1)

		reachability1.PrintSummary(true, true, true)
	}()

}

/**
ginkgo.It("should enforce multiple egress policies with egress allow-all policy taking precedence [Feature:NetworkPolicy]", func() {
ginkgo.It("should enforce policies to check ingress and egress policies can be controlled independently based on PodSelector [Feature:NetworkPolicy]", func() {
ginkgo.It("should enforce egress policy allowing traffic to a server in a different namespace based on PodSelector and NamespaceSelector [Feature:NetworkPolicy]", func() {
*/
func TestEgressAndIngressIntegration(k8s *utils.Kubernetes, stacked bool) []*Stack {
	// ingress policies stack
	builder1 := &utils.NetworkPolicySpecBuilder{}
	builder1 = builder1.SetName("x", "deny-all").SetPodSelector(map[string]string{"pod": "a"})
	builder1.SetTypeIngress()
	builder1.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, nil, nil, nil)
	policy1 := builder1.Get()
	m1 := utils.NewReachableMatrix(false, pods, namespaces)
	reachability1 := utils.NewReachability(allPods)
	m1.ExpectAllIngress("x", "a", false)
	m1.Expect("x", "b", "x", "a", true)
	m1.Expect("y", "b", "x", "a", true)
	m1.Expect("z", "b", "x", "a", true)
	m1.Expect("x", "a", "x", "a", true)

	// egress policies stack w pod selector and ns selector
	builder2 := &utils.NetworkPolicySpecBuilder{}
	builder2 = builder2.SetName("x", "deny-all").SetPodSelector(map[string]string{"pod": "a"})
	builder2.SetTypeEgress().AddEgress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, map[string]string{"ns": "y"}, nil, nil)
	policy2 := builder1.Get()
	m2 := utils.NewReachableMatrix(false, pods, namespaces)
	reachability2 := utils.NewReachability(allPods)
	// copied from m1
	m2.ExpectAllIngress("x", "a", false)
	m2.Expect("x", "b", "x", "a", true)
	m2.Expect("y", "b", "x", "a", true)
	m2.Expect("z", "b", "x", "a", true)
	m2.Expect("x", "a", "x", "a", true)

	// new egress rule.
	m2.Expect("x", "a", "y", "b", true)

	builder3 := &utils.NetworkPolicySpecBuilder{}
	// by preserving the same name, this policy will also serve to test the 'updated policy' scenario.
	builder3 = builder2.SetName("x", "allow-all").SetPodSelector(map[string]string{"pod": "a"})
	builder3.SetTypeEgress()
	builder3.AddEgress(nil, &p80, nil, nil, nil, nil, nil, nil)

	policy3 := builder2.Get()
	m3 := utils.NewReachableMatrix(true, pods, namespaces)
	reachability3 := utils.NewReachability(allPods)

	return []*Stack{
		&Stack{
			m1,
			reachability1,
			policy1,
			p80,
		},
		&Stack{
			m2,
			reachability2,
			policy2,
			p80,
		},
		&Stack{
			m3,
			reachability3,
			policy3,
			p80,
		},
	}
}

// should enforce multiple ingress policies with ingress allow-all policy taking precedence [Feature:NetworkPolicy]"
func TestAllowAllPrecedenceIngress(k8s *utils.Kubernetes, stackedOrUpdated bool) []*Stack {
	if !stackedOrUpdated {
		panic("this test always true")
	}

	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("x", "deny-all").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress()
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{}, nil, nil, nil)

	policy1 := builder.Get()
	m1 := utils.NewReachableMatrix(true, pods, namespaces)
	m1.ExpectAllIngress("x", "a", false)
	m1.Expect("x", "a", "x", "a", true)
	reachability1 := utils.NewReachability(allPods)

	builder2 := &utils.NetworkPolicySpecBuilder{}
	// by preserving the same name, this policy will also serve to test the 'updated policy' scenario.
	builder2 = builder2.SetName("x", "allow-all").SetPodSelector(map[string]string{"pod": "a"})
	builder2.SetTypeIngress()
	builder2.AddIngress(nil, &p80, nil, nil, nil, nil, nil, nil)

	policy2 := builder2.Get()
	m2 := utils.NewReachableMatrix(true, pods, namespaces)
	reachability2 := utils.NewReachability(allPods)

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
	}
}

// should allow egress access on one named port [Feature:NetworkPolicy]
func testEgressOnNamedPort(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability) {
	namedPorts := "serve-80"
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("x", "allow-client-a-via-named-port-egress-rule").SetPodSelector(map[string]string{"pod": "a"})

	// note egress DNS isnt necessary to test egress over a named port.
	builder.SetTypeEgress().WithEgressDNS().AddEgress(nil, nil, &namedPorts, nil, nil, nil, nil, nil)

	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := utils.NewReachableMatrix(true, pods, namespaces)
	reachability := utils.NewReachability(allPods)

	// TODO, maybe add validation that 81 doesn't work as well?
	return m, reachability
}

// should allow ingress access from namespace on one named port [Feature:NetworkPolicy]
func testNamedPortWNamespace(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability) {
	namedPorts := "serve-80"
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("x", "allow-client-a-via-named-port-ingress-rule").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress().AddIngress(nil, nil, &namedPorts, nil, nil, map[string]string{"ns": "x"}, nil, nil)

	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := utils.NewReachableMatrix(true, pods, namespaces)
	reachability := utils.NewReachability(allPods)
	m.ExpectAllIngress("x", "a", false)
	m.Expect("x", "a", "x", "a", true)
	m.Expect("x", "b", "x", "a", true)
	m.Expect("x", "c", "x", "a", true)

	// TODO, add validation that 81 doesn't work.
	return m, reachability
}

// testNamedPort should allow ingress access on one named port [Feature:NetworkPolicy]
func testNamedPort(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability) {
	namedPorts := "serve-80"
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("x", "allow-client-a-via-named-port-ingress-rule").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress().AddIngress(nil, &p80, &namedPorts, nil, nil, nil, nil, nil)

	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := utils.NewReachableMatrix(true, pods, namespaces)
	// No egress rules because we're deny all !
	reachability := utils.NewReachability(allPods)

	// TODO, add validation that 81 doesn't work.
	return m, reachability
}

// testAllowAll should support allow-all policy [Feature:NetworkPolicy]
func testAllowAll(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability) {
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("x", "default-deny").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, nil, nil, nil, nil)
	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := utils.NewReachableMatrix(true, pods, namespaces)
	// No egress rules because we're deny all !
	reachability := utils.NewReachability(allPods)
	return m, reachability
}

// This covers two test cases: stacked policy's and updated policies.
// 1) should enforce policy based on Ports [Feature:NetworkPolicy] (disallow 80) (stacked == false)
// 2) should enforce updated policy (stacked == true)
func testPortsPoliciesStackedOrUpdated(k8s *utils.Kubernetes, stackInsteadOfUpdate bool) []*Stack {
	blocked := func() *utils.ReachableMatrix {
		x := utils.NewReachableMatrix(true, pods, namespaces)
		x.ExpectAllIngress("x", "a", false)
		x.Expect("x", "a", "x", "a", true)
		return x
	}
	unblocked := func() *utils.ReachableMatrix {
		x := utils.NewReachableMatrix(true, pods, namespaces)
		return x
	}

	/***
	Initially, only whitelist port 80, and verify 81 is blocked.
	*/
	policyName := "policy-that-will-update-for-ports"
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("x", policyName).SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress()
	builder.AddIngress(nil, &p80, nil, nil, nil, nil, nil, nil)
	policy1 := builder.Get()

	/***
	  Now, whitelist port 81, and verify 81 it is open.
	*/
	// using false makes this a test for 'updated' policies...
	if stackInsteadOfUpdate {
		policyName = "policy-that-will-update-for-ports-2"
	}
	builder2 := &utils.NetworkPolicySpecBuilder{}
	// by preserving the same name, this policy will also serve to test the 'updated policy' scenario.
	builder2 = builder2.SetName("x", policyName).SetPodSelector(map[string]string{"pod": "a"})
	builder2.SetTypeIngress()
	builder2.AddIngress(nil, &p81, nil, nil, nil, nil, nil, nil)
	policy2 := builder2.Get()

	// The first policy was on port 80, which was whitelisted, while 81 wasn't.
	// The second policy was on port 81, which was whitelisted.
	// At this point, if we stacked, make sure 80 is still unblocked
	// Whereas if we DIDNT stack, make sure 80 is blocked.
	s3 := &Stack{
		blocked(),
		utils.NewReachability(allPods),
		nil, // nil policy wont be created, this is just a 2nd validation, this time, of port 81.
		80,
	}
	if stackInsteadOfUpdate {
		s3.ReachableMatrix = unblocked()
	}
	return []*Stack{
		&Stack{
			blocked(), // 81 blocked
			utils.NewReachability(allPods),
			policy1,
			81,
		},
		&Stack{
			unblocked(), // 81 open now
			utils.NewReachability(allPods),
			policy2,
			81,
		},
		s3,
	}
}

// "should enforce policy based on Ports [Feature:NetworkPolicy] (disallow 80)
func testPortsPolicies(k8s *utils.Kubernetes) {
	bootstrap(k8s)
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("x", "allow-port-81-not-port-80").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress()
	// anyone on port 81 is ok...
	builder.AddIngress(nil, &p81, nil, nil, nil, nil, nil, nil)

	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())

	m80 := utils.NewReachableMatrix(true, pods, namespaces)
	m80.ExpectAllIngress("x", "a", false)
	m80.Expect("x", "a", "x", "a", true)
	validate(k8s, m80, nil, 80)
	s, p := m80.Summary()
	fmt.Println(s, p)

	fmt.Println("***** port 81 *****")
	m81 := utils.NewReachableMatrix(true, pods, namespaces)
	m81.ExpectAllIngress("x", "a", true)
	validate(k8s, m81, nil, 81)
	s, p = m81.Summary()
	fmt.Println(s, p)

}

// should enforce policy to allow traffic only from a pod in a different namespace based on PodSelector and NamespaceSelector [Feature:NetworkPolicy]
// should enforce policy based on PodSelector and NamespaceSelector [Feature:NetworkPolicy]
func testEnforcePodAndNSSelector(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability) {
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("x", "allow-x-via-pod-and-ns-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress()
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, map[string]string{"ns": "y"}, nil, nil)

	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := utils.NewReachableMatrix(true, pods, namespaces)
	reachability := utils.NewReachability(allPods)
	m.ExpectAllIngress("x", "a", false)
	m.Expect("y", "b", "x", "a", true)
	m.Expect("x", "a", "x", "a", true)

	return m, reachability
}

// should enforce policy based on PodSelector or NamespaceSelector [Feature:NetworkPolicy]
func testEnforcePodOrNSSelector(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability) {
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("x", "allow-x-via-pod-or-ns-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress()
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, nil, nil, nil)
	builder.AddIngress(nil, &p80, nil, nil, nil, map[string]string{"ns": "y"}, nil, nil)

	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := utils.NewReachableMatrix(true, pods, namespaces)
	reachability := utils.NewReachability(allPods)
	m.ExpectAllIngress("x", "a", false)
	m.Expect("y", "a", "x", "a", true)
	m.Expect("y", "b", "x", "a", true)
	m.Expect("y", "c", "x", "a", true)
	m.Expect("x", "b", "x", "a", true)
	m.Expect("y", "b", "x", "a", true)
	//m.Expect("z", "b", "x", "a", true)
	m.Expect("x", "a", "x", "a", true)

	return m, reachability
}

// should enforce policy based on NamespaceSelector with MatchExpressions[Feature:NetworkPolicy]
func testNamespaceSelectorMatchExpressions(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability) {
	builder := &utils.NetworkPolicySpecBuilder{}
	selector := []metav1.LabelSelectorRequirement{{
		Key:      "ns",
		Operator: metav1.LabelSelectorOpIn,
		Values:   []string{"y"},
	}}
	builder = builder.SetName("x", "allow-a-via-ns-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, nil, nil, &selector, nil)
	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := utils.NewReachableMatrix(true, pods, namespaces)
	reachability := utils.NewReachability(allPods)
	m.ExpectAllIngress("x", "a", false)
	m.Expect("y", "a", "x", "a", true)
	m.Expect("y", "b", "x", "a", true)
	m.Expect("y", "c", "x", "a", true)
	m.Expect("x", "a", "x", "a", true)

	return m, reachability
}

// testPodSelectorMatchExpressions should enforce policy based on PodSelector with MatchExpressions[Feature:NetworkPolicy]
func testPodSelectorMatchExpressions(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability) {
	builder := &utils.NetworkPolicySpecBuilder{}
	selector := []metav1.LabelSelectorRequirement{{
		Key:      "pod",
		Operator: metav1.LabelSelectorOpIn,
		Values:   []string{"b"},
	}}
	builder = builder.SetName("x", "allow-client-b-via-pod-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, nil, nil, &selector, nil)
	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := utils.NewReachableMatrix(true, pods, namespaces)
	reachability := utils.NewReachability(allPods)
	m.ExpectAllIngress("x", "a", false)
	m.Expect("x", "b", "x", "a", true)
	m.Expect("y", "b", "x", "a", true)
	m.Expect("z", "b", "x", "a", true)
	m.Expect("x", "a", "x", "a", true)

	return m, reachability
}

// testInnerNamespaceTraffic should enforce policy to allow traffic from pods within server namespace based on PodSelector [Feature:NetworkPolicy]
func testIntraNamespaceTrafficOnly(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability) {
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("x", "allow-client-b-via-pod-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, nil, map[string]string{"ns": "y"}, nil, nil)
	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := utils.NewReachableMatrix(true, pods, namespaces)
	reachability := utils.NewReachability(allPods)
	m.ExpectAllIngress("x", "a", false)
	m.Expect("y", "a", "x", "a", true)
	m.Expect("y", "b", "x", "a", true)
	m.Expect("y", "c", "x", "a", true)

	return m, reachability
}

// testInnerNamespaceTraffic should enforce policy to allow traffic from pods within server namespace, based on PodSelector [Feature:NetworkPolicy]
// note : network policies are applied to a namespace by default, meaning that you need a specific policy to select pods in external namespaces.
// thus in this case, we don't expect y/b -> x/a, because even though it is labelled 'b', it is in a different namespace.
func testInnerNamespaceTraffic(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability) {
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("x", "allow-client-b-via-pod-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, nil, nil, nil)
	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := utils.NewReachableMatrix(true, pods, namespaces)
	reachability := utils.NewReachability(allPods)
	m.ExpectAllIngress("x", "a", false)
	m.Expect("x", "b", "x", "a", true)
	m.Expect("x", "a", "x", "a", true)

	return m, reachability
}

func TestDefaultDeny(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability) {
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("x", "default-deny")
	builder.SetTypeIngress() //	.AddIngress(nil, &p80, nil, nil, nil, nil, nil, nil)
	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())
	m := utils.NewReachableMatrix(true, pods, namespaces)
	m.ExpectAllIngress("x", "a", false)
	m.ExpectAllIngress("x", "b", false)
	m.ExpectAllIngress("x", "c", false)

	m.Expect("x", "a", "x", "a", true)
	m.Expect("x", "b", "x", "b", true)
	m.Expect("x", "c", "x", "c", true)

	// No egress rules because we're deny all !
	reachability := utils.NewReachability(allPods)

	return m, reachability
}

func TestPodLabelWhitelistingFromBToA(k8s *utils.Kubernetes) (*utils.ReachableMatrix, *utils.Reachability) {
	builder := &utils.NetworkPolicySpecBuilder{}
	builder = builder.SetName("x", "allow-client-a-via-pod-selector").SetPodSelector(map[string]string{"pod": "a"})
	builder.SetTypeIngress()
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, map[string]string{"ns": "x"}, nil, nil)
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, map[string]string{"ns": "y"}, nil, nil)
	builder.AddIngress(nil, &p80, nil, nil, map[string]string{"pod": "b"}, map[string]string{"ns": "z"}, nil, nil)

	k8s.CreateOrUpdateNetworkPolicy("x", builder.Get())

	m := utils.NewReachableMatrix(true, pods, namespaces)
	m.ExpectAllIngress("x", "a", false)
	m.Expect("x", "b", "x", "a", true)
	m.Expect("y", "b", "x", "a", true)
	m.Expect("z", "b", "x", "a", true)
	m.Expect("x", "a", "x", "a", true)

	reachability := utils.NewReachability(allPods)
	reachability.ExpectAllIngress(utils.NewPod("x", "a"), false)
	reachability.Expect(utils.NewPod("x", "b"), utils.NewPod("x", "a"), true)
	reachability.Expect(utils.NewPod("y", "b"), utils.NewPod("x", "a"), true)
	reachability.Expect(utils.NewPod("z", "b"), utils.NewPod("x", "a"), true)
	reachability.Expect(utils.NewPod("x", "a"), utils.NewPod("x", "a"), true)

	// TODO move this to a unit test !
	if m.Expected["z_c"]["x_a"] {
		panic("expectations are wrong")
	}
	return m, reachability
}
