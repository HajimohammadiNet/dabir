DROP INDEX IF EXISTS idx_letters_year_serial_unique;

ALTER TABLE letters
DROP COLUMN IF EXISTS letter_serial,
DROP COLUMN IF EXISTS letter_year_suffix,
DROP COLUMN IF EXISTS letter_year;

DROP TABLE IF EXISTS letter_number_counters;