package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Sleeps17/linker-client/internal/models"
	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/olekukonko/tablewriter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

const (
	methodPostTopic   = "post_topic"
	methodDeleteTopic = "delete_topic"
	methodListTopics  = "list_topics"
	methodPostLink    = "post_link"
	methodPickLink    = "pick_link"
	methodListLinks   = "list_links"
	methodDeleteLink  = "delete_links"
	methodHelp        = "help"

	argTopic = "topic"
	argLink  = "link"
	argAlias = "alias"

	linkerUsername = "LINKER_USERNAME"
)

var (
	availableMethods = []string{
		methodPostTopic,
		methodListTopics,
		methodDeleteTopic,
		methodPostLink,
		methodPickLink,
		methodListLinks,
		methodDeleteLink,
		methodHelp,
	}

	serverAddress = "linker-sleeps17.amvera.io"
	serverTimeout = 30 * time.Second
)

type MethodInfo struct {
	Name         string
	Description  string
	RequiredArgs []string
	OptionalArgs []string
}

var Help = []MethodInfo{
	{
		Name:         methodPostTopic,
		Description:  "Создает новый топик",
		RequiredArgs: []string{"topic"},
		OptionalArgs: []string{},
	},
	{
		Name:         methodListTopics,
		Description:  "Выводит список топиков",
		RequiredArgs: []string{},
		OptionalArgs: []string{},
	},
	{
		Name:         methodDeleteTopic,
		Description:  "Удаляет топик",
		RequiredArgs: []string{"topic"},
		OptionalArgs: []string{},
	},
	{
		Name:         methodPostLink,
		Description:  "Создает новую ссылку в топике",
		RequiredArgs: []string{"topic", "link"},
		OptionalArgs: []string{"alias"},
	},
	{
		Name:         methodPickLink,
		Description:  "Выводит ссылку из топика по алиасу",
		RequiredArgs: []string{"topic", "alias"},
	},
	{
		Name:         methodListLinks,
		Description:  "Выводит список ссылок из топика",
		RequiredArgs: []string{"topic"},
	},
	{
		Name:         methodDeleteLink,
		Description:  "Удаляет ссылку из топика",
		RequiredArgs: []string{"topic", "alias"},
	},
}

var (
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

func MustCreateUsername() string {
	username := os.Getenv(linkerUsername)

	if username == "" {
		fmt.Println("Вы не зарегистрировались в клиенте")
		fmt.Println("Установите переменную окружения LINKER_USERNAME")
		os.Exit(1)
	}

	return username
}

func RequireContains(args ...string) bool {
	for _, arg := range args {
		ok := false
		for _, cmdArg := range os.Args {
			if strings.Contains(cmdArg, arg) {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}

	return true
}

func ParseArg(arg string) string {
	for _, cmdArg := range os.Args {
		if strings.HasPrefix(cmdArg, fmt.Sprintf("--%s=", arg)) {
			return strings.TrimPrefix(cmdArg, fmt.Sprintf("--%s=", arg))
		}
	}

	return ""
}

func PrintResponseForPostTopic(resp models.PostTopicResponse, err error) {
	var output string
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			output = yellow(fmt.Sprintf("Неверные входные данные: %s", st.Message()))
		} else {
			output = red(fmt.Sprintf("Произошла ошибка при выполнении запроса: %s\n", st.Message()))
		}
	} else {
		output = green(fmt.Sprintf("Запрос успешно выполнен, TopicId = %d\n", resp.TopicID))
	}
	fmt.Println(output)
}

func PrintResponseForDeleteTopic(resp models.DeleteTopicResponse, err error) {
	var output string
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			output = yellow(fmt.Sprintf("Неверные входные данные: %s", st.Message()))
		} else {
			output = red(fmt.Sprintf("Произошла ошибка при выполнении запроса: %s\n", st.Message()))
		}
	} else {
		output = green(fmt.Sprintf("Запрос успешно выполнен, TopicId = %d\n", resp.TopicID))
	}
	fmt.Println(output)
}

func PrintResponseForListTopics(resp models.ListTopicsResponse, err error) {
	var output string
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			output = yellow(fmt.Sprintf("Неверные входные данные: %s", st.Message()))
		} else {
			output = red(fmt.Sprintf("Произошла ошибка при выполнении запроса: %s\n", st.Message()))
		}
		fmt.Println(output)
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Topic"})
	for i, topic := range resp.Topics {
		table.Append([]string{green(fmt.Sprint(i + 1)), green(topic)})
	}

	fmt.Println(green("Запрос успешно выполнен"))
	table.Render()
}

