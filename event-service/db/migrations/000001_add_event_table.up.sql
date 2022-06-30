CREATE TABLE event(
										id VARCHAR(36) PRIMARY KEY,
										scope TEXT NOT NULL,
										object TEXT NOT NULL,
										meta JSON,
										event_type TEXT NOT NULL,
										event_family TEXT NOT NULL,
										created_at TIMESTAMP NOT NULL,
										payload JSON,
										debug JSON
);
