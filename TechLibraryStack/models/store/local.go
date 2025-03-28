package store

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/TruthHun/gotil/filetil"
)

//本地存储
type Local struct{}

//判断文件对象是否存在
//@param                object              文件对象
//@return               err                 错误，nil表示文件存在，否则表示文件不存在，并告知错误信息
func (this *Local) IsObjectExist(object string) (err error) {
	_, err = os.Stat(object)
	return
}

//文件存储[如果是图片文件，不要使用gzip压缩，否则在使用阿里云OSS自带的图片处理功能无法处理图片]
//@param            tmpfile          临时文件
//@param            save             存储文件，不建议与临时文件相同，特别是IsDel参数值为true的时候
//@param            IsDel            文件上传后，是否删除临时文件
func (this *Local) MoveToStore(tmpfile, save string) (err error) {
	save = strings.TrimLeft(save, "/")
	//"./a.png"与"a.png"是相同路径
	if strings.HasPrefix(tmpfile, "./") || strings.HasPrefix(save, "./") {
		tmpfile = strings.TrimPrefix(tmpfile, "./")
		save = strings.TrimPrefix(save, "./")
	}
	if strings.ToLower(tmpfile) != strings.ToLower(save) { //不是相同文件路径

		os.MkdirAll(filepath.Dir(save), os.ModePerm)
		// 不使用rename，因为在docker中会挂载外部卷，导致错误
		// 见https://gocn.vip/article/178
		if b, err := ioutil.ReadFile(tmpfile); err == nil {
			ioutil.WriteFile(save, b, os.ModePerm)
		}
		os.Remove(tmpfile)
	}
	return
}

//从OSS中删除文件
//@param           object                     文件对象
//@param           IsPreview                  是否是预览的OSS
func (this *Local) DelFiles(object ...string) error {
	for _, file := range object {
		os.Remove(strings.TrimLeft(file, "/"))
	}
	return nil
}

//删除文件夹
func (this *Local) DelFromFolder(folder string) (err error) {
	return os.RemoveAll(folder)
}

func (this *Local) CopyFile(sourceFile, targetFile string) (err error) {

	return
}

// 拷贝文件夹
func (this *Local) CopyDir(sourceDir, targetDir string) (err error) {
	var fl []filetil.FileList

	sourceDir = strings.TrimLeft(sourceDir, "./")
	targetDir = strings.TrimLeft(targetDir, "./")

	if fl, err = filetil.ScanFiles(sourceDir); err != nil {
		return
	}
	if len(fl) == 0 {
		return
	}

	for _, f := range fl {
		if f.IsDir {
			continue
		}

		// 不能用defer，因为在文件多的时候，会打开跟多文件句柄
		targetFile := strings.ReplaceAll(filepath.Join(targetDir, strings.TrimPrefix(f.Path, sourceDir)), "\\", "/")
		src, errOpen := os.Open(f.Path)
		if errOpen != nil {
			return errOpen
		}

		os.MkdirAll(filepath.Dir(targetFile), os.ModePerm)
		target, errOpen := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0644)
		if errOpen != nil {
			src.Close()
			return errOpen
		}

		if _, err = io.Copy(target, src); err != nil {
			target.Close()
			src.Close()
			return
		}
		target.Close()
		src.Close()
	}
	return
}
