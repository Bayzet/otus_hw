package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	internalhttp "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/server/http"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/consts"
	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/logger"
)

var configFile string
var Logger *logger.Logger

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

func helloWorld(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), consts.StartRequestTime, time.Now())
		ctx = context.WithValue(ctx, consts.ResponseStatusCode, http.StatusOK)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))

		r = r.WithContext(ctx)

		h(w, r)
	}
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
		fmt.Printf("Ошибка чтения файла: %v", err.Error())
		return
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		fmt.Printf("Ошибка маршелинга: %v", err.Error())
		return
	}

	Logger, err = logger.New(config.Logger.Level, config.Logger.File)
	if err != nil {
		fmt.Println(err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloWorld(internalhttp.LoggingMiddleware(Logger)))

	addr := fmt.Sprintf("%v:%v", config.HTTP.Host, config.HTTP.Port)

	server := internalhttp.NewServer(addr, mux)

	ctx := context.Background()
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			Logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	Logger.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		Logger.Error("failed to start http server: " + err.Error())
		os.Exit(1)
	}
}
