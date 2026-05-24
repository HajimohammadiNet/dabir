package settings

import "time"

type Setting struct {
	Key       string
	Value     []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	KeyOrganizationName = "organization_name"
	KeyLetterConfig     = "letter_config"
)

type NumberingMode string

const (
	NumberingModeFixedPrefix  NumberingMode = "fixed_prefix"
	NumberingModeJalaliYearly NumberingMode = "jalali_yearly"
	NumberingModeManual       NumberingMode = "manual"
)

type YearSource string

const (
	YearSourceLetterDate YearSource = "letter_date"
	YearSourceCreatedAt  YearSource = "created_at"
)

type LetterConfig struct {
	NumberingMode NumberingMode `json:"numbering_mode"`

	NumberPrefix  string `json:"number_prefix"`
	NumberPadding int    `json:"number_padding"`

	YearlyPrefixDigits  int        `json:"yearly_prefix_digits"`
	YearlySerialPadding int        `json:"yearly_serial_padding"`
	YearlySeparator     string     `json:"yearly_separator"`
	YearSource          YearSource `json:"year_source"`
}

func DefaultLetterConfig() LetterConfig {
	return LetterConfig{
		NumberingMode: NumberingModeFixedPrefix,

		NumberPrefix:  "DABIR",
		NumberPadding: 6,

		YearlyPrefixDigits:  3,
		YearlySerialPadding: 4,
		YearlySeparator:     "-",
		YearSource:          YearSourceLetterDate,
	}
}

func NormalizeLetterConfig(config LetterConfig) LetterConfig {
	defaultConfig := DefaultLetterConfig()

	if config.NumberingMode == "" {
		config.NumberingMode = defaultConfig.NumberingMode
	}

	if config.NumberingMode != NumberingModeFixedPrefix &&
		config.NumberingMode != NumberingModeJalaliYearly &&
		config.NumberingMode != NumberingModeManual {
		config.NumberingMode = defaultConfig.NumberingMode
	}

	if config.NumberPrefix == "" {
		config.NumberPrefix = defaultConfig.NumberPrefix
	}

	if config.NumberPadding <= 0 {
		config.NumberPadding = defaultConfig.NumberPadding
	}

	if config.YearlyPrefixDigits <= 0 {
		config.YearlyPrefixDigits = defaultConfig.YearlyPrefixDigits
	}

	if config.YearlySerialPadding <= 0 {
		config.YearlySerialPadding = defaultConfig.YearlySerialPadding
	}

	if config.YearlySeparator == "" {
		config.YearlySeparator = defaultConfig.YearlySeparator
	}

	if config.YearSource == "" {
		config.YearSource = defaultConfig.YearSource
	}

	if config.YearSource != YearSourceLetterDate &&
		config.YearSource != YearSourceCreatedAt {
		config.YearSource = defaultConfig.YearSource
	}

	return config
}
