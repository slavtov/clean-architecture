package repository

var (
	getUserQuery  = `SELECT * FROM users WHERE id = $1`
	getUsersQuery = `SELECT id, email, updated_at, created_at 
								FROM users ORDER BY created_at DESC`
	createUserQuery = `INSERT INTO users (email, "password") 
								VALUES ($1, $2) RETURNING *`
	updateUserQuery = `UPDATE users 
								SET email = COALESCE(NULLIF($1, ''), email), 
									"password" = COALESCE(NULLIF($2, ''), "password"), 
									updated_at = now() 
								WHERE id = $3 RETURNING *`
	deleteUserQuery      = `DELETE FROM users WHERE id = $1`
	findUserByEmailQuery = `SELECT * FROM users WHERE email = $1`
)
