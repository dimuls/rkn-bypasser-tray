package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/getlantern/systray"
	"github.com/sirupsen/logrus"
)

var (
	torCmd         *exec.Cmd
	rknBypasserCmd *exec.Cmd
)

func main() {
	// Should be called at the very beginning of main().
	systray.Run(onReady, onQuit)
}

func onReady() {

	iconBytes, err := ioutil.ReadFile("./icon.ico")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load icon.ico")
	}

	systray.SetIcon(iconBytes)
	systray.SetTitle("RKN Bypasser")
	systray.SetTooltip("RKN Bypasser")

	mQuit := systray.AddMenuItem("Выход", "Выйти и выключить прокси-сервер")

	onStart()

	go func() {
		select {
		case <-mQuit.ClickedCh:
			onQuit()
			systray.Quit()
		}
	}()
}

func onStart() {

	go func() {
		for {
			torCmd = exec.Command("./tor.exe", "--quiet")
			err := torCmd.Run()
			if err != nil {
				logrus.WithError(err).Error("Failed to run tor.exe")
				time.Sleep(time.Second)
			} else {
				return
			}
		}
	}()

	go func() {
		rknBypasserCmd = exec.Command("./rkn-bypasser.exe",
			"--bind-addr", "127.0.0.1:8000")
		err := rknBypasserCmd.Run()
		if err != nil {
			logrus.WithError(err).Error("Failed to run rkn-bypasser.exe")
			time.Sleep(time.Second)
		} else {
			return
		}
	}()
}

func onQuit() {
	rknBypasserCmd.Process.Signal(os.Kill)
	torCmd.Process.Signal(os.Kill)
}
