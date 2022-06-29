package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/ioutil"

	_ "github.com/jackc/pgx/stdlib"
	goose "github.com/pressly/goose/v3"
	"github.com/spf13/pflag"
	yaml "gopkg.in/yaml.v3"
)

//go:embed *.sql
var embedMigrations embed.FS

var configFile string

func init() {
	pflag.StringVar(&configFile, "config", "config.yaml", "Path to configuration file")
}

func main() {
	pflag.Parse()

	var config struct {
		Storage struct {
			Driver string `yaml:"driver"`
			DSN    string `yaml:"dsn"`
		} `yaml:"storage"`
	}

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

	var db *sql.DB
	db, err = sql.Open(config.Storage.Driver, config.Storage.DSN)
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("pgx"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "."); err != nil {
		panic(err)
	}
}
