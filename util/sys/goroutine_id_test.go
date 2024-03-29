package sys

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSmokeGoroutineId(t *testing.T) {
	Convey("检查在同一routine的多级方法获取的routineid是否一致", t, func(c C) {
		baseId := CurGoroutineID()
		block := make(chan bool, 1)
		go func() {
			mainId := CurGoroutineID()
			c.So(baseId, ShouldNotEqual, mainId)
			func() {
				dept1 := CurGoroutineID()
				c.So(mainId, ShouldEqual, dept1)
				func() {
					dept2 := CurGoroutineID()
					c.So(dept1, ShouldEqual, dept2)
					close(block)
				}()
			}()
		}()
		<-block
	})
}

func TestSmokeRunGoroutine(t *testing.T) {
	Convey("检查RunRoutine是否可以正确设置父子关系", t, func(c C) {
		parentId := CurGoroutineID()
		var childId uint64
		block := make(chan bool, 1)
		RunRoutine(func() {
			childId = CurGoroutineID()
			parentIdFound := FindRootRoutineId()
			c.So(parentId, ShouldEqual, parentIdFound)
			close(block)
		})
		<-block
		_, found := routineMap[childId]
		c.So(found, ShouldBeFalse)
	})
}
