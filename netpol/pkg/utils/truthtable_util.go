package utils

type ReachableMatrix struct {
	DefaultExpect bool
	Pods []string
	Namespaces []string
	Expected map[string]map[string]int
	Observed map[string]map[string]int
}

func addToMap(m map[string]map[string]int, ns string, pod string, ns2 string, pod2 string, connectionTime int) {
	if m == nil {
		m = map[string]map[string]int{}
	}
	if m[ns+"_"+pod] == nil {
		m[ns+"_"+pod] = map[string]int{}
	}
	m[ns+"_"+pod][ns2+"_"+pod2] = connectionTime
}

func (r *ReachableMatrix) init(){
	if r.Expected != nil {
		return
	}

	// create datastructures if not created yet.
	r.Expected = map[string]map[string]int{}
	if r.pods == nil || r.ns == nil {
		panic("pods/ns must be existent in the struct.")
	}
	time := -1
	if r.DefaultExpect == true {
		time = 5 // expect connect within 5 seconds
	}
	// create the matrix
	for _, n := range r.ns {
		for _, p := range r.pods {
			for _, nn := range r.ns {
				for _, pp := range r.pods {
					addToMap(r.Expected, n, p, nn, pp, time)
				}
			}
		}
	}
}

func (r *ReachableMatrix) Expect(ns string, pod string, ns2 string, pod2 string, connectionTime int) {
	r.init()

	addToMap(r.Expected, ns, pod, ns2, pod2, connectionTime)
}

func (r *ReachableMatrix) Observe(ns string, pod string, ns2 string, pod2 string, connectionTime int) {
	r.init()

	addToMap(r.Observed, ns, pod, ns2, pod2, connectionTime)
}

func (r *ReachableMatrix) GetExpectedObserverd(ns string, pod string, ns2 string, pod2 string) (int, int) {
	r.init()
	return r.Expected[ns+"_"+pod][ns2+"_"+pod2], r.Observed[ns+"_"+pod][ns2+"_"+pod2]
}

func (r *ReachableMatrix) LengthExpectedObserved() (int, int) {
	r.init()
	return len(r.Expected), len(r.Observed)
}
