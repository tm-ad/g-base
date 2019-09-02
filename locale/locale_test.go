package locale

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSmokeLWithArgs(t *testing.T) {
	Convey("检查内置的语言快捷调用是否生效，有参数", t, func() {
		pattern := "this is a test for %s, %s"
		args := []interface{}{
			"locale resource pack",
			"hello",
		}

		want := "this is a test for locale resource pack, hello"

		So(L("some-catalog", "some-key", pattern, args...), ShouldEqual, want)
	})
}

func TestSmokeLWithNoArgs(t *testing.T) {
	Convey("检查内置的语言快捷调用是否生效，无参数", t, func() {
		pattern := "this is a test for"
		want := "this is a test for"

		So(L("some-catalog", "some-key", pattern), ShouldEqual, want)
	})
}

type myLPack struct {
}

func (m *myLPack) Localize(catalog, key, reserved string, args ...interface{}) string {
	return fmt.Sprintf("hello %s", args...)
}

func TestSmokeLPack(t *testing.T) {
	Convey("冒烟测试对LPack的调用实现", t, func() {
		lp := &myLPack{}
		SetLPack(lp)

		want := "hello world"

		So(L("some-catalog", "some-key", "%s", "world"), ShouldEqual, want)
	})
}

func TestSmokeSetLPackOK(t *testing.T) {
	Convey("冒烟测试对SetLPack的有效设置", t, func() {
		lp := &myLPack{}
		So(func() { SetLPack(lp) }, ShouldNotPanic)
	})
}

func TestSmokeSetLPackNotOK(t *testing.T) {
	Convey("冒烟测试对SetLPack的无效设置，设置为空", t, func() {
		So(func() { SetLPack(nil) }, ShouldPanicWith, "locale resource pack must be specified")
	})
}
