package plugin

import (
	"os"
	"path"

	"github.com/livexy/plugins/plugin/osser"
)

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
func (o *LocalFS) Upload(objname, localfile string) error {
	objfile := path.Join(o.data, objname)
	err := os.Rename(localfile, objfile)
	if err != nil {
		return err
	}
	return nil
}
func (o *LocalFS) Delete(objname string) error {
	objfile := path.Join(o.data, objname)
	err := os.Remove(objfile)
	if err != nil {
		return err
	}
	return nil
}
