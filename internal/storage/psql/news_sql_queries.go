package psql

const (
	createNews = `INSERT INTO news (author_id, title, content, image_url, category, created_at)
				VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), now())
				RETURNING *`

	updateNews = `UPDATE news
				SET title = COALESCE(NULLIF($1, ''), title),
					content = COALESCE(NULLIF($2, ''), content),
					image_url = COALESCE(NULLIF($3, ''), image_url),
					category = COALESCE(NULLIF($4, ''), category),
					updated_at = now()
				WHERE news_id = $5
				RETURNING *`

	deleteNews = `DELETE FROM news WHERE news_id = $1`
)