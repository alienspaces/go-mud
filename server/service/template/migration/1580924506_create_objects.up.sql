-- template
CREATE TABLE "template" (
    "id"         UUID CONSTRAINT template_pk PRIMARY KEY DEFAULT gen_random_uuid(),
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "deleted_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);