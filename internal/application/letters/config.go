package letters

import (
	"context"
	"encoding/json"

	domainsettings "github.com/hajimohammadinet/dabir/internal/domain/settings"
)

type LetterConfigProvider struct {
	settingsRepo domainsettings.Repository
}

func NewLetterConfigProvider(settingsRepo domainsettings.Repository) *LetterConfigProvider {
	return &LetterConfigProvider{
		settingsRepo: settingsRepo,
	}
}

func (p *LetterConfigProvider) Get(ctx context.Context) LetterNumberConfig {
	domainConfig := domainsettings.DefaultLetterConfig()

	setting, err := p.settingsRepo.Get(ctx, domainsettings.KeyLetterConfig)
	if err == nil && setting != nil {
		var stored domainsettings.LetterConfig
		if err := json.Unmarshal(setting.Value, &stored); err == nil {
			domainConfig = domainsettings.NormalizeLetterConfig(stored)
		}
	}

	return LetterNumberConfig{
		Mode: NumberingMode(domainConfig.NumberingMode),

		Prefix:  domainConfig.NumberPrefix,
		Padding: domainConfig.NumberPadding,

		YearlyPrefixDigits:  domainConfig.YearlyPrefixDigits,
		YearlySerialPadding: domainConfig.YearlySerialPadding,
		YearlySeparator:     domainConfig.YearlySeparator,
		YearSource:          domainConfig.YearSource,
	}
}
