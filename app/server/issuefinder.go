package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	"issuefinder/infra/config"
	"issuefinder/interface/web"
)

func main() {
	app := &cli.App{
		Name: "issuefinder",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "show debug information",
			},
			&cli.IntFlag{
				Name:    "port",
				Value:   0,
				Aliases: []string{"p"},
				Usage:   "port where to run the service (default: random)",
			},
			&cli.BoolFlag{
				Name:  "no-webpage",
				Value: false,
				Usage: "do not automatically open the web page",
			},
		},
		Action: RunWebServer,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Println("Got error")
		log.Fatal(err)
	}
}

func RunWebServer(c *cli.Context) error {
	config := config.NewHandler()
	catchSignals(config)

	address := fmt.Sprintf("127.0.0.1:%d", c.Int("port"))
	webserver := web.NewWebServer(address, config)

	webserver.SetDebug(c.Bool("debug"))

	// launch the browser and open the page to this
	go func() {
		<-time.After(100 * time.Millisecond)
		if !c.Bool("no-webpage") {
			open(fmt.Sprintf("http://%s/index.tml", address))
		}
	}()
	return webserver.Start()
}

// open tells the operating system to launch the default application for the provided URL
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func catchSignals(cfg config.Handler) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		for s := range sigs {
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT:
				fmt.Println("got signal and try to exit: ", s)
				cfg.SaveConfig()
				os.Exit(0)
			default:
				fmt.Println("other: ", s)
			}
		}
	}()
}
