package rsa_test

import (
	"encoding/base64"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/tm-ad/g-base/crypto/rsa"
	"testing"
)

func Test_GenRsaKey(t *testing.T) {
	Convey("验证是否能正确生成RSA密钥", t, func() {
		pem, key, err := GenRsaKey(1024)
		fmt.Println(string(pem))
		fmt.Println(string(key))
		So(err, ShouldBeNil)
		// 获取 key
		publicKey, err := GetPublicKey(pem)
		So(err, ShouldBeNil)
		So(publicKey.Size(), ShouldEqual, 1024/8)

		privateKey, err := GetPrivateKey(key)
		So(err, ShouldBeNil)
		So(privateKey.Size(), ShouldEqual, 1024/8)
	})
}

func Test_PrivateEncrypt_To_PublicDecrypt(t *testing.T) {
	Convey("验证私钥加密公钥解密的正确性", t, func() {
		// 获取key
		pem, key, err := GenRsaKey(1024)
		fmt.Println(string(pem))
		fmt.Println(string(key))
		So(err, ShouldBeNil)
		// 获取 key
		publicKey, err := GetPublicKey(pem)
		So(err, ShouldBeNil)
		So(publicKey.Size(), ShouldEqual, 1024/8)

		privateKey, err := GetPrivateKey(key)
		So(err, ShouldBeNil)
		So(privateKey.Size(), ShouldEqual, 1024/8)

		ori := input
		pwd, err := PrivateEncrypt([]byte(ori), privateKey)
		So(err, ShouldBeNil)
		fmt.Println(string(pwd))
		d, err := PublicDecrypt(pwd, publicKey)
		So(ori, ShouldEqual, string(d))
		fmt.Println(string(d))
	})
}

func Test_PrivateEncrypt_To_PublicDecrypt_longtext(t *testing.T) {
	Convey("验证私钥加密公钥解密的正确性", t, func() {
		// 获取key
		pem, key, err := GenRsaKey(1024)
		fmt.Println(string(pem))
		fmt.Println(string(key))
		So(err, ShouldBeNil)
		// 获取 key
		publicKey, err := GetPublicKey(pem)
		So(err, ShouldBeNil)
		So(publicKey.Size(), ShouldEqual, 1024/8)

		privateKey, err := GetPrivateKey(key)
		So(err, ShouldBeNil)
		So(privateKey.Size(), ShouldEqual, 1024/8)

		ori := long_text_input
		pwd, err := PrivateEncrypt([]byte(ori), privateKey)
		So(err, ShouldBeNil)
		fmt.Println(string(pwd))
		d, err := PublicDecrypt(pwd, publicKey)
		So(ori, ShouldEqual, string(d))
		fmt.Println(string(d))
	})
}

func Test_PublicEncrypt_To_PrivateDecrypt_longtext(t *testing.T) {
	Convey("验证公钥加密私钥解密的正确性", t, func() {
		// 获取key
		pem, key, err := GenRsaKey(1024)
		fmt.Println(string(pem))
		fmt.Println(string(key))
		So(err, ShouldBeNil)
		// 获取 key
		publicKey, err := GetPublicKey(pem)
		So(err, ShouldBeNil)
		So(publicKey.Size(), ShouldEqual, 1024/8)

		privateKey, err := GetPrivateKey(key)
		So(err, ShouldBeNil)
		So(privateKey.Size(), ShouldEqual, 1024/8)

		ori := long_text_input
		pwd, err := PublicEncrypt([]byte(ori), publicKey)
		So(err, ShouldBeNil)
		fmt.Println(base64.StdEncoding.EncodeToString(pwd))
		d, err := PrivateDecrypt(pwd, privateKey)
		So(ori, ShouldEqual, string(d))
		fmt.Println(string(d))
	})
}

func Test_Sha1_Sign(t *testing.T) {
	Convey("验证啥sha1签名及验证", t, func() {
		// 获取key
		pem, key, err := GenRsaKey(1024)
		So(err, ShouldBeNil)
		// 获取 key
		publicKey, err := GetPublicKey(pem)
		So(err, ShouldBeNil)

		privateKey, err := GetPrivateKey(key)
		So(err, ShouldBeNil)

		ori := long_text_input
		pwd, err := PublicEncrypt([]byte(ori), publicKey)
		So(err, ShouldBeNil)
		sign, err := SignSha1WithRsa(pwd, privateKey)
		// fmt.Println(base64.StdEncoding.EncodeToString(sign))
		// err = VerifySignSha1WithRsa([]byte(ori), sign, publicKey)
		err = VerifySignSha1WithRsa(pwd, sign, publicKey)
		So(err, ShouldBeNil)
	})
}

func Test_Sha256_Sign(t *testing.T) {
	Convey("验证啥sha1签名及验证", t, func() {
		// 获取key
		pem, key, err := GenRsaKey(1024)
		So(err, ShouldBeNil)
		// 获取 key
		publicKey, err := GetPublicKey(pem)
		So(err, ShouldBeNil)

		privateKey, err := GetPrivateKey(key)
		So(err, ShouldBeNil)

		ori := long_text_input
		pwd, err := PublicEncrypt([]byte(ori), publicKey)
		So(err, ShouldBeNil)
		sign, err := SignSha256WithRsa(pwd, privateKey)
		// fmt.Println(base64.StdEncoding.EncodeToString(sign))
		// err = VerifySignSha1WithRsa([]byte(ori), sign, publicKey)
		err = VerifySignSha256WithRsa(pwd, sign, publicKey)
		So(err, ShouldBeNil)
	})
}

var input = `马云啊is a human～24#mayun是个人`
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
