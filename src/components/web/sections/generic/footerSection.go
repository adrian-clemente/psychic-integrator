package generic

type FooterSection struct {}

func (page FooterSection)GetTemplateName() string {
	return "generic/footerTemplate.tmpl"
}
