export type NumberingMode = "fixed_prefix" | "jalali_yearly" | "manual";
export type YearSource = "letter_date" | "created_at";

export type PublicSettings = {
  organization_name: string;
  letter_config: {
    numbering_mode: NumberingMode;

    number_prefix: string;
    number_padding: number;

    yearly_prefix_digits: number;
    yearly_serial_padding: number;
    yearly_separator: string;
    year_source: YearSource;
  };
};