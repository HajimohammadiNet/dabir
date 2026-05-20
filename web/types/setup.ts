export type SetupStatus = {
  initialized: boolean;
  setup_needed: boolean;
};

export type InitializeSetupInput = {
  organization_name: string;
  superuser: {
    username: string;
    full_name: string;
    password: string;
  };
  letter_config: {
    number_prefix: string;
    number_padding: number;
  };
};

export type InitializeSetupResponse = {
  initialized: boolean;
  user_id: string;
  username: string;
};