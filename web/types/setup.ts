export type SetupStatus = {
  initialized: boolean;
  setup_needed: boolean;
};

export type NumberingMode = "fixed_prefix" | "jalali_yearly";
export type YearSource = "letter_date" | "created_at";

export type InitializeSetupInput = {
  organization_name: string;
  superuser: {
    username: string;
    full_name: string;
    password: string;
  };
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

export type InitializeSetupResponse = {
  initialized: boolean;
  user_id: string;
  username: string;
};