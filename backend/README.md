# Quotes App - Веб-приложение для цитат

RESTful API для управления цитатами с возможностью комментирования, лайков и фильтрации.

## 🚀 Технологии

- **Backend**: Go (Gin framework)
- **Database**: PostgreSQL
- **ORM**: Gorm
- **Authentication**: JWT
- **Containerization**: Docker & Docker Compose

## 📋 Предварительные требования

- Docker & Docker Compose
- Go 1.25 (для локальной разработки)

## 🏃‍♂️ Быстрый старт

### Запуск с Docker (рекомендуется)

```bash
# Клонируйте репозиторий
git clone <repository-url>
cd quotes-app/backend

# Запустите приложение
docker-compose up --build

# Приложение будет доступно по http://localhost:8080
```

### Локальная разработка

```bash
cd backend

# Установите зависимости
go mod download

# Запустите базу данных
docker-compose up db -d

# Запустите приложение
go run main.go
```

## 🗄️ Структура базы данных

Миграции применяются автоматически при запуске контейнера:

- `001_create_tables.sql` - создание таблиц
- `002_insert_initial_data.sql` - начальные данные
- `003_create_indexes.sql` - индексы производительности
- `004_add_constraints.sql` - ограничения целостности

## 🔐 Аутентификация

API использует JWT токены. Добавьте токен в заголовок запроса:
```
Authorization: Bearer <your_jwt_token>
```

## 📚 API Endpoints

### 🔓 Публичные эндпоинты

#### 🔑 Аутентификация

