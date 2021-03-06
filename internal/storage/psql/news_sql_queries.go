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

	getTotalNewsCount = `SELECT COUNT(news_id) FROM news`

	getNews = `SELECT news_id, author_id, title, content, image_url, category, updated_at, created_at 
			FROM news
			WHERE news_id < (news_id + $1)
			ORDER BY news_id DESC, created_at, updated_at
			LIMIT $2`

	getNewsByID = `SELECT n.news_id,
				n.title,
				n.content,
				n.updated_at,
				n.image_url,
				n.category,
				CONCAT(u.first_name, ' ', u.last_name) as author,
				u.user_id as author_id
			FROM news n
				LEFT JOIN users u on u.user_id = n.author_id
			WHERE news_id = $1`

	findByTitle = `SELECT  news_id, author_id, title, content, image_url, category, updated_at, created_at
				FROM news	
				WHERE title ILIKE '%' || $1 || '%' 
					and news_id < (news_id + $2)
				ORDER BY news_id DESC, title, created_at, updated_at
				LIMIT $3`

	getTitleCount = `SELECT COUNT(title)
					FROM news
					WHERE title ILIKE '%' || $1 || '%'`
)