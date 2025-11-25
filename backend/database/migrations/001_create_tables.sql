-- Пользователи
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(50) UNIQUE NOT NULL,
                                     email VARCHAR(100) UNIQUE NOT NULL,
                                     password_hash VARCHAR(255) NOT NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Категории цитат
CREATE TABLE IF NOT EXISTS categories (
                                          id SERIAL PRIMARY KEY,
                                          name VARCHAR(100) UNIQUE NOT NULL,
                                          description TEXT
);

-- Цитаты
CREATE TABLE IF NOT EXISTS quotes (
                                      id SERIAL PRIMARY KEY,
                                      content TEXT NOT NULL,
                                      author VARCHAR(100),
                                      user_id INTEGER,
                                      category_id INTEGER,
                                      likes_count INTEGER DEFAULT 0,
                                      dislikes_count INTEGER DEFAULT 0,
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
                                      FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

-- Лайки/дизлайки цитат
CREATE TABLE IF NOT EXISTS quote_likes (
                                           id SERIAL PRIMARY KEY,
                                           quote_id INTEGER NOT NULL,
                                           user_id INTEGER NOT NULL,
                                           type VARCHAR(10),
                                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                           FOREIGN KEY (quote_id) REFERENCES quotes(id) ON DELETE CASCADE,
                                           FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Комментарии
CREATE TABLE IF NOT EXISTS comments (
                                        id SERIAL PRIMARY KEY,
                                        content TEXT NOT NULL,
                                        quote_id INTEGER NOT NULL,
                                        user_id INTEGER,
                                        likes_count INTEGER DEFAULT 0,
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        FOREIGN KEY (quote_id) REFERENCES quotes(id) ON DELETE CASCADE,
                                        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Лайки комментариев
CREATE TABLE IF NOT EXISTS comment_likes (
                                             id SERIAL PRIMARY KEY,
                                             comment_id INTEGER NOT NULL,
                                             user_id INTEGER NOT NULL,
                                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                             FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
                                             FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);