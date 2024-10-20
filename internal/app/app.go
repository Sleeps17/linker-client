package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sleeps17/linker-client/internal/clients/linker"
	httpclient "github.com/Sleeps17/linker-client/internal/clients/linker/http"
	"github.com/Sleeps17/linker-client/internal/config"
	"github.com/Sleeps17/linker-client/internal/models"
	"github.com/Sleeps17/linker-client/internal/utils/formatter"
	"github.com/Sleeps17/linker-client/internal/utils/parser"
	"os"
)

const (
	linkerUsernameEnv = "LINKER_USERNAME"
)

type App struct {
	ctx      context.Context
	stopChan chan os.Signal
	cfg      *config.Config
	parser   *parser.Parser
	format   *formatter.Formatter
	client   linker.Client
}

func New(
	ctx context.Context,
	stopChan chan os.Signal,
	cfg *config.Config,
	parser *parser.Parser,
	format *formatter.Formatter,
) *App {
	return &App{
		ctx:      ctx,
		stopChan: stopChan,
		cfg:      cfg,
		parser:   parser,
		format:   format,
		client:   httpclient.New(&cfg.Client),
	}
}

func (a *App) MustRun() {
	username := a.getUsername()

	method, err := a.parser.ParseArgs(os.Args)
	if err != nil {
		if errors.Is(err, parser.ErrEmptyMethod) {
			a.format.Error("Кажется вы забыли указать метод")
		}

		if errors.Is(err, parser.ErrUnknownMethod) {
			a.format.Errorf("Метод %s не поддерживается", os.Args[1])
		}

		if errors.Is(err, parser.ErrMissRequiredArg) {
			a.format.Error("Кажется вы указали не все необходимые аргументы, используйте help")
		}
	}

	switch method {
	case models.PostTopic:
		topic := a.parser.Arg(models.Topic).(string)
		id, err := a.client.PostTopic(a.ctx, username, topic)
		if err != nil {
			a.format.Error(err.Error())
		}
		a.format.Successf("Топик успешно создан id = %d", id)
	case models.DeleteTopic:
		topic := a.parser.Arg(models.Topic).(string)
		id, err := a.client.DeleteTopic(a.ctx, username, topic)
		if err != nil {
			a.format.Error(err.Error())
		}
		a.format.Successf("Топик успешно удален id = %d", id)
	case models.ListTopics:
		topics, err := a.client.ListTopics(a.ctx, username)
		if err != nil {
			a.format.Error(err.Error())
		}

		headers := []string{"ID", "Topic"}
		values := make([][]string, 0)
		for idx, topic := range topics {
			values = append(values, []string{fmt.Sprint(idx + 1), a.format.SuccessString(topic)})
		}

		a.format.SuccessTable(headers, values...)
	case models.PostLink:
		topic := a.parser.Arg(models.Topic).(string)
		link := a.parser.Arg(models.Link).(string)
		alias := a.parser.Arg(models.Alias).(string)

		alias, err := a.client.PostLink(a.ctx, username, topic, alias, link)
		if err != nil {
			a.format.Error(err.Error())
		}

		a.format.Successf("Ссылка успешно создана, alias = %s", alias)
	case models.PickLink:
		topic := a.parser.Arg(models.Topic).(string)
		alias := a.parser.Arg(models.Alias).(string)

		link, err := a.client.PickLink(a.ctx, username, topic, alias)
		if err != nil {
			a.format.Error(err.Error())
		}

		a.format.Successf("Ссылка успешно получена, link = %s", link)
	case models.DeleteLink:
		topic := a.parser.Arg(models.Topic).(string)
		alias := a.parser.Arg(models.Alias).(string)

		alias, err := a.client.DeleteLink(a.ctx, username, topic, alias)
		if err != nil {
			a.format.Error(err.Error())
		}

		a.format.Successf("Ссылка успешно удалена, alias = %s", alias)
	case models.ListLinks:
		topic := a.parser.Arg(models.Topic).(string)

		links, err := a.client.ListLinks(a.ctx, username, topic)
		if err != nil {
			a.format.Error(err.Error())
		}

		headers := []string{"ID", "Alias", "Link"}
		values := make([][]string, 0)

		idx := 0
		for alias, link := range links {
			values = append(values, []string{fmt.Sprint(idx + 1), a.format.SuccessString(alias), a.format.SuccessString(link)})
			idx++
		}

		a.format.SuccessTable(headers, values...)
	case models.Help:
	}

	a.stopChan <- os.Interrupt
}

func (a *App) Stop(cancel context.CancelFunc) {
	//a.format.Warning("Остановка приложения")
	cancel()
}

func (a *App) getUsername() string {
	username := os.Getenv(linkerUsernameEnv)

	if username == "" {
		a.format.Error("Вы не зарегистрировались в клиенте, установите переменную окружения LINKER_USERNAME")
	}

	return username
}
