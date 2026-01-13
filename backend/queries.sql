-- name: CreateUser :one
INSERT INTO users (username, hashed_password)
VALUES ($1, $2)
RETURNING user_id, username;

-- name: GetUserByUsername :one
SELECT user_id,username, hashed_password
FROM users
WHERE username = $1;

-- name: GetUserByID :one
SELECT user_id,username,hashed_password
FROM users
WHERE user_id = $1;

-- name: CreateFoodItem :one
INSERT INTO food(user_id,food_name,calories_100,protein_100,carbs_100,fats_100)
VALUES($1,$2,$3,$4,$5,$6)
RETURNING user_id,food_name,calories_100,protein_100,carbs_100,fats_100;

-- name: CreateFoodCacheItem :one
INSERT INTO food_Cache(user_id,food_name,calories_100,protein_100,carbs_100,fats_100)
VALUES($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: LogFoodItem :one
INSERT INTO food_entries (
    user_id,
    food_id,
    recipe_id,
    calories,
    total_grams,
    protein,
    carbs,
    fats
) VALUES (
    $1,  -- user_id (BIGINT, NOT NULL)
    $2,  -- food_id (BIGINT, can be NULL)
    $3,  -- recipe_id (BIGINT, can be NULL) 
    $4,  -- calories (DOUBLE PRECISION, NOT NULL)
    $5,  -- total_grams (DOUBLE PRECISION, NOT NULL)
    $6,  -- protein (DOUBLE PRECISION, NOT NULL)
    $7,  -- carbs (DOUBLE PRECISION, NOT NULL)
    $8   -- fats (DOUBLE PRECISION, NOT NULL)
)
RETURNING *;

-- name: ViewFood :many
SELECT calories, protein, carbs, fats
FROM food_entries
WHERE user_id = $1 
  AND created_at BETWEEN $2 AND $3
ORDER BY created_at;
;

-- name: ViewFoodTotal :one
SELECT 
  SUM(calories)::float as total_calories,
  SUM(protein)::float as total_protein,
  SUM(carbs)::float as total_carbs,
  SUM(fats)::float as total_fats
FROM food_entries
WHERE user_id = $1 
  AND created_at BETWEEN $2 AND $3;


-- name: LogExercise :one
INSERT INTO exercise_entries(user_id,exercise_name,weight,sets,reps,rpe,notes)
VALUES($1,$2,$3,$4,$5,$6,$7)
RETURNING *;

