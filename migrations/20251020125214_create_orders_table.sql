-- +goose Up
CREATE TABLE "orders"(
    "order_id" SERIAL PRIMARY KEY,
    "buyer_id" INTEGER NOT NULL REFERENCES "users"("user_id"),
    "seller_id" INTEGER NOT NULL REFERENCES "users"("user_id"),
    "quantity" INTEGER NOT NULL CHECK ("quantity" > 0),
    "product_id" INTEGER NOT NULL REFERENCES "products"("product_id"),
    "order_date" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "shipping_address" TEXT NOT NULL,
    "shipping_cost" DECIMAL(10, 2) NOT NULL CHECK ("shipping_cost" >= 0),
    "total_amount" DECIMAL(10, 2) NOT NULL CHECK ("total_amount" >= 0),
    "tracking_number" VARCHAR(50),
    "shipped_at" TIMESTAMP WITH TIME ZONE,
    "delivered_at" TIMESTAMP WITH TIME ZONE,
    "status" VARCHAR(20) CHECK ("status" IN('pending', 'processing', 'completed', 'cancelled')) NOT NULL DEFAULT 'pending',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "idx_orders_buyer" ON "orders"("buyer_id");
CREATE INDEX "idx_orders_seller" ON "orders"("seller_id");
CREATE INDEX "idx_orders_status" ON "orders"("status");
CREATE INDEX "idx_orders_date" ON "orders"("order_date");

-- +goose Down
DROP INDEX "idx_orders_date";
DROP INDEX "idx_orders_status";
DROP INDEX "idx_orders_buyer";
DROP INDEX "idx_orders_seller";
DROP TABLE "orders";
