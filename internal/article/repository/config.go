package repository

var (
	getArticleQuery    = `SELECT * FROM articles WHERE id = $1`
	getArticlesQuery   = `SELECT * FROM articles ORDER BY created_at DESC`
	createArticleQuery = `INSERT INTO articles (author_id, title, "desc") 
									VALUES ($1, $2, $3) RETURNING *`
	updateArticleQuery = `UPDATE articles 
									SET title = COALESCE(NULLIF($1, ''), title), 
										"desc" = COALESCE(NULLIF($2, ''), "desc"), 
										updated_at = now() 
									WHERE id = $3 RETURNING *`
	deleteArticleQuery = `DELETE FROM articles WHERE id = $1 
									AND author_id = $2`
)
