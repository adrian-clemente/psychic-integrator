package page

import "components/web/section/generic"
import "components/web/section/release"
import "components/printer"
import "fmt"

type ReleasePerformedPage struct {
	ProjectName string
	ReleaseResult bool
}

func (page *ReleasePerformedPage)GetContent() string {

	resultHeader := fmt.Sprintf("%v release summary", page.ProjectName)
	result := "Release submit has failed"
	if page.ReleaseResult {
		result = "Release submit has finished correctly"
	}


	headerSection := generic.HeaderSection{"Release manager"}
	bodySection := release.BodyReleaseResultSection{resultHeader, result}
	footerSection := generic.FooterSection{}
	mainSection := generic.MainSection{headerSection, bodySection, footerSection}
	printerPage := printer.PrinterPage{}

	content, _ := printerPage.PrintContent(mainSection);

	return content;
}
