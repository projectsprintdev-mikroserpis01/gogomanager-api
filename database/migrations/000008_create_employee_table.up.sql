CREATE TABLE employees (
    id SERIAL PRIMARY KEY, 
    identity_number VARCHAR(33) NOT NULL CHECK (LENGTH(identity_number) BETWEEN 5 AND 33), 
    name VARCHAR(33) NOT NULL CHECK (LENGTH(name) BETWEEN 4 AND 33), 
    employee_image_uri TEXT, 
    gender VARCHAR(6) NOT NULL CHECK (gender IN ('male', 'female')), 
    department_id INT NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "employees"
ADD CONSTRAINT fk_department
FOREIGN KEY (department_id)
REFERENCES departments(id);