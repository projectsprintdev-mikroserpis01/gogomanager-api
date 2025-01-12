CREATE TABLE managers (
	id SERIAL PRIMARY KEY,
	email VARCHAR(255) NOT NULL,
	password VARCHAR(255) NOT NULL,
	name VARCHAR(255),
	user_image_uri VARCHAR(4096),
	company_name VARCHAR(255),
	company_image_uri VARCHAR(4096),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
