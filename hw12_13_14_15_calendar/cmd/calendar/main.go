//nolint:varnamelen
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	sqlstorage "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/server/grpc"

	"github.com/spf13/pflag"
	yaml "gopkg.in/yaml.v3"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/consts"
	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

type StorageType string

func (st StorageType) Validate() error {
	switch st {
	case consts.StorageTypeSQL, consts.StorageTypeMemory:
		return nil
	default:
		return logger.ErrStorageTypeNotValid
	}
}

func init() {
	pflag.StringVar(&configFile, "config", "/calendar/config.yaml", "Path to configuration file")
}

//nolint:funlen
func main() {
	pflag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	var config Config
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Ошибка чтения файла: %v", err.Error())
		return
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		fmt.Printf("Ошибка маршелинга: %v", err.Error())
		return
	}

	Logger, err := logger.New(config.Logger.Level, config.Logger.File)
	if err != nil {
		fmt.Println(err.Error())
	}

	var storage app.Repository
	if config.Storage.Type == consts.StorageTypeMemory {
		storage = memorystorage.New()
	} else {
		storage = sqlstorage.New(config.Storage.Driver, config.Storage.DSN)
	}

	calendar := app.New(storage)
	router := mux.NewRouter()

	httpServer := internalhttp.NewServer(router, config.HTTP.Host, config.HTTP.Port, Logger)
	calendarAPI := internalhttp.CalendarApp{App: calendar}
	calendarAPI.RegisterHTTPHandlers(router)

	grpcServer := grpc.NewServer(calendar, config.GRPC.Host, config.GRPC.Port, Logger)

	ctx := context.Background()
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(ctx, time.Second*3)
		defer cancel()

		if err := httpServer.Stop(ctx); err != nil {
			Logger.Error("failed to stop http httpServer: " + err.Error())
		}
	}()

	go func() {
		if err := httpServer.Start(ctx); err != nil {
			Logger.Error("failed to start http server: " + err.Error())
			os.Exit(1)
		}
	}()

	go func() {
		if err := grpcServer.Run(ctx); err != nil {
			Logger.Error("failed to start gRPC server: " + err.Error())
			os.Exit(1)
		}
	}()
	sChan := make(chan os.Signal, 1)
	signal.Notify(sChan, syscall.SIGINT, syscall.SIGHUP)

	for range sChan {
		grpcServer.Stop(ctx)
		_ = httpServer.Stop(ctx)

		os.Exit(1)
	}
}
