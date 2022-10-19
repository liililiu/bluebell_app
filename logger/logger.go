package logger

import (
	"bluebell_app/settings"
	"bluebell_app/utils"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lg *zap.Logger

// 注意这里的lg暂时没有用全局zap.L()代替

// InitLogger 初始化Logger

func Init(cfg *settings.LogConfig, mode string) (err error) {
	// core 需要三部分

	// 第二部分：encoder
	encoder := getEncoder()

	// 第三部分： level
	var l = new(zapcore.Level)
	// 将配置文件中的level类型转换为zapCore的level
	err = l.UnmarshalText([]byte(
		//viper.GetString("" + "log.level"),
		cfg.Level,
	))

	if err != nil {
		fmt.Println("将配置文件中的level类型转换为zapCore的level失败，err：", err)
		return
	}
	// 第一部分：writeSyncer
	writeSyncer := getLogWriter(l.String())
	var core zapcore.Core
	fmt.Printf("日志模式mode的值：%v\n", mode)
	if mode == "dev" {
		// 开发模式,日志打印到终端;根据需要选择是否打印到文件
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
			//zapcore.NewCore(zapcore.NewCore(encoder, writeSyncer, l))
		)
		fmt.Println("日志记录模式：dev")
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
		fmt.Println("日志记录模式：非dev")
	}

	lg = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	return
}

// Encoder 的函数式配置

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// WriteSyncer 的函数式配置

func getLogWriter(level string) zapcore.WriteSyncer {
	// 检查日志目录
	if ok, _ := utils.PathExists(settings.GlobalConfig.LogConfig.Dir); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", settings.GlobalConfig.LogConfig.Dir)
		_ = os.Mkdir(settings.GlobalConfig.LogConfig.Dir, os.ModePerm)
	}
	fileWriter, _ := rotatelogs.New(
		// 生成日志文件 日志文件名可以按时间来创建
		path.Join(settings.GlobalConfig.LogConfig.Dir, "%Y-%m-%d", level+".log"),
		// 日志文件相关的时间 根本本地时间
		rotatelogs.WithClock(rotatelogs.Local),
		// 文件存活有效期
		rotatelogs.WithMaxAge(time.Duration(settings.GlobalConfig.LogConfig.MaxAge)*24*time.Hour), // 日志留存时间
		// 多久创建一次新的日志文件
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	// 根据条件选择是否同时打印控制台和文件日志
	//if viper.GetBool("log.logInConsole") {
	//	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter))
	//}
	//if settings.GlobalConfig.LogConfig.LogInConsole {
	//	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter))
	//}
	return zapcore.AddSync(fileWriter)
}

// 两个中间件

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)
		lg.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
