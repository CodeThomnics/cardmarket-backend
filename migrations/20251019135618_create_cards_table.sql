-- +goose Up
CREATE TABLE "cards"(
    "card_id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "image_url" VARCHAR(500),
    "description" TEXT,
    "set_name" VARCHAR(100),
    "card_number" VARCHAR(50),
    "rarity" VARCHAR(50),
    "tcg_game_id" INTEGER NOT NULL REFERENCES "tcg_games"("tcg_game_id"),
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "idx_cards_tcg_game" ON "cards"("tcg_game_id");
CREATE INDEX "idx_cards_name" ON "cards"("name");

-- +goose Down
DROP INDEX "idx_cards_tcg_game";
DROP INDEX "idx_cards_name";
DROP TABLE "cards";
