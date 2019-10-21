package logger

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	test *LoggerCapturer
)

//
type LoggerCapturer struct {
	Buffer *bytes.Buffer
	Writer *bufio.Writer
}

func newLoggerCapturer() *LoggerCapturer {
	buffer := &bytes.Buffer{}
	capture := &LoggerCapturer{
		Buffer: buffer,
		Writer: bufio.NewWriter(buffer),
	}

	logger.SetOutput(capture.Writer)

	return capture
}

func (c *LoggerCapturer) Result() string {
	c.Writer.Flush()

	message := strings.TrimSuffix(c.Buffer.String(), "\n")

	c.Buffer.Reset()

	return message
}

func TestLogger(t *testing.T) {
	Convey("Logging test set", t, func() {
		loggerCapturer := newLoggerCapturer()

		Convey("Format message method", func() {
			got := FormatMessage(TRACE, "Hello %s", "format")
			expected := formats[TRACE] + outputColors[TRACE]("Hello format")

			So(got, ShouldEqual, expected)
		})

		Convey("Write method", func() {
			Write(FATAL, "Hello %s", "fatal")
			So(loggerCapturer.Result(), ShouldEqual, FormatMessage(FATAL, "Hello %s", "fatal"))
		})

		Convey("Basic method", func() {
			Trace("Hello %s", "trace")
			So(loggerCapturer.Result(), ShouldEqual, FormatMessage(TRACE, "Hello %s", "trace"))

			Info("Hello %s", "info")
			So(loggerCapturer.Result(), ShouldEqual, FormatMessage(INFO, "Hello %s", "info"))

			Warn("Hello %s", "warn")
			So(loggerCapturer.Result(), ShouldEqual, FormatMessage(WARN, "Hello %s", "warn"))

			Error("Hello %s", "error")
			So(loggerCapturer.Result(), ShouldEqual, FormatMessage(ERROR, "Hello %s", "error"))
		})

		Convey("Fatal method", func() {
			// Store the original exit method and restore at the end
			originalOsExit := osExit
			defer func() {
				osExit = originalOsExit
			}()

			// Overwrite the original exit process by empty method
			osExit = func(code int) {}

			Fatal("Hello %s", "fatal")
			So(loggerCapturer.Result(), ShouldEqual, FormatMessage(FATAL, "Hello %s", "fatal"))
		})
	})
}

func TestLoggerFatal(t *testing.T) {
	if os.Getenv("MUST_CRASHER") == "1" {
		Fatal("Hello %s", "fatal")
		return
	}

	Convey("Exit should be OK", t, func() {
		cmd := exec.Command(os.Args[0], "-test.run=TestLoggerFatal")
		cmd.Env = append(os.Environ(), "MUST_CRASHER=1")

		exitError, ok := cmd.Run().(*exec.ExitError)

		So(ok, ShouldBeTrue)
		So(exitError.Success(), ShouldBeFalse)
	})
}
