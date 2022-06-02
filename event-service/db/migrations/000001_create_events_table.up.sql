CREATE TABLE event(
  id BINARY(16) PRIMARY KEY,
  company TEXT NOT NULL,
  team TEXT NOT NULL,
  environment TEXT NOT NULL,
  event_type TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  message TEXT,
  debug JSON
);
