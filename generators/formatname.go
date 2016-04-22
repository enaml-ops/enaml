package generators

import "strings"

func FormatName(unformattedName string) (formattedName string) {
	formattedName = ConvertToCamelCase(unformattedName)
	return
}

func ConvertToCamelCase(name string) string {
	f := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == '-'
	})
	for i := range f {
		f[i] = strings.ToUpper(f[i][:1]) + f[i][1:]
	}
	return strings.Join(f, "")
}
