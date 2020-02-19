package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Lesson struct {
	Title   string
	Content string
}

func createSampleLessons() {
	home, err := os.UserHomeDir()
	if err != nil {
		errorHandling(err)
	}

	files, err := ioutil.ReadDir(home + "/.config/gotypist/lessons")
	if err != nil {
		errorHandling(err)
	}

	// check whether there are already lessons
	for _, file := range files {
		splits := strings.Split(file.Name(), ".")
		if splits[len(splits)-1] == "yaml" {
			return
		}
	}

	// create sample lessons if no lessons exist
	files, err = ioutil.ReadDir("data/sample_lessons")
	if err != nil {
		errorHandling(err)
	}

	for _, lesson := range files {
		source_file, err := os.Open(fmt.Sprintf("data/sample_lessons/%s", lesson.Name()))
		if err != nil {
			errorHandling(err)
		}
		target_file, err := os.Create(home + fmt.Sprintf("/.config/gotypist/lessons/%s", lesson.Name()))
		if err != nil {
			errorHandling(err)
		}
		io.Copy(target_file, source_file)

		source_file.Close()
		target_file.Close()
	}
}
