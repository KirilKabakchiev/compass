BEGIN;

DROP TYPE IF EXISTS system_auth_access;

ALTER TABLE system_auths
    DROP COLUMN IF EXISTS access_level;

DROP TABLE IF EXISTS system_auth_restrictions;

ALTER TABLE applications
    ADD COLUMN integration_system_id varchar(256) REFERENCES integration_systems(id);

COMMIT;