package logging

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// FileLogger class.
type FileLogger struct {
	dir                string
	prefix             string
	suffix             string
	maxSizeInBytes     uint64
	rollover           bool
	haveDateInFileName bool
	haveTimeInFileName bool
	maxFileCount       uint16
	// Below fields are dynamic
	currenDate       string
	currenTime       string
	currentFileCount uint16
	currentFileName  string
	currentFileSize  uint64
	file             *os.File
}

func (filelogger *FileLogger) buildNewFileName() (string, error) {
	var filename string
	t := time.Now()
	if strings.HasSuffix(filelogger.dir, "/") {
		filename = filelogger.dir
	} else {
		filename = filelogger.dir
		filename += "/"
	}
	filename += filelogger.prefix
	if filelogger.haveDateInFileName {
		filename += "_"
		filename += fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
	}
	if filelogger.haveTimeInFileName {
		filename += "_"
		filename += fmt.Sprintf("%02d%02d%02d", t.Hour(), t.Minute(), t.Second())
	}
	if filelogger.rollover {
		if filelogger.currentFileCount < filelogger.maxFileCount {
			filename += "_"
			filelogger.currentFileCount++
			filename += fmt.Sprintf("%02d", filelogger.currentFileCount)
		} else {
			return "", fmt.Errorf("Reached max file count: [%d] for file: [%s]",
				filelogger.currentFileCount, filelogger.currentFileName)
		}

	}
	filename += filelogger.suffix
	return filename, nil
}

func New(path, filePrefix, fileSuffix string,
	fileMaxSizeInMegaBytes uint64,
	fileToRollover, appendDateInFileName, appendTimeInFileName bool,
	maximumFileCount uint16) (*FileLogger, error) {
	filelogger := FileLogger{
		dir:                path,
		prefix:             filePrefix,
		suffix:             fileSuffix,
		maxSizeInBytes:     fileMaxSizeInMegaBytes * 1024 * 1024,
		rollover:           fileToRollover,
		haveDateInFileName: appendDateInFileName,
		haveTimeInFileName: appendTimeInFileName,
		maxFileCount:       maximumFileCount,
		currentFileCount:   0,
	}
	var err error
	filelogger.currentFileName, err = filelogger.buildNewFileName()
	if err != nil {
		return nil, err
	}
	filelogger.file, err = os.OpenFile(filelogger.currentFileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	return &filelogger, nil
}

func (filelogger *FileLogger) Write(b []byte) (int, error) {
	currentFileSize, maxSize, size := filelogger.currentFileSize, filelogger.maxSizeInBytes, uint64(len(b))
	newFileSize := currentFileSize + size
	if maxSize <= newFileSize {
		filename, err := filelogger.buildNewFileName()
		if err != nil {
			return -1, err
		}
		if err := filelogger.file.Close(); err != nil {
			fmt.Printf("FileLogger:Warning - Unable to close previous file [%s] due to error: [%s]\n",
				filelogger.currentFileName, err.Error())
		}
		filelogger.currentFileName = filename
		newFileSize = size
		filelogger.file, err = os.OpenFile(filelogger.currentFileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return -1, err
		}
	}

	filelogger.currentFileSize = newFileSize
	return filelogger.file.Write(b)
}

func (filelogger *FileLogger) Close() error {
	if filelogger.file != nil {
		return filelogger.file.Close()
	}
	return nil
}
