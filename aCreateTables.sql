CREATE TABLE IF NOT EXISTS objects (
    id SERIAL PRIMARY KEY,
    name TEXT,
    clock INTEGER CHECK (clock <= 14 AND clock >= -12)
);

CREATE EXTENSION IF NOT EXISTS btree_gist;
CREATE TABLE IF NOT EXISTS tech_windows (
    id SERIAL PRIMARY KEY,
    id_object INTEGER REFERENCES objects(id) ON DELETE CASCADE,
    duration TSRANGE,
    EXCLUDE USING GIST (id_object WITH =, duration WITH &&)
);