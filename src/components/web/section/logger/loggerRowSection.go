package logger

type LoggerRowSection struct {
	LevelColor string
	CreationTime string
	Level string
	Text string
}

func (page LoggerRowSection)GetTemplateName() string {
	return "web/logger/loggerRowTemplate.tmpl"
}
