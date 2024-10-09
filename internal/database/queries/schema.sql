CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    type INTEGER CHECK (type >= 0 AND type <= 9) NOT NULL,
    value INTEGER CHECK (value >= 0 AND value <= 99) NOT NULL,
    state task_state NOT NULL,
    creationtime BIGINT NOT NULL,
    lastupdatetime BIGINT NOT NULL
);