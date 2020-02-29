package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Lesson implements Viewport
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

	// Check whether there are already lessons.
	for _, file := range files {
		splits := strings.Split(file.Name(), ".")
		if splits[len(splits)-1] == "yaml" {
			return
		}
	}

	// Copy bags of words.
	files, err = ioutil.ReadDir("data/bags_of_words")
	if err != nil {
		errorHandling(err)
	}

	for _, bag_of_words := range files {
		sourceFile, err := os.Open(fmt.Sprintf("data/bags_of_words/%s", bag_of_words.Name()))
		if err != nil {
			errorHandling(err)
		}
		targetFile, err := os.Create(home + fmt.Sprintf("/.config/gotypist/bags_of_words/%s", bag_of_words.Name()))
		if err != nil {
			errorHandling(err)
		}
		io.Copy(targetFile, sourceFile)

		sourceFile.Close()
		targetFile.Close()
	}

	// Create sample lessons if no lessons exist.
	files, err = ioutil.ReadDir("data/sample_lessons")
	if err != nil {
		errorHandling(err)
	}

	for _, lesson := range files {
		sourceFile, err := os.Open(fmt.Sprintf("data/sample_lessons/%s", lesson.Name()))
		if err != nil {
			errorHandling(err)
		}
		targetFile, err := os.Create(home + fmt.Sprintf("/.config/gotypist/lessons/%s", lesson.Name()))
		if err != nil {
			errorHandling(err)
		}
		io.Copy(targetFile, sourceFile)

		sourceFile.Close()
		targetFile.Close()
	}
}
