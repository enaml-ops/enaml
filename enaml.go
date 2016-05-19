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

//Deploy -
func Deploy(deployment Deployment) (result string, err error) {
	return Paint(deployment)
}

//Cloud -
func Cloud(config CloudConfig) (result string, err error) {
	var cloudConfigManifest CloudConfigManifest
	cloudConfigManifest = config.GetManifest()
	var ccfgYaml []byte
	if ccfgYaml, err = yaml.Marshal(cloudConfigManifest); err == nil {
		result = string(ccfgYaml)
	}
	return
}
