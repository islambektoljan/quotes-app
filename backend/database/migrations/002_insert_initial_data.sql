INSERT INTO categories (name, description) VALUES
                                               ('Мотивация', 'Вдохновляющие цитаты для мотивации'),
                                               ('Юмор', 'Смешные и ироничные цитаты'),
                                               ('Философия', 'Глубокие философские мысли'),
                                               ('Любовь', 'Цитаты о любви и отношениях'),
                                               ('Успех', 'Цитаты об успехе и достижениях'),
                                               ('Жизнь', 'Цитаты о жизни и её смысле')
ON CONFLICT (name) DO NOTHING;

-- Добавляем тестового пользователя
INSERT INTO users (username, email, password_hash) VALUES
    ('admin', 'admin@quotes.app', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi') -- password: password
ON CONFLICT (email) DO NOTHING;