CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_users_name_trgm ON users USING gin (name gin_trgm_ops);
CREATE INDEX idx_users_surname_trgm ON users USING gin (surname gin_trgm_ops);
CREATE INDEX idx_users_patronymic_trgm ON users USING gin (patronymic gin_trgm_ops);
--indexes will work only with records >10000