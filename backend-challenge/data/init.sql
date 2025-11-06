-- Database initialization SQL
-- Run: sqlite3 data/store.db < data/init.sql

DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS valid_coupons;

CREATE TABLE products (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    price REAL NOT NULL,
    image_thumbnail TEXT,
    image_mobile TEXT,
    image_tablet TEXT,
    image_desktop TEXT
);

CREATE TABLE valid_coupons (
    code TEXT PRIMARY KEY
);

-- Insert valid coupons (from coupon processing)
INSERT INTO valid_coupons (code) VALUES ('BIRTHDAY');
INSERT INTO valid_coupons (code) VALUES ('BUYGETON');
INSERT INTO valid_coupons (code) VALUES ('FIFTYOFF');
INSERT INTO valid_coupons (code) VALUES ('FREEZAAA');
INSERT INTO valid_coupons (code) VALUES ('GNULINUX');
INSERT INTO valid_coupons (code) VALUES ('HAPPYHRS');
INSERT INTO valid_coupons (code) VALUES ('OVER9000');
INSERT INTO valid_coupons (code) VALUES ('SIXTYOFF');

-- Insert products
INSERT INTO products (id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop) VALUES ("1", "Waffle with Berries", "Waffle", 6.5, "https://orderfoodonline.deno.dev/public/images/image-waffle-thumbnail.jpg", "https://orderfoodonline.deno.dev/public/images/image-waffle-mobile.jpg", "https://orderfoodonline.deno.dev/public/images/image-waffle-tablet.jpg", "https://orderfoodonline.deno.dev/public/images/image-waffle-desktop.jpg");
INSERT INTO products (id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop) VALUES ("2", "Vanilla Bean Crème Brûlée", "Crème Brûlée", 7, "https://orderfoodonline.deno.dev/public/images/image-creme-brulee-thumbnail.jpg", "https://orderfoodonline.deno.dev/public/images/image-creme-brulee-mobile.jpg", "https://orderfoodonline.deno.dev/public/images/image-creme-brulee-tablet.jpg", "https://orderfoodonline.deno.dev/public/images/image-creme-brulee-desktop.jpg");
INSERT INTO products (id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop) VALUES ("3", "Macaron Mix of Five", "Macaron", 8, "https://orderfoodonline.deno.dev/public/images/image-macaron-thumbnail.jpg", "https://orderfoodonline.deno.dev/public/images/image-macaron-mobile.jpg", "https://orderfoodonline.deno.dev/public/images/image-macaron-tablet.jpg", "https://orderfoodonline.deno.dev/public/images/image-macaron-desktop.jpg");
INSERT INTO products (id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop) VALUES ("4", "Classic Tiramisu", "Tiramisu", 5.5, "https://orderfoodonline.deno.dev/public/images/image-tiramisu-thumbnail.jpg", "https://orderfoodonline.deno.dev/public/images/image-tiramisu-mobile.jpg", "https://orderfoodonline.deno.dev/public/images/image-tiramisu-tablet.jpg", "https://orderfoodonline.deno.dev/public/images/image-tiramisu-desktop.jpg");
INSERT INTO products (id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop) VALUES ("5", "Pistachio Baklava", "Baklava", 4, "https://orderfoodonline.deno.dev/public/images/image-baklava-thumbnail.jpg", "https://orderfoodonline.deno.dev/public/images/image-baklava-mobile.jpg", "https://orderfoodonline.deno.dev/public/images/image-baklava-tablet.jpg", "https://orderfoodonline.deno.dev/public/images/image-baklava-desktop.jpg");
INSERT INTO products (id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop) VALUES ("6", "Lemon Meringue Pie", "Pie", 5, "https://orderfoodonline.deno.dev/public/images/image-meringue-thumbnail.jpg", "https://orderfoodonline.deno.dev/public/images/image-meringue-mobile.jpg", "https://orderfoodonline.deno.dev/public/images/image-meringue-tablet.jpg", "https://orderfoodonline.deno.dev/public/images/image-meringue-desktop.jpg");
INSERT INTO products (id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop) VALUES ("7", "Red Velvet Cake", "Cake", 4.5, "https://orderfoodonline.deno.dev/public/images/image-cake-thumbnail.jpg", "https://orderfoodonline.deno.dev/public/images/image-cake-mobile.jpg", "https://orderfoodonline.deno.dev/public/images/image-cake-tablet.jpg", "https://orderfoodonline.deno.dev/public/images/image-cake-desktop.jpg");
INSERT INTO products (id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop) VALUES ("8", "Salted Caramel Brownie", "Brownie", 4.5, "https://orderfoodonline.deno.dev/public/images/image-brownie-thumbnail.jpg", "https://orderfoodonline.deno.dev/public/images/image-brownie-mobile.jpg", "https://orderfoodonline.deno.dev/public/images/image-brownie-tablet.jpg", "https://orderfoodonline.deno.dev/public/images/image-brownie-desktop.jpg");
INSERT INTO products (id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop) VALUES ("9", "Vanilla Panna Cotta", "Panna Cotta", 6.5, "https://orderfoodonline.deno.dev/public/images/image-panna-cotta-thumbnail.jpg", "https://orderfoodonline.deno.dev/public/images/image-panna-cotta-mobile.jpg", "https://orderfoodonline.deno.dev/public/images/image-panna-cotta-tablet.jpg", "https://orderfoodonline.deno.dev/public/images/image-panna-cotta-desktop.jpg");
