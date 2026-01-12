package eftprocess

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Processor handles the end-to-end EFT notification flow.
type Processor struct {
	fileConf FileProcessorConfig
	notifier Notifier
}

// NewProcessor creates and initializes a new EFT processor.
func NewProcessor(fileConf FileProcessorConfig, notifier Notifier) (*Processor, error) {
	if err := ensureMandatoryDirsExist(fileConf.InputDir, fileConf.DoneDir, fileConf.ErrorDir); err != nil {
		return nil, err
	}
	return &Processor{
		fileConf: fileConf,
		notifier: notifier,
	}, nil
}

func ensureMandatoryDirsExist(dirs ...string) error {
	var allErrs []error
	for _, dir := range dirs {
		info, err := os.Stat(dir)
		if os.IsNotExist(err) {
			allErrs = append(allErrs, fmt.Errorf("mandatory dir not found: %s. Hint: create it manually", dir))
			continue
		}
		if err != nil {
			allErrs = append(allErrs, fmt.Errorf("error accessing path %s: %v", dir, err))
			continue
		}
		if !info.IsDir() {
			allErrs = append(allErrs, fmt.Errorf("path is not a directory: %s", dir))
		}
	}
	if len(allErrs) > 0 {
		return errors.Join(allErrs...)
	}
	return nil
}

// FilesMatch finds files matching the configured glob pattern that are ready for processing.
func (p *Processor) FilesMatch() ([]InputFileInfo, error) {
	matchGlob := filepath.Join(p.fileConf.InputDir, p.fileConf.GlobPath)
	
	log.Printf("Searching for files matching: %s", matchGlob)
	
	var files []InputFileInfo
	err := filepath.Walk(p.fileConf.InputDir, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		matched, err := filepath.Match(matchGlob, walkPath)
		if err != nil {
			return err
		}

		if matched {
			inputFileInfo := InputFileInfo{Path: walkPath, Info: info}
			if p.isFileReadyForProcessing(inputFileInfo) {
				files = append(files, inputFileInfo)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking input directory: %w", err)
	}
	return files, nil
}

func (p *Processor) isFileReadyForProcessing(fileInfo InputFileInfo) bool {
	if !isOlderThanSecs(fileInfo.Info.ModTime(), p.fileConf.OlderThanSeconds) {
		return false
	}
	return p.isBankFileUploaded(fileInfo)
}

func isOlderThanSecs(fileTime time.Time, olderSec int) bool {
	return time.Since(fileTime) > time.Duration(olderSec)*time.Second
}

func (p *Processor) isBankFileUploaded(fileInfo InputFileInfo) bool {
	bankfilepath := strings.ReplaceAll(fileInfo.Path, "566", "565")
	if fileExists(bankfilepath) {
		log.Printf("Bank file %s exists. Skipping processing for %s", bankfilepath, fileInfo.Info.Name())
		return false
	}
	return true
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil && !info.IsDir()
}

// Process parses the CSV file and sends emails.
func (p *Processor) Process(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		errStr := fmt.Sprintf("Error opening input file %s: %v", filePath, err)
		log.Println(errStr)
		_ = p.notifier.Error("Error: Markham Notification - EFT", errStr)
		return err
	}

	eftInfos, err := GetEftInfosFromCSV(string(data))
	if err != nil {
		errStr := fmt.Sprintf("Error parsing input file %s: %v", filePath, err)
		log.Println(errStr)
		_ = p.notifier.Error("Error: Markham Notification - EFT", errStr)
		return err
	}

	if err := p.notifier.NotifyAll(eftInfos); err != nil {
		log.Printf("Error sending emails for %s: %v", filePath, err)
		return err
	}

	log.Printf("Successfully processed %s. Emails sent: %d", filePath, len(eftInfos.EftInfos))
	return nil
}

// PostProcess moves the file to the target directory with a timestamp.
func (p *Processor) PostProcess(fileInfo InputFileInfo, targetPath string) {
	timestamp := time.Now().Format("2006-01-02_150405.000")
	newFileName := fileInfo.Info.Name() + "_" + timestamp
	newFullName := filepath.Join(targetPath, newFileName)

	if err := os.Rename(fileInfo.Path, newFullName); err != nil {
		log.Fatalf("Critical error: failed to move processed file %s to %s: %v", fileInfo.Path, newFullName, err)
	}
}
