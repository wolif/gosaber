package jsonrpc20

import (
	"fmt"
)

// client error ----------------------------------------------------------------

type ClientError string

const BatchCallOverSize = ClientError("batch call too many")

func ClientErrorf(format string, a ...interface{}) ClientError {
	return ClientError(fmt.Sprintf(format, a...))
}

func (e ClientError) Error() string {
	return fmt.Sprintf("jsonrpc client, %s", string(e))
}

func IsClientError(e error) bool {
	_, ok := e.(ClientError)
	return ok
}

func IsBatchCallOverSize(e error) bool {
	err, ok := e.(ClientError)
	if !ok {
		return false
	}
	return err == BatchCallOverSize
}

// server error ----------------------------------------------------------------

type ServerError string

func ServerErrorf(format string, a ...interface{}) ServerError {
	return ServerError(fmt.Sprintf(format, a...))
}

func (e ServerError) Error() string {
	return fmt.Sprintf("jsonrpc server, %s", string(e))
}

func IsServerError(e error) bool {
	_, ok := e.(ServerError)
	return ok
}
