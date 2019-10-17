package log

import (
	"errors"
	"fmt"
	"github.com/tm-ad/g-base/util/fs"
	"github.com/tm-ad/g-base/util/option"
	"github.com/tm-ad/g-base/util/strftime"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// region rotate option

const (
	optkeyMaxAge       = "max-age"
	optkeyRotationTime = "rotation-time"
)

// WithMaxAge creates a new Option that sets the
// max age of a log file before it gets purged from
// the file system.
func WithMaxAge(d time.Duration) Option {
	return option.New(optkeyMaxAge, d)
}

// WithRotationTime creates a new Option that sets the
// time between rotation.
func WithRotationTime(d time.Duration) Option {
	return option.New(optkeyRotationTime, d)
}

func defaultRotatePattern(pattern string) string {
	if pattern == "" {
		return "%Y-%m-%d"
	}

	return pattern
}

func defaultRotationTime(rotationTime time.Duration) time.Duration {
	if int(rotationTime) <= 0 {
		return 24 * time.Hour
	}
	return rotationTime
}

func defaultMaxAge(maxAge time.Duration) time.Duration {
	if int(maxAge) <= 0 {
		return 7 * 24 * time.Hour
	}
	return maxAge
}

func defaultName(lname string) string {
	if lname == "" {
		return "log"
	}

	return lname
}

func defaultLevel(llvl string) string {
	if llvl == "" {
		return "info"
	}

	return llvl
}

// end region option

// Option is used to pass optional arguments to
// the RotateLogs constructor
type Option interface {
	Name() string
	Value() interface{}
}

// Clock is the interface used by the RotateLogs
// object to determine the current time
type Clock interface {
	Now() time.Time
}
type clockFn func() time.Time

// UTC is an object satisfying the Clock interface, which
// returns the current time in UTC
var UTC = clockFn(func() time.Time { return time.Now().UTC() })

func (c clockFn) Now() time.Time {
	return c()
}

// Local is an object satisfying the Clock interface, which
// returns the current time in the local timezone
var Local = clockFn(time.Now)

var patternConversionRegexps = []*regexp.Regexp{
	regexp.MustCompile(`%[%+A-Za-z]`),
	regexp.MustCompile(`\*+`),
}

type cleanupGuard struct {
	enable bool
	fn     func()
	mutex  sync.Mutex
}

func (g *cleanupGuard) Enable() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.enable = true
}
func (g *cleanupGuard) Run() {
	g.fn()
}

// NewRotateFileLog 创建一个根据时间周期切分的文件日志
// 	root: 日志存储的根路径，不能为空
//	name: 日志的主文件名，可认为为 prefix
//	lvl: 日志的记录的等级
//	pattern: 文件名的格式化字符串，如 %Y-%m-%d
//	rotationTime: 文件切分的间隔
//	maxAge: 文件最大的时间有效期
func NewRotateFileLog(root, name, lvl, pattern string, rotationTime, maxAge time.Duration) (*Logger, error) {
	l := New()
	l.SetLevel(defaultLevel(lvl))

	// 检查并创建日志根目录
	if err := fs.Mkdir(root); err != nil {
		TipInDevelopment(fmt.Sprintf(`log root initialize failed: %v \n`, err))
	}
	baseLogName := path.Join(root, defaultName(name))
	// 构建 rotate log file
	writer, err := NewRotateWriter(
		baseLogName+defaultRotatePattern(pattern)+".log",
		WithMaxAge(defaultMaxAge(maxAge)),
		WithRotationTime(defaultRotationTime(rotationTime)),
	)

	if err != nil {
		return nil, err
	}

	l.AddOutput(writer)

	return l, nil
}

// RotateWriter represents a log file that gets
// automatically rotated as you write to it.
type RotateWriter struct {
	clock        Clock
	curFn        string
	curBaseFn    string
	globPattern  string
	generation   int
	maxAge       time.Duration
	mutex        sync.RWMutex
	outFh        *os.File
	pattern      *strftime.Strftime
	rotationTime time.Duration
	forceNewFile bool
}

// NewRotateWriter creates a new RotateLogs object. A log filename pattern
// must be passed. Optional `Option` parameters may be passed
func NewRotateWriter(p string, options ...Option) (*RotateWriter, error) {
	globPattern := p
	for _, re := range patternConversionRegexps {
		globPattern = re.ReplaceAllString(globPattern, "*")
	}

	pattern, err := strftime.New(p)
	if err != nil {
		return nil, errors.New(`invalid strftime pattern`)
	}

	var clock Clock = Local
	rotationTime := 24 * time.Hour
	var maxAge time.Duration
	var forceNewFile bool

	for _, o := range options {
		switch o.Name() {
		case optkeyMaxAge:
			maxAge = o.Value().(time.Duration)
			if maxAge < 0 {
				maxAge = 0
			}
		case optkeyRotationTime:
			rotationTime = o.Value().(time.Duration)
			if rotationTime < 0 {
				rotationTime = 0
			}
		}
	}

	return &RotateWriter{
		clock:        clock,
		globPattern:  globPattern,
		maxAge:       maxAge,
		pattern:      pattern,
		rotationTime: rotationTime,
		forceNewFile: forceNewFile,
	}, nil
}

// Write satisfies the io.Writer interface. It writes to the
// appropriate file handle that is currently being used.
// If we have reached rotation time, the target file gets
// automatically rotated, and also purged if necessary.
func (rl *RotateWriter) Write(p []byte) (n int, err error) {
	// Guard against concurrent writes
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	out, err := rl.getWriter_nolock(false, false)
	if err != nil {
		return 0, errors.New(`failed to acquite target io.Writer`)
	}

	return out.Write(p)
}

