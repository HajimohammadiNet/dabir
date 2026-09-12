DROP INDEX IF EXISTS idx_letters_direction_created_at;
DROP INDEX IF EXISTS idx_letters_direction_active_date;
DROP INDEX IF EXISTS idx_letters_display_letter_number_unique;
DROP INDEX IF EXISTS idx_letters_year_serial_unique;
DROP INDEX IF EXISTS idx_letters_fixed_number_unique;

-- A rollback keeps outgoing rows as legacy letters. Reassign their internal
-- numbers and yearly fields so the pre-direction unique indexes can be restored
-- without deleting user data.
WITH numbered AS (
    SELECT
        id,
        COALESCE((SELECT MAX(letter_number) FROM letters WHERE direction = 'incoming'), 0)
            + ROW_NUMBER() OVER (ORDER BY created_at, id) AS new_number
    FROM letters
    WHERE direction = 'outgoing'
)
UPDATE letters AS l
SET
    letter_number = numbered.new_number,
    letter_year = NULL,
    letter_year_suffix = NULL,
    letter_serial = NULL
FROM numbered
WHERE l.id = numbered.id;

WITH duplicates AS (
    SELECT outgoing.id
    FROM letters AS outgoing
    WHERE outgoing.direction = 'outgoing'
      AND outgoing.is_deleted = false
      AND outgoing.display_letter_number IS NOT NULL
      AND outgoing.display_letter_number <> ''
      AND EXISTS (
          SELECT 1
          FROM letters AS incoming
          WHERE incoming.direction = 'incoming'
            AND incoming.is_deleted = false
            AND incoming.display_letter_number = outgoing.display_letter_number
      )
)
UPDATE letters AS l
SET display_letter_number = 'OUTGOING-' || l.id::TEXT || '-' || l.display_letter_number
FROM duplicates
WHERE l.id = duplicates.id;

ALTER TABLE letter_number_counters
DROP CONSTRAINT IF EXISTS letter_number_counters_pkey;

UPDATE letter_number_counters AS incoming
SET last_number = merged.last_number,
    updated_at = NOW()
FROM (
    SELECT jalali_year, MAX(last_number) AS last_number
    FROM letter_number_counters
    GROUP BY jalali_year
) AS merged
WHERE incoming.direction = 'incoming'
  AND incoming.jalali_year = merged.jalali_year;

INSERT INTO letter_number_counters (direction, jalali_year, last_number, updated_at)
SELECT 'incoming', jalali_year, MAX(last_number), NOW()
FROM letter_number_counters
GROUP BY jalali_year
HAVING NOT EXISTS (
    SELECT 1
    FROM letter_number_counters AS incoming
    WHERE incoming.direction = 'incoming'
      AND incoming.jalali_year = letter_number_counters.jalali_year
);

DELETE FROM letter_number_counters
WHERE direction = 'outgoing';

ALTER TABLE letter_number_counters
DROP CONSTRAINT IF EXISTS letter_number_counters_direction_check;

ALTER TABLE letter_number_counters
DROP COLUMN direction;

ALTER TABLE letter_number_counters
ADD PRIMARY KEY (jalali_year);

ALTER TABLE letters
DROP CONSTRAINT IF EXISTS letters_direction_check;

ALTER TABLE letters
DROP COLUMN direction;

CREATE UNIQUE INDEX idx_letters_fixed_number_unique
ON letters (letter_number)
WHERE letter_year IS NULL;

CREATE UNIQUE INDEX idx_letters_year_serial_unique
ON letters (letter_year, letter_serial)
WHERE letter_year IS NOT NULL
  AND letter_serial IS NOT NULL;

CREATE UNIQUE INDEX idx_letters_display_letter_number_unique
ON letters (display_letter_number)
WHERE display_letter_number IS NOT NULL
  AND display_letter_number <> ''
  AND is_deleted = false;

SELECT setval(
    'letter_number_seq',
    GREATEST(COALESCE((SELECT MAX(letter_number) FROM letters), 1), 1),
    EXISTS (SELECT 1 FROM letters)
);

DROP SEQUENCE IF EXISTS outgoing_letter_number_seq;
