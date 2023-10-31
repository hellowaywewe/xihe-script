package score

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/opensourceways/xihe-script/config"
	"github.com/opensourceways/xihe-script/domain/score"
	"github.com/opensourceways/xihe-script/infrastructure/message"
	"github.com/opensourceways/xihe-script/utils"
)

type evaluateImpl struct {
	evaluate string
}

func NewEvaluateScore(evaluate string) score.EvaluateScore {
	return &evaluateImpl{
		evaluate: evaluate,
	}
}

func (s *evaluateImpl) Evaluate(col *message.MatchFields, cfg *config.OBSConfig) (data []byte, err error) {
	args := []string{"python3", s.evaluate, "--pred_path", col.Path, "--true_path", col.AnswerPath, "--cls", strconv.Itoa(col.Cls), "--pos", strconv.Itoa(col.Pos)}
	input := fmt.Sprintf("%s+%s+%s+%s\n", cfg.AccessKey, cfg.SecretKey, cfg.Bucket, cfg.Endpoint)
	var b bytes.Buffer
	b.Write([]byte(input))
	c, err, _ := utils.RunCmd(&b, args...)
	if err != nil {
		return
	}

	data = []byte(c)
	data = bytes.ReplaceAll(bytes.TrimSpace(data), []byte(`'`), []byte(`"`))

	return
}
