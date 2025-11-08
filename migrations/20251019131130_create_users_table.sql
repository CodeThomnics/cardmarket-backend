-- +goose Up
CREATE TABLE "users"(
    "user_id" SERIAL PRIMARY KEY,
    "username" VARCHAR(50) UNIQUE NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "password_hash" TEXT NOT NULL,
    "first_name" VARCHAR(100) NOT NULL,
    "last_name" VARCHAR(100) NOT NULL,
    "street_name" VARCHAR(255) NOT NULL,
    "street_number" VARCHAR(20) NOT NULL,
    "city" VARCHAR(100) NOT NULL,
    "state" VARCHAR(100) NOT NULL,
    "zip_code" VARCHAR(20) NOT NULL,
    "seller_type" VARCHAR(20) CHECK ("seller_type" IN('private', 'professional', 'powerseller')) NOT NULL DEFAULT 'private',
    "country_id" INTEGER NOT NULL  REFERENCES "countries"("country_id"),
    "language_id" INTEGER NOT NULL REFERENCES "languages"("language_id"),
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE "users";