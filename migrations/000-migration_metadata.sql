CREATE TABLE migration_metadata (
    migration_id int not null,
    migrated_at_utc TEXT default CURRENT_TIMESTAMP
);

CREATE VIEW latest_migration_id AS
SELECT migration_id
FROM migration_metadata
WHERE migrated_at_utc = (
    SELECT MAX(migrated_at_utc) AS latest_date
    FROM migration_metadata
);
