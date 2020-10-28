package utils

import (
	//"archive/tar"
	//"fmt"
	//"io"
	//"io/ioutil"
	//"os"
	"testing"
)

func TestTar(t *testing.T) {
	Tar(".git", "a.tar", false)
}

/*
func TestTar(t *testing.T) {
	filename := "a.tar"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		println(err.Error())
		return
	}
	tw, err := NewTar(f)
	if err != nil {
		println("new tar error:", err.Error())
		return
	}
	info, err := ioutil.ReadDir(".")
	if err != nil {
		println(err.Error())
		return
	}
	for _, v := range info {
		if v.Name() == filename {
			continue
		}
		F, err := os.Open(v.Name())
		if err != nil {
			println(err.Error())
			return
		}
		TarAddFile(F, "", tw)
		F.Close()
	}
}
*/
