CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);

CREATE TABLE users_profile(
    user_id INTEGER PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
    date_of_birth DATE,
    email VARCHAR(255) UNIQUE,
    height DECIMAL,
    weight DECIMAL,
    is_trainer   BOOLEAN NOT NULL DEFAULT FALSE,
    is_vip       BOOLEAN NOT NULL DEFAULT FALSE,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE food (
    food_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    food_name VARCHAR(255) NOT NULL,
    calories_100 DECIMAL NOT NULL,
    protein_100 DECIMAL NOT NULL,
    carbs_100 DECIMAL NOT NULL,
    fats_100 DECIMAL NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE food_Cache (
    food_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    food_name VARCHAR(255) NOT NULL,
    is_solid BOOLEAN,
    calories_100 DECIMAL NOT NULL,
    protein_100 DECIMAL NOT NULL,
    carbs_100 DECIMAL NOT NULL,
    fats_100 DECIMAL NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE recipes (
    recipe_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    recipe_name VARCHAR(255) NOT NULL,
    instructions TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE recipe_ingredients (
    ingredient_id SERIAL PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES recipes(recipe_id),
    food_id INTEGER NOT NULL REFERENCES food(food_id),
    total_grams DECIMAL NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE food_entries (
    nutrition_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ,
    food_id INTEGER REFERENCES food(food_id),
    recipe_id INTEGER REFERENCES recipes(recipe_id),
    calories DECIMAL NOT NULL,
    total_grams DECIMAL NOT NULL,
    protein DECIMAL NOT NULL,
    carbs DECIMAL NOT NULL,
    fats DECIMAL NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);



CREATE TABLE training_routine (
    routine_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    training_routine_name VARCHAR(255),

    notes VARCHAR(255)
);

CREATE TABLE training (
    training_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ,
    routine_id INTEGER REFERENCES training_routine(routine_id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    exercise_name VARCHAR(255),
    notes VARCHAR(255)
);

CREATE TABLE training_ingredients (
    training_entry_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id),
    routine_id INTEGER REFERENCES training_routine(routine_id),
    exercise_name VARCHAR(255),
    weight_ DECIMAL NOT NULL,
    sets_ INTEGER NOT NULL,
    reps INTEGER NOT NULL,
    notes VARCHAR(255)
);