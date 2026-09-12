CREATE SEQUENCE IF NOT EXISTS outgoing_letter_number_seq
    START WITH 1
    INCREMENT BY 1
    MINVALUE 1
    NO MAXVALUE
    CACHE 1;

ALTER TABLE letters
ADD COLUMN direction VARCHAR(20) NOT NULL DEFAULT 'incoming';

ALTER TABLE letters
ADD CONSTRAINT letters_direction_check
CHECK (direction IN ('incoming', 'outgoing'));

ALTER TABLE letter_number_counters
DROP CONSTRAINT IF EXISTS letter_number_counters_pkey;

ALTER TABLE letter_number_counters
ADD COLUMN direction VARCHAR(20) NOT NULL DEFAULT 'incoming';

ALTER TABLE letter_number_counters
ADD CONSTRAINT letter_number_counters_direction_check
CHECK (direction IN ('incoming', 'outgoing'));

ALTER TABLE letter_number_counters
ADD PRIMARY KEY (direction, jalali_year);

DROP INDEX IF EXISTS idx_letters_fixed_number_unique;
DROP INDEX IF EXISTS idx_letters_year_serial_unique;
DROP INDEX IF EXISTS idx_letters_display_letter_number_unique;

CREATE UNIQUE INDEX idx_letters_fixed_number_unique
ON letters (direction, letter_number)
WHERE letter_year IS NULL;

CREATE UNIQUE INDEX idx_letters_year_serial_unique
ON letters (direction, letter_year, letter_serial)
WHERE letter_year IS NOT NULL
  AND letter_serial IS NOT NULL;

CREATE UNIQUE INDEX idx_letters_display_letter_number_unique
ON letters (direction, display_letter_number)
WHERE display_letter_number IS NOT NULL
  AND display_letter_number <> ''
  AND is_deleted = false;

CREATE INDEX idx_letters_direction_active_date
ON letters (direction, is_deleted, letter_date DESC, created_at DESC);

CREATE INDEX idx_letters_direction_created_at
ON letters (direction, created_at DESC);
