package main

import (
	"log"
	"os"
	"time"

	"path/filepath"
)

func main12() {
	// basePath := "/Volumes/data01/projects/projects_go/go-mail/inputfile*"
	basePath := "C:\\data\\projects_go\\go-mail\\input\\"
	conf := FileProcessorConfig{GlobPath: "*556*.csv", InputDir: basePath, OlderThanSeconds: 2}
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
func filesMatch(conf FileProcessorConfig) (files []InputFileInfo, err error) {
	matchGlob := conf.InputDir + conf.GlobPath //path.Join - Does not work for windows(see log snippet below), going back to +
	log.Println("Find match for:", matchGlob)
	// 2020/09/26 21:31:59 Find match for: C:\data\projects_go\go-mail\input/*556*.csv
	// 2020/09/26 21:31:59 FIle : C:\data\projects_go\go-mail\input/*556*.csv C:\data\projects_go\go-mail\input
	// 2020/09/26 21:31:59 FIle : C:\data\projects_go\go-mail\input/*556*.csv C:\data\projects_go\go-mail\input\input556.csv
	// 2020/09/26 21:31:59 []  <- hmm no files found! input566.csv must be found.

	// var files []os.FileInfo
	err = filepath.Walk(conf.InputDir,
		func(walkPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// log.Println("File :", matchGlob, walkPath)
			if matched, _ := filepath.Match(matchGlob, walkPath); matched == true { //info.Mode().IsRegular()

				if isOlderThanSecs(info.ModTime(), conf.OlderThanSeconds) {
					log.Println("Match: ", walkPath)
				}
				myNonsenseFinfoWithPath := InputFileInfo{path: walkPath, info: info}
				files = append(files, myNonsenseFinfoWithPath) //fileInfo is expensive, just return and reuse
			} else {
				//fmt.Println("Skip: ", walkPath, info.Size())
			}

			return nil
		}) //filewalk
	return
}

type InputFileInfo struct {
	path string
	info os.FileInfo
}

// func filesMatch(conf FileFilterConfig) (files []os.FileInfo, err error) { //<- chan string
// 	now := time.Now()
// 	path.Split(conf.globPath)
// 	fileInfos, _ := glob(conf.globPath)
// 	for finf := range fileInfos {
// 		if finf.Mode().IsRegular() {
// 			if isOlderThanSecs(finf.ModTime(), conf.olderThanSeconds) {
// 				files = append(files, finf)
// 			}else{
// 				log.Println("Skipped", finf.Name())
// 			}
// 		}
// 	}
// 	return
// }

// copy of filpath.Glob but directly returning FileInfo
// func glob(dir string, ext string) ([]os.FileInfo, error) {
// 	files := []os.FileInfo
// 	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
// 		if filepath.Ext(path) == ext {
// 			files = append(files, f)
// 		}
// 		return nil
// 	})
// 	return files, err
// }
