package encoder

import (
  "reflect"
  "fmt"
  "strings"
)

func ConfigEncoder(title string, config interface{}) string {
	// This will read the fields of the interface and create a config file

	t := reflect.TypeOf(config)
	values := reflect.ValueOf(config)

	maxIndex := t.NumField()
	output := fmt.Sprintf("[%s]\n", title)

	for i := 0; i < maxIndex; i++ {
		f := t.Field(i)

		value := values.FieldByName(f.Name).String()
		name := ToSnakeCase(f.Name)

		if len(value) != 0 {
			output += fmt.Sprintf("aws%s = %s\n", name, value)
		}
	}

	return output
}

func ToSnakeCase(input string) string {
	output := ""
	for _, v := range input {
		val := string(v)
		lowerV := strings.ToLower(val)
		if lowerV != val {
			output += "_"
		}
		output += lowerV
	}

	return output
}