func (rl *RotateWriter) genFilename() string {
	now := rl.clock.Now()

	// XXX HACK: Truncate only happens in UTC semantics, apparently.
	// observed values for truncating given time with 86400 secs:
	//
	// before truncation: 2018/06/01 03:54:54 2018-06-01T03:18:00+09:00
	// after  truncation: 2018/06/01 03:54:54 2018-05-31T09:00:00+09:00
	//
	// This is really annoying when we want to truncate in local time
	// so we hack: we take the apparent local time in the local zone,
	// and pretend that it's in UTC. do our math, and put it back to
	// the local zone
	var base time.Time
	if now.Location() != time.UTC {
		base = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC)
		base = base.Truncate(time.Duration(rl.rotationTime))
		base = time.Date(base.Year(), base.Month(), base.Day(), base.Hour(), base.Minute(), base.Second(), base.Nanosecond(), base.Location())
	} else {
		base = now.Truncate(time.Duration(rl.rotationTime))
	}
	return rl.pattern.FormatString(base)
}

// must be locked during this operation
func (rl *RotateWriter) getWriter_nolock(bailOnRotateFail, useGenerationalNames bool) (io.Writer, error) {
	generation := rl.generation
	// previousFn := rl.curFn
	// This filename contains the name of the "NEW" filename
	// to log to, which may be newer than rl.currentFilename
	baseFn := rl.genFilename()
	filename := baseFn
	var forceNewFile bool
	if baseFn != rl.curBaseFn {
		generation = 0
		// even though this is the first write after calling New(),
		// check if a new file needs to be created
		if rl.forceNewFile {
			forceNewFile = true
		}
	} else {
		if !useGenerationalNames {
			// nothing to do
			return rl.outFh, nil
		}
		forceNewFile = true
		generation++
	}
	if forceNewFile {
		// A new file has been requested. Instead of just using the
		// regular strftime pattern, we create a new file name using
		// generational names such as "foo.1", "foo.2", "foo.3", etc
		var name string
		for {
			if generation == 0 {
				name = filename
			} else {
				name = fmt.Sprintf("%s.%d", filename, generation)
			}
			if _, err := os.Stat(name); err != nil {
				filename = name
				break
			}
			generation++
		}
	}
	// make sure the dir is existed, eg:
	// ./foo/bar/baz/hello.log must make sure ./foo/bar/baz is existed
	dirname := filepath.Dir(filename)
	if err := os.MkdirAll(dirname, 0755); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create directory %s", dirname))
	}
	// if we got here, then we need to create a file
	fh, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to open file %s: %s", rl.pattern, err))
	}

	if err := rl.rotate_nolock(filename); err != nil {
		err = errors.New("failed to rotate")
		if bailOnRotateFail {
			// Failure to rotate is a problem, but it's really not a great
			// idea to stop your application just because you couldn't rename
			// your log.
			//
			// We only return this error when explicitly needed (as specified by bailOnRotateFail)
			//
			// However, we *NEED* to close `fh` here
			if fh != nil { // probably can't happen, but being paranoid
				fh.Close()
			}
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	rl.outFh.Close()
	rl.outFh = fh
	rl.curBaseFn = baseFn
	rl.curFn = filename
	rl.generation = generation

	//if h := rl.eventHandler; h != nil {
	//	go h.Handle(&FileRotatedEvent{
	//		prev:    previousFn,
	//		current: filename,
	//	})
	//}
	return fh, nil
}

func (rl *RotateWriter) rotate_nolock(filename string) error {
	lockfn := filename + `_lock`
	fh, err := os.OpenFile(lockfn, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		// Can't lock, just return
		return err
	}

	var guard cleanupGuard
	guard.fn = func() {
		fh.Close()
		os.Remove(lockfn)
	}
	defer guard.Run()

	//if rl.linkName != "" {
	//	tmpLinkName := filename + `_symlink`
	//	if err := os.Symlink(filename, tmpLinkName); err != nil {
	//		return errors.Wrap(err, `failed to create new symlink`)
	//	}
	//
	//	if err := os.Rename(tmpLinkName, rl.linkName); err != nil {
	//		return errors.Wrap(err, `failed to rename new symlink`)
	//	}
	//}
	//
	//if rl.maxAge <= 0 && rl.rotationCount <= 0 {
	//	return errors.New("panic: maxAge and rotationCount are both set")
	//}

	matches, err := filepath.Glob(rl.globPattern)
	if err != nil {
		return err
	}

	cutoff := rl.clock.Now().Add(-1 * rl.maxAge)
	var toUnlink []string
	for _, path := range matches {
		// Ignore lock files
		if strings.HasSuffix(path, "_lock") || strings.HasSuffix(path, "_symlink") {
			continue
		}

		fi, err := os.Stat(path)
		if err != nil {
			continue
		}
		//_, err := os.Lstat(path)
		//if err != nil {
		//	continue
		//}

		if rl.maxAge > 0 && fi.ModTime().After(cutoff) {
			continue
		}

		//if rl.rotationCount > 0 && fl.Mode()&os.ModeSymlink == os.ModeSymlink {
		//	continue
		//}
		toUnlink = append(toUnlink, path)
	}

	//if rl.rotationCount > 0 {
	//	// Only delete if we have more than rotationCount
	//	if rl.rotationCount >= uint(len(toUnlink)) {
	//		return nil
	//	}
	//
	//	toUnlink = toUnlink[:len(toUnlink)-int(rl.rotationCount)]
	//}

	if len(toUnlink) <= 0 {
		return nil
	}

	guard.Enable()
	go func() {
		// unlink files on a separate goroutine
		for _, path := range toUnlink {
			os.Remove(path)
		}
	}()

	return nil
}
