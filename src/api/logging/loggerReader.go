package logging

import "os"
import "bufio"
import "time"
import "code.google.com/p/go.exp/inotify"
import "log"
import "io"

func ProcessFile(path string, fileName string, output chan<- string) {

	lineNum := ReadFile(path, fileName, 0, output)

	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Watch(path)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case ev := <-watcher.Event:
			if ev.Mask == inotify.IN_MODIFY {
				time.Sleep(1 * time.Second)
				lineNum = ReadFile(path, fileName, lineNum, output)
			}
		case err := <-watcher.Error:
			log.Println("error:", err)
		}

	}
}

func ReadFile(path string, fileName string, startingPoint int, output chan<- string) (lastLine int) {
	file, err := os.Open(path + "/" + fileName)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	lineNum := 0

	for {
		line, err := r.ReadString('\n')

		if (startingPoint < lineNum) {
			output <- line + " " + fileName
			time.Sleep(1 * time.Second)
		}
		if err == io.EOF {
			return lineNum
		}
		lineNum++
	}
}

