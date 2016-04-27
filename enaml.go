package enaml

import "gopkg.in/yaml.v2"

//Paint -
func Paint(deployment Deployment) (result string, err error) {
	var deploymentManifest DeploymentManifest
	deploymentManifest = deployment.GetDeployment()
	var dmYaml []byte
	if dmYaml, err = yaml.Marshal(deploymentManifest); err == nil {
		result = string(dmYaml)
	}
	return
}
