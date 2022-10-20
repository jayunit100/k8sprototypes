package tanzuclimodel

import (
	"fmt"
	"strings"
)

func GetPrometheusPackage() *Package {
	return &Package{
		Name:      "prometheus",
		Package:   "prometheus.tanzu.vmware.com",
		Version:   "2.36.2+vmware.1-tkg.1",
		Namespace: "default",
	}
}

type Package struct {
	Name      string
	Package   string
	Version   string
	Namespace string
	Install   Command
}

type PackageInstallCommand struct {
	Package
}

// Commands

func (p *PackageInstallCommand) String(c *Cluster) []string {
	tmpl := "tanzu package install %v -- package-name %v --version %v --namespace %v"
	cmd := fmt.Sprintf(tmpl, c.Name, p.Package.Name, p.Package.Version, p.Package.Namespace)
	return strings.SplitAfter(" ", cmd)
}
func (p *PackageInstallCommand) Configure(inputs map[string]string) error {
	return nil
}
