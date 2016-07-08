package generators

import "strings"

//StructName - gets struct name for property
func (v *Record) StructName(i int, packagename string, properties []string) (structname string) {
	if i > 0 {
		currentNode := v.Slice[i-1]
		structname = FormatName(currentNode)
		if i > 1 {
			parentNames := v.FindAllParentsOfSameNamedElement(currentNode, properties)
			if len(parentNames) > 1 {
				structname = FormatName(v.Slice[i-2] + "_" + currentNode)
			}
		}
	} else {
		structname = FormatName(packagename + "_job")
	}
	return
}

//TypeName - returns valid type name for a give record
func (v *Record) TypeName(i int, properties []string) (typename string) {
	if i+1 < v.Length {
		currentNode := v.Slice[i]
		typename = "*" + FormatName(currentNode)
		if i >= 1 {
			parentNames := v.FindAllParentsOfSameNamedElement(currentNode, properties)
			if len(parentNames) > 1 {
				typename = "*" + FormatName(v.Slice[i-1]+"_"+currentNode)
			}
		}
	} else {
		typename = "interface{}"
	}
	return
}

func (v *Record) FindAllParentsOfSameNamedElement(currentNode string, properties []string) (parentNames []string) {
	parentNameMap := make(map[string]string)
	for _, property := range properties {
		if strings.Contains(property, currentNode) {
			parts := strings.Split(property, ".")
			for i := 0; i < len(parts); i++ {
				if currentNode == parts[i] {
					x := i - 1
					if x >= 0 {
						parentNameMap[parts[i-1]] = parts[i-1]
					} else {
						parentNameMap[""] = ""
					}
					break
				}
			}
		}
	}

	for parent, _ := range parentNameMap {
		parentNames = append(parentNames, parent)
	}
	return
}
