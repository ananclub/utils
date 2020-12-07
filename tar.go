package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"

	//"io/ioutil"
	"os"
	"strings"
)

type TarFile struct {
	gw *gzip.Writer
	tw *tar.Writer
}

func NewGzipFile(dest string, deleteIfExist bool) (gw *gzip.Writer, err error) {
	d, err := os.Create(dest)
	if err == os.ErrExist {
		if err = os.Remove(dest); err != nil {
			return
		}
	} else {
		return
	}
	gw = NewGzip(d)
	return
}
func NewGzip(w io.Writer) (gw *gzip.Writer) {
	gw = gzip.NewWriter(w)
	return
}

func NewTarFile(dest string, deleteIfExist bool) (tw *tar.Writer, err error) {
	d, err := os.Create(dest)
	if err == os.ErrExist {
		if err = os.Remove(dest); err != nil {
			return
		}
	} else {
		return
	}
	tw = tar.NewWriter(d)
	return
}
func NewTar(w io.Writer) (tw *tar.Writer, err error) {

	tw = tar.NewWriter(w)
	return
}

//压缩 使用gzip压缩成tar.gz
func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	gw := gzip.NewWriter(d)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for _, file := range files {
		err := TarAddFile(tw, file, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func TarAddFile(tw *tar.Writer, file *os.File, prefix string) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = TarAddFile(tw, f, prefix)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//解压 tar.gz
func DeCompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := dest + hdr.Name
		file, err := createFile(filename)
		if err != nil {
			return err
		}
		io.Copy(file, tr)
	}
	return nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
func Tar(filesource, filetarget string, deleteIfExist bool) (err error) {
	var tarfile *os.File
	if deleteIfExist {
		tarfile, err = os.OpenFile(filetarget, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
	} else {
		tarfile, err = os.OpenFile(filetarget, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
	}

	defer tarfile.Close()
	tw, err := NewTar(tarfile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	f, err := os.Open(filesource)
	if err != nil {
		fmt.Println(err)
		return err
	}
	TarAddFile(tw, f, "")
	tw.Flush()
	return nil
}
