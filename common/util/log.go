package util

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tigapilarmandiri/perkakas/common/constant"
	"go.elastic.co/apm/module/apmzerolog/v2"
)

var (
	onceLog sync.Once
	Log     zerolog.Logger
)

func init() {
	onceLog.Do(func() {
		Log = log.With().Caller().Logger()
	})
}

func InitElk() {
	if os.Getenv(constant.IS_USE_ELK) == "true" {
		// apmzerolog.Writer will send log records with the level error or greater to Elastic APM.
		Log = zerolog.New(zerolog.MultiLevelWriter(os.Stdout, &apmzerolog.Writer{})).With().Caller().Logger()
		return
	}
}
