package printer

import "io/ioutil"
import "text/template"
import "bytes"
import "components"
import "components/web/section/generic"
import "reflect"

const VIEW_PATH = "templates/"

type Printer interface {
	PrintContent(page *generic.MainSection) (string, error)
}

type PrinterPage struct {
	Printer
}

func (printerPage *PrinterPage) PrintContent(component components.Component) (string, error) {
	typeComponent := reflect.TypeOf(component)
	valueComponent := reflect.ValueOf(component)

	dataTemplate := make(map[string]string)

	componentContent := ""

	if typeComponent.Kind() == reflect.Struct {
		for index := 0; index < typeComponent.NumField(); index++ {
			field := typeComponent.Field(index)
			if !field.Anonymous {

				fieldValue := reflect.Indirect(valueComponent).FieldByName(field.Name)
				fieldInterface := fieldValue.Interface()

				if fieldValue.Kind() == reflect.Slice {

					var buffer bytes.Buffer
					s := reflect.ValueOf(fieldInterface)

					for i := 0; i < s.Len(); i++ {
						sliceValue := s.Index(i).Interface()
						if component, ok := sliceValue.(components.Component); ok {
							componentContentField, _ := printerPage.PrintContent(component)
							buffer.WriteString(componentContentField)
						}
					}

					dataTemplate[field.Name] = buffer.String()

				} else {
					if component, ok := fieldInterface.(components.Component); ok {

						componentContentField, _ := printerPage.PrintContent(component)
						dataTemplate[field.Name] = componentContentField
					} else {
						dataTemplate[field.Name] = fieldValue.String()
					}
				}
			}
		}
		componentContent, _ = getContent(component.GetTemplateName(), dataTemplate)
	}
	return componentContent, nil
}

func getContent(templateName string, component interface{}) (string, error) {

	filepath := getFilePath(templateName)
	content, error := ioutil.ReadFile(filepath)
	if error!= nil {
		return "", error
	}

	t := template.New("template") //create a new template with some name
	t, _ = t.Parse(string(content))

	var doc bytes.Buffer

	error = t.Execute(&doc, component)

	if error != nil {
		return "", error
	}

	return doc.String(), nil
}

func getFilePath(filename string) string {
	return VIEW_PATH + filename
}
