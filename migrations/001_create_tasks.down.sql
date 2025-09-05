CREATE TABLE IF NOT EXISTS tasks (
                                     id UUID PRIMARY KEY,
                                     title TEXT NOT NULL,
                                     description TEXT,
                                     status TEXT NOT NULL,
                                     created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
    );

DROP TABLE IF EXISTS tasks;
