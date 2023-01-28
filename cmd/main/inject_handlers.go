package main

import (
	"github.com/google/wire"
	"telegram-api/internal/infrastructure/handler"
)

var handlerSet = wire.NewSet(
	handler.NewMessageFormer,
	handler.NewCommandHandler,
	handler.NewInlineMessageHandler,
	handler.NewCustomMessageHandler,
)
