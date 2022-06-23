package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/consts"

	// "context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/logger"

	// "os"
	// "os/signal"
	// "syscall"
	// "time"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
	// internalhttp "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/server/http"
	// memorystorage "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
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
	pflag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func LoggingMiddleware(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rTime := ctx.Value("start-request-time").(time.Time)
	ip := strings.Split(r.RemoteAddr, ":")

	msg := fmt.Sprintf("%v [%v] %v %v %v %v %v %v",
		ip[0],
		rTime.Format(time.RFC822Z),
		r.Method,
		r.RequestURI,
		r.Proto,
		r.Context().Value("response-status-code"),
		time.Since(rTime),
		r.Header.Get("User-Agent"),
	)
	Logger.Info(msg)
}

func HelloWorld(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "start-request-time", time.Now())
		ctx = context.WithValue(ctx, "response-status-code", http.StatusOK)

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
		fmt.Printf("Ошибка чтения файла: %w", err)
		return
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		fmt.Printf("Ошибка маршелинга: %w", err)
		return
	}

	Logger, err = logger.New(config.Logger.Level, config.Logger.File)
	if err != nil {
		fmt.Println(err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", HelloWorld(LoggingMiddleware))

	log.Println("Запуск веб-сервера на " + fmt.Sprintf("%v:%v", config.Http.Host, config.Http.Port))
	err = http.ListenAndServe(fmt.Sprintf("%v:%v", config.Http.Host, config.Http.Port), mux)
	log.Fatal(err)
	//storage := memorystorage.New()
	//calendar := app.New(logg, storage)

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
