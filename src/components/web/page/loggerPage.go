package page

import "components/web/section/logger"
import "components/web/section/generic"
import "components/printer"
import "api/logging"

type LoggerPage struct {
	LoggerRows *logging.LoggerLines
}

func (page *LoggerPage)GetContent() string {

	var loggerRowSections []logger.LoggerRowSection  // an empty list

	for c := page.LoggerRows.List.Front(); c != nil; c = c.Next() {
		loggerLine := (c.Value.(*logging.LoggerLine))

		loggerRowSection := logger.LoggerRowSection{loggerLine.Level, loggerLine.CreationTime, loggerLine.Level, loggerLine.Text}
		loggerRowSections = append(loggerRowSections, loggerRowSection)
	}

	headerSection := generic.HeaderSection{"Logger"}
	bodySection := logger.BodyLoggerSection{loggerRowSections}
	footerSection := generic.FooterSection{}
	mainSection := generic.MainSection{headerSection, bodySection, footerSection}
	printerPage := printer.PrinterPage{}

	content, _ := printerPage.PrintContent(mainSection);

	return content;
}
