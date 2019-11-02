package base

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/huoxuhuoxu/UseGoexPackaging/diyLog"
	"github.com/huoxuhuoxu/UseGoexPackaging/systemRp"
)

type MainControl struct {
	Output    *diyLog.DiyLog     // 日志输出
	Sr        *systemRp.SystemRp // 重启对象
	Ctx       context.Context    // 上下文对象
	CtxCancel context.CancelFunc // 上下文通知结束
	Interrupt chan os.Signal     // 信号监控
}

type MainCtrlParams struct {
	LogPrefix      string        // 日志前缀
	SrIntervalTime time.Duration // 重启间隔时长
	IsDebug        bool          // debug模式
}

func NewMainCtrl(taskType string, mcp *MainCtrlParams) (*MainControl, error) {
	var (
		mc  *MainControl
		err error
	)

	func() {
		defer func() {
			if tmpErr := recover(); tmpErr != nil {
				errMsg := fmt.Sprintf("initial task failed, %s", tmpErr)
				err = errors.New(errMsg)
			}
		}()

		mc = &MainControl{}

		// Log
		if mcp.IsDebug {
			mc.Output, _ = diyLog.NewDiyLog(diyLog.DEBUG, diyLog.STDOUT)
		} else {
			mc.Output, _ = diyLog.NewDiyLog(diyLog.DEBUG, diyLog.FILE_RECORD)
			mc.Output.SetOutputParams("logs", taskType+"-"+mcp.LogPrefix+"-")
		}

		// Process restart
		mc.Sr = systemRp.NewSystemRp(mcp.SrIntervalTime)
		mc.Sr.SubscribeBeforeFunc(mc.Output.Close)

		// 主控制流
		mc.Ctx, mc.CtxCancel = context.WithCancel(context.Background())

		// 监听信号, SIGINT, SIGTERM
		mc.Interrupt = make(chan os.Signal, 1)
		signal.Notify(mc.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	}()

	return mc, err
}
