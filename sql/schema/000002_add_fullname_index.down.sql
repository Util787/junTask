DROP INDEX idx_users_name_trgm;
DROP INDEX idx_users_surname_trgm;
DROP INDEX idx_users_patronymic_trgm;
DROP EXTENSION IF EXISTS pg_trgm;