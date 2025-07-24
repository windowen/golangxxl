package mctx

import (
	"context"
	"strconv"

	"liveJob/pkg/constant"
	"liveJob/pkg/tools/errs"
	"liveJob/pkg/tools/utils"
)

func HaveOpUser(ctx context.Context) bool {
	return ctx.Value(constant.RpcOpUserID) != nil
}

func Check(ctx context.Context) (int, int32, error) {
	opUserId, ok := ctx.Value(constant.RpcOpUserID).(int)
	if !ok || opUserId == 0 {
		return 0, 0, errs.ErrNoPermission.Wrap("opuser_id_empty")
	}

	opUserTypeArr, ok := ctx.Value(constant.RpcOpUserType).([]string)
	if !ok {
		return 0, 0, errs.ErrNoPermission.Wrap("missing_user_type")
	}
	if len(opUserTypeArr) == 0 {
		return 0, 0, errs.ErrNoPermission.Wrap("user type empty")
	}
	userType, err := strconv.Atoi(opUserTypeArr[0])
	if err != nil {
		return 0, 0, errs.ErrNoPermission.Wrap("user type invalid " + err.Error())
	}
	if !(userType == constant.AdminUser || userType == constant.NormalUser) {
		return 0, 0, errs.ErrNoPermission.Wrap("user type invalid")
	}
	return opUserId, int32(userType), nil
}

func CheckUser(ctx context.Context) (int, error) {
	userID, userType, err := Check(ctx)
	if err != nil {
		return 0, err
	}
	if userType != constant.NormalUser {
		return 0, errs.ErrNoPermission.Wrap("not user")
	}
	return userID, nil
}

func GetOpUserID(ctx context.Context) int {
	userID, _ := ctx.Value(constant.OpUserId).(int)
	return userID
}

func GetCountryCode(ctx context.Context) string {
	countryCode, _ := ctx.Value(constant.CountryCode).(string)
	return countryCode
}

func GetLanguage(ctx context.Context) string {
	language, _ := ctx.Value(constant.Language).(string)
	return language
}

func GetUserType(ctx context.Context) (int, error) {
	userTypeArr, _ := ctx.Value(constant.RpcOpUserType).([]string)
	userType, err := strconv.Atoi(userTypeArr[0])
	if err != nil {
		return 0, errs.ErrNoPermission.Wrap("user type invalid " + err.Error())
	}
	return userType, nil
}

func WithOpUserID(ctx context.Context, opUserID string, userType int) context.Context {
	headers, _ := ctx.Value(constant.RpcCustomHeader).([]string)
	ctx = context.WithValue(ctx, constant.RpcOpUserID, opUserID)
	ctx = context.WithValue(ctx, constant.RpcOpUserType, []string{strconv.Itoa(userType)})
	if utils.IndexOf(constant.RpcOpUserType, headers...) < 0 {
		ctx = context.WithValue(ctx, constant.RpcCustomHeader, append(headers, constant.RpcOpUserType))
	}
	return ctx
}

func WithApiToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, constant.CtxApiToken, token)
}
