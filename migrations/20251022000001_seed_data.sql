-- +goose Up
-- Reference/Lookup Tables
INSERT INTO languages (language_code, language_name) VALUES 
('EN', 'English'),
('DE', 'German'),
('FR', 'French'),
('ES', 'Spanish'),
('IT', 'Italian'),
('JP', 'Japanese');

INSERT INTO countries (country_name, country_code) VALUES 
('United States', 'US'),
('Germany', 'DE'),
('France', 'FR'),
('United Kingdom', 'GB'),
('Italy', 'IT'),
('Spain', 'ES'),
('Japan', 'JP'),
('Canada', 'CA'),
('Netherlands', 'NL');

INSERT INTO tcg_games (name) VALUES 
('Magic: The Gathering'),
('Pokemon'),
('Yu-Gi-Oh!'),
('Flesh and Blood'),
('Digimon'),
('One Piece Card Game');

-- Users (mix of buyers and sellers)
INSERT INTO users (username, email, password_hash, first_name, last_name, street_name, street_number, city, state, zip_code, seller_type, country_id, language_id) VALUES
('magicdealer', 'dealer@example.com', '$2a$12$tPpK2HKTiNJAhWPX9zK9NOD0txn0OZKGb1kTF8TCr3MmTNP4gHzcu', 'Thomas', 'Wilson', 'Oak Street', '42', 'Berlin', 'Berlin', '10115', 'professional',
	(SELECT country_id FROM countries WHERE country_code = 'DE'),
	(SELECT language_id FROM languages WHERE language_code = 'DE')
),
('cardcollector', 'collector@example.com', '$2a$12$CLkNQFCY9fPEGBN1FgZ6ZugGmCzZxNBsRUmYOiVDa7jBvX/geB.zS', 'Marie', 'Dubois', 'Rue de Seine', '15', 'Paris', 'Île-de-France', '75006', 'private',
	(SELECT country_id FROM countries WHERE country_code = 'FR'),
	(SELECT language_id FROM languages WHERE language_code = 'FR')
),
('powertcg', 'power@example.com', '$2a$12$xKlPVnpyN8/rS6VdZQu9n.36dNvD0d1XRdEhqXZF4KC0H.CKWleca', 'Alex', 'Johnson', 'High Street', '27', 'London', 'England', 'EC1V 7JN', 'powerseller',
	(SELECT country_id FROM countries WHERE country_code = 'GB'),
	(SELECT language_id FROM languages WHERE language_code = 'EN')
),
('casualplayer', 'casual@example.com', '$2a$12$hW6mK9Oq4RzOXF5ciV2Rue9OBG1YB4kA8nMhnA.jQ99S4WKrTmzne', 'Carlos', 'Gonzalez', 'Calle Mayor', '8', 'Madrid', 'Madrid', '28013', 'private',
	(SELECT country_id FROM countries WHERE country_code = 'ES'),
	(SELECT language_id FROM languages WHERE language_code = 'ES')
),
('mtgbuyer', 'buyer@example.com', '$2a$12$T.b4WAjmc6c/1xKgz9MsSOcbkZgRJHJ7UzgKjKJQbPA5S8D6/slKi', 'Emma', 'Brown', 'Maple Avenue', '103', 'Toronto', 'Ontario', 'M5V 2A4', 'private',
	(SELECT country_id FROM countries WHERE country_code = 'CA'),
	(SELECT language_id FROM languages WHERE language_code = 'EN')
),
('rarefinds', 'rare@example.com', '$2a$12$5YhKx1ksZat7KklT.LMWweo8A3vBSY96LP0B1NeRNEMxiBPZZ1gsS', 'Takashi', 'Yamamoto', 'Sakura Street', '5-2', 'Tokyo', 'Tokyo', '150-0043', 'professional',
	(SELECT country_id FROM countries WHERE country_code = 'JP'),
	(SELECT language_id FROM languages WHERE language_code = 'JP')
);

