package utils

import (
	"io/ioutil"
	"os"
)

// ReadShader ReadShader
func ReadShader(filename string) (shader string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return string(bytes) + "\x00"
}
