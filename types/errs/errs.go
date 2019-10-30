package errs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	_ sdk.CodeType = iota
	CodeNotFound
	CodeAlreadyExists
	CodeInvalidArgument
	CodeMarshalFailure
	CodeUnmarshalFailure

	CodespaceUEX sdk.CodespaceType = "dex-demo"
)

func newErrWithUEXCodespace(code sdk.CodeType, msg string) sdk.Error {
	return sdk.NewError(CodespaceUEX, code, msg)
}

func ErrNotFound(msg string) sdk.Error {
	return newErrWithUEXCodespace(CodeNotFound, msg)
}

func ErrAlreadyExists(msg string) sdk.Error {
	return newErrWithUEXCodespace(CodeAlreadyExists, msg)
}

func ErrInvalidArgument(msg string) sdk.Error {
	return newErrWithUEXCodespace(CodeInvalidArgument, msg)
}

func ErrMarshalFailure(msg string) sdk.Error {
	return newErrWithUEXCodespace(CodeMarshalFailure, msg)
}

func ErrUnmarshalFailure(msg string) sdk.Error {
	return newErrWithUEXCodespace(CodeUnmarshalFailure, msg)
}

func ErrOrBlankResult(err sdk.Error) sdk.Result {
	if err == nil {
		return sdk.Result{}
	}

	return err.Result()
}
