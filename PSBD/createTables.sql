CREATE TABLE IF NOT EXISTS objects (
    id SERIAL PRIMARY KEY,
    name TEXT,
    clock INTEGER CHECK (clock < 24 AND clock >= 0)
);

CREATE EXTENSION IF NOT EXISTS btree_gist;
CREATE TABLE IF NOT EXISTS tech_windows (
    id SERIAL PRIMARY KEY,
    id_object INTEGER REFERENCES objects(id) ON DELETE CASCADE,
    duration TSRANGE,
    EXCLUDE USING GIST (id_object WITH =, duration WITH &&)
);

-- INSERT INTO objects(name, clock) VALUES ('cc', 3);

-- INSERT INTO tech_windows(id_object, duration) VALUES(1, '[2021-01-01 00:01, 2021-02-02 00:01)');