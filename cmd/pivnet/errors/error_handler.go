package errors

import (
	"errors"
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/pivotal-cf-experimental/go-pivnet"
	"github.com/pivotal-cf-experimental/go-pivnet/cmd/pivnet/printer"
)

var (
	ErrAlreadyHandled = errors.New("error already handled")
	RedFunc           = color.New(color.FgRed).SprintFunc()
)

//go:generate counterfeiter . ErrorHandler

type ErrorHandler interface {
	HandleError(err error) error
}

type errorHandler struct {
	format       string
	outputWriter io.Writer
	logWriter    io.Writer
}

func NewErrorHandler(
	format string,
	outputWriter io.Writer,
	logWriter io.Writer,
) ErrorHandler {
	return &errorHandler{
		format:       format,
		outputWriter: outputWriter,
		logWriter:    logWriter,
	}
}

func (h errorHandler) HandleError(err error) error {
	if err == nil {
		return nil
	}

	var message string

	switch err.(type) {
	case pivnet.ErrUnauthorized:
		message = fmt.Sprintf("Failed to authenticate - please provide valid API token")
	case pivnet.ErrNotFound:
		message = fmt.Sprintf("Pivnet error: %s", err.Error())
	default:
		message = err.Error()
	}

	coloredMessage := fmt.Sprintf(RedFunc(message))

	switch h.format {
	case printer.PrintAsJSON:
		e := h.printLogln(coloredMessage)
		if e != nil {
			return e
		}

		return ErrAlreadyHandled

	case printer.PrintAsYAML:
		e := h.printLogln(coloredMessage)
		if e != nil {
			return e
		}

		return ErrAlreadyHandled

	default:
		h.println(coloredMessage)
		return ErrAlreadyHandled
	}
}

func (h errorHandler) println(message string) error {
	_, err := h.outputWriter.Write([]byte(fmt.Sprintln(message)))
	return err
}

func (h errorHandler) printLogln(message string) error {
	_, err := h.logWriter.Write([]byte(fmt.Sprintln(message)))
	return err
}
