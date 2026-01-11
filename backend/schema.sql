-- Add NOT NULL to user_id foreign keys
CREATE TABLE users (
    user_id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);

CREATE TABLE users_profile(
    user_id BIGINT PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE NOT NULL,  -- ✅ NOT NULL
    date_of_birth DATE,
    email VARCHAR(255) UNIQUE,
    height DECIMAL,
    weight DECIMAL,
    is_trainer   BOOLEAN NOT NULL DEFAULT FALSE,
    is_vip       BOOLEAN NOT NULL DEFAULT FALSE,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE food (
    food_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(user_id) NOT NULL,  -- ✅ NOT NULL
    food_name VARCHAR(255) NOT NULL,
    calories_100 DOUBLE PRECISION NOT NULL,
    protein_100 DOUBLE PRECISION NOT NULL,
    carbs_100 DOUBLE PRECISION NOT NULL,
    fats_100 DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE food_Cache (
    food_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(user_id) NOT NULL,  -- ✅ NOT NULL
    food_name VARCHAR(255) NOT NULL,
    calories_100 DOUBLE PRECISION NOT NULL,
    protein_100 DOUBLE PRECISION NOT NULL,
    carbs_100 DOUBLE PRECISION NOT NULL,
    fats_100 DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE recipes (
    recipe_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(user_id) NOT NULL,  -- ✅ NOT NULL
    recipe_name VARCHAR(255) NOT NULL,
    instructions TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE recipe_ingredients (
    ingredient_id BIGSERIAL PRIMARY KEY,
    recipe_id BIGINT NOT NULL REFERENCES recipes(recipe_id) NOT NULL,  -- ✅ NOT NULL
    food_id BIGINT NOT NULL REFERENCES food(food_id) NOT NULL,  -- ✅ NOT NULL
    total_grams DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE food_entries (
    nutrition_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id) NOT NULL,  -- ✅ NOT NULL
    food_id BIGINT REFERENCES food(food_id),  -- ❓ Can be NULL if using recipe_id
    recipe_id BIGINT REFERENCES recipes(recipe_id),  -- ❓ Can be NULL if using food_id
    calories DOUBLE PRECISION NOT NULL,
    total_grams DOUBLE PRECISION NOT NULL,
    protein DOUBLE PRECISION NOT NULL,
    carbs DOUBLE PRECISION NOT NULL,
    fats DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Add constraint: either food_id OR recipe_id must be set
    CONSTRAINT chk_food_or_recipe CHECK (
        (food_id IS NOT NULL AND recipe_id IS NULL) OR 
        (food_id IS NULL AND recipe_id IS NOT NULL) 
    )
);

CREATE TABLE training_routine (
    routine_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id) NOT NULL,  -- ✅ NOT NULL
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    training_routine_name VARCHAR(255),
    notes VARCHAR(255)
);

CREATE TABLE training (
    training_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id) NOT NULL,  -- ✅ NOT NULL
    routine_id BIGINT REFERENCES training_routine(routine_id),  -- ❓ Can be NULL if standalone
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    exercise_name VARCHAR(255),
    notes VARCHAR(255)
);

CREATE TABLE training_ingredients (
    training_entry_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id) NOT NULL,  -- ✅ NOT NULL
    routine_id BIGINT REFERENCES training_routine(routine_id),  -- ❓ Can be NULL
    exercise_name VARCHAR(255),
    weight_ DOUBLE PRECISION NOT NULL,
    sets_ INTEGER NOT NULL,
    reps INTEGER NOT NULL,
    notes VARCHAR(255)
);