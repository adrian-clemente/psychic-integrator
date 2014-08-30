package page

import "components/web/element"
import "components/web/section/generic"
import "components/web/section/release"
import "components/printer"

type ReleaseHandlerPage struct {
	ReleaseCommitsElements []element.CommitElement
	ReleaseProjects []element.OptionElement
	Project string
	Version string
}

func (page *ReleaseHandlerPage)GetContent() string {

	headerSection := generic.HeaderSection{"Release manager"}
	commitSection := release.CommitsSection{page.ReleaseCommitsElements, page.Project, page.Version}
	bodySection := release.BodyReleaseMainSection{page.ReleaseProjects, commitSection }

	footerSection := generic.FooterSection{}
	mainSection := generic.MainSection{headerSection, bodySection, footerSection}
	printerPage := printer.PrinterPage{}

	content, _ := printerPage.PrintContent(mainSection);

	return content;
}
