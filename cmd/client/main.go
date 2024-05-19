package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"slices"
	"strings"
	"time"

	linkerV2 "github.com/Sleeps17/linker-protos/gen/go/linker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	methodPostTopic   = "post_topic"
	methodDeleteTopic = "delete_topic"
	methodListTopics  = "list_topics"
	methodPostLink    = "post_link"
	methodPickLink    = "pick_link"
	methodListLinks   = "list_links"
	methodDeleteLink  = "delete_links"

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
	}

	serverAddress = "localhost:4404"
	serverTimeout = 5 * time.Second
)

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

func PrintResponseForPostTopic(resp *linkerV2.PostTopicResponse, err error) {
	var output string
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			output = yellow(fmt.Sprintf("Неверные входные данные: %s", st.Message()))
		} else {
			output = red(fmt.Sprintf("Произошла ошибка при выполнении запроса: %s\n", st.Message()))
		}
	} else {
		output = green(fmt.Sprintf("Запрос успешно выполнен, TopicId = %d\n", resp.TopicId))
	}
	fmt.Println(output)
}

func PrintResponseForDeleteTopic(resp *linkerV2.DeleteTopicResponse, err error) {
	var output string
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			output = yellow(fmt.Sprintf("Неверные входные данные: %s", st.Message()))
		} else {
			output = red(fmt.Sprintf("Произошла ошибка при выполнении запроса: %s\n", st.Message()))
		}
	} else {
		output = green(fmt.Sprintf("Запрос успешно выполнен, TopicId = %d\n", resp.TopicId))
	}
	fmt.Println(output)
}

func PrintResponseForListTopics(resp *linkerV2.ListTopicsResponse, err error) {
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

func PrintResponseForPostLink(resp *linkerV2.PostLinkResponse, err error) {
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

func PrintResponseForPickLink(resp *linkerV2.PickLinkResponse, err error) {
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

func PrintResponseForDeleteLink(resp *linkerV2.DeleteLinkResponse, err error) {
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

func PrintResponseForListLinks(topic string, resp *linkerV2.ListLinksResponse, err error) {
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

func main() {
	username := MustCreateUsername()

	if len(os.Args) < 2 {
		fmt.Printf(red("Кажется вы забыли указать метод\n"))
		os.Exit(1)
	}

	method := os.Args[1]

	if !slices.Contains(availableMethods, method) {
		fmt.Println(fmt.Sprintf(red("Метод %s не поддерживается\n", method)))
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), serverTimeout)
	defer cancel()

	dial, err := grpc.DialContext(ctx, serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(red("Что-то пошло не так"))
		os.Exit(1)
	}

	client := linkerV2.NewLinkerClient(dial)

	switch method {
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

		request := &linkerV2.PostTopicRequest{
			Username: username,
			Topic:    topic,
		}

		resp, err := client.PostTopic(ctx, request)
		PrintResponseForPostTopic(resp, err)
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

		request := &linkerV2.DeleteTopicRequest{
			Username: username,
			Topic:    topic,
		}

		resp, err := client.DeleteTopic(ctx, request)
		PrintResponseForDeleteTopic(resp, err)
	case methodListTopics:
		request := &linkerV2.ListTopicsRequest{
			Username: username,
		}

		resp, err := client.ListTopics(ctx, request)
		PrintResponseForListTopics(resp, err)
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

		request := &linkerV2.PostLinkRequest{
			Username: username,
			Topic:    topic,
			Link:     link,
			Alias:    alias,
		}

		resp, err := client.PostLink(ctx, request)
		PrintResponseForPostLink(resp, err)
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

		request := &linkerV2.PickLinkRequest{
			Username: username,
			Topic:    topic,
			Alias:    alias,
		}

		resp, err := client.PickLink(ctx, request)
		PrintResponseForPickLink(resp, err)
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

		request := &linkerV2.ListLinksRequest{
			Username: username,
			Topic:    topic,
		}

		resp, err := client.ListLinks(ctx, request)
		PrintResponseForListLinks(topic, resp, err)
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

		request := &linkerV2.DeleteLinkRequest{
			Username: username,
			Topic:    topic,
			Alias:    alias,
		}

		resp, err := client.DeleteLink(ctx, request)
		PrintResponseForDeleteLink(resp, err)
	}
}
