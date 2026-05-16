package letters

import (
	"context"
	"encoding/json"

	domainsettings "github.com/hajimohammadinet/dabir/internal/domain/settings"
)

type LetterConfigProvider struct {
	settingsRepo domainsettings.Repository
}

type storedLetterConfig struct {
	NumberPrefix  string `json:"number_prefix"`
	NumberPadding int    `json:"number_padding"`
}

func NewLetterConfigProvider(settingsRepo domainsettings.Repository) *LetterConfigProvider {
	return &LetterConfigProvider{
		settingsRepo: settingsRepo,
	}
}

func (p *LetterConfigProvider) Get(ctx context.Context) LetterNumberConfig {
	cfg := LetterNumberConfig{
		Prefix:  "DABIR",
		Padding: 6,
	}

	setting, err := p.settingsRepo.Get(ctx, domainsettings.KeyLetterConfig)
	if err != nil || setting == nil {
		return cfg
	}

	var stored storedLetterConfig
	if err := json.Unmarshal(setting.Value, &stored); err != nil {
		return cfg
	}

	if stored.NumberPrefix != "" {
		cfg.Prefix = stored.NumberPrefix
	}

	if stored.NumberPadding > 0 {
		cfg.Padding = stored.NumberPadding
	}

	return cfg
}
