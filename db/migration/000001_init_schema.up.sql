CREATE TYPE product_type AS ENUM (
  'sayuran',
  'protein',
  'buah',
  'snack'
);

CREATE TABLE product (
  id serial PRIMARY KEY,
  name varchar(100) UNIQUE NOT NULL,
  product_code varchar(100) UNIQUE NOT NULL,
  product_type product_type NOT NULL,
  description text,
  price integer NOT NULL,
  register_date date DEFAULT CURRENT_DATE
);

