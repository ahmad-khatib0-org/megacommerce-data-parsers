package main

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Db struct {
		Dsn string `yaml:"dsn"`
	} `yaml:"db"`
}

var config *Config

func initConfig() error {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	config = &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return err
	}

	return nil
}

// The OpenDB() function returns a  *pgxpool.Pool postgres connection pool.
func OpenDB() (*pgxpool.Pool, error) {
	d, _ := pgxpool.New(context.Background(), config.Db.Dsn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := d.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return d, nil
}
