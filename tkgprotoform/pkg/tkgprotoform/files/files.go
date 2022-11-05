package files

import (
	"fmt"
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
	//go:embed 1.6/cluster.yaml
	ClusterConfig string

	//go:embed 1.6/2_management_cluster.sh
	ManagementClusterInstall string

	// //go:embed image_builder.sh
	ImageBuilderScript string
)

// files has all the files we are going to write to disk, so that
// they can be hacked up by the user after running tkgprotoform init the first time.
func files(version string) map[string]string {
	return map[string]string{
		"cluster.yaml":            ClusterConfig,
		"2_management_cluster.sh": ManagementClusterInstall,
		"image_builder.sh":        ImageBuilderScript,
	}
}

// WriteAllToLocal writes the shell scripts and yaml files
// to your local directory.  It takes a config.yaml as input.
func WriteAllToLocal(conf *config.Config, version string) []string {
	klog.Infof("Writing out %v static files to local directory.", len(files(version)))
	returnVal := []string{}

	for file, contents := range files(version) {
		outputFileLoc := func() string {
			return conf.OutputFilesPath + "/" + file
		}
		fmt.Println(outputFileLoc())
		if util.FileExists(outputFileLoc()) {
			klog.Infof("File exists %v , not writing...", file)
			returnVal = append(returnVal, fmt.Sprintf("skip %v", file))
		} else {
			klog.Infof("File not exists %v, writing...", file)

			if file == ImageBuilderScript {
				contents = GetImageBuilderSubstituted(conf, contents)
			}
			output := util.WriteStringToFile(contents, outputFileLoc())
			if output != nil {
				fmt.Println("ERRORrrr ", output)
				returnVal = append(returnVal, "ERROR")
			} else {
				returnVal = append(returnVal, "success_"+outputFileLoc())
			}
		}
	}
	return returnVal
}

func GetImageBuilderSubstituted(conf *config.Config, contents string) string {
	contents = strings.ReplaceAll(contents, "blah_iso_xyz", conf.ImageBuilderInputs["iso_path"])
	return contents
}
