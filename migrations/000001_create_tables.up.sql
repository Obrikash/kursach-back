CREATE TABLE pools (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    type VARCHAR(20) CHECK (type IN ('Спортивный', 'Оздоровительный', 'комбинированный')) NOT NULL
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL CHECK (name IN ('тренер', 'клиент', 'админ'))
);

CREATE TABLE group_category (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL CHECK (name IN ('начинающие', 'подростки', 'взрослые', 'спортсмены')) 
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    full_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    role_id INT REFERENCES roles(id) NOT NULL
);

CREATE TABLE trainers (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) UNIQUE,
    pool_id INT REFERENCES pools(id) NOT NULL
);

CREATE TABLE training_groups (
    id SERIAL PRIMARY KEY,
    pool_id INT REFERENCES pools(id) NOT NULL,
    category_id INT REFERENCES group_category(id) NOT NULL,
    trainer_id INT REFERENCES trainers(id) NOT NULL
);

CREATE TABLE user_groups (
    user_id INT REFERENCES users(id),
    group_id INT REFERENCES training_groups(id),
    PRIMARY KEY (user_id, group_id)
);

CREATE TABLE schedules (
    id SERIAL PRIMARY KEY,
    group_id INT REFERENCES training_groups(id) NOT NULL,
    day_of_week INT CHECK (day_of_week BETWEEN 1 AND 7) NOT NULL,
    time_of_day TIME WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE abonements (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    visits_per_week INT NOT NULL CHECK (visits_per_week IN (1, 2, 3, 5))
);

CREATE TABLE user_abonements (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    abonement_id INT REFERENCES abonements(id) ON DELETE RESTRICT,
    price DECIMAL(10, 2) NOT NULL,
    date_start TIMESTAMP(0) NOT NULL,
    date_end TIMESTAMP(0) NOT NULL,
    PRIMARY KEY (user_id, abonement_id)
);

