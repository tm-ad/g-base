package util_test

import (
	base642 "encoding/base64"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/tm-ad/g-base/util"
	"testing"
)

func Test_Smoke_SubStr(t *testing.T) {
	Convey("验证substr是否能正确拆分", t, func() {
		ori := `马云啊is a human～24#mayun是个人`
		expected := []string{
			`马云啊i`,
			`s a `,
			`huma`,
			`n～24`,
			`#may`,
			`un是个`,
			`人`,
		}

		So(SubStr(ori, 0, 4), ShouldEqual, expected[0])
		So(SubStr(ori, 4, 4), ShouldEqual, expected[1])
		So(SubStr(ori, 8, 4), ShouldEqual, expected[2])
		So(SubStr(ori, 12, 4), ShouldEqual, expected[3])
		So(SubStr(ori, 16, 4), ShouldEqual, expected[4])
		So(SubStr(ori, 20, 4), ShouldEqual, expected[5])
		So(SubStr(ori, 24, 4), ShouldEqual, expected[6])
	})
}

func Test_SubStr_IsOk_with_empty(t *testing.T) {
	Convey("验证substr是否能对空字符串正确处理", t, func() {
		So(SubStr("", 0, 4), ShouldEqual, "")
		So(SubStr("", 2, 200), ShouldEqual, "")
		So(SubStr("n", 2, 4), ShouldEqual, "")
		So(SubStr("n", 2, -1), ShouldEqual, "")
		So(SubStr("n", -1, 1), ShouldEqual, "")
	})
}

func Test_WebSafeBase64(t *testing.T) {
	Convey("验证websafebase64的互转", t, func() {
		base64 := base642.StdEncoding.EncodeToString([]byte(long_text_input))
		fmt.Println(base64)
		wsBase64 := ToWebSafeBase64(base64)
		fmt.Println(wsBase64)
		fromBase64 := FromWebSafeBase64(wsBase64)
		fmt.Println(fromBase64)
		So(fromBase64, ShouldEqual, base64)
	})
}

var long_text_input = `golang的一种特殊的加密解密算法AES/ECB/PKCS5,但是算法并没有包含在标准库中,经...
https://www.cnblogs.com/lavin/...  - 百度快照
golang base64解码碰到的坑 - u014270740的专栏 - CSDN博客
2019年6月6日 -    url safe 将+/字符串转化成_-    no padding is add  末尾...golang操作base64例子。简单的编码、解码功能。 x 下载 base64的图片编码转...
CSDN技术社区 - 百度快照
关于WebSafeBase64的加密和解密,求解答。。。 - Golang 中国
enc_price = WebSafeBase64Decode(final_message)(iv, p, sig) = dec_...看了java的源码,不知道如何用golang实现。如果能提供下思路的话更好,非常感谢...
www.golangtc.com/t/56f...  - 百度快照
golang base64加密与解密 - Go语言中文网 - Golang中文社区
2016年4月27日 - 查看原文:golang base64加密与解密入群交流(和以上内容无关):Go中文网 QQ 交流群:798786647 或加微信入微信群:274768166 备注:入群;关注公众号:Go语言...
https://studygolang.com/articl...  - 百度快照
Golang实现的Base64加密 - Go语言中文网 - Golang中文社区
2016年2月5日 - 也欢迎加入知识星球 Go粉丝们(免费)base64加密是我们经常看到的一种加密方法,比如ESMTP的验证过程和二进制文件的网际传输等都会用到这种编码。 base64...
https://studygolang.com/articl...  - 百度快照
[Golang] base64加密与解密 - Kirai - 博客园
2016年8月7日 - [Golang] base64加密与解密 首先解释以下什么是base64(来自百度百科): Base64是网络上最常见的用于传输8Bit字节代码的编码方式之一,大家可以查看RFC20...
https://www.cnblogs.com/kirai/...  - 百度快照
golang base64函数基本用法 - 简书
2018年4月10日 - golang base64函数基本用法 base64主要两个函数编码和解码。 编码:把一段字节buffer翻译成base64格式字符串。func EncodeToString...
简书社区 - 百度快照
go语言base64加密解密的方法_Golang_脚本之家
2015年3月2日 - 这篇文章主要介绍了go语言base64加密解密的方法,实例分析了Go语言base64加密解密的技巧,需要的朋友可以参考下
脚本之家 - 百度快照
golang基础学习-base64使用 - Keil - SegmentFault 思否
2019年8月21日 - 在近期的项目开发中对图片进行base64编码,简单使用了golang的base64包。 1.使用方法 1.1 引入包 import "encoding/base64" 1.2 base64使用 这里所有的...
https://segmentfault.com/a/119...  - 百度快照
1
2
3
4
5
6
7
8
9
10
下一页>`
