package logging

import (
	"time"
	"container/list"
	"strings"
	"regexp"
)

const (
	WARN = "WARN"
	ERROR = "ERROR"
	INFO = "INFO"
	DEBUG = "DEBUG"
	WARN_CLASS = "warn"
	ERROR_CLASS = "error"
	INFO_CLASS = "info"
	DEBUG_CLASS = "debug"
	CARTERO_SERVICE = "cartero"
	FAVOR_SERVICE = "favor"
)

type LoggerLines struct {
	List *list.List
}

type LoggerLine struct {
	Text string
	Level string
	CreationTime string
	LoggerService string
}

var loggerLines LoggerLines = LoggerLines{list.New()}

func (loggerLine *LoggerLine) ToString() string {
	return loggerLine.CreationTime + " " + loggerLine.Level + " " + loggerLine.Text
}


func (loggerLines *LoggerLines) addLoggerLine(loggerLine *LoggerLine) {
	loggerLines.List.PushBack(loggerLine)
}

func RetrieveLoggerLines() (*LoggerLines) {
	return &loggerLines;
}

func GetLevelColor(stringLevel string) string {
	stringLevel = strings.TrimSpace(stringLevel)
	colors := map[string]string{
		WARN: WARN_CLASS,
		ERROR: ERROR_CLASS,
		INFO: INFO_CLASS,
		DEBUG: DEBUG_CLASS,
	}
	return colors[stringLevel]
}

func PrintLoggerLine(inputChan <- chan string) {
	
	for {
		loggerLine := formatRawText(<- inputChan)
		loggerLines.addLoggerLine(&loggerLine)
		//fmt.Println(loggerLine);
		time.Sleep(time.Second * 1)
	}
}

func formatRawText(rawText string) LoggerLine {
	//.Println(rawText)
	rawTextTokens := strings.Split(rawText, " ")

	dateRegex := regexp.MustCompile("\\d{8}")
	timeRegex := regexp.MustCompile("\\d{2}:\\d{2}:\\d{2}:\\d{3}")
	timeAltRegex := regexp.MustCompile("\\d{8}T\\d{6}.\\d{3}\\+\\d{4}")
	levelRegex := regexp.MustCompile("INFO|WARN|ERROR|DEBUG")
	serviceRegex := regexp.MustCompile("loggerCartero|loggerFavor")

	loggerEventTime := ""
	loggerEventLevel := ""
	loggerEventText := ""
	loggerEventService := ""

	for _, rawTextToken := range rawTextTokens {

		rawTextTokenByte := []byte(rawTextToken)

		if (dateRegex.Match(rawTextTokenByte) && loggerEventTime == "") || timeRegex.Match(rawTextTokenByte) ||
				timeAltRegex.Match(rawTextTokenByte) {
			loggerEventTime = loggerEventTime + " " + rawTextToken
		} else if levelRegex.Match(rawTextTokenByte) {
			loggerEventLevel = loggerEventLevel+" "+rawTextToken
		} else if serviceRegex.Match(rawTextTokenByte) {
			loggerEventService = rawTextToken;
		} else {
			loggerEventText = loggerEventText + " " + rawTextToken;
		}
	}

	return LoggerLine{strings.TrimSpace(loggerEventText), strings.TrimSpace(loggerEventLevel),
		strings.TrimSpace(loggerEventTime), strings.TrimSpace(loggerEventService)};
}
