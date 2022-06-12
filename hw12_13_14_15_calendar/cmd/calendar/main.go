package main

import (
	// "context"
	"flag"
	"fmt"
	"io/ioutil"

	memorystorage "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"

	// "os"
	// "os/signal"
	// "syscall"
	// "time"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	// "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/logger"
	// internalhttp "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/server/http"
	// memorystorage "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	pflag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	pflag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	var config Config

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Ошибка чтения файла: %w", err)
		return
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		fmt.Printf("Ошибка маршелинга: %w", err)
		return
	}

	fmt.Printf("%#v", config)

	Logg, err := logger.New(config.Logger.Level, config.Logger.File)
	if err != nil {
		fmt.Println(err.Error())
	}

	storage := memorystorage.New()
	// calendar := app.New(logg, storage)

	// server := internalhttp.NewServer(logg, calendar)

	// ctx, cancel := signal.NotifyContext(context.Background(),
	// 	syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	// defer cancel()

	// go func() {
	// 	<-ctx.Done()

	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// 	defer cancel()

	// 	if err := server.Stop(ctx); err != nil {
	// 		logg.Error("failed to stop http server: " + err.Error())
	// 	}
	// }()

	// logg.Info("calendar is running...")

	// if err := server.Start(ctx); err != nil {
	// 	logg.Error("failed to start http server: " + err.Error())
	// 	cancel()
	// 	os.Exit(1) //nolint:gocritic
	// }
}
