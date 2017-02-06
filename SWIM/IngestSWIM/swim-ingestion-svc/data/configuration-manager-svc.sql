DROP ROLE IF EXISTS "configuration-manager-svc";
CREATE ROLE "configuration-manager-svc" LOGIN PASSWORD 'password';
GRANT SELECT, UPDATE, INSERT, DELETE ON ALL TABLES IN SCHEMA "public" TO "configuration-manager-svc";

DROP TABLE IF EXISTS dynamic;

CREATE TABLE dynamic (
  noun VARCHAR(1024),
  id   CHARACTER(36),
  data BYTEA,
  PRIMARY KEY (noun, id)
);