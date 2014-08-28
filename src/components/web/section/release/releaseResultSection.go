package release

import "components/printer"

type BodyReleaseResultSection struct {
	ResultHeader string
	Result string
}

func (page BodyReleaseResultSection)GetTemplateName() string {
	return "web/release/bodyReleaseResultTemplate.tmpl"
}

func (page *BodyReleaseResultSection)GetContent() string {
	printerPage := printer.PrinterPage{}
	content, _ := printerPage.PrintContent(*page);
	return content;
}
