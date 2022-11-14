package plugin

import (
	"path"
	gplugin "plugin"
	"strings"

	cmap "github.com/orcaman/concurrent-map/v2"
)

var pluginmap = cmap.New[gplugin.Symbol]()

func Load[T any](plugpath, name string) T {
	if len(plugpath) == 0 {
		plugpath = "./plugins"
	}
	name = strings.ToLower(name)
	obj, ok := pluginmap.Get(name)
	if ok {
		return obj.(T)
	}
	so := path.Join(plugpath, name+".so")
	plug, err := gplugin.Open(so)
	if err != nil {
		panic("插件 " + name + " 加载失败！" + err.Error())
	}
	obj, err = plug.Lookup("Plugin")
	if err != nil {
		panic("插件 " + name + " 加载，无Plugin符号！" + err.Error())
	}
	t, ok := obj.(T)
	if !ok {
		panic("插件 " + name + " 转换T失败！")
	}
	pluginmap.Set(name, t)
	return t
}

func Reload[T any](plugpath, name string) T {
	pluginmap.Remove(name)
	return Load[T](plugpath, name)
}
