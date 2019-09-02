// Package locale 用于提供本地语言包处理接口定义和快捷调用入口
package locale

import "fmt"

// Dash 是保留的语言包目录名称
const Dash = "-"

// _lpk 当前正在使用的语言包
var _lpk ILPack

// ILPack 定义一个语言包的行为接口
type ILPack interface {
	// Localize 用于输出当前语言包下指定的目录中的指定key的语言文本
	Localize(catalog, key, reserved string, args ...interface{}) string
}

// L 返回指定预演目录中指定key的语言文本，当无法找到指定文本时使用reserved作为代替文本
func L(catalog, key, reserved string, args ...interface{}) string {
	if _lpk == nil {
		if len(args) == 0 {
			return reserved
		} else {
			return fmt.Sprintf(reserved, args...)
		}
	}

	return _lpk.Localize(catalog, key, reserved, args...)
}

// SetLPack 设置当前正在使用的语言包
func SetLPack(pack ILPack) {
	if pack == nil {
		panic("locale resource pack must be specified")
	}
	_lpk = pack
}
