package tanzuclimodel

//  tkgprotoform.com/tkgprotoform/tanzumodel

// VC7 returns two clusters, one mgmt, one worker
func GetVC7(name string, ip1 string, ip2 string) (*Cluster, *Cluster) {
	mgmt := NewVSphereManagementCluster(name+"-mgmt", ip1)
	worker := NewVsphereWorkloadCluster(name+"-worker", ip2)
	return mgmt, worker
}

// NewVsphereCluster creates a new vsphere cluster.  IF you send a Vip, its
// kube-vip based....
func NewVSphereCluster(name string, vip string) *Cluster {
	return &Cluster{
		Parent: nil,
		LoadbalancerConfigs: map[string]string{
			"kube-vip": vip,
		},
		Name:     name,
		Packages: []*Package{},
	}
}

func NewVSphereManagementCluster(name string, vip string) *Cluster {
	c := NewVSphereCluster(name, vip)
	return c
}

func NewVsphereWorkloadCluster(name string, vip string) *Cluster {
	c := NewVSphereCluster(name, vip)
	c.Packages = []*Package{
		GetPrometheusPackage(),
	}
	return c
}
