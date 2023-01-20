package multithread

import (
	"strings"

	"github.com/HakimMouchquelita/GO-F22ssh/src/fs"
)

func SplitFiles(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func CopyFolder(ending chan bool, folders []string) {
	for _, folder := range folders {
		s := strings.Split(folder, ":")
		fs.CopyFolder(s[0], s[1])
	}

	ending <- true
}

func CopyFile(ending chan bool, files []string) {
	for _, file := range files {
		s := strings.Split(file, ":")
		fs.CopyFile(s[0], s[1])
	}

	ending <- true
}
