package pprofserver

import (
	"context"
	"github.com/chemax/url-shorter/logger"
)

func ExampleInit() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lg, err := logger.Init()
	if err != nil {
		return
	}
	Init(ctx, lg)
}
