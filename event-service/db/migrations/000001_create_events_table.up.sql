CREATE TABLE event(
  id VARCHAR(36) PRIMARY KEY,
  company TEXT NOT NULL,
  team TEXT NOT NULL,
  environment TEXT NOT NULL,
  event_type TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  payload JSON,
  debug JSON
);
