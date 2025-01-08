INSERT INTO roles (name) VALUES
('Superadmin'),
('Lead Admin'),
('Admin'),
('User')
ON CONFLICT (id) DO NOTHING;
