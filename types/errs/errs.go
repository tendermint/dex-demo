package errs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/dex-demo/storeutil"
)

const (
	_ sdk.CodeType = iota
	CodeNotFound
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

func ErrInvalidArgument(msg string) sdk.Error {
	return newErrWithUEXCodespace(CodeInvalidArgument, msg)
}

func ErrMarshalFailure(msg string) sdk.Error {
	return newErrWithUEXCodespace(CodeMarshalFailure, msg)
}

func ErrUnmarshalFailure(msg string) sdk.Error {
	return newErrWithUEXCodespace(CodeUnmarshalFailure, msg)
}

func WrapNotFound(err error) sdk.Error {
	if err == nil {
		return nil
	}
	if err == storeutil.ErrStoreKeyNotFound {
		return ErrNotFound(err.Error())
	}
	return sdk.ErrInternal(err.Error())
}

func WrapOrNil(err error) sdk.Error {
	if err == nil {
		return nil
	}

	return sdk.ErrInternal(err.Error())
}

func ErrOrBlankResult(err sdk.Error) sdk.Result {
	if err == nil {
		return sdk.Result{}
	}

	return err.Result()
}
