package app

import (
	"bytes"
	"encoding/json"

	"github.com/opensourceways/xihe-script/config"
	"github.com/opensourceways/xihe-script/domain/score"
	"github.com/opensourceways/xihe-script/infrastructure/message"
)

type scoreService struct {
	e   score.EvaluateScore
	cfg *config.Configuration
}

type EvaluateService interface {
	Evaluate(*message.MatchFields, *message.ScoreRes) error
}

func NewEvaluateService(e score.EvaluateScore, cfg *config.Configuration) EvaluateService {
	return &scoreService{
		e:   e,
		cfg: cfg,
	}
}

func (s *scoreService) Evaluate(col *message.MatchFields, res *message.ScoreRes) error {
	bys, err := s.e.Evaluate(col, &s.cfg.OBSConfig)
	if err != nil {
		return err
	}

	return json.NewDecoder(bytes.NewBuffer(bys)).Decode(res)
}
