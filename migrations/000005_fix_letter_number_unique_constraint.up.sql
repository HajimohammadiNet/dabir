ALTER TABLE letters
DROP CONSTRAINT IF EXISTS letters_letter_number_key;

CREATE UNIQUE INDEX IF NOT EXISTS idx_letters_fixed_number_unique
ON letters (letter_number)
WHERE letter_year IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_letters_year_serial_unique
ON letters (letter_year, letter_serial)
WHERE letter_year IS NOT NULL
  AND letter_serial IS NOT NULL;