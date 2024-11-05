package contexts

import (
	"context"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

const tmaInitDataKey CtxKey = "tma-init-data"

func WithTMAInitData(ctx context.Context, initData initdata.InitData) context.Context {
	return context.WithValue(ctx, tmaInitDataKey, initData)
}

func TMAInitData(ctx context.Context) (initdata.InitData, bool) {
	initData, ok := ctx.Value(tmaInitDataKey).(initdata.InitData)
	return initData, ok
}
