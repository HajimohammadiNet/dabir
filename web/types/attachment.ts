export type LetterAttachment = {
  id: string;
  letter_id: string;
  file_name: string;
  content_type: string;
  size_bytes: number;
  is_deleted: boolean;
  created_at: string;
  updated_at: string;
};

export type AttachmentDownloadURLResponse = {
  url: string;
};