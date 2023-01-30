package model

import "context"

const ContextUserKey string = "ContextUserKey"
const ContextChatIDKey string = "ContextChatIDKey"
const ContextMessageIDKey string = "ContextMessageIDKey"

func UserExist(ctx context.Context) bool {
	return ctx.Value(ContextUserKey) != nil
}

func GetCurrentUser(ctx context.Context) User {
	return ctx.Value(ContextUserKey).(User)
}

func GetCurrentChatID(ctx context.Context) int64 {
	return ctx.Value(ContextChatIDKey).(int64)
}

func GetCurrentMessageID(ctx context.Context) int {
	return ctx.Value(ContextMessageIDKey).(int)
}
