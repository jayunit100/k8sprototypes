# What a network policy API should express

- Define policy scopes (ns or pod)
- Default deny all traffic for entities in a scope
- Define policies that whitelist traffic between entities

type PolicyGroup {
	Spec {
		Name string
		
		// many ways to define a policy group...
		// redundant is ok, will be normalized
		// via status
		Namespaces []string
		Pods []string
		NSLabels map[string]string
		PodLables map[string]string
		PodIPs []string
	}

	// Managed by controller
	Status {
		Items []api.Pod
	}
}

type PolicyEdge {
	// edges can allow, or deny.
	// The policy controller makes
	// sure there is no trampling, i.e., 
	// that there is only one edge between
	// 2 Policy groups
	Allow boolean
	Policy incoming
	Policy outgoing
}

// Example usages
