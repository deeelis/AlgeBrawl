CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    login TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS equation_sets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS equations (
    id UUID PRIMARY KEY,
    set_id UUID NOT NULL REFERENCES equation_sets(id),
    equation_text TEXT NOT NULL,
    root1 FLOAT NOT NULL,
    root2 FLOAT NOT NULL,
    user_answer1 FLOAT,
    user_answer2 FLOAT,
    solved_correctly BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_equation_sets_user_id ON equation_sets(user_id);
CREATE INDEX IF NOT EXISTS idx_equations_set_id ON equations(set_id);
