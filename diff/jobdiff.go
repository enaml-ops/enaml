package diff

import (
	"github.com/kr/pretty"
	"github.com/xchapter7x/enaml"
	"gopkg.in/yaml.v2"
)

func JobPropertiesDiff(a, b []byte) []string {
	var objA enaml.JobManifest
	var objB enaml.JobManifest
	yaml.Unmarshal(a, &objA)
	yaml.Unmarshal(b, &objB)
	mp := pretty.Diff(objA, objB)
	return mp
}
