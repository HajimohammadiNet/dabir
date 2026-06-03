CREATE TABLE IF NOT EXISTS letter_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    letter_id UUID NOT NULL REFERENCES letters(id) ON DELETE CASCADE,

    file_name TEXT NOT NULL,
    object_key TEXT NOT NULL,
    content_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,

    uploaded_by UUID NOT NULL REFERENCES users(id),

    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_by UUID NULL REFERENCES users(id),
    deleted_at TIMESTAMPTZ NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_letter_attachments_object_key_unique
ON letter_attachments (object_key)
WHERE object_key IS NOT NULL
  AND object_key <> '';

CREATE INDEX IF NOT EXISTS idx_letter_attachments_letter_id
ON letter_attachments (letter_id)
WHERE is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_letter_attachments_created_at
ON letter_attachments (created_at DESC);