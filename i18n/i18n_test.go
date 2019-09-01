package i18n

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLocalizeWithLocale(t *testing.T) {
	Convey("检查内置的直接输出Provider是否正常工作", t, func() {
		locale := "zh-CN"
		catalog := "catalog"
		key := "key"
		pattern := "this is a test for %s, %s"
		args := []interface{}{
			"i18n",
			"hello",
		}

		want := "this is a test for i18n, hello"

		So(LocalizeWithLocale(locale, catalog, key, pattern, args...), ShouldEqual, want)
	})
}
