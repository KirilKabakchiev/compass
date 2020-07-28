BEGIN;

CREATE TYPE system_auth_access_level AS ENUM (
    'RESTRICTED',
    'UNRESTRICTED'
    );

ALTER TABLE system_auths
    ADD COLUMN access_level system_auth_access_level DEFAULT 'UNRESTRICTED' ::system_auth_access_level NOT NULL;

CREATE TABLE system_auth_restrictions(
  id uuid,
  system_auth_id uuid REFERENCES system_auths(id) ON DELETE CASCADE,
  app_id uuid REFERENCES applications(id) ON DELETE CASCADE,
  runtime_id uuid REFERENCES runtimes(id) ON DELETE CASCADE,
  integration_system_id uuid REFERENCES integration_systems(id) ON DELETE CASCADE,
  application_template_id uuid REFERENCES app_templates(id) ON DELETE CASCADE
);

ALTER TABLE applications
    DROP COLUMN  integration_system_id;

COMMIT;