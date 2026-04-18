package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/portwatch/internal/supervisor"
g, err := loadConfig()
	if err != nil : %v", err)
	}

	logger := log.New(os.Stdout, "[portwatch] ", log.LstdFlags)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	w, err := watcher.New(cfg.Interval, cfg.Protocol, logger)
	if err != nil {
		log.Fatalf("watcher: %v", err)
	}

	policy := supervisor.RestartPolicy{
		MaxRetries: cfg.MaxRestarts,
		Delay:      2 * time.Second,
	}
	sup := supervisor.New(policy, logger)

	if err := sup.Run(ctx, func(ctx context.Context) error {
		return w.Run(ctx)
	}); err != nil && err != context.Canceled {
		log.Fatalf("fatal: %v", err)
	}
}
