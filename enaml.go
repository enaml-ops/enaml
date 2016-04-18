package enaml

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v1"
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
		deploymentManifest = deployment.VSphere()
	}

	if dmYaml, err := yaml.Marshal(deploymentManifest); err != nil {
		panic(fmt.Sprintf("couldnt parse deployment manifest: ", err))
	} else {
		fmt.Println(string(dmYaml))
	}
}
