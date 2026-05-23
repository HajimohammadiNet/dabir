CREATE TABLE IF NOT EXISTS letter_number_counters (
    jalali_year INT PRIMARY KEY,
    last_number BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE letters
ADD COLUMN IF NOT EXISTS letter_year INT,
ADD COLUMN IF NOT EXISTS letter_year_suffix VARCHAR(8),
ADD COLUMN IF NOT EXISTS letter_serial BIGINT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_letters_year_serial_unique
ON letters (letter_year, letter_serial)
WHERE letter_year IS NOT NULL
  AND letter_serial IS NOT NULL;