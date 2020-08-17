package utils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func GetDir(dirpath string) ([]string,error) {
	dirs:=make([]string,0)
	dir,err:=ioutil.ReadDir(dirpath)
	if err!=nil {
		return nil,err
	}
	pth:=string(os.PathSeparator)
	for _,f:=range dir {
		if f.IsDir() {
			dirs=append(dirs,dirpath+pth+f.Name())
		}
	}
	return dirs,nil
}
func GetFileList(dirpath string) ([]string,error) {
	files:=make([]string,0)
	dir,err:=ioutil.ReadDir(dirpath)
	if err!=nil {
		return nil,err
	}
	pth:=string(os.PathSeparator)
	for _,f:=range dir {
		if !f.IsDir() {
			if strings.HasSuffix(f.Name() ,".txt"){
				files=append(files,dirpath+pth+f.Name())
			}
		}
	}
	return files,nil
}
func FileIsExisted(path string) bool {
	_,err:=os.Stat(path)
	if err!=nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func CopyDir(srcPath, desPath string) error {
	//检查目录是否正确
	if srcInfo, err := os.Stat(srcPath); err != nil {
		return err
	} else {
		if !srcInfo.IsDir() {
			return errors.New("源路径不是一个正确的目录！")
		}
	}

	if desInfo, err := os.Stat(desPath); err != nil {
		return err
	} else {
		if !desInfo.IsDir() {
			return errors.New("目标路径不是一个正确的目录！")
		}
	}

	if strings.TrimSpace(srcPath) == strings.TrimSpace(desPath) {
		return errors.New("源路径与目标路径不能相同！")
	}

	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		//复制目录是将源目录中的子目录复制到目标路径中，不包含源目录本身
		if path == srcPath {
			return nil
		}

		//生成新路径
		destNewPath := strings.Replace(path, srcPath, desPath, -1)

		if !f.IsDir() {
			CopyFile(path, destNewPath)
		} else {
			if !FileIsExisted(destNewPath) {
				return MakeDir(destNewPath)
			}
		}

		return nil
	})

	return err
}

func CopyFile(src, des string) (written int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	//获取源文件的权限
	fi, _ := srcFile.Stat()
	perm := fi.Mode()

	//desFile, err := os.Create(des)  //无法复制源文件的所有权限
	desFile, err := os.OpenFile(des, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)  //复制源文件的所有权限
	if err != nil {
		return 0, err
	}
	defer desFile.Close()

	return io.Copy(desFile, srcFile)
}

func MakeDir(dir string) error {
	if !FileIsExisted(dir) {
		if err := os.MkdirAll(dir, 0777); err != nil { //os.ModePerm
			fmt.Println("MakeDir failed:", err)
			return err
		}
	}
	return nil
}
