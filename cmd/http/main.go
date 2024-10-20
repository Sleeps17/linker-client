package main

import (
	"context"
	"github.com/Sleeps17/linker-client/internal/app"
	"github.com/Sleeps17/linker-client/internal/config"
	"github.com/Sleeps17/linker-client/internal/models"
	"github.com/Sleeps17/linker-client/internal/utils/formatter"
	"github.com/Sleeps17/linker-client/internal/utils/parser"
	"os"
	"os/signal"
	"syscall"
)

type MethodInfo struct {
	Name         models.Method
	Description  string
	RequiredArgs []string
	OptionalArgs []string
}

var Help = []MethodInfo{
	{
		Name:         models.PostTopic,
		Description:  "Создает новый топик",
		RequiredArgs: []string{"topic"},
		OptionalArgs: []string{},
	},
	{
		Name:         models.ListTopics,
		Description:  "Выводит список топиков",
		RequiredArgs: []string{},
		OptionalArgs: []string{},
	},
	{
		Name:         models.DeleteTopic,
		Description:  "Удаляет топик",
		RequiredArgs: []string{"topic"},
		OptionalArgs: []string{},
	},
	{
		Name:         models.PostLink,
		Description:  "Создает новую ссылку в топике",
		RequiredArgs: []string{"topic", "link"},
		OptionalArgs: []string{"alias"},
	},
	{
		Name:         models.PickLink,
		Description:  "Выводит ссылку из топика по алиасу",
		RequiredArgs: []string{"topic", "alias"},
	},
	{
		Name:         models.ListLinks,
		Description:  "Выводит список ссылок из топика",
		RequiredArgs: []string{"topic"},
	},
	{
		Name:         models.DeleteLink,
		Description:  "Удаляет ссылку из топика",
		RequiredArgs: []string{"topic", "alias"},
	},
}

func main() {
	format := formatter.New(os.Stdout)

	cfg, err := config.Load()
	if err != nil {
		format.Error(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Client.Timeout)
	defer cancel()

	//validate := corevalidator.New[interface{}](validator.New())
	parse := parser.New(cfg.AvailableMethods)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	application := app.New(ctx, stop, cfg, parse, format)
	go application.MustRun()

	<-stop
	application.Stop(cancel)
}
