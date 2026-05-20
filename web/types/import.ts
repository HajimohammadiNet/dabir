export type ImportError = {
  row: number;
  field: string;
  message: string;
};

export type ImportedLetterRow = {
  row_number: number;
  letter_number: number;
  title: string;
  letter_date: string;
  sender: string;
  receiver: string;
};

export type ImportJob = {
  id: string;
  type: "letters";
  status: "previewed" | "committed" | "failed";
  file_name: string;

  total_rows: number;
  valid_rows: number;
  invalid_rows: number;

  max_letter_number?: number | null;

  detected_columns?: Record<string, string>;
  preview_data?: ImportedLetterRow[];
  errors?: ImportError[];

  created_by: string;
  committed_by?: string | null;

  created_at: string;
  committed_at?: string | null;
};

export type CommitImportResponse = {
  import_id: string;
  imported_rows: number;
  skipped_rows: number;
  next_letter_number: number;
};