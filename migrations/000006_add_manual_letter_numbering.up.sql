ALTER TABLE letters
ADD COLUMN IF NOT EXISTS display_letter_number TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_letters_display_letter_number_unique
ON letters (display_letter_number)
WHERE display_letter_number IS NOT NULL
  AND display_letter_number <> '';