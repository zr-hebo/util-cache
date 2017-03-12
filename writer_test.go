package bufferedwriter

import (
	"os"
	"testing"
)

func Test_Writer(t *testing.T) {
	bufferSize := 1024 * 1024 * 5
	file, err := os.Create("test.txt")
	if err != nil {
		t.Log(err.Error())
	}

	writer, err := NewBufferedWriter(file, bufferSize)
	if err != nil {
		t.Log(err.Error())
		return
	}

	for i := 0; i < 500; i++ {
		var rowData []byte
		for j := 0; j < 256*1024; j++ {
			rowData = append(rowData, []byte("hehe\n")...)
		}

		t.Log("write one line")
		if err = writer.Write(rowData); err != nil {
			t.Log(err.Error())
			return
		}
	}

	if err = writer.WaitClean(); err != nil {
		t.Log(err.Error())
		return
	}
}
