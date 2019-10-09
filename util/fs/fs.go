package fs

import (
	. "github.com/tm-ad/g-base/util"
	"os"
	"path"
	"path/filepath"
)

// WorkingPath 获取当前应用的执行目录
func WorkingPath(offset string) string {
	if Development() {
		wp := os.Getenv("GO_WORKING_PATH")

		if wp != "" {
			wp = filepath.ToSlash(path.Clean(wp))
			return path.Join(wp, offset)
		}
	}

	return path.Join(Dir(os.Args[0]), offset)
}

// Dir 获取路径所在文件夹
func Dir(path string) string {
	dir, err := filepath.Abs(filepath.Dir(path))

	if nil != err {
		return ""
	}

	return dir
}

// PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Mkdir 无错误创建文件夹
func Mkdir(path string) error {
	exists, err := PathExists(path)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	os.Mkdir(path, os.ModePerm)

	return nil
}
