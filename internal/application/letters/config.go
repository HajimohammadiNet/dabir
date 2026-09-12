package letters

import (
	"context"
	"encoding/json"

	"github.com/hajimohammadinet/dabir/internal/domain/letter"
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

func (p *LetterConfigProvider) Get(ctx context.Context, direction letter.Direction) LetterNumberConfig {
	domainConfig := domainsettings.DefaultLetterConfig()
	settingKey := domainsettings.KeyLetterConfig
	if direction == letter.DirectionOutgoing {
		settingKey = domainsettings.KeyOutgoingLetterConfig
	}

	setting, err := p.settingsRepo.Get(ctx, settingKey)
	if err == nil && setting == nil && direction == letter.DirectionOutgoing {
		setting, err = p.settingsRepo.Get(ctx, domainsettings.KeyLetterConfig)
	}
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
