package psql

const (
	createComments = `INSER INTO comments (author_id, news_id, message)
					VALUES ($1, $2, $3)
					RETURNING *`
	
	deleteComment = `DELETE FROM comments WHERE comment_id = $1`

	updateComment = `UPDATE comments
				SET message = $1, updated_at = CURRENT_TIMESTAMP 
				WHERE comment_id = $2 
				RETURNING *`
)