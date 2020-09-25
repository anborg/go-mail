package main
import (
	"path/filepath"
	"log"
	"os"
	"time"
)

func main(){
	var conf = FileFilterConfig{globPath: "/Volumes/data01/projects/projects_go/go-mail/inputfile*",olderThanSeconds: 2}
	files, _ := filesMatch(conf)
	log.Println(files)
	for file := range files {
		log.Println(file)
	}
}

type FileFilterConfig struct {
	globPath string
	olderThanSeconds uint32
}

func filesMatch(conf FileFilterConfig) (files []os.FileInfo, err error){//<- chan string
	now := time.Now()
	tempfiles, _ := glob(conf.globPath)
	for f := range tempfiles {
		if f.Mode().IsRegular(){
			if isOlderThanSecs(f.ModTime(),conf.olderThanSeconds){
				files = append(files, f)
			}
		}
	}
	return 
}

func isOlderThanSecs(fileTime, olderSec int32) bool{
	return time.Now().Sub(fileTime) > olderSec*time.Second
}

// copy of filpath.Glob but directly returning FileInfo
func glob(dir string, ext string) ( []os.FileInfo,  error) {
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
	  if filepath.Ext(path) == ext {
		files = append(files, f)
	  }
	  return nil
	})
  	return files, err
  }