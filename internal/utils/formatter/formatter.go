package formatter

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"io"
	"os"
)

const (
	errorCode = 1
)

var (
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

type Formatter struct {
	out io.Writer
}

func New(out io.Writer) *Formatter {
	return &Formatter{
		out: out,
	}
}

func (f *Formatter) SuccessString(msg string) string {
	return green(msg)
}

func (f *Formatter) WarningString(msg string) string {
	return yellow(msg)
}

func (f *Formatter) ErrorString(msg string) string {
	return red(msg)
}

func (f *Formatter) Success(msg string) {
	_, _ = fmt.Fprintln(f.out, green(msg))
}

func (f *Formatter) Warning(msg string) {
	_, _ = fmt.Fprintln(f.out, yellow(msg))
}

func (f *Formatter) Error(msg string) {
	_, _ = fmt.Fprintln(f.out, red(msg))
	os.Exit(errorCode)
}

func (f *Formatter) Successf(format string, a ...any) {
	_, _ = fmt.Fprintf(f.out, green(format+"\n"), a...)
}

func (f *Formatter) Warningf(format string, a ...any) {
	_, _ = fmt.Fprintf(f.out, yellow(format+"\n"), a...)
}

func (f *Formatter) Errorf(format string, a ...any) {
	_, _ = fmt.Fprintf(f.out, red(format+"\n"), a...)
	os.Exit(errorCode)
}

func (f *Formatter) SuccessTable(headers []string, values ...[]string) {
	table := tablewriter.NewWriter(f.out)
	table.SetHeader(headers)
	table.AppendBulk(values)
	table.Render()
}
