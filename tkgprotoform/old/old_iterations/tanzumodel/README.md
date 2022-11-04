# Tanzu Model classes

This models the **tanzu** user experience as a series of classes.

Each "command" in tanzu cli maps to datastructures in this folder,
for example...

```
type Cluster struct {
	// something like "type=kubevip, vip=1.2.3.4"... can make a struct later
	Parent              *Cluster
	LoadbalancerConfigs map[string]string
	Name                string
	Packages            []*Package
	CreateCommand       Command
	DeleteCommand       Command
}
```

The Cluster struct models what a user thinks about when making a cluster.  

All of the components (Clusters, Packages, and so on) that a user interacts with
implement a common interface, mostly for code uniformity:

```
type Command interface {
	String(c *Cluster) []string
	Configure(c *Cluster, inputs map[string]string) bool
}
```

The purpose of this interface is to enable

- simulating and generation of tanzu commands
- installing tanzu components and initializing them by other packages.