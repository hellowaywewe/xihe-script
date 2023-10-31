package score

import (
	"github.com/opensourceways/xihe-script/config"
	"github.com/opensourceways/xihe-script/infrastructure/message"
)

type EvaluateScore interface {
	Evaluate(*message.MatchFields, *config.OBSConfig) ([]byte, error)
}
