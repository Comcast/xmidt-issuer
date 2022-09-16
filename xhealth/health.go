package xhealth

import (
	"context"
	"errors"

	health "github.com/InVisionApp/go-health"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/xmidt-org/themis/xlog"
)

// Options holds the available configuration options for the health infrastructure
type Options struct {
	// DisableLogging controls whether the created health service logs anything.  The default is false,
	// meaning logging is enabled.
	DisableLogging bool

	// Custom is an optional map passed to NewHandler that is included in all responses to health checks
	Custom map[string]interface{}
}

// New constructs an IHealth instance for the given environment.  If either the DisableLogging option field
// is set or the given logger is nil, logging will be disabled on the returned health object.  The listener
// is optional.
func New(o Options, logger log.Logger, listener health.IStatusListener) (health.IHealth, error) {
	h := health.New()
	if o.DisableLogging || logger == nil {
		h.DisableLogging()
	} else {
		h.Logger = NewHealthLoggerAdapter(logger)
	}

	h.StatusListener = listener
	return h, nil
}

// OnStart returns an uber/fx Lifecycle hook for startup
func OnStart(logger log.Logger, h health.IHealth) func(context.Context) error {
	return func(_ context.Context) error {
		logger.Log(
			level.Key(), level.InfoValue(),
			xlog.MessageKey(), "health service starting",
		)

		return h.Start()
	}
}

// OnStop returns an uber/fx Lifecycle hook for shutdown
func OnStop(logger log.Logger, h health.IHealth) func(context.Context) error {
	return func(_ context.Context) error {
		logger.Log(
			level.Key(), level.InfoValue(),
			xlog.MessageKey(), "health service stopping",
		)

		err := h.Stop()
		if errors.Is(err, health.ErrAlreadyStopped) {
			logger.Log(
				level.Key(), level.InfoValue(),
				xlog.MessageKey(), "health service already stopped or not running",
			)

			return nil
		}

		return err
	}
}
