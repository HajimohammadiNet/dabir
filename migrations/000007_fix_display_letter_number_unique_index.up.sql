DROP INDEX IF EXISTS idx_letters_display_letter_number_unique;

CREATE UNIQUE INDEX idx_letters_display_letter_number_unique
ON letters (display_letter_number)
WHERE display_letter_number IS NOT NULL
  AND display_letter_number <> ''
  AND is_deleted = false;