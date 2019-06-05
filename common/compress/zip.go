package compress

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Compress(filePath, destFile string) error {
	filePath = filepath.Clean(filePath)
	if strings.HasSuffix(filePath, "/") {
		filePath = filePath[0 : len(filePath)-1]
	}

	// file write
	fw, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer fw.Close()

	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	baseDir := filepath.Dir(filePath)
	// 递归遍历文件夹并压缩
	err = writeFile2Gzip(baseDir, "", fileInfo, tw)
	if err != nil {
		return err
	}
	return nil
}

func getAbBaseDir(abBaseDir, fileName string) string {
	if len(strings.TrimSpace(abBaseDir)) == 0 {
		return fileName
	}
	return abBaseDir + "/" + fileName
}
func writeFile2Gzip(realBaseDir, abBaseDir string, fileInfo os.FileInfo, tw *tar.Writer) error {
	// 逃过文件夹, 我这里就不递归了
	if fileInfo.IsDir() {
		// 打开文件夹
		dir, err := os.Open(realBaseDir + "/" + fileInfo.Name())
		if err != nil {
			return err
		}
		defer dir.Close()

		// 读取文件列表
		fis, err := dir.Readdir(0)
		if err != nil {
			return err
		}

		for _, fi := range fis {
			err = writeFile2Gzip(dir.Name(), getAbBaseDir(abBaseDir, fileInfo.Name()), fi, tw)
			if err != nil {
				return err
			}
		}
		return nil
	}

	// 写信息头
	err := writeHeader(abBaseDir, fileInfo, tw)
	if err != nil {
		return err
	}

	// 打开文件
	fr, err := os.Open(realBaseDir + "/" + fileInfo.Name())
	if err != nil {
		return err
	}
	defer fr.Close()

	// 写文件
	_, err = io.Copy(tw, fr)
	if err != nil {
		return err
	}
	return nil
}

func writeHeader(baseDir string, fileInfo os.FileInfo, tw *tar.Writer) error {
	h := new(tar.Header)
	h.Name = getAbBaseDir(baseDir, fileInfo.Name())
	h.Size = fileInfo.Size()
	h.Mode = int64(fileInfo.Mode())
	h.ModTime = fileInfo.ModTime()

	return tw.WriteHeader(h)
}

func Decompress(gzipFile, destPath string) error {
	// file read
	fr, err := os.Open(gzipFile)
	if err != nil {
		return err
	}
	defer fr.Close()

	// gzip read
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return err
	}
	defer gr.Close()

	// tar read
	tr := tar.NewReader(gr)

	// 读取文件
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if h.FileInfo().IsDir() {
			dir := filepath.Dir(destPath + "/" + h.Name)
			os.MkdirAll(dir, h.FileInfo().Mode())
		} else {
			// 打开文件
			fw, err := os.OpenFile(destPath+"/"+h.Name, os.O_CREATE|os.O_WRONLY, h.FileInfo().Mode() /*os.FileMode(h.Mode)*/)
			if err != nil {
				return err
			}
			defer fw.Close()

			// 写文件
			_, err = io.Copy(fw, tr)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
