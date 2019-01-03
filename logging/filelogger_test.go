package logging

import (
	"fmt"
	"testing"
)

func TestBuildFileName(t *testing.T) {
	filelogger, err := New("../../", "TestFile", ".dat",
		60,
		true, true, true,
		10)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Log(filelogger.currentFileName)
	for i := 0; i < 100000000; i++ {
		s, err := filelogger.Write([]byte("This is a very good line, it has been written by a great scientist!!\n"))
		if s <= 0 || err != nil {
			fmt.Printf("Write failed:[%d][%s]\n", s, err)
			return
		}
	}
	filelogger.Close()
}
