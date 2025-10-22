-- +goose Up
CREATE TABLE "order_items"(
    "order_item_id" SERIAL PRIMARY KEY,
    "sub_order_id" INTEGER NOT NULL REFERENCES "sub_orders"("sub_order_id"),
    "product_id" INTEGER NOT NULL REFERENCES "products"("product_id"),
    "quantity" INTEGER NOT NULL CHECK ("quantity" > 0),
    "unit_price" DECIMAL(10, 2) NOT NULL CHECK("unit_price" >= 0),
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "idx_order_items_sub_order" ON "order_items"("sub_order_id");
CREATE INDEX "idx_order_items_product" ON "order_items"("product_id");

-- +goose Down
DROP INDEX "idx_order_items_sub_order";
DROP INDEX "idx_order_items_product";
DROP TABLE "order_items";
