package main

import (
	"log"
	"os"
	"strings"
	"time"

	"path/filepath"
)

func main12() {
	// basePath := "/Volumes/data01/projects/projects_go/go-mail/inputfile*"
	//basePath := "C:\\data\\projects\\projects_go\\go-mail\\input\\"
	basePath := "\\\\markham.ca\\data\\Cayenta79_APDATA_TEST\\"
	conf := FileProcessorConfig{GlobPath: "*566*.csv", InputDir: basePath, OlderThanSeconds: 2}
	files, err := filesMatch(conf)
	if err != nil {
		log.Println(err)
	}
	log.Println(files)

}
func isOlderThanSecs(fileTime time.Time, olderSec int) bool {
	now := time.Now()
	diff := now.Sub(fileTime)
	cutoff := time.Duration(olderSec) * time.Second
	//log.Println("Now:", time.Now(), ", Cutoff:", cutoff, ", diff:", diff)
	return diff > cutoff
}

func isReadyForCayantaEFT(fileInfo InputFileInfo, conf FileProcessorConfig) bool {
	if isOlderThanSecs(fileInfo.info.ModTime(), conf.OlderThanSeconds) == false {
		return false //too new, pass ...let time go..
	}
	return isBankFileUploaded(fileInfo)
}

func isBankFileUploaded(fileInfo InputFileInfo) bool {
	bankfilepath := strings.ReplaceAll(fileInfo.path, "566", "565")
	//log.Println("Check for this bank file : ", bankfilepath)
	if fileExists(bankfilepath) {
		log.Println("Bank file not processet YET! Skip: ", fileInfo.info.Name(), bankfilepath)
		return false
	}
	return true
}

func filesMatch(conf FileProcessorConfig) (files []InputFileInfo, err error) {
	matchGlob := conf.InputDir + conf.GlobPath //path.Join - Does not work for windows(see log snippet below), going back to +
	log.Println("Find match for:", matchGlob)
	err = filepath.Walk(conf.InputDir,
		func(walkPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// log.Println("File :", matchGlob, walkPath)
			if matched, _ := filepath.Match(matchGlob, walkPath); matched == true { //info.Mode().IsRegular()
				//fmt.Println("Yes glob match : ", walkPath, info.Size())
				inputFileInfo := InputFileInfo{path: walkPath, info: info}
				if isBankFileUploaded(inputFileInfo) == true {
					files = append(files, inputFileInfo) //fileInfo is expensive, just return and reuse
				}
			} else {
				//fmt.Println("No glob match - Skip: ", walkPath, info.Size())
			}

			return nil
		}) //filewalk
	return
}

// InputFileInfo wrapper to store path and os file info
type InputFileInfo struct {
	path string
	info os.FileInfo
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	//log.Println("fileExists(): info=", info)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

