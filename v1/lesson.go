package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Lesson implements Viewport
type Lesson struct {
	Title   string
	Content string
	Tag     ModeTag
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
	err = createLessonFromBagOfWords(home)
	if err != nil {
		errorHandling(err)
	}

	// Create sample lessons if no lessons exist.
	err = createLessonsFromSampleLessons(home)
	if err != nil {
		errorHandling(err)
	}
}

func createLessonFromBagOfWords(home string) error {
	files, err := ioutil.ReadDir("data/bags_of_words")
	if err != nil {
		return err
	}
	for _, bag_of_words := range files {
		sourceFilePath := fmt.Sprintf("data/bags_of_words/%s", bag_of_words.Name())
		targetFilePath := home + fmt.Sprintf("/.config/gotypist/bags_of_words/%s", bag_of_words.Name())
		err := createLesson(sourceFilePath, targetFilePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func createLessonsFromSampleLessons(home string) error{
	files, err := ioutil.ReadDir("data/sample_lessons")
	if err != nil {
		return err
	}

	for _, lesson := range files {
		sourceFilePath:= fmt.Sprintf("data/sample_lessons/%s", lesson.Name())
		targetFilePath := home + fmt.Sprintf("/.config/gotypist/lessons/%s", lesson.Name())
		err := createLesson(sourceFilePath, targetFilePath)
		if err != nil {
			return err
		}

	}
	return nil
}

func createLesson(sourceFilePath, targetFilePath string) error {
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	err = os.MkdirAll(filepath.Dir(targetFilePath), os.ModePerm)
	if err != nil {
		return err
	}
	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, sourceFile)
	return err
}
