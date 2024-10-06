CREATE TABLE  IF NOT EXISTS tasks (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    type INTEGER CHECK (type >= 0 AND type <= 9) NOT NULL,
    value INTEGER CHECK (value >= 0 AND value <= 99) NOT NULL,
    state VARCHAR(20) CHECK (state IN ('RECEIVED', 'PROCESSING', 'DONE')) NOT NULL,
    creationtime INTEGER NOT NULL,
    lastupdatetime INTEGER NOT NULL
);
