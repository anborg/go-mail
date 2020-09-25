package main

import (
	"log"
	"os"

	// "path"
	"path/filepath"
)

func main1() {
	// macPath := "/Volumes/data01/projects/projects_go/go-mail/inputfile*"
	winPath := "C:/data/projects_go/go-mail/inputfile*"
	var conf = FileFilterConfig{globPath: winPath, olderThanSeconds: 2}
	//files, _ := filesMatch(conf)

	err := filepath.Walk(conf.globPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			log.Println(path, info.Size())
			return nil
		}) //filewalk
	if err != nil {
		log.Println(err)
	}

	// log.Println(files)
	// for file := range files {
	// 	log.Println(file)
	// }
}

type FileFilterConfig struct {
	globPath         string
	olderThanSeconds uint32
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

// func isOlderThanSecs(fileTime, olderSec int32) bool {
// 	return time.Now().Sub(fileTime) > olderSec*time.Second
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