func PrintResponseForPostLink(resp models.PostLinkResponse, err error) {
	var output string
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			output = yellow(fmt.Sprintf("Неверные входные данные: %s", st.Message()))
		} else {
			output = red(fmt.Sprintf("Произошла ошибка при выполнении запроса: %s\n", st.Message()))
		}
	} else {
		output = green(fmt.Sprintf("Запрос успешно выполнен, Alias = %s\n", resp.Alias))
	}

	fmt.Println(output)
}

func PrintResponseForPickLink(resp models.PickLinkResponse, err error) {
	var output string
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			output = yellow(fmt.Sprintf("Неверные входные данные: %s", st.Message()))
		} else {
			output = red(fmt.Sprintf("Произошла ошибка при выполнении запроса: %s\n", st.Message()))
		}
	} else {
		output = green(fmt.Sprintf("Запрос успешно выполнен, Link = %s\n", resp.Link))
	}
	fmt.Println(output)
}

func PrintResponseForDeleteLink(resp models.DeleteLinkResponse, err error) {
	var output string
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			output = yellow(fmt.Sprintf("Неверные входные данные: %s", st.Message()))
		} else {
			output = red(fmt.Sprintf("Произошла ошибка при выполнении запроса: %s\n", st.Message()))
		}
	} else {
		output = green(fmt.Sprintf("Запрос успешно выполнен, Alias = %s\n", resp.Alias))
	}
	fmt.Println(output)
}

func PrintResponseForListLinks(topic string, resp models.ListLinksResponse, err error) {
	var output string
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			output = yellow(fmt.Sprintf("Неверные входные данные: %s", st.Message()))
		} else {
			output = red(fmt.Sprintf("Произошла ошибка при выполнении запроса: %s\n", st.Message()))
		}
		fmt.Println(output)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Alias", "Link"})

	for i := range resp.Links {
		table.Append([]string{green(fmt.Sprint(i + 1)), green(resp.Aliases[i]), green(resp.Links[i])})
	}

	fmt.Println(green("Запрос успешно выполнен"))
	fmt.Println(fmt.Sprintf("Список ссылок в топике %s:", topic))
	table.Render()
}

func PrintResponseForHelp() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name", "Description", "Required_Args", "Option_Args"})

	for i := range Help {
		table.Append([]string{
			green(fmt.Sprint(i + 1)),
			green(Help[i].Name),
			green(Help[i].Description),
			green(Help[i].RequiredArgs),
			green(Help[i].OptionalArgs),
		})
	}

	table.Render()
}

