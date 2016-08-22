package generators

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/enaml-ops/enaml"
	"github.com/xchapter7x/lo"
	"gopkg.in/yaml.v2"
)

//Generate - used to generate a struct for a given job
func Generate(packagename string, fileBytes []byte, outputDir string) {
	b := preprocessJobManifest(fileBytes)
	objects := make(map[string]map[string]ObjectField)
	var properties []string
	for _, v := range b.recs {
		properties = append(properties, v.Orig)
	}
	for i := 0; i < b.max; i++ {

		for _, v := range b.recs {
			if v.Length-1 >= i {

				var structname = v.StructName(i, packagename, properties)
				var typeName = v.TypeName(i, properties)
				elementName := v.Slice[i]
				attributeName := FormatName(elementName)

				if _, ok := objects[structname]; !ok {
					objects[structname] = make(map[string]ObjectField)
				}

				if previousElement, ok := objects[structname][attributeName]; !ok {
					lo.G.Debug("Adding", attributeName, "to", structname, "with type", typeName)
					objects[structname][attributeName] = ObjectField{
						ElementName:       attributeName,
						ElementType:       typeName,
						ElementAnnotation: createElementAnnotation(elementName),
						Meta:              v.Yaml,
					}
				} else {
					if previousElement.ElementAnnotation != createElementAnnotation(elementName) {
						lo.G.Warning("******** Recommended creating custom yaml marshaller on", structname, "for", attributeName, " ********")
						previousElement.ElementAnnotation = "`yaml:\"-\"`"
						objects[structname][attributeName] = previousElement
					}
				}
			}
		}
	}
	structs := generateStructs(objects, packagename)
	writeStructsToDisk(structs, outputDir)
}

func createElementAnnotation(elementName string) string {
	return fmt.Sprintf("`yaml:\"%s,omitempty\"`", elementName)
}

func writeStructsToDisk(structs []jobStructTemplate, outputDir string) {
	for _, v := range structs {
		tmpl, err := template.New("job").Parse(structTemplate)
		if err != nil {
			panic(err)
		}

		os.MkdirAll(outputDir, 0700)
		jobPath := path.Join(outputDir, strings.ToLower(FormatName(v.JobName))+".go")
		f, _ := os.Create(jobPath)
		err = tmpl.Execute(f, v)
		if err != nil {
			panic(err)
		}
	}
}

func generateStructs(objects map[string]map[string]ObjectField, packagename string) (structList []jobStructTemplate) {

	for k, v := range objects {
		tmpJob := jobStructTemplate{
			PackageName: packagename,
			JobName:     k,
		}
		for _, v1 := range v {
			tmpJob.Elements = append(tmpJob.Elements, elementStruct{
				ElementName:     v1.ElementName,
				ElementType:     v1.ElementType,
				ElementYamlName: v1.ElementAnnotation,
				ElementComments: fmt.Sprintf("Descr: %v Default: %v\n", v1.Meta.Description, v1.Meta.Default),
			})
		}
		structList = append(structList, tmpJob)
	}
	return
}

func preprocessJobManifest(jobmanifest []byte) (proc processing) {
	manifestYaml := enaml.JobManifest{}
	yaml.Unmarshal(jobmanifest, &manifestYaml)

	for k, v := range manifestYaml.Properties {
		rec := CreateNewRecord(k, v)
		proc.recs = append(proc.recs, rec)
		if proc.max < rec.Length {
			proc.max = rec.Length
		}
	}
	return
}

//CreateNewRecord - creates a record from a given period delimited property and enaml.JobManifestProperty
func CreateNewRecord(property string, yaml enaml.JobManifestProperty) (record Record) {
	elementArray := strings.Split(property, ".")
	record = Record{
		Length: len(elementArray),
		Orig:   property,
		Slice:  elementArray,
		Yaml:   yaml,
	}
	return
}
