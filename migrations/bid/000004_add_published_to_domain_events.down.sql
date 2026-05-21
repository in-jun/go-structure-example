DROP INDEX IF EXISTS idx_domain_events_unpublished;
ALTER TABLE domain_events DROP COLUMN IF EXISTS published;