func SendRequest[T, U any](ctx context.Context, request T, method string, path string) (int, U) {
	data, err := jsoniter.Marshal(request)
	if err != nil {
		fmt.Println(red("Ошибка при парсинге json"))
		os.Exit(1)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("http://%s/%s", serverAddress, path), bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(red("Ошибка при создании запроса"))
		os.Exit(1)
	}

	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(red("Ошибка при выполнении запроса"))
		fmt.Println(err)
		os.Exit(1)
	}

	var response U
	if err := jsoniter.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println(red("Ошибка при парсинге json"))
		os.Exit(1)
	}

	return resp.StatusCode, response
}
func main() {
	username := MustCreateUsername()

	if len(os.Args) < 2 {
		fmt.Printf(red("Кажется вы забыли указать метод\n"))
		os.Exit(1)
	}

	method := os.Args[1]

	if !slices.Contains(availableMethods, method) {
		fmt.Println(red(fmt.Sprintf("Метод %s не поддерживается", method)))
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), serverTimeout)
	defer cancel()

	switch method {
	case methodHelp:
		PrintResponseForHelp()
	case methodPostTopic:
		if !RequireContains(argTopic) {
			fmt.Println(red("Для выполнения метода post вы должны указать флаг topic"))
			os.Exit(1)
		}

		topic := ParseArg(argTopic)
		if topic == "" {
			fmt.Println(red("Выполнение запроса post с пустым значением topic невозможно"))
			os.Exit(1)
		}

		request := models.PostTopicRequest{
			Username: username,
			Topic:    topic,
		}

		_, resp := SendRequest[models.PostTopicRequest, models.PostTopicResponse](
			ctx, request,
			http.MethodPost, "topics",
		)

		PrintResponseForPostTopic(resp, nil)
	case methodDeleteTopic:
		if !RequireContains(argTopic) {
			fmt.Println(red("Для выполнения метода delete вы должны указать флаг topic"))
			os.Exit(1)
		}

		topic := ParseArg(argTopic)
		if topic == "" {
			fmt.Println(red("Выполнение запроса delete с пустым значением topic невозможно"))
			os.Exit(1)
		}

		request := models.DeleteTopicRequest{
			Username: username,
			Topic:    topic,
		}

		_, resp := SendRequest[models.DeleteTopicRequest, models.DeleteTopicResponse](
			ctx, request,
			http.MethodDelete, "topics",
		)

		PrintResponseForDeleteTopic(resp, nil)
	case methodListTopics:
		request := models.ListTopicsRequest{
			Username: username,
		}

		_, resp := SendRequest[models.ListTopicsRequest, models.ListTopicsResponse](
			ctx, request,
			http.MethodGet, "topics",
		)

		PrintResponseForListTopics(resp, nil)
	case methodPostLink:
		if !RequireContains(argTopic, argLink) {
			fmt.Println(red("Для выполнения метода post вы должны указать флаг link"))
			os.Exit(1)
		}

		link := ParseArg(argLink)
		if link == "" {
			fmt.Println(red("Выполнение запроса post с пустым значением link невозможно"))
			os.Exit(1)
		}

		topic := ParseArg(argTopic)
		if topic == "" {
			fmt.Println(red("Выполнение запроса post с пустым значением topic невозможно"))
			os.Exit(1)
		}

		alias := ParseArg(argAlias)

		request := models.PostLinkRequest{
			Username: username,
			Topic:    topic,
			Link:     link,
			Alias:    alias,
		}

		_, resp := SendRequest[models.PostLinkRequest, models.PostLinkResponse](
			ctx, request,
			http.MethodPost, "links",
		)

		PrintResponseForPostLink(resp, nil)
	case methodPickLink:
		if !RequireContains(argTopic, argAlias) {
			fmt.Println(red("Для выполнения метода pick вы должны указать флаг alias"))
			os.Exit(1)
		}

		topic := ParseArg(argTopic)
		if topic == "" {
			fmt.Println(red("Выполнение запроса pick с пустым значением topic невозможно"))
			os.Exit(1)
		}

		alias := ParseArg(argAlias)
		if alias == "" {
			fmt.Println(red("Выполнение запроса pick с пустым значением alias невозможно"))
			os.Exit(1)
		}

		request := models.PickLinkRequest{
			Username: username,
			Topic:    topic,
			Alias:    alias,
		}

		_, resp := SendRequest[models.PickLinkRequest, models.PickLinkResponse](
			ctx, request,
			http.MethodGet, "links",
		)

		PrintResponseForPickLink(resp, nil)
	case methodListLinks:
		if !RequireContains(argTopic) {
			fmt.Println(red("Для выполнения метода list вы должны указать флаг topic"))
			os.Exit(1)
		}

		topic := ParseArg(argTopic)
		if topic == "" {
			fmt.Println(red("Выполнение запроса list с пустым значением topic невозможно"))
			os.Exit(1)
		}

		request := models.ListLinksRequest{
			Username: username,
			Topic:    topic,
		}

		_, resp := SendRequest[models.ListLinksRequest, models.ListLinksResponse](
			ctx, request,
			http.MethodGet, "links/list",
		)

		PrintResponseForListLinks(topic, resp, nil)
	case methodDeleteLink:
		if !RequireContains(argTopic, argAlias) {
			fmt.Println(red("Для выполнения метода delete вы должны указать флаг alias"))
			os.Exit(1)
		}

		topic := ParseArg(argTopic)
		if topic == "" {
			fmt.Println(red("Выполнение запроса delete с пустым значением topic невозможно"))
			os.Exit(1)
		}

		alias := ParseArg(argAlias)
		if alias == "" {
			fmt.Println(red("Выполнение запроса delete с пустым значением alias невозможно"))
			os.Exit(1)
		}

		request := models.DeleteLinkRequest{
			Username: username,
			Topic:    topic,
			Alias:    alias,
		}

		_, resp := SendRequest[models.DeleteLinkRequest, models.DeleteLinkResponse](
			ctx, request,
			http.MethodDelete, "links",
		)

		PrintResponseForDeleteLink(resp, nil)
	}
}
