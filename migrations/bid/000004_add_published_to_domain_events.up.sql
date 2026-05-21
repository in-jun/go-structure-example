ALTER TABLE domain_events ADD COLUMN published BOOLEAN NOT NULL DEFAULT FALSE;
CREATE INDEX idx_domain_events_unpublished ON domain_events(published) WHERE published = FALSE;
