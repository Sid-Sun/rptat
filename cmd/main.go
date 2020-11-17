package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"

	"github.com/sid-sun/rptat/app"
	"github.com/sid-sun/rptat/cmd/config"
)

func main() {
	argc := len(os.Args)
	errMsg := func() {
		fmt.Println("Run rptat with either -s reload to reload for path to connfig file")
	}
	switch argc {
	case 2:
		// Stat and start
		info, err := os.Stat(os.Args[1])
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Invalid path")
			}
			return
		}
		if info.IsDir() {
			fmt.Println("Such wow")
			return
		}
		break
	case 3:
		if os.Args[1] != "-s" && os.Args[2] != "reload" {
			errMsg()
			return
		}

		path := "/var/run/rptat"
		_, err := os.Stat(path + "/rptat.pid")
		if err != nil {
			if !os.IsPermission(err) && !os.IsNotExist(err) {
				panic(err)
			}

			path = "/var/run/user/" + strconv.Itoa(os.Getuid()) + "/rptat"
			_, err := os.Stat(path + "/rptat.pid")

			if err != nil {
				if os.IsPermission(err) {
					fmt.Println("Permission error")
					return
				}
				if os.IsNotExist(err) {
					fmt.Println("Could not find an RPTAT process")
					return
				}
				panic(err)
			}
		}
		data, err := ioutil.ReadFile(path + "/rptat.pid")
		if err != nil {
			panic(err)
		}

		pid, err := strconv.Atoi(string(data))
		if err != nil {
			panic(err)
		}

		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Println("Process not found")
			return
		}

		err = process.Signal(syscall.SIGUSR1)
		if err != nil {
			panic(err)
		}
		return
		// Reload
	default:
		errMsg()
		return
	}
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	initLogger(cfg.GetEnv())
	app.StartServer(cfg, logger)
}
