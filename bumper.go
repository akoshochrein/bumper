package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var basePath = "/Users/akos.hochrein/workspace"
var reqFileName = "requirements.txt"

type Package struct {
	name string
	version string
}

func main() {
	p := Package{
	"test-package",
	"1.0.4",
	}

	files, err := ioutil.ReadDir(basePath)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		go ProcessDirectoryCandidate(file, p)
	}
}

func ProcessDirectoryCandidate(file os.FileInfo, p Package) {
	if file.IsDir() {
		repoName := filepath.Join(basePath, file.Name())
		reqFilePath := filepath.Join(repoName, reqFileName)

		if _, err := os.Stat(reqFilePath); !os.IsNotExist(err) {
			ProcessRequirementsFile(reqFilePath, p)
		}
	}
}

func ProcessRequirementsFile(reqFilePath string, p Package) {
	reqFile, err := os.Open(reqFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer reqFile.Close()

	BuildNewRequirementsFile(reqFile, err, reqFilePath, p)
}

func BuildNewRequirementsFile(reqFile *os.File, err error, newReqFilePath string, p Package) *os.File {
	s := bufio.NewScanner(reqFile)
	newReqF, err := os.Create(newReqFilePath + "2")
	if err != nil {
		log.Fatal(err)
	}

	for s.Scan() {
		t := s.Text()
		if strings.Contains(t, p.name) {
			re := regexp.MustCompile("[0-9]+\\.[0-9]+\\.[0-9]+")
			newReqF.Write([]byte(re.ReplaceAllString(t, p.version)))
		} else {
			newReqF.Write([]byte(t))
		}
		newReqF.Write([]byte("\n"))
	}

	return newReqF
}
