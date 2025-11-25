-- CHECK constraint для типа реакции (like/dislike)
ALTER TABLE quote_likes
    DROP CONSTRAINT IF EXISTS check_quote_like_type;

ALTER TABLE quote_likes
    ADD CONSTRAINT check_quote_like_type
        CHECK (type IN ('like', 'dislike'));

-- NOT NULL constraints для обязательных полей
ALTER TABLE quotes
    ALTER COLUMN content SET NOT NULL;

ALTER TABLE quotes
    ALTER COLUMN author SET NOT NULL;

ALTER TABLE comments
    ALTER COLUMN content SET NOT NULL;

-- Добавляем ограничение на минимальную длину контента (если поддерживается)
ALTER TABLE quotes
    ADD CONSTRAINT min_content_length CHECK (length(content) >= 5);

ALTER TABLE comments
    ADD CONSTRAINT min_comment_length CHECK (length(content) >= 1);