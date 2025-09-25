-- 002_seed.sql
INSERT INTO users (name, email, password_hash)
VALUES ('Admin User', 'admin@kerjadiluar.test', '$2a$10$INVALIDPLACEHOLDERFORBCRYPT'); -- later replace with real hash or register via API
