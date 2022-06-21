package psql

const (
	createUserQuery = `INSERT INTO users 
			(first_name, last_name, 
			email, password, role, avatar, 
			phone_number, address, city, country, 
			postcode, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, 
				$7, $8, $9, $10, $11, now(), now()) RETURNING *`
)