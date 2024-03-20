CREATE TYPE product_type AS ENUM (
  'sayuran',
  'protein',
  'buah',
  'snack'
);

CREATE TABLE product (
                     id serial PRIMARY KEY,
  name varchar(100) UNIQUE NOT NULL,
  type product_type NOT NULL,
  description text,
  price DECIMAL(10, 2) NOT NULL,
  date_added date NOT NULL
);