-- Cards
INSERT INTO cards (name, image_url, description, set_name, card_number, rarity, tcg_game_id) VALUES
('Black Lotus', 'https://example.com/black_lotus.jpg', 'Adds 3 mana of any single color to your mana pool, then is discarded.', 'Alpha', '232', 'Mythic Rare',
	(SELECT tcg_game_id FROM tcg_games WHERE name = 'Magic: The Gathering')
),
('Charizard', 'https://example.com/charizard.jpg', 'Stage 2 Pokémon. Evolved from Charmeleon.', 'Base Set', '4/102', 'Holo Rare',
	(SELECT tcg_game_id FROM tcg_games WHERE name = 'Pokemon')
),
('Blue-Eyes White Dragon', 'https://example.com/blue_eyes.jpg', 'This legendary dragon is a powerful engine of destruction.', 'Legend of Blue Eyes', 'LOB-001', 'Ultra Rare',
	(SELECT tcg_game_id FROM tcg_games WHERE name = 'Yu-Gi-Oh!')
),
('Mox Ruby', 'https://example.com/mox_ruby.jpg', 'Adds one red mana to your mana pool.', 'Beta', '265', 'Rare',
	(SELECT tcg_game_id FROM tcg_games WHERE name = 'Magic: The Gathering')
),
('Pikachu', 'https://example.com/pikachu.jpg', 'Mouse Pokémon., Weight: 13 lbs.', 'Base Set', '58/102', 'Common',
	(SELECT tcg_game_id FROM tcg_games WHERE name = 'Pokemon')
),
('Dark Magician', 'https://example.com/dark_magician.jpg', 'The ultimate wizard in terms of attack and defense.', 'Legend of Blue Eyes', 'LOB-005', 'Ultra Rare',
	(SELECT tcg_game_id FROM tcg_games WHERE name = 'Yu-Gi-Oh!')
),
('Briar', 'https://example.com/briar.jpg', 'Embodiment of earth and lightning', 'Tales of Aria', '1', 'Legendary',
	(SELECT tcg_game_id FROM tcg_games WHERE name = 'Flesh and Blood')
),
('Omnimon', 'https://example.com/omnimon.jpg', 'DNA digivolved from WarGreymon and MetalGarurumon', 'Release Special', 'BT1-084', 'Secret Rare',
	(SELECT tcg_game_id FROM tcg_games WHERE name = 'Digimon')
),
('Monkey D. Luffy', 'https://example.com/luffy.jpg', 'Captain of the Straw Hat Pirates', 'Romance Dawn', 'OP01-001', 'Leader Rare',
	(SELECT tcg_game_id FROM tcg_games WHERE name = 'One Piece Card Game')
),
('Jace, the Mind Sculptor', 'https://example.com/jace.jpg', 'Powerful planeswalker with multiple abilities', 'Worldwake', '31', 'Mythic Rare',
	(SELECT tcg_game_id FROM tcg_games WHERE name = 'Magic: The Gathering')
);

-- Products
INSERT INTO products (seller_id, card_id, price, condition, quantity, is_available, language_id) VALUES
((SELECT user_id FROM users WHERE username = 'magicdealer'),
 (SELECT card_id FROM cards WHERE name = 'Black Lotus' AND set_name = 'Alpha' LIMIT 1),
 25000.00, 'good', 1, true, (SELECT language_id FROM languages WHERE language_code = 'EN')),
((SELECT user_id FROM users WHERE username = 'powertcg'),
 (SELECT card_id FROM cards WHERE name = 'Charizard' AND set_name = 'Base Set' LIMIT 1),
 5000.00, 'near mint', 1, true, (SELECT language_id FROM languages WHERE language_code = 'EN')),
((SELECT user_id FROM users WHERE username = 'magicdealer'),
 (SELECT card_id FROM cards WHERE name = 'Blue-Eyes White Dragon' AND set_name = 'Legend of Blue Eyes' LIMIT 1),
 120.00, 'excellent', 2, true, (SELECT language_id FROM languages WHERE language_code = 'EN')),
((SELECT user_id FROM users WHERE username = 'powertcg'),
 (SELECT card_id FROM cards WHERE name = 'Mox Ruby' AND set_name = 'Beta' LIMIT 1),
 4500.00, 'light_played', 1, true, (SELECT language_id FROM languages WHERE language_code = 'EN')),
