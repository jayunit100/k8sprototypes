package tanzuclimodel

type Command interface {
	String(c *Cluster) []string
	Configure(c *Cluster, inputs map[string]string) bool
}
