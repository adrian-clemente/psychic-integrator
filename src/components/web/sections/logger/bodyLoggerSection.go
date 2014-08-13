package logger

type BodyLoggerSection struct {
	LoggerRows []LoggerRowSection
}

func (page BodyLoggerSection)GetTemplateName() string {
	return "logger/bodyLoggerTemplate.tmpl"
}
