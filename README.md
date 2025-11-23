# PR Reviewer

Простой сервис для управления Pull Request-ами и ревьюерами в командах.

## Технологии

- Go 1.27
- PostgreSQL 15
- Docker + Docker Compose
- migrate/migrate для миграций

## Быстрый старт

```bash
# 1. Клонируем и заходим
git clone https://github.com/vysotskaya-a/pr-reviewer
cd pr-reviewer

docker-compose up --build