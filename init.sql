CREATE DATABASE allworkusers;

\c allworkusers;

CREATE TYPE user_type AS ENUM ('customer', 'service_provider');

CREATE TABLE IF NOT EXISTS user_account (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  phone_number VARCHAR(20) NOT NULL,
  address VARCHAR(255) NOT NULL,
  user_type user_type NOT NULL
);