package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-service/internal/app"
	"auth-service/internal/container"
	"auth-service/internal/logger"
)

// ------------------------------------------------------------------------------
// local
// ------------------------------------------------------------------------------
var (
	ct            *container.Container
	server        *app.Application = nil
	isShutDownApp bool             = false
	terminate     bool             = false
)

// ------------------------------------------------------------------------------
// sigHandler
// ------------------------------------------------------------------------------
func sigHandler(chSig chan os.Signal) {
	logger.Log.Print(2, "[server]sigHandler")
	for {
		signal := <-chSig
		str := fmt.Sprintf("[server] Accept Signal : %d", signal)
		logger.Log.Print(2, "%s", str)
		switch signal {
		case syscall.SIGHUP:
			logger.Log.Print(2, "[server]SIGHUP(%d)\n", signal)
		case syscall.SIGINT:
			logger.Log.Print(2, "[server]SIGINT(%d)\n", signal)
			shudownApp()
			terminate = true
			// os.Exit(0)
		case syscall.SIGTERM:
			logger.Log.Print(2, "SIGTERM(%d)\n", signal)
			terminate = true
			// os.Exit(0)
		case syscall.SIGKILL:
			logger.Log.Print(2, "SIGKILL(%d)\n", signal)
			terminate = true
		case syscall.SIGUSR1:
			logger.Log.Print(2, "SIGUSR1(%d)\n", signal)
			go shudownApp()
			// os.Exit(0)
		default:
			logger.Log.Print(2, "Unknown signal(%d)\n", signal)
			// panic(signal)
		}
	}
}

// ------------------------------------------------------------------------------
// initEnvVaiable
// ------------------------------------------------------------------------------
func initEnvVaiable() bool {
	//
	return true
}

// ------------------------------------------------------------------------------
// initContainer
// ------------------------------------------------------------------------------
func initContainer() bool {
	var err error = nil
	ct, err = container.NewContainer()
	if err != nil {
		logger.Log.Print(2, "[server]initContainer err.. %v \n", err)
		return false
	}

	return true
}

// ------------------------------------------------------------------------------
// initSignal
// ------------------------------------------------------------------------------
func initSignal() {
	logger.Log.Print(2, "[server]initSignal...")
	// signal handler
	ch_signal := make(chan os.Signal, 10)
	signal.Notify(ch_signal, syscall.SIGSEGV, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGUSR1)
	go sigHandler(ch_signal)
}

// ------------------------------------------------------------------------------
// initEnv
// ------------------------------------------------------------------------------
func initEnv() bool {
	initEnvVaiable()

	// container
	if !initContainer() {
		return false
	}

	// setdebug level
	logger.Log.SetLevel(ct.Config.DebugLv)

	// signal
	initSignal()
	return true
}

// ------------------------------------------------------------------------------
// shudownApp
// ------------------------------------------------------------------------------
func shudownApp() {
	if isShutDownApp {
		return
	}
	isShutDownApp = true
	logger.Log.Print(2, "[server]shudownApp..")

	server.Shutdown()
}

// ------------------------------------------------------------------------------
// clearEnv
// ------------------------------------------------------------------------------
func clearEnv() {
}

func main() {
	process_mode := flag.String("mode", "debug", "프로세스 실행 모드를 선택")
	flag.Parse()
	logger.Log.Print(2, "process mode : %s", *process_mode)
	logger.Log.Print(2, "남은 인자들:", flag.Args())

	ok := initEnv()
	defer clearEnv()

	logger.Log.Print(2, "init state : %v", ok)

	if !ok {
		logger.Log.Error("initEnv Error...")
	}

	var ch_terminate chan bool = make(chan bool)
	// NewApplication()
	if ok {
		server = app.NewApplication(ct, ch_terminate)
		go server.Start()
	}

	for ok {
		select {
		case <-ch_terminate:
			logger.Log.Print(2, "Server shutdown ok.")
			shudownApp()
			// request manage-service
			terminate = true
		default:
		}

		if terminate {
			logger.Log.Print(2, "Quit auth-service .. ")
			break
		}

		time.Sleep(time.Millisecond * 1000)
		// check manage process
	}
}
