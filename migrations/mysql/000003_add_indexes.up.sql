CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_todos_user_created ON todos(user_id, created_at DESC);
