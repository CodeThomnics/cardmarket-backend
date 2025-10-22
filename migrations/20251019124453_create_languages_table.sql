-- +goose Up
CREATE TABLE "languages"(
    "language_id" SERIAL PRIMARY KEY,
    "language_code" CHAR(2) UNIQUE NOT NULL,
    "language_name" VARCHAR(100) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE "languages";
