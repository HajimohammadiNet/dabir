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
	OrganizationName string             `json:"organization_name"`
	LetterConfig     PublicLetterConfig `json:"letter_config"`
}

type PublicLetterConfig struct {
	NumberPrefix  string `json:"number_prefix"`
	NumberPadding int    `json:"number_padding"`
}

func NewGetPublicSettingsUseCase(settingsRepo domainsettings.Repository) *GetPublicSettingsUseCase {
	return &GetPublicSettingsUseCase{
		settingsRepo: settingsRepo,
	}
}

func (uc *GetPublicSettingsUseCase) Execute(ctx context.Context) (*PublicSettingsOutput, error) {
	output := &PublicSettingsOutput{
		OrganizationName: "Dabir",
		LetterConfig: PublicLetterConfig{
			NumberPrefix:  "DABIR",
			NumberPadding: 6,
		},
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
		if err := json.Unmarshal(letterSetting.Value, &output.LetterConfig); err != nil {
			return nil, fmt.Errorf("failed to unmarshal letter config: %w", err)
		}
	}

	return output, nil
}
