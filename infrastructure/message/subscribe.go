package message

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/sirupsen/logrus"
)

const (
	StyleCompetition = "1"
	TextCompetition  = "2"
	ImageCompetition = "3"
	LearnText        = "learn23-text"
	LeanImg          = "learn23-img"

	CompetitionPhaseFinal       = "final"
	CompetitionPhasePreliminary = "preliminary"
)

type Config struct {
	kfklib.Config

	Topics topics `json:"topics"  required:"true"`
}

type topics struct {
	Match string `json:"submission" required:"true"`
}

func Subscribe(ctx context.Context, impl MatchImpl, cfg *Config, log *logrus.Entry) error {
	if err := kfklib.Init(&cfg.Config, log, nil, "", true); err != nil {
		return err
	}

	defer kfklib.Exit()

	h := handler{impl}

	err := kfklib.Subscribe("xihe-script", h.handle, []string{cfg.Topics.Match})
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

type handler struct {
	impl MatchImpl
}

func (h *handler) handle(payload []byte, header map[string]string) error {
	body := MatchMessage{}
	if err := json.Unmarshal(payload, &body); err != nil {
		return err
	}

	m := h.impl.GetMatch(body.CompetitionId)
	if m == nil {
		return fmt.Errorf("unknown competition id: %s", body.CompetitionId)
	}

	switch m.GetCompetitionId() {
	case TextCompetition, ImageCompetition, LearnText, LeanImg:
		go h.evaluate(&body, m)

	case StyleCompetition:
		go h.calculate(&body, m)

	default:
		logrus.Errorf("unknown competition id: %s", m.GetCompetitionId())
	}

	return nil
}

func (h *handler) evaluate(body *MatchMessage, m MatchFieldImpl) {
	c := MatchFields{
		Path: m.GetPrefix() + "/" + strings.TrimPrefix(body.Path, "/"),
		Cls:  m.GetCls(),
		Pos:  m.GetPos(),
	}

	switch body.Phase {
	case CompetitionPhaseFinal:
		c.AnswerPath = m.GetAnswerFinalPath()

	case CompetitionPhasePreliminary:
		c.AnswerPath = m.GetAnswerPreliminaryPath()

	default:
		logrus.Errorf("unknown competition phase: %s", body.Phase)
	}

	if err := h.impl.Evaluate(body, &c); err != nil {
		logrus.Errorf(
			"evaluate failed, competition id:%s,user:%v",
			body.CompetitionId, body.UserId,
		)
	}
}

func (h *handler) calculate(body *MatchMessage, m MatchFieldImpl) {
	c := MatchFields{Path: m.GetPrefix() + "/" + strings.TrimPrefix(body.Path, "/")}

	switch body.Phase {
	case CompetitionPhaseFinal:
		c.FidWeightsPath = m.GetFidWeightsFinalPath()
		c.RealPath = m.GetRealFinalPath()

	case CompetitionPhasePreliminary:
		c.FidWeightsPath = m.GetFidWeightsPreliminaryPath()
		c.RealPath = m.GetRealPreliminaryPath()

	default:
		logrus.Errorf("unknown competition phase: %s", body.Phase)
	}

	if err := h.impl.Calculate(body, &c); err != nil {
		logrus.Errorf(
			"evaluate failed, competition id:%s,user:%v",
			body.CompetitionId, body.UserId,
		)
	}
}
