package element

type OptionElement struct {
	OptionValue string
	OptionName string
}

func (page OptionElement)GetTemplateName() string {
	return "web/element/optionTemplate.tmpl"
}
