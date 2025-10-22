-- +goose Up
CREATE TABLE "countries"(
    "country_id" SERIAL PRIMARY KEY,
    "country_name" VARCHAR(100) NOT NULL,
    "country_code" CHAR(2) UNIQUE NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE "countries";
