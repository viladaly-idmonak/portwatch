package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/port/portwatch/internal/scannerg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	s := scanner.New(cfg.Protocol, cfg.Hosts)
	l := logger.New(os.Stdout)
	w := watcher.New(time.Duration(cfg.IntervalSecs)*time.Second, s, l)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		w.Stop()
	}()

	fmt.Printf("portwatch starting (interval=%ds, protocol=%s, hosts=%v)\n",
		cfg.IntervalSecs, cfg.Protocol, cfg.Hosts)

	if err := w.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "watcher error: %v\n", err)
		os.Exit(1)
	}
}
