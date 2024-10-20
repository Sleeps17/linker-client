package parser

import (
	"errors"
	"fmt"
	"github.com/Sleeps17/linker-client/internal/models"
	"slices"
	"strings"
)

var (
	ErrEmptyMethod     = errors.New("empty method")
	ErrUnknownMethod   = errors.New("unknown method")
	ErrMissRequiredArg = errors.New("miss required arg")
)

const (
	emptyMethod = ""
)

var args = map[models.Arg]interface{}{}

var requiredAgs = map[models.Method][]models.Arg{
	models.PostTopic:   {models.Topic},
	models.DeleteTopic: {models.Topic},
	models.ListTopics:  {},
	models.PostLink:    {models.Topic, models.Link},
	models.PickLink:    {models.Topic, models.Alias},
	models.ListLinks:   {models.Topic},
	models.DeleteLink:  {models.Topic, models.Alias},
	models.Help:        {},
}

var optionalArgs = map[models.Method][]models.Arg{
	models.PostTopic:   {},
	models.DeleteTopic: {},
	models.ListTopics:  {},
	models.PostLink:    {models.Alias},
	models.PickLink:    {},
	models.ListLinks:   {},
	models.DeleteLink:  {},
	models.Help:        {},
}

type Parser struct {
	//corevalidator.Validator[any]
	availableMethods []string
}

func New(availableMethods []string) *Parser {
	return &Parser{
		availableMethods: availableMethods,
	}
}

func (p *Parser) ParseArgs(args []string) (models.Method, error) {
	if len(args) < 2 { //nolint:mnd
		return emptyMethod, ErrEmptyMethod
	}

	method := args[1]
	if !slices.Contains(p.availableMethods, method) {
		return emptyMethod, ErrUnknownMethod
	}

	if err := p.parseArgs(method, args[2:]); err != nil {
		return emptyMethod, err
	}

	return models.Method(method), nil
}

func (p *Parser) Arg(argName models.Arg) interface{} {
	return args[argName]
}

func (p *Parser) parseArgs(method string, cmdArgs []string) error {
	requiredArgs := requiredAgs[models.Method(method)]
	optionalArgs := optionalArgs[models.Method(method)]

	for _, requiredArg := range requiredArgs {
		found := false
		for _, cmdArg := range cmdArgs {
			if strings.HasPrefix(cmdArg, fmt.Sprintf("--%s=", requiredArg)) {
				args[requiredArg] = strings.TrimPrefix(cmdArg, fmt.Sprintf("--%s=", requiredArg))
				found = true
			}
		}

		if !found {
			return ErrMissRequiredArg
		}
	}

	for _, optionalArg := range optionalArgs {
		for _, cmdArg := range cmdArgs {
			if strings.HasPrefix(cmdArg, fmt.Sprintf("--%s=", optionalArg)) {
				args[optionalArg] = strings.TrimPrefix(cmdArg, fmt.Sprintf("--%s=", optionalArg))
			}
		}
	}

	return nil
}
