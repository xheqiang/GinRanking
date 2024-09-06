package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var LogBaseDir = "./runtime/log"

func init() {
	/* logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}) */
	// logrus.SetReportCaller(false) // 用来指示日志是否要报告调用者的信息
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func Write(msg string, filename string) {
	setOutPutFile(logrus.InfoLevel, filename)
	logrus.Info(msg)
}

func Debug(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.DebugLevel, "debug")
	logrus.WithFields(fields).Debug(args...)
}

func Info(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.InfoLevel, "info")
	logrus.WithFields(fields).Info(args...)
}

func Warn(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.WarnLevel, "warn")
	logrus.WithFields(fields).Warn(args...)
}

func Fatal(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.FatalLevel, "fatal")
	logrus.WithFields(fields).Fatal(args...)
}

func Error(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.ErrorLevel, "error")
	logrus.WithFields(fields).Error(args...)
}

func Panic(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.PanicLevel, "panic")
	logrus.WithFields(fields).Panic(args...)
}

func Trace(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.TraceLevel, "trace")
	logrus.WithFields(fields).Trace(args...)
}

func setOutPutFile(level logrus.Level, logName string) {

	if _, err := os.Stat(LogBaseDir); os.IsNotExist(err) {
		err = os.MkdirAll(LogBaseDir, 0777)
		if err != nil {
			panic(fmt.Errorf("create log dir '%s' error: %s ", LogBaseDir, err))
		}
	}

	timeStr := time.Now().Format("20060102")
	fileName := path.Join(LogBaseDir, logName+"."+timeStr+".log")

	var err error
	os.Stderr, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Errorf("open log file '%s' error: %s ", fileName, err))
	}

	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(level)
}

func LoggerToFile() gin.LoggerConfig {
	if _, err := os.Stat(LogBaseDir); os.IsNotExist(err) {
		err = os.MkdirAll(LogBaseDir, 0777)
		if err != nil {
			panic(fmt.Errorf("create log dir '%s' error: %s ", LogBaseDir, err))
		}
	}

	timeStr := time.Now().Format("20060102")
	fileName := path.Join(LogBaseDir, "access_log."+timeStr+".log")

	os.Stderr, _ = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	var conf = gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - %s \"%s %s %s %d %s \"%s\" %s\"\n",
				param.TimeStamp.Format("2006-01-02 15:04:05"),
				param.ClientIP,
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		},
		Output: io.MultiWriter(os.Stdout, os.Stderr),
	}
	return conf
}

func Recover(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			if _, errDir := os.Stat(LogBaseDir); os.IsNotExist(errDir) {
				errDir = os.MkdirAll(LogBaseDir, 0777)
				if errDir != nil {
					panic(fmt.Errorf("create log dir '%s' error: %s ", LogBaseDir, errDir))
				}
			}

			timeStr := time.Now().Format("2006-01-02")
			fileName := path.Join(LogBaseDir, "panic_."+timeStr+".log")

			f, errFile := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if errFile != nil {
				fmt.Println(errFile)
			}

			timeFileStr := time.Now().Format("2006-01-02 15:04:05")
			f.WriteString("panic error time: " + timeFileStr + "\n")
			f.WriteString(fmt.Sprintf("%v", err) + "\n")
			f.WriteString("stacktrace from panic:" + string(debug.Stack()) + "\n")
			f.Close()

			ctx.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  fmt.Sprintf("%v", err),
			})
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			ctx.Abort()
		}
	}()
	ctx.Next()
}