((SELECT user_id FROM users WHERE username = 'rarefinds'),
 (SELECT card_id FROM cards WHERE name = 'Pikachu' AND set_name = 'Base Set' LIMIT 1),
 25.00, 'mint', 4, true, (SELECT language_id FROM languages WHERE language_code = 'EN')),
((SELECT user_id FROM users WHERE username = 'magicdealer'),
 (SELECT card_id FROM cards WHERE name = 'Dark Magician' AND set_name = 'Legend of Blue Eyes' LIMIT 1),
 85.00, 'excellent', 3, true, (SELECT language_id FROM languages WHERE language_code = 'EN')),
((SELECT user_id FROM users WHERE username = 'rarefinds'),
 (SELECT card_id FROM cards WHERE name = 'Briar' AND set_name = 'Tales of Aria' LIMIT 1),
 150.00, 'near mint', 2, true, (SELECT language_id FROM languages WHERE language_code = 'EN')),
((SELECT user_id FROM users WHERE username = 'powertcg'),
 (SELECT card_id FROM cards WHERE name = 'Omnimon' AND set_name = 'Release Special' LIMIT 1),
 200.00, 'mint', 1, true, (SELECT language_id FROM languages WHERE language_code = 'EN')),
((SELECT user_id FROM users WHERE username = 'rarefinds'),
 (SELECT card_id FROM cards WHERE name = 'Monkey D. Luffy' AND set_name = 'Romance Dawn' LIMIT 1),
 45.00, 'excellent', 5, true, (SELECT language_id FROM languages WHERE language_code = 'EN')),
((SELECT user_id FROM users WHERE username = 'magicdealer'),
 (SELECT card_id FROM cards WHERE name = 'Jace, the Mind Sculptor' AND set_name = 'Worldwake' LIMIT 1),
 250.00, 'near mint', 2, true, (SELECT language_id FROM languages WHERE language_code = 'DE')),
((SELECT user_id FROM users WHERE username = 'powertcg'),
 (SELECT card_id FROM cards WHERE name = 'Black Lotus' AND set_name = 'Alpha' LIMIT 1),
 30000.00, 'near mint', 1, true, (SELECT language_id FROM languages WHERE language_code = 'EN')),
((SELECT user_id FROM users WHERE username = 'rarefinds'),
 (SELECT card_id FROM cards WHERE name = 'Charizard' AND set_name = 'Base Set' LIMIT 1),
 6500.00, 'mint', 1, true, (SELECT language_id FROM languages WHERE language_code = 'EN'));

-- Orders
INSERT INTO orders (buyer_id, seller_id, quantity, product_id, order_date, shipping_address, shipping_cost, total_amount, tracking_number, shipped_at, delivered_at, status) VALUES
-- Order 1: cardcollector buys Black Lotus from magicdealer
((SELECT user_id FROM users WHERE username = 'cardcollector'),
 (SELECT user_id FROM users WHERE username = 'magicdealer'),
 1,
 (SELECT p.product_id FROM products p 
  JOIN users us ON us.user_id = p.seller_id 
  JOIN cards c ON c.card_id = p.card_id 
  WHERE us.username = 'magicdealer' AND c.name = 'Black Lotus' AND c.set_name = 'Alpha' LIMIT 1),
 '2025-09-15 10:30:00+00',
 '15 Rue de Seine, 75006 Paris, France',
 85.00,
 25085.00,
 'TRK7829456321',
 '2025-09-16 11:45:00+00',
 '2025-09-19 14:30:00+00',
 'completed'),
-- Order 2: casualplayer buys Blue-Eyes White Dragon from magicdealer
((SELECT user_id FROM users WHERE username = 'casualplayer'),
 (SELECT user_id FROM users WHERE username = 'magicdealer'),
 1,
 (SELECT p.product_id FROM products p 
  JOIN users us ON us.user_id = p.seller_id 
  JOIN cards c ON c.card_id = p.card_id 
  WHERE us.username = 'magicdealer' AND c.name = 'Blue-Eyes White Dragon' AND c.set_name = 'Legend of Blue Eyes' LIMIT 1),
 '2025-09-28 14:15:00+00',
 '8 Calle Mayor, 28013 Madrid, Spain',
 5.00,
 125.00,
 'TRK4832957612',
 '2025-09-29 10:20:00+00',
 '2025-10-02 15:45:00+00',
 'completed'),
