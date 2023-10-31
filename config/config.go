package config

import (
	"os"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/opensourceways/xihe-script/infrastructure/message"
	"github.com/opensourceways/xihe-script/utils"
)

type Configuration struct {
	Matchs    []Match        `json:"matchs"    required:"true"`
	Message   message.Config `json:"message"`
	Endpoint  string         `json:"endpoint"  required:"true"`
	MaxRetry  int            `json:"max_retry" required:"true"`
	OBSConfig OBSConfig      `json:"obs_config" required:"true"`
}

func (cfg *Configuration) validate() error {
	if err := utils.CheckConfig(cfg, ""); err != nil {
		return err
	}

	return cfg.Message.Validate()
}

func (cfg *Configuration) setDefault() {
	if cfg.MaxRetry <= 0 {
		cfg.MaxRetry = 10
	}
}

type Match struct {
	Id                        string `json:"competition_id"                required:"true"`
	AnswerFinalPath           string `json:"answer_final_path"`
	AnswerPreliminaryPath     string `json:"answer_preliminary_path"`
	FidWeightsFinalPath       string `json:"fid_weights_final_path"`
	FidWeightsPreliminaryPath string `json:"fid_weights_preliminary_path"`
	RealFinalPath             string `json:"real_final_path"`
	RealPreliminaryPath       string `json:"real_preliminary_path"`
	Pos                       int    `json:"pos"`
	Cls                       int    `json:"cls"`
	Prefix                    string `json:"prefix" required:"true"`
}

type OBSConfig struct {
	Prefix    string `json:"prefix"`
	Bucket    string `json:"bucket"         required:"true"`
	Endpoint  string `json:"endpoint"       required:"true"`
	AccessKey string `json:"access_key"     required:"true"`
	SecretKey string `json:"secret_key"     required:"true"`
}

func (m *Match) GetAnswerFinalPath() string {
	return m.AnswerFinalPath
}

func (m *Match) GetAnswerPreliminaryPath() string {
	return m.AnswerPreliminaryPath
}

func (m *Match) GetPrefix() string {
	return m.Prefix
}

func (m *Match) GetFidWeightsFinalPath() string {
	return m.FidWeightsFinalPath
}

func (m *Match) GetFidWeightsPreliminaryPath() string {
	return m.FidWeightsPreliminaryPath
}

func (m *Match) GetRealFinalPath() string {
	return m.RealFinalPath
}

func (m *Match) GetRealPreliminaryPath() string {
	return m.RealPreliminaryPath
}

func (m *Match) GetPos() int {
	return m.Pos
}

func (m *Match) GetCls() int {
	return m.Cls
}

func (m *Match) GetCompetitionId() string {
	return m.Id
}

func (cfg *Configuration) GetMatch(id string) *Match {
	for k := range cfg.Matchs {
		m := &cfg.Matchs[k]
		if strings.EqualFold(m.Id, id) {
			return m
		}
	}

	return nil
}

func loadFromYaml(path string, cfg interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, cfg)
}

func LoadConfig(path string, cfg *Configuration) error {
	if err := loadFromYaml(path, cfg); err != nil {
		return err
	}

	cfg.setDefault()

	return cfg.validate()
}
