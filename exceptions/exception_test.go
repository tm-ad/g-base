package exceptions

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIsException_built_in_error_is_not_a_exception(t *testing.T) {
	Convey("常规error不是Exception类型", t, func() {
		e := errors.New("我是一个异常")

		So(IsException(e), ShouldBeFalse)
	})
}

type myExp struct {
}

func (e *myExp) Error() string {
	return ""
}

func (e *myExp) Message() string {
	return ""
}

func (e *myExp) Code() int {
	return 0
}

func (e *myExp) Causes() []error {
	return nil
}

func TestIsException_exception_is_exception(t *testing.T) {
	Convey("自定义实现的exception检测为exception", t, func() {
		e := &myExp{}

		So(IsException(e), ShouldBeTrue)
	})
}

func TestNew_exception_message_common_code_are_wanted(t *testing.T) {
	Convey("基本New方法是否可以正确创建文字相同和code为常规code的错误", t, func() {
		msg := "exception123"
		e := New(msg)

		So(e.Message(), ShouldEqual, msg)
		So(e.Code(), ShouldEqual, CommonExceptionCode)
	})
}

func TestCode_exception_message_code_are_wanted(t *testing.T) {
	Convey("基本Code方法是否可以正确创建文字相同和code的错误", t, func() {
		msg := "exception123"
		c := 7788
		e := Code(c, msg)

		So(e.Message(), ShouldEqual, msg)
		So(e.Code(), ShouldEqual, c)
	})
}

func TestCode_exception_error_stack_are_wanted(t *testing.T) {
	Convey("检查指定cause是否为预期", t, func() {
		err1 := "this is 1st error"
		err2 := "this is 2nd error"

		msg := "exception123"
		c := 7788
		e := Code(c, msg, errors.New(err1), errors.New(err2))

		So(e.Message(), ShouldEqual, msg)
		So(e.Code(), ShouldEqual, c)
		So(len(e.Causes()), ShouldEqual, 2)
		So(e.Causes()[0].Error(), ShouldEqual, err1)
		So(e.Causes()[1].Error(), ShouldEqual, err2)
	})
}

func TestNew_Error_out_json(t *testing.T) {
	Convey("检查指定json字符串错误信息输出是否为预期", t, func() {
		err1 := "this is 1st error"
		err2 := "this is 2nd error"

		msg := "exception123"
		c := 7788
		e := Code(c, msg, errors.New(err1), errors.New(err2))

		So(e.Error(), ShouldEqual, "{\"code\":7788,\"message\":\"exception123\",\"causes\":[\"this is 1st error\",\"this is 2nd error\"]}")
	})
}
