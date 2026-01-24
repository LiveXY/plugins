package plugin

import (
	"io"
	"os"
	"path"

	"github.com/livexy/plugin/osser"
)

// NewLocalFS 创建一个新的本地文件系统实例
func NewLocalFS(cfg osser.OSSConfig) (osser.OSSer, error) {
	fs := &LocalFS{}
	err := fs.setup(cfg)
	return fs, err
}

type LocalFS struct {
	data string
}

func (o *LocalFS) setup(cfg osser.OSSConfig) error {
	o.data = cfg.Endpoint
	return nil
}

// Upload 上传文件（将本地文件移动到对象存储目录）
func (o *LocalFS) Upload(objname, localfile string) error {
	objfile := path.Join(o.data, objname)
	// 尝试直接重命名（原子操作，性能最好）
	err := os.Rename(localfile, objfile)
	if err == nil {
		return nil
	}

	// 如果重命名失败（可能是跨设备移动），尝试 复制+删除
	// 这里不区分具体错误类型，作为一种通用的回退策略
	return o.moveFileByCopy(localfile, objfile)
}

// moveFileByCopy 通过复制后删除原文件的方式移动文件
// 用于解决跨设备移动文件失败的问题
func (o *LocalFS) moveFileByCopy(source, dest string) error {
	inputFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	// 确保文件写入并关闭，再进行删除操作
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return err
	}

	// 关闭文件确保数据刷盘（虽然defer会关，但这里显式关闭更好处理错误）
	if err := outputFile.Close(); err != nil {
		return err
	}
	// 关闭输入文件以便删除
	if err := inputFile.Close(); err != nil {
		return err
	}

	// 删除原文件
	return os.Remove(source)
}

// Delete 删除指定的文件
func (o *LocalFS) Delete(objname string) error {
	objfile := path.Join(o.data, objname)
	err := os.Remove(objfile)
	if err != nil {
		return err
	}
	return nil
}
