package enaml

import (
	"fmt"

	"gopkg.in/yaml.v1"
)

func Paint(deployment Deployment) {
	var deploymentManifest DeploymentManifest
	deploymentManifest = deployment.GetDeployment()

	if dmYaml, err := yaml.Marshal(deploymentManifest); err != nil {
		panic(fmt.Sprintf("couldnt parse deployment manifest: ", err))
	} else {
		fmt.Println(string(dmYaml))
	}
}
