ALTER TABLE applications
ADD COLUMN client_id VARCHAR(32) NOT NULL,
ADD COLUMN client_secret VARCHAR(64) NOT NULL DEFAULT '';
