package processor

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"muni/go-mail/internal/config"
	"muni/go-mail/internal/mail"
	"muni/go-mail/internal/model"
	"muni/go-mail/internal/parser"
)

func FilesMatch(conf config.FileProcessorConfig) (files []model.InputFileInfo, err error) {
	matchGlob := conf.InputDir + conf.GlobPath //path.Join - Does not work for windows(see log snippet below), going back to +
	log.Println("Find match for:", matchGlob)
	err = filepath.Walk(conf.InputDir,
		func(walkPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// log.Println("File :", matchGlob, walkPath)
			// Note: filepath.Match checks against the name, not the full path usually, unless pattern contains separators?
			// But here we matched full path against full glob?
			// The original code passed `matchGlob` (dir + glob) to `filepath.Match` against `walkPath` (full path).
			// This works if matchGlob matches the full structure.
			if matched, _ := filepath.Match(matchGlob, walkPath); matched == true { //info.Mode().IsRegular()
				//fmt.Println("Yes glob match : ", walkPath, info.Size())
				inputFileInfo := model.InputFileInfo{Path: walkPath, Info: info}
				if isReadyForCayantaEFT(inputFileInfo, conf) == true {
					files = append(files, inputFileInfo) //fileInfo is expensive, just return and reuse
				}
			} else {
				//fmt.Println("No glob match - Skip: ", walkPath, info.Size())
			}

			return nil
		}) //filewalk
	return
}

func isReadyForCayantaEFT(fileInfo model.InputFileInfo, conf config.FileProcessorConfig) bool {
	if isOlderThanSecs(fileInfo.Info.ModTime(), conf.OlderThanSeconds) == false {
		return false //too new, pass ...let time go..
	}
	return isBankFileUploaded(fileInfo)
}

func isOlderThanSecs(fileTime time.Time, olderSec int) bool {
	now := time.Now()
	diff := now.Sub(fileTime)
	cutoff := time.Duration(olderSec) * time.Second
	//log.Println("Now:", time.Now(), ", Cutoff:", cutoff, ", diff:", diff)
	return diff > cutoff
}

func isBankFileUploaded(fileInfo model.InputFileInfo) bool {
	//Original logic: bankfilepath := strings.ReplaceAll(fileInfo.path, "566", "565")
	bankfilepath := strings.ReplaceAll(fileInfo.Path, "566", "565")
	//log.Println("Check for this bank file : ", bankfilepath)
	if fileExists(bankfilepath) {
		log.Println("Bank file not processet YET! Skip: ", fileInfo.Info.Name(), bankfilepath)
		return false
	}
	return true
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	//log.Println("fileExists(): info=", info)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func EnsureMandatoryDirsExist(dirs ...string) error {
	for _, dir := range dirs {
		info, err := os.Stat(dir)
		if os.IsNotExist(err) {
			return errors.New("mandatory dir not found. Hint: create necessary folders manually before executing: " + dir)
		}
		if !info.IsDir() {
			return errors.New("mandatory dir not found. Hint: create necessary folders manually before executing: " + dir)
		}
	}
	return nil
}

func Process(filePath string, mailConf config.MailServerConfig) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		var errStr = "Error opening input file: " + filePath + ", " + err.Error()
		log.Println(errStr) //Go to next file. Email?
		_ = mail.SendErrorAlert(mailConf, errStr)
		return err

	}
	input := string(data)

	eftInfos, err := parser.GetEftFromCSV(input)
	//eftInfos.EftInfos[1].Invoices
	if err != nil {
		var errStr = "Error parsing input file: " + filePath + ", " + err.Error()
		log.Println(errStr) //Go to next file. Email?
		err2 := mail.SendErrorAlert(mailConf, errStr)
		if err2 != nil {
			log.Println(err2)
		}
		return err
	}
	//send mails
	err1 := mail.BatchSendMail(mailConf, eftInfos)
	if err1 != nil {
		log.Println("Error sending emails for input file:", filePath, err1) //Go to next file. Email?
		return err1
	}
	log.Println("Processed done:", filePath, "Emails sent #: ", len(eftInfos.EftInfos)) //Go to next file. Email?

	return nil
}

func PostProcess(fileInfo model.InputFileInfo, targetPath string) {
	currentfileName := fileInfo.Info.Name()
	newFileName := currentfileName + time.Now().Format("2006-01-02_150405.000") // Fixed format to standard Go time format string (original was 2020-01-31 which is not standard layout usage, Go uses Mon Jan 2 15:04:05 MST 2006 reference time)
	// Original was: time.Now().Format("2020-01-31_154560.555") 
	// The reference time is 2006-01-02 15:04:05.
	// 2020-01-31... looks like they HARDCODED a specific format using non-standard chars or they misunderstood Go layouts.
	// "2006-01-02" is the standard layout. 
	// If I use their string literally "2020-01-31_154560.555", it might just print static text if it doesn't match reference.
	// But it seems they wanted a timestamp.
	// I'll assume they wanted a timestamp. "2006-01-02_150405.000" is a safe bet for a unique timestamp.
	
	newfullName := targetPath + newFileName
	if e := os.Rename(fileInfo.Path, newfullName); e != nil {
		log.Fatal("Error moving processed file to target dir: ", newfullName, e)
	} //
}
