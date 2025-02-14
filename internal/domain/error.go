package domain

import (
	resp "avito_shop/pkg/http/responses"
	pkgerr "avito_shop/pkg/pkgerror"
	"errors"
)

var (
	ErrBadRequest       = errors.New("bad request")
	ErrUserNotFound     = errors.New("user not found")
	ErrLowBalance       = errors.New("not enough coins on the balance")
	ErrSelfSending      = errors.New("can't send coins to yourself")
	ErrMerchNotFound    = errors.New("merch not found")
	ErrUserExists       = errors.New("user already exist")
	ErrInvalidAuthToken = errors.New("invalid auth token")
	ErrUnauthorized     = errors.New("unauthorized")
)

func HandleResult(err error, r any) resp.Response {
	if err == nil {
		return resp.OK(r)
	}

	err = pkgerr.UnwrapAll(err)

	switch {
	case errors.Is(err, ErrUnauthorized),
		errors.Is(err, ErrInvalidAuthToken),
		errors.Is(err, ErrUserExists):
		return resp.Unauthorized(err)
	case errors.Is(err, ErrBadRequest),
		errors.Is(err, ErrLowBalance),
		errors.Is(err, ErrMerchNotFound),
		errors.Is(err, ErrUserNotFound),
		errors.Is(err, ErrSelfSending):
		return resp.BadRequest(err)
	default:
		return resp.Unknown(err)
	}
}
