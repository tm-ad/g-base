package strftime_test

import (
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/tm-ad/g-base/util/strftime"
	"testing"
	"time"
)

func TestStrftime_FormatString(t *testing.T) {
	Convey("简单验证formatstring是否正确", t, func() {
		const expected = `apm-test/logs/apm.log.01000101`
		p, _ := New("apm-test/logs/apm.log.%Y%m%d")
		dt := time.Date(100, 1, 1, 1, 0, 0, 0, time.UTC)
		So(p.FormatString(dt), ShouldEqual, expected)
	})
}