-- Order 3: casualplayer buys Charizard from powertcg
((SELECT user_id FROM users WHERE username = 'casualplayer'),
 (SELECT user_id FROM users WHERE username = 'powertcg'),
 1,
 (SELECT p.product_id FROM products p 
  JOIN users us ON us.user_id = p.seller_id 
  JOIN cards c ON c.card_id = p.card_id 
  WHERE us.username = 'powertcg' AND c.name = 'Charizard' AND c.set_name = 'Base Set' LIMIT 1),
 '2025-09-28 14:15:00+00',
 '8 Calle Mayor, 28013 Madrid, Spain',
 15.00,
 5015.00,
 'TRK3918274650',
 '2025-09-29 09:15:00+00',
 '2025-10-03 12:30:00+00',
 'completed'),
-- Order 4: mtgbuyer buys Mox Ruby from powertcg
((SELECT user_id FROM users WHERE username = 'mtgbuyer'),
 (SELECT user_id FROM users WHERE username = 'powertcg'),
 1,
 (SELECT p.product_id FROM products p 
  JOIN users us ON us.user_id = p.seller_id 
  JOIN cards c ON c.card_id = p.card_id 
  WHERE us.username = 'powertcg' AND c.name = 'Mox Ruby' AND c.set_name = 'Beta' LIMIT 1),
 '2025-10-10 09:45:00+00',
 '103 Maple Avenue, M5V 2A4 Toronto, Canada',
 25.00,
 4525.00,
 'TRK9274619385',
 '2025-10-12 13:40:00+00',
 NULL,
 'completed'),
-- Order 5: mtgbuyer buys Briar from rarefinds
((SELECT user_id FROM users WHERE username = 'mtgbuyer'),
 (SELECT user_id FROM users WHERE username = 'rarefinds'),
 1,
 (SELECT p.product_id FROM products p 
  JOIN users us ON us.user_id = p.seller_id 
  JOIN cards c ON c.card_id = p.card_id 
  WHERE us.username = 'rarefinds' AND c.name = 'Briar' AND c.set_name = 'Tales of Aria' LIMIT 1),
 '2025-10-10 09:45:00+00',
 '103 Maple Avenue, M5V 2A4 Toronto, Canada',
 10.00,
 160.00,
 NULL,
 NULL,
 NULL,
 'processing'),
-- Order 6: cardcollector buys Briar from rarefinds
((SELECT user_id FROM users WHERE username = 'cardcollector'),
 (SELECT user_id FROM users WHERE username = 'rarefinds'),
 1,
 (SELECT p.product_id FROM products p 
  JOIN users us ON us.user_id = p.seller_id 
  JOIN cards c ON c.card_id = p.card_id 
  WHERE us.username = 'rarefinds' AND c.name = 'Briar' AND c.set_name = 'Tales of Aria' LIMIT 1),
 '2025-10-18 16:20:00+00',
 '15 Rue de Seine, 75006 Paris, France',
 8.00,
 158.00,
 NULL,
 NULL,
 NULL,
 'pending'),
-- Order 7: cardcollector buys Monkey D. Luffy from rarefinds
((SELECT user_id FROM users WHERE username = 'cardcollector'),
 (SELECT user_id FROM users WHERE username = 'rarefinds'),
 1,
 (SELECT p.product_id FROM products p 
  JOIN users us ON us.user_id = p.seller_id 
  JOIN cards c ON c.card_id = p.card_id 
  WHERE us.username = 'rarefinds' AND c.name = 'Monkey D. Luffy' AND c.set_name = 'Romance Dawn' LIMIT 1),
 '2025-10-18 16:20:00+00',
 '15 Rue de Seine, 75006 Paris, France',
 8.00,
 53.00,
 NULL,
 NULL,
 NULL,
 'pending');

-- +goose Down
-- Clean up all data in reverse order of dependencies
DELETE FROM orders;
DELETE FROM products;
DELETE FROM cards;
DELETE FROM users;
DELETE FROM tcg_games;
DELETE FROM countries;
DELETE FROM languages;
