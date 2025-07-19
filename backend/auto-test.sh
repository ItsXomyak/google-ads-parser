#!/bin/bash

set -e

BASE_URL="http://localhost:8080"

# Очистка базы
echo "Очищаем базу данных..."
curl -s -X POST "$BASE_URL/domains/clear"
echo

# Первый парсинг — должен занять до 20 сек
echo "Запрос #1: parse batch [example.com] (первый раз)"
START=$(date +%s)
response1=$(curl -s -X POST "$BASE_URL/parse/batch" \
  -H "Content-Type: application/json" \
  -d '{"domains": ["example.com"]}')
END=$(date +%s)
DURATION=$((END - START))

echo "Время: $DURATION сек"
echo "Ответ: $response1"
echo

# Повторный запрос — должен быть из базы, быстро
echo "Запрос #2: parse batch [example.com] (повторно, из БД)"
START=$(date +%s)
response2=$(curl -s -X POST "$BASE_URL/parse/batch" \
  -H "Content-Type: application/json" \
  -d '{"domains": ["example.com"]}')
END=$(date +%s)
DURATION=$((END - START))

echo "Время: $DURATION сек"
echo "Ответ: $response2"
echo

# Проверка фильтрации и пагинации
echo " Проверка фильтрации + пагинации"
curl -s "$BASE_URL/domains?domain=example&limit=5&page=1"
echo

# Получение всех доменов
echo "Получаем все домены"
curl -s "$BASE_URL/domains"
echo

# Удаление домена по ID (допустим ID=1)
echo "🗑 Удаляем домен с ID=1"
curl -s -X DELETE "$BASE_URL/domains/1"
echo

# Проверка, что домен удалён
echo "Список доменов после удаления"
curl -s "$BASE_URL/domains"
echo

# Некорректный домен
echo "Тест на невалидный домен: badsite.abcdef"
curl -s -X POST "$BASE_URL/parse/batch" \
  -H "Content-Type: application/json" \
  -d '{"domains": ["badsite.abcdef"]}'
echo

echo "Все тесты завершены успешно"
