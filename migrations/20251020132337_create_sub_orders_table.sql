-- +goose Up
CREATE TABLE "sub_orders"(
    "sub_order_id" SERIAL PRIMARY KEY,
    "order_id" INTEGER NOT NULL REFERENCES "orders"("order_id") ON DELETE CASCADE,
    "seller_id" INTEGER NOT NULL REFERENCES "users"("user_id"),
    "subtotal" DECIMAL(10, 2) NOT NULL CHECK ("subtotal" >= 0),
    "shipping_cost" DECIMAL(10, 2) NOT NULL DEFAULT 0 CHECK ("shipping_cost" >= 0),
    "status" VARCHAR(20) CHECK ("status" IN('pending', 'processing', 'shipped', 'delivered', 'cancelled')) NOT NULL DEFAULT 'pending',
    "tracking_number" VARCHAR(100),
    "shipped_at" TIMESTAMP WITH TIME ZONE,
    "delivered_at" TIMESTAMP WITH TIME ZONE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "idx_sub_orders_order" ON "sub_orders"("order_id");
CREATE INDEX "idx_sub_orders_seller" ON "sub_orders"("seller_id");
CREATE INDEX "idx_sub_orders_status" ON "sub_orders"("status");

-- +goose Down
DROP INDEX "idx_sub_orders_order";
DROP INDEX "idx_sub_orders_seller";
DROP INDEX "idx_sub_orders_status";
DROP TABLE "sub_orders";
