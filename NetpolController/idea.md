# What a network policy API should express

- Define groups of assets which we want to protect
- Define connectivity between those assets
- Define default policies for new assets which don't have groups
- Provide definitive (boolean) information about wether two assets can communicate

# What we express now

- Ways to select individual assets and block communication to them
- Ways to layer allowed communication on to blocked assets

# What users have to do

- Layer multiple allowed policies in order to determine security posture
- Understand about nil vs empty policies and how that amounts to deny/allow 
- Use CNI provider to determine ultimate connectivity decisions
- Manually inspect/verify that selectors match assets (pods) properly

# What users should be able to do

A policy should be built in 3 steps

- Define policies either in terms of allow or deny, depending on perspective.
- Define proper defaults explicitly (Allow or deny all).
- Define groups using either regex or explicitly, confirm the groups are correctly matching assets.
- Define connectivity across groups, with assurance the connectivity isnt conflicting.
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
