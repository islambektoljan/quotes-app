-- Индексы для таблицы quotes
CREATE INDEX IF NOT EXISTS idx_quotes_category_id ON quotes(category_id);
CREATE INDEX IF NOT EXISTS idx_quotes_user_id ON quotes(user_id);
CREATE INDEX IF NOT EXISTS idx_quotes_created_at ON quotes(created_at);
CREATE INDEX IF NOT EXISTS idx_quotes_likes_count ON quotes(likes_count);

-- Индексы для таблицы comments
CREATE INDEX IF NOT EXISTS idx_comments_quote_id ON comments(quote_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at);

-- Индексы для таблицы users
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- Уникальные индексы для предотвращения дублирования лайков
CREATE UNIQUE INDEX IF NOT EXISTS idx_quote_likes_user_quote ON quote_likes(user_id, quote_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_comment_likes_user_comment ON comment_likes(user_id, comment_id);