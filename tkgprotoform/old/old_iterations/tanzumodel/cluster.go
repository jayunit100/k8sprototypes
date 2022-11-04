package tanzuclimodel

type Cluster struct {
	// something like "type=kubevip, vip=1.2.3.4"... can make a struct later
	Parent              *Cluster
	LoadbalancerConfigs map[string]string
	Name                string
	Packages            []*Package
	CreateCommand       Command
	DeleteCommand       Command
}

func (c *Cluster) IsWorkload() bool {
	return c.Parent != nil
}
func (c *Cluster) IsManagement() bool {
	return !c.IsWorkload()
}

// Command: Create

type Create struct{}
type Delete struct{}

func (p *Create) String(c *Cluster) []string {
	return []string{"tanzu", "cluster", "create", "-f", c.Name}
}
func (p *Create) Configure(c *Cluster, inputs map[string]string) error {
	return nil
}

func (p *Delete) String(c *Cluster) []string {
	return []string{"tanzu", "cluster", "delete", c.Name, "-y"}
}

func (p *Delete) Configure(c *Cluster, inputs map[string]string) error {
	return nil
}
