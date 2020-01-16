package main

import (
	"flag"
	"log"
	"stock_analysis/base"
	"time"
)

var (
	mc      *base.MainControl
	isDebug bool
	err     error
)

func init() {
	flag.BoolVar(&isDebug, "debug", false, "debug模式")

	flag.Parse()
	log.SetFlags(0)
}

func main() {
	mc, err = base.NewMainCtrl(&base.MainCtrlParams{
		LogPrefix:      "rt",
		SrIntervalTime: time.Hour * 12,
		IsDebug:        isDebug,
	})
	if err != nil {
		log.Fatal("Level-High Exception", err)
	}
	defer func() {
		mc.Output.Log("Exit ...")
		mc.Output.Close()
	}()

	rt, err := NewRealTime(mc, isDebug)
	if err != nil {
		mc.Output.Error(err)
		return
	}
	rt.Running()

	for {
		select {
		case <-mc.Interrupt:
			mc.Output.Log("Interrupt")
			return
		case <-mc.Ctx.Done():
			mc.Output.Log("Ctx.Done")
			return
		}
	}
}
