-- +goose Up
CREATE TABLE "products"(
    "product_id" SERIAL PRIMARY KEY,
    "seller_id" INTEGER NOT NULL REFERENCES "users"("user_id"),
    "card_id" INTEGER NOT NULL REFERENCES "cards"("card_id"),
    "price" DECIMAL(10, 2) NOT NULL CHECK ("price" >= 0),
    "condition" VARCHAR(20) CHECK ("condition" IN ('mint', 'near mint', 'excellent', 'good', 'light_played', 'played', 'poor')) NOT NULL,
    "quantity" INTEGER NOT NULL DEFAULT 1 CHECK ("quantity" >= 0),
    "is_available" BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "idx_products_seller" ON "products"("seller_id");
CREATE INDEX "idx_products_card" ON "products"("card_id");
CREATE INDEX "idx_products_available" ON "products"("is_available") WHERE "is_available" = TRUE;

-- +goose Down
DROP INDEX "idx_products_seller";
DROP INDEX "idx_products_card";
DROP INDEX "idx_products_available";
DROP TABLE "products";