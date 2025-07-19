#!/bin/bash

set -e

BASE_URL="http://localhost:8080"

echo " Запрос #1: parse example.com (первый раз, максимум 20 сек)"
START=$(date +%s)
response1=$(curl -s --max-time 20 "$BASE_URL/parse?domain=example.com")
END=$(date +%s)
DURATION=$((END - START))

echo "  Время: $DURATION сек"
echo " Ответ: $response1"
echo

echo "🔄 Запрос #2: parse example.com (повторно, должно быть мгновенно)"
START=$(date +%s)
response2=$(curl -s "$BASE_URL/parse?domain=example.com")
END=$(date +%s)
DURATION=$((END - START))

echo " Время: $DURATION сек"
echo " Ответ: $response2"
echo

echo " Проверка фильтрации + пагинации"
curl -s "$BASE_URL/domains?domain=example&limit=5&page=1"
echo

echo " Получаем все домены"
curl -s "$BASE_URL/domains"
echo

echo " Удаляем домен id=1"
curl -s -X DELETE "$BASE_URL/domains/1"
echo

echo " Проверяем, что домен удалён"
curl -s "$BASE_URL/domains"
echo

echo " Тест на несуществующий домен (badsite.abcdef)"
curl -s "$BASE_URL/parse?domain=badsite.abcdef"
echo

echo " Тест на очистку всей базы"
curl -s -X POST "$BASE_URL/domains/clear"

echo
echo "Программа завершена успешно"