
CREATE TABLE IF NOT EXISTS domains (
    id SERIAL PRIMARY KEY,
    domain TEXT UNIQUE NOT NULL,
    legal_name TEXT NOT NULL,
    country TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

DELETE FROM domains 
WHERE legal_name LIKE '%Ошибка получени%данных%'
   OR country LIKE '%Ошибка получени%данных%'
   OR legal_name = 'Ошибка получения данных' 
   OR country = 'Ошибка получения данных';

-- Создаем простое ограничение на уровне таблицы
ALTER TABLE domains 
ADD CONSTRAINT check_no_error_data 
CHECK (
    legal_name NOT LIKE '%Ошибка получени%' 
    AND country NOT LIKE '%Ошибка получени%'
);