ALTER TABLE letters
ADD COLUMN sender VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN receiver VARCHAR(255) NOT NULL DEFAULT '';

UPDATE letters
SET receiver = destination
WHERE receiver = '';

CREATE INDEX idx_letters_sender ON letters(sender);
CREATE INDEX idx_letters_receiver ON letters(receiver);