package enaml

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-yaml/yaml"
)

func Paint(deployment Deployment) {
	var deploymentManifest DeploymentManifest
	switch strings.ToLower(os.Args[1]) {
	case "vsphere":
		deploymentManifest = deployment.VSphere()
	case "aws":
		deploymentManifest = deployment.AWS()
	case "azure":
		deploymentManifest = deployment.Azure()
	case "openstack":
		deploymentManifest = deployment.OpenStack()
	default:
		panic("please define a target iaas (vsphere, aws, openstack, azure)")
	}

	if dmYaml, err := yaml.Marshal(deploymentManifest); err != nil {
		panic(fmt.Sprintf("couldnt parse deployment manifest: ", err))
	} else {
		fmt.Println(string(dmYaml))
	}
}
