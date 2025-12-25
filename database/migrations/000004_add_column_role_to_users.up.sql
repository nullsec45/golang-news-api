CREATE TYPE user_role AS ENUM ('admin', 'writer', 'editor', 'viewer');

ALTER TABLE users ADD COLUMN role user_role DEFAULT 'viewer';