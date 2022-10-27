package xlog

import (
	"bytes"
	"testing"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestProvide(t *testing.T) {
	var (
		assert = assert.New(t)

		expected = log.NewJSONLogger(new(bytes.Buffer))

		actual *zap.Logger
		app    = fxtest.New(
			t,
			fx.Provide(Provide(expected)),
			fx.Populate(&actual),
		)
	)

	app.RequireStart()
	assert.Equal(expected, actual)
	app.RequireStop()
}
