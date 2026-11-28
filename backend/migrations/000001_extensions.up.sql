-- citext: a text type that compares case-insensitively. users.email uses it, so
-- "Ali@x.com" and "ali@x.com" count as the same address.
-- It is a "trusted" extension, so the database owner may create it.
CREATE EXTENSION IF NOT EXISTS citext;
