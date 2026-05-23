DROP INDEX IF EXISTS idx_letters_fixed_number_unique;
DROP INDEX IF EXISTS idx_letters_year_serial_unique;

ALTER TABLE letters
ADD CONSTRAINT letters_letter_number_key UNIQUE (letter_number);