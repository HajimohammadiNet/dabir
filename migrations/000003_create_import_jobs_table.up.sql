CREATE TABLE import_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,

    file_name TEXT NOT NULL,

    total_rows INT NOT NULL DEFAULT 0,
    valid_rows INT NOT NULL DEFAULT 0,
    invalid_rows INT NOT NULL DEFAULT 0,

    max_letter_number BIGINT,

    detected_columns JSONB,
    preview_data JSONB,
    errors JSONB,

    created_by UUID NOT NULL REFERENCES users(id),
    committed_by UUID REFERENCES users(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    committed_at TIMESTAMPTZ
);

CREATE INDEX idx_import_jobs_type ON import_jobs(type);
CREATE INDEX idx_import_jobs_status ON import_jobs(status);
CREATE INDEX idx_import_jobs_created_by ON import_jobs(created_by);
CREATE INDEX idx_import_jobs_created_at ON import_jobs(created_at);