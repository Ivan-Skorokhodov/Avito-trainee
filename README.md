# Avito-trainee

Для запуска просто выполнить `docker compose up --build` .

Сделал то, что успел:
- Разнес микросервис по слоям;
- Добавил кастомные middleware для логов и panic recover;
- "Вынес" .env в docker compose и собираю их в кофиг на старте app;
- Накатываю на Postgres миграции при запуске;
- Добавил линтер (golangci-lint);
- Вынес ошибки и json decode/encode в отдельные пакеты;
- Сгенерировал mock и написал unit-тесты для usecase.