BEGIN;

DELETE FROM users;
DROP TABLE users_transactions;
DROP TABLE users_data;
DROP TABLE users;
DROP TABLE schema_migrations;
DROP TYPE "e_role";

COMMIT;
