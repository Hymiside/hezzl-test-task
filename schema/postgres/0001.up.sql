CREATE TABLE campaigns (
    id SERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER PRIMARY KEY REFERENCES campaigns (id),
    name TEXT,
    description TEXT,
    priority INTEGER,
    removed BOOLEAN,
    created_at timestamp
);

INSERT INTO campaigns (name) VALUES ('DnD');