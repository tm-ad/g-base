package sys

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
)

var littleBuf = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 64)
		return &buf
	},
}

// parseUintBytes is like strconv.ParseUint, but using a []byte.
func parseUintBytes(s []byte, base int, bitSize int) (n uint64, err error) {
	var cutoff, maxVal uint64

	if bitSize == 0 {
		bitSize = int(strconv.IntSize)
	}
	s0 := s
	switch {
	case len(s) < 1:
		err = strconv.ErrSyntax
		goto Error

	case 2 <= base && base <= 36:
		// valid base; nothing to do

	case base == 0:
		// Look for octal, hex prefix.
		switch {
		case s[0] == '0' && len(s) > 1 && (s[1] == 'x' || s[1] == 'X'):
			base = 16
			s = s[2:]
			if len(s) < 1 {
				err = strconv.ErrSyntax
				goto Error
			}
		case s[0] == '0':
			base = 8
		default:
			base = 10
		}

	default:
		err = errors.New("invalid base " + strconv.Itoa(base))
		goto Error
	}

	n = 0
	cutoff = cutoff64(base)
	maxVal = 1<<uint(bitSize) - 1

	for i := 0; i < len(s); i++ {
		var v byte
		d := s[i]
		switch {
		case '0' <= d && d <= '9':
			v = d - '0'
		case 'a' <= d && d <= 'z':
			v = d - 'a' + 10
		case 'A' <= d && d <= 'Z':
			v = d - 'A' + 10
		default:
			n = 0
			err = strconv.ErrSyntax
			goto Error
		}
		if int(v) >= base {
			n = 0
			err = strconv.ErrSyntax
			goto Error
		}

		if n >= cutoff {
			// n*base overflows
			n = 1<<64 - 1
			err = strconv.ErrRange
			goto Error
		}
		n *= uint64(base)

		n1 := n + uint64(v)
		if n1 < n || n1 > maxVal {
			// n+v overflows
			n = 1<<64 - 1
			err = strconv.ErrRange
			goto Error
		}
		n = n1
	}

	return n, nil

Error:
	return n, &strconv.NumError{Func: "ParseUint", Num: string(s0), Err: err}
}

// Return the first number n such that n*base >= 1<<64.
func cutoff64(base int) uint64 {
	if base < 2 {
		return 0
	}
	return (1<<64-1)/uint64(base) + 1
}

var goroutineSpace = []byte("goroutine ")

func CurGoroutineID() uint64 {
	bp := littleBuf.Get().(*[]byte)
	defer littleBuf.Put(bp)
	b := *bp
	b = b[:runtime.Stack(b, false)]
	// Parse the 4707 out of "goroutine 4707 ["
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}
	b = b[:i]
	n, err := parseUintBytes(b, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %q: %v", b, err))
	}
	return n
}

var routineMapLocker sync.Mutex
var routineMap = make(map[uint64]uint64)
var routineDirectMap = make(map[uint64]uint64)

// FindRootRoutineId 寻找最顶级父routine的id
func FindRootRoutineId() uint64 {
	curRoutineId := CurGoroutineID()
	var parentId uint64
	var found bool
	parentId, found = routineMap[curRoutineId]
	// 特例处理首次就没找到，则认为此 routine 为顶级
	if !found {
		return curRoutineId
	}

	// 增加从缓存中获取
	cachedParentId, cacheFound := routineDirectMap[curRoutineId]
	if cacheFound {
		// 1/20 几率重新计算
		// TODO: 需要尝试压测
		rnd := rand.Intn(20)
		if rnd == 5 {
			cacheFound = false
		}
	}
	if cacheFound {
		return cachedParentId
	}

	for {
		if !found {
			cachedParentId = parentId
			break
		}
		parentId, found = routineMap[parentId]
	}

	return cachedParentId
}

// RunRoutine 提供一个封装的 go goroutine 启动方法
// 	该方法会在运行环境中建立 parent goroutine id 和 current goroutine id 的 map 映射
func RunRoutine(rFunc func()) {
	parentRoutineId := CurGoroutineID()
	go func() {
		curRoutineId := CurGoroutineID()
		func() {
			// 检查是否存在当前的routine id
			defer routineMapLocker.Unlock()
			routineMapLocker.Lock()
			preParentRoutineId, found := routineMap[curRoutineId]
			fmt.Println(preParentRoutineId)
			if !found || preParentRoutineId != parentRoutineId {
				routineMap[curRoutineId] = parentRoutineId
			}
			// routineMapLocker.Unlock()
		}()

		// 退出后推出
		defer func() {
			defer routineMapLocker.Unlock()
			routineMapLocker.Lock()
			delete(routineMap, curRoutineId)
			delete(routineDirectMap, curRoutineId)
			// routineMapLocker.Unlock()
		}()
		// 执行
		rFunc()
	}()
}
