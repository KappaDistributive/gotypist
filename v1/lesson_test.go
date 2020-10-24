package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func Test_createLesson_ShouldNotReturnErrorWhenAbleToCreateLesson(t *testing.T) {
	closingActions := setupSource(t)
	defer closingActions()
	targetDirectory := "./tmp/target/"
	sourceFilePath := "./tmp/source/source-file"
	targetFile := "target-file"
	targetFilePath := targetDirectory + targetFile

	err := createLesson(sourceFilePath, targetFilePath)

	assert.NoError(t, err)
	files, err := ioutil.ReadDir(targetDirectory)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Contains(t, files[0].Name(), targetFile)
}

func Test_createLesson_ShouldReturnErrorWhenUnableToCreateLesson(t *testing.T) {
	targetDirectory := "./tmp/target/"
	sourceFilePath := "./tmp/source/source-file"
	targetFile := "target-file"
	targetFilePath := targetDirectory + targetFile

	err := createLesson(sourceFilePath, targetFilePath)

	assert.Error(t, err)
}
func setupSource(t *testing.T) func(){
	err := os.MkdirAll("./tmp/source/", os.ModePerm)
	fileHandle, err := os.Create("./tmp/source/source-file")
	defer fileHandle.Close()

	assert.NoError(t, err)

	return func() {
		err := os.RemoveAll("./tmp")
		assert.NoError(t, err)
	}
}