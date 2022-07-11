package psql

const (
	createUserQuery = `INSERT INTO users 
				(first_name, last_name, 
				email, password, role, avatar, 
				phone_number, address, city, country, 
				postcode, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, 
				$7, $8, $9, $10, $11, now(), now()) 
			RETURNING *`

	updateUserQuery = `UPDATE users 
					SET first_name = COALESCE(NULLIF($1, ''), first_name),
						last_name = COALESCE(NULLIF($2, ''), last_name),
						email = COALESCE(NULLIF($3, ''), email),
						role = COALESCE(NULLIF($4, ''), role),
						avatar = COALESCE(NULLIF($5, ''), avatar),
						phone_number = COALESCE(NULLIF($6, ''), phone_number),
						address = COALESCE(NULLIF($7, ''), address),
						city = COALESCE(NULLIF($8, ''), city),
						country = COALESCE(NULLIF($9, ''), country),
						postcode = COALESCE(NULLIF($10, 0), postcode),
						updated_at = now()
					WHERE user_id = $12
					RETURNING *`

	deleteUserQuery = `DELETE FROM users WHERE user_id = $1`

	getUserByID = `SELECT user_id, first_name, last_name, 
					email, password, role, avatar, 
					phone_number, address, city, country, 
					postcode, created_at, updated_at
				FROM users
				WHERE user_id = $1`

	findUsersByName = `SELECT first_name, last_name, 
						email, password, role, avatar, 
						phone_number, address, city, country, 
						postcode, created_at, updated_at
					FROM users
					WHERE first_name ILIKE '%' $1 '%' or last_name ILIKE '%' $1 '%'
					ORDER BY first_name, last_name`

	getUsers = `SELECT first_name, last_name, 
				email, password, role, avatar, 
				phone_number, address, city, country, 
				postcode, created_at, updated_at
			FROM users
			WHERE user_id < (user_id + $1)
			ORDER BY user_id DESC, COALESCE(NULLIF($2, ''), first_name)
			LIMIT $3
			`

	getTotal = `SELECT COUNT(user_id) FROM users`

	getTotalCount = `SELECT COUNT(user_id) 
					FROM users 
					WHERE first_name ILIKE '%' || $1 || '%' 
						or last_name ILIKE '%' || $1 || '%'`

	findUserByEmail = `SELECT irst_name, last_name, 
						email, password, role, avatar, 
						phone_number, address, city, country, 
						postcode, created_at, updated_at
					FROM users
					WHERE email = $1`
)