**Регистрация пользователя**
- **URL**: `POST /register`
- **Body**:
```json
{
  "username": "string (3-50 chars)",
  "email": "string (valid email)",
  "password": "string (min 6 chars)"
}
```
- **Response** (201):
```json
{
  "token": "jwt_token",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

**Логин пользователя**
- **URL**: `POST /login`
- **Body**:
```json
{
  "email": "string",
  "password": "string"
}
```
- **Response** (200):
```json
{
  "token": "jwt_token",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

#### 📖 Цитаты

**Получить список цитат**
- **URL**: `GET /quotes`
- **Query Parameters**:
    - `page` - номер страницы (default: 1)
    - `limit` - количество на странице (default: 10)
    - `category_id` - фильтр по категории
    - `author` - поиск по автору
    - `content` - поиск по содержанию
    - `sort` - поле для сортировки (default: created_at)
    - `order` - порядок сортировки (asc/desc, default: desc)
- **Response** (200):
```json
{
  "quotes": [
    {
      "id": 1,
      "content": "Цитата текст...",
      "author": "Автор",
      "user": {"id": 1, "username": "user1"},
      "category": {"id": 1, "name": "Мотивация"},
      "likes_count": 5,
      "dislikes_count": 1,
      "created_at": "2023-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "pages": 10
  }
}
```

**Получить цитату по ID**
- **URL**: `GET /quotes/:id`
- **Response** (200):
```json
{
  "id": 1,
  "content": "Цитата текст...",
  "author": "Автор",
  "user": {"id": 1, "username": "user1"},
  "category": {"id": 1, "name": "Мотивация"},
  "likes_count": 5,
  "dislikes_count": 1,
  "comments": [
    {
      "id": 1,
      "content": "Комментарий текст...",
      "user": {"id": 2, "username": "user2"},
      "likes_count": 2,
      "created_at": "2023-01-01T00:00:00Z"
    }
  ],
  "created_at": "2023-01-01T00:00:00Z"
}
```

#### 📂 Категории

**Получить все категории**
- **URL**: `GET /categories`
- **Response** (200):
```json
[
  {
    "id": 1,
    "name": "Мотивация",
    "description": "Вдохновляющие цитаты для мотивации"
  }
]
```

#### 💬 Комментарии

**Получить комментарии цитаты**
- **URL**: `GET /quotes/:id/comments`
- **Response** (200):
```json
[
  {
    "id": 1,
    "content": "Комментарий текст...",
    "user": {"id": 1, "username": "user1"},
    "likes_count": 2,
    "created_at": "2023-01-01T00:00:00Z"
  }
]
```

### 🔒 Защищенные эндпоинты (требуют JWT токен)

#### ✍️ Цитаты

**Создать цитату**
- **URL**: `POST /quotes`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
```json
{
  "content": "string (1-1000 chars)",
  "author": "string (1-100 chars)",
  "category_id": 1
}
```
- **Response** (201): Объект цитаты

**Обновить цитату**
- **URL**: `PUT /quotes/:id`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
```json
{
  "content": "string (optional)",
  "author": "string (optional)",
  "category_id": 1
}
```
- **Response** (200): Обновленный объект цитаты

**Удалить цитату**
- **URL**: `DELETE /quotes/:id`
- **Headers**: `Authorization: Bearer <token>`
- **Response** (200):
```json
{
  "message": "Quote deleted successfully"
}
```

**Лайк цитаты**
- **URL**: `POST /quotes/:id/like`
- **Headers**: `Authorization: Bearer <token>`
- **Response** (200):
```json
{
  "message": "Reaction updated successfully"
}
```

**Дизлайк цитаты**
- **URL**: `POST /quotes/:id/dislike`
- **Headers**: `Authorization: Bearer <token>`
- **Response** (200):
```json
{
  "message": "Reaction updated successfully"
}
```

#### 💭 Комментарии

**Добавить комментарий**
- **URL**: `POST /quotes/:id/comments`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
```json
{
  "content": "string (1-500 chars)"
}
```
- **Response** (201): Объект комментария

**Лайк комментария**
- **URL**: `POST /comments/:id/like`
- **Headers**: `Authorization: Bearer <token>`
- **Response** (200):
```json
{
  "message": "Like updated successfully"
}
```

**Обновить комментарий**
- **URL**: `PUT /comments/:id`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
```json
{
  "content": "string (1-500 chars)"
}
```
- **Response** (200): Обновленный объект комментария

**Удалить комментарий**
- **URL**: `DELETE /comments/:id`
- **Headers**: `Authorization: Bearer <token>`
- **Response** (200):
```json
{
  "message": "Comment deleted successfully"
}
```

### 🩺 Системные эндпоинты

**Проверка здоровья**
- **URL**: `GET /health`
- **Response** (200):
```json
{
  "status": "OK",
  "database": "connected"
}
```

**Проверка базы данных**
- **URL**: `GET /db-check`
- **Response** (200):
```json
{
  "status": "Database accessible",
  "data": {
    "users_count": 10,
    "quotes_count": 50,
    "categories_count": 6
  }
}
```

## 🐛 Тестирование API

### Примеры с использованием curl:

**Регистрация:**
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```

**Создание цитаты:**
```bash
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{"content":"Life is beautiful","author":"Unknown","category_id":1}'
```

**Получение цитат:**
```bash
curl "http://localhost:8080/quotes?page=1&limit=5&category_id=1"
```

## 🔧 Конфигурация

Создайте файл `.env` для настройки окружения:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=quotes_db
DB_PORT=5432
JWT_SECRET=your_super_secret_jwt_key_here
```

## 📁 Структура проекта

```
backend/
├── config/           # Конфигурация БД и JWT
├── database/
│   └── migrations/   # SQL миграции
├── handlers/         # Обработчики HTTP запросов
├── middleware/       # Промежуточное ПО
├── models/           # Модели данных
└── main.go          # Точка входа
```

## 🐳 Docker команды

```bash
# Запуск
docker-compose up

# Запуск в фоновом режиме
docker-compose up -d

# Просмотр логов
docker-compose logs -f app

# Остановка
docker-compose down

# Пересборка и запуск
docker-compose up --build
```

## 🤝 Разработка

### Добавление новых миграций:

1. Создайте файл в `database/migrations/` с префиксом номера версии
2. Файлы выполняются в алфавитном порядке
3. Используйте `IF NOT EXISTS` для идемпотентности
