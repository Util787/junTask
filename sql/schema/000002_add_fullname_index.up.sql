CREATE INDEX idx_users_fullname ON users (name, surname, patronymic);
--TODO:check with explain analyze and change index logic if necessarry