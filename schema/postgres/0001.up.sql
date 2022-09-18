CREATE TABLE campaigns (
    id SERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaigns (id),
    name TEXT,
    description TEXT,
    priority SERIAL,
    removed BOOLEAN,
    created_at timestamp
);

INSERT INTO campaigns (name) VALUES ('DnD');

CREATE UNIQUE INDEX ON campaigns (id);
CREATE UNIQUE INDEX ON items (id);
CREATE UNIQUE INDEX ON items (campaign_id);
CREATE INDEX ON items (name);

