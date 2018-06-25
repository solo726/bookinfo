package global

import (
	"os"
	"io/ioutil"
	"strconv"
)

func GenPid(fileName string) int {
	pid := os.Getpid()

	ioutil.WriteFile(fileName, []byte(strconv.Itoa(pid)), os.ModePerm)

	return pid
}
