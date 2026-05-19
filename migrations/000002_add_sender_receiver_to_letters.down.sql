DROP INDEX IF EXISTS idx_letters_sender;
DROP INDEX IF EXISTS idx_letters_receiver;

ALTER TABLE letters
DROP COLUMN IF EXISTS sender,
DROP COLUMN IF EXISTS receiver;