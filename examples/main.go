package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/ophum/mcclient"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

var config Config

func init() {
	path := flag.String("config", "config.yaml", "-config config.yaml")
	flag.Parse()

	f, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		log.Fatal(err)
	}
}
func main() {
	client, err := mcclient.New(config.Host, config.Port, config.Password)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Whitelist().Add(ctx, "hum_op"); err != nil {
		log.Fatal(err)
	}
	log.Println("added hum_op")

	users, err := client.Whitelist().List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range users {
		log.Println(u)
	}

	if err := client.Whitelist().Remove(ctx, "hum_op"); err != nil {
		log.Fatal(err)
	}
	log.Println("Removed hum_op")

	users, err = client.Whitelist().List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range users {
		log.Println(u)
	}

}
