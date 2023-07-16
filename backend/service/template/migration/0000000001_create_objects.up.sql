-- template
CREATE TABLE template
(
    id           UUID
        CONSTRAINT template_pk PRIMARY KEY DEFAULT gen_random_uuid(),
    extra_column TEXT NULL,
    created_at   timestamp with time zone  DEFAULT current_timestamp,
    updated_at   timestamp with time zone  DEFAULT NULL,
    deleted_at   timestamp with time zone  DEFAULT NULL
);