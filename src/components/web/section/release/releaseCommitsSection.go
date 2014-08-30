package release

import "components/web/element"
import "components/printer"

type CommitsSection struct {
	ReleaseCommitsElements []element.CommitElement
	Project string
	Version string
}

func (page CommitsSection)GetTemplateName() string {
	return "web/release/commitsReleaseTemplate.tmpl"
}

func (page *CommitsSection)GetContent() string {
	printerPage := printer.PrinterPage{}
	content, _ := printerPage.PrintContent(*page);
	return content;
}
