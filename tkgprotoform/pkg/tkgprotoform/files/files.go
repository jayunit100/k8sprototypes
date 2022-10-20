package files

import (
	"k8s.io/klog/v2"
	"strings"
	config "tkgprotoform.com/protoform/pkg/config"
	"tkgprotoform.com/protoform/pkg/util"
)

// PURPOSE: To write yaml and shell scripts that are custom installation artifacts to disk
// so that we can replicate the work our CI is doing locally, easily, without parsing python in our heads, or
// copying auto generated log dumps into shell scripts.
// INPUT: compile time, scans all yaml and sh files
// OUTPUT: writes the files to cwd, when you call WriteAllToLocal

import (
	_ "embed"
)

var (
	//go:embed cluster.yaml
	ClusterConfig string

	//go:embed management-cluster.sh
	ManagementClusterInstall string

	//go:embed
	ImageBuilderScript string
)

// files has all the files we are going to write to disk, so that
// they can be hacked up by the user after running tkgprotoform init the first time.
func files() map[string]string {
	return map[string]string{
		"cluster.yaml":          ClusterConfig,
		"management_cluster.sh": ManagementClusterInstall,
		"image_builder.sh":      ImageBuilderScript,
	}
}

// WriteAllToLocal writes the shell scripts and yaml files
// to your local directory.  It takes a config.yaml as input... that way
// if you have a specific home directory or whatever, it eventually
// can infer that.
func WriteAllToLocal(conf *config.Config) {
	klog.Infof("Writing out %v static files to local directory.", len(files()))
	for file, contents := range files() {
		if util.FileExists(file) {
			klog.Infof("File exists %v , not writing...", file)
		} else {
			klog.Infof("File not exists %v, writing...", file)

			if file == ImageBuilderScript {
				contents = ImageBuilderSubstitutions(conf, contents)
			}
			util.StringToFile(contents, conf.OutputFilesPath+"/"+file)
		}
	}
}

func ImageBuilderSubstitutions(conf *config.Config, contents string) string {
	contents = strings.ReplaceAll(contents, "blah_iso_xyz", conf.ImageBuilderInputs["iso_path"])
	return contents
}
