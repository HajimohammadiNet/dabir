export type Letter = {
  id: string;
  letter_number: number;
  formatted_letter_number: string;

  title: string;
  letter_date: string;
  letter_date_jalali: string;

  registrar_name: string;
  sender: string;
  receiver: string;

  description?: string | null;

  created_by: string;
  updated_by?: string | null;
  deleted_by?: string | null;

  is_deleted: boolean;

  created_at: string;
  updated_at: string;
  deleted_at?: string | null;
};

export type ListLettersResponse = {
  items: Letter[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
};