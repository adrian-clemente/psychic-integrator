package config

import "os"
import "log"
import "bufio"
import "strings"

var systemProperties SystemProperties

type SystemProperties struct {
	Values map[string]string
}

func LoadConfig() {

	log.Print("Loading system properties...")

	systemProps := SystemProperties{make(map[string]string)}

	file, err := os.Open("config/properties.conf")
	if err != nil {
		log.Fatal(err)
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			propertyKeyValue := strings.Split(scanner.Text(), "=")
			if len(propertyKeyValue) == 2 {
				systemProps.Values[propertyKeyValue[0]] = propertyKeyValue[1]
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	log.Print("Finished loading system properties...")

	systemProperties = systemProps;
}

func GetProperty(key string) string {
	return systemProperties.Values[key]
}
