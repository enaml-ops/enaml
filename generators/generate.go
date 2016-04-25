package generators

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/xchapter7x/enaml"
	"gopkg.in/yaml.v2"
)

func generate(packagename string, fileBytes []byte, outputDir string) {
	b := preprocessJobManifest(fileBytes)
	objects := make(map[string]map[string]ObjectField)

	for i := 0; i < b.max; i++ {

		for _, v := range b.recs {

			if v.Length-1 >= i {
				var structname = newStructName(i, v, packagename)
				var typeName = newTypeName(i, v)

				if _, ok := objects[structname]; !ok {
					objects[structname] = make(map[string]ObjectField)
				}

				if !nestedTypeDefinition(structname, typeName, FormatName(v.Slice[i])) {

					objects[structname][v.Slice[i]] = ObjectField{
						ElementName:       FormatName(v.Slice[i]),
						ElementType:       typeName,
						ElementAnnotation: "`yaml:\"" + v.Slice[i] + ",omitempty\"`",
						Meta:              v.Yaml,
					}
				}
			}
		}
	}
	structs := generateStructs(objects, packagename)
	writeStructsToDisk(structs, outputDir)
}

func nestedTypeDefinition(structname, typeName, elementName string) bool {
	return (structname == typeName && typeName == elementName)
}

func newStructName(i int, v record, packagename string) (structname string) {
	if i > 0 {
		structname = FormatName(strings.Join(v.Slice[i-1:i], ""))

	} else {
		structname = FormatName(packagename)
	}
	return
}

func newTypeName(i int, v record) (typename string) {
	if i+1 < v.Length {
		typename = FormatName(v.Slice[i])
	} else {
		typename = "interface{}"
	}
	return
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
		rec := record{
			Length: len(strings.Split(k, ".")),
			Orig:   k,
			Slice:  strings.Split(k, "."),
			Yaml:   v,
		}
		proc.recs = append(proc.recs, rec)
		if proc.max < rec.Length {
			proc.max = rec.Length
		}
	}
	return
}
