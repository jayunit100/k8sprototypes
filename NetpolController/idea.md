This markdown file describes an idea I've had recently to declaratively define policies, and then generate network policy implementations on the fly.


How current network policy works (requires understanding data model to use it): 
![alt text](https://github.com/jayunit100/k8sprototypes/raw/master/NetpolController/Kubernetes:%20An%20Illustrated%20API.png "v1 net pol diagram")


IMO, the existing NetworkPolicy API in Kubernetes is very hard to understand.  The reason is b/c of the corner cases around nil, empty, missing policies, and what 'deny all' really means.  For example, a deny all policy is currently a policy that selects a pod, but has an emptyIngress rule on it.  If that empty ingress rule then has one `{}` entry in it, the policy now becomes 'allow all' since an empty ingress from rule is effectively a match to any pod.

These subtle changes resulting in wildly different semantics make network policies counterintuitive to use and communicate across security channels without a higher level construct.

This MD file explains one approach to such a construct: A NetworkPolicy Controller, which uses 

- PolicyGroups
- DAG's with +/- polarity (allow/deny) between those groups
- Default allow / deny rules for pods *not* in a DAG

Using these 3 concepts policy's can be written the same way they are designed intuitively: By a DAG over a graph of edges and nodes, where multiple edges can exist between nodes (pods).

# What a network policy API should express

- Declare groups of assets (IP, pod, ns) which we want to protect 
- Declare connectivity between those assets 
- Declare default policies for new assets which have low overall priority
- Provide definitive (boolean) information about wether two specific containers can communicate

Note that none of these IMO should require a CNI provider.

# What users have to do due to our lack of declarative policies in K8's networkPolicy API

- Calculate effective policy from layers
- Understand datastructure properties (nil vs empty) how those amount to deny/allow
- Use CNI provider logs or metrics to determine ultimate connectivity decisions
- Manually inspect/verify that selectors match assets (pods) properly

# What users should be able to do

A policy should be built in a sequence of steps which are independently verifiable:

- Define groups of assets using either regex or explicitly, confirm the groups are correctly matching assets.
- Define default policies on those assets, either in terms of allow or deny, depending on perspective.
- Define specific connectivity between asset groups, with assurance the connectivity isnt conflicting.
- Layer policies in a pinch i.e. 'deny all'  or 'allow all' based on perspective.

```
type PolicyGroupSpec {
	Name string
	includes *PodSelector
	exceptions *PodSelector	
}
type PolicyGroupStatus {
	Pods []api.Pod
}
type PolicyGroup struct {
	Spec PolicyGroupSpec
	Status PolicyGroupStatus // updated by controller
}

type PodSelector {
	Spec *PodSelectorSpec
	Status *PodSelectorStatus
}
type PodSelectorSpec {
	Namespaces []string
	Pods []string
	NSLabels map[string]string
	PodLables map[string]string
	IPs []string
	images []string
}
type PodSelectorStatus {
	Pods []api.Pod
}
type PolicyEdge {
	// edges can allow, or deny.
	// The policy controller makes
	// sure there is no trampling, i.e., 
	// that there is only one edge between
	// 2 Policy groups
	Priority int
	Polarity bool // true = positive
	PortRange []int

	PolicyGroup incoming
	incomingRegex string

	PolicyGroup outgoing
	outgoingRegex string
}

type PodCommMap {
	// adjacency list of all pod communications
	// ultimately this is the only thing CNI's
	// really need in order to make a decision
	map[string]map[string]bool
}

/** Example usages

	xNs := PodSelector{
		Spec: {
			Namespace: []string{"x"}
		}
	}

	pgX := PolicyGroup{
		Spec: {
			Name: "x",
			includes: xNs,
			excludes: nil
		}		
	}
	pgX := PolicyGroup{
		Spec: {
			Name: "allow-all-x",
			includes: xNs,
			excludes: nil
		}		
	}


	pe := PolicyEdge {
		Priority: 1
		Allow: true
		PortRange: []int{0,10000}
		incoming: nil
		incomingRegex: string[]{"x"} // <- ns x
	
		outgoing: PolicyGroup[]
		outgoingRegex: "*"
	}
**/

```
