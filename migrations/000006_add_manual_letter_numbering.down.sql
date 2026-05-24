DROP INDEX IF EXISTS idx_letters_display_letter_number_unique;

ALTER TABLE letters
DROP COLUMN IF EXISTS display_letter_number;