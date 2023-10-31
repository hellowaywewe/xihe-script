package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/opensourceways/xihe-grpc-protocol/grpc/client"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-script/app"
	"github.com/opensourceways/xihe-script/config"
	"github.com/opensourceways/xihe-script/infrastructure/message"
	"github.com/opensourceways/xihe-script/infrastructure/score"
)

type options struct {
	service     liboptions.ServiceOptions
	enableDebug bool
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) (options, error) {
	var o options

	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.enableDebug, "enable_debug", false,
		"whether to enable debug model.",
	)

	err := fs.Parse(args)

	return o, err
}

func main() {
	logrusutil.ComponentInit("xihe")
	log := logrus.NewEntry(logrus.StandardLogger())

	o, err := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err != nil {
		logrus.Fatalf("new options failed, err:%s", err.Error())
	}

	if err := o.Validate(); err != nil {
		logrus.Fatalf("Invalid options, err:%s", err.Error())
	}

	if o.enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled.")
	}

	cfg := new(config.Configuration)
	if err := config.LoadConfig(o.service.ConfigFile, cfg); err != nil {
		logrus.Fatalf("load config, err:%s", err.Error())
	}

	if err := os.Remove(o.service.ConfigFile); err != nil {
		logrus.Fatalf("config file delete failed, err:%s", err.Error())
	}

	cli, err := client.NewCompetitionClient(cfg.Endpoint)
	if err != nil {
		logrus.Errorf("init rpc server err: %v", err)

		return
	}

	defer cli.Disconnect()

	run(newHandler(cfg, cli, log), &cfg.Message, log)
}

func newHandler(cfg *config.Configuration, cli *client.CompetitionClient, log *logrus.Entry) *handler {
	return &handler{
		maxRetry: cfg.MaxRetry,
		log:      log,
		evaluate: app.NewEvaluateService(score.NewEvaluateScore(os.Getenv("EVALUATE")), cfg),
		match:    cfg,
		cli:      cli,
	}
}

func run(h *handler, cfg *message.Config, log *logrus.Entry) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	defer wg.Wait()

	called := false
	ctx, done := context.WithCancel(context.Background())

	defer func() {
		if !called {
			called = true
			done()
		}
	}()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			log.Info("receive done. exit normally")
			return

		case <-sig:
			log.Info("receive exit signal")
			done()
			called = true
			os.Exit(1)
			return
		}
	}(ctx)

	if err := message.Subscribe(ctx, h, cfg, log); err != nil {
		log.Errorf("subscribe failed, err:%v", err)
	}
}
