package settings

import (
	"context"
	"encoding/json"
	"fmt"

	domainsettings "github.com/hajimohammadinet/dabir/internal/domain/settings"
)

type GetPublicSettingsUseCase struct {
	settingsRepo domainsettings.Repository
}

type PublicSettingsOutput struct {
	OrganizationName string                      `json:"organization_name"`
	LetterConfig     domainsettings.LetterConfig `json:"letter_config"`
}

func NewGetPublicSettingsUseCase(settingsRepo domainsettings.Repository) *GetPublicSettingsUseCase {
	return &GetPublicSettingsUseCase{
		settingsRepo: settingsRepo,
	}
}

func (uc *GetPublicSettingsUseCase) Execute(ctx context.Context) (*PublicSettingsOutput, error) {
	output := &PublicSettingsOutput{
		OrganizationName: "Dabir",
		LetterConfig:     domainsettings.DefaultLetterConfig(),
	}

	orgSetting, err := uc.settingsRepo.Get(ctx, domainsettings.KeyOrganizationName)
	if err != nil {
		return nil, err
	}

	if orgSetting != nil {
		if err := json.Unmarshal(orgSetting.Value, &output.OrganizationName); err != nil {
			return nil, fmt.Errorf("failed to unmarshal organization name: %w", err)
		}
	}

	letterSetting, err := uc.settingsRepo.Get(ctx, domainsettings.KeyLetterConfig)
	if err != nil {
		return nil, err
	}

	if letterSetting != nil {
		var storedConfig domainsettings.LetterConfig

		if err := json.Unmarshal(letterSetting.Value, &storedConfig); err != nil {
			return nil, fmt.Errorf("failed to unmarshal letter config: %w", err)
		}

		output.LetterConfig = domainsettings.NormalizeLetterConfig(storedConfig)
	}

	return output, nil
}
