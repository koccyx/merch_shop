### Установка и запуск

#### Предварительные требования

- Go 1.23.3 или выше
- Docker и Docker Compose
- PostgreSQL (если запуск без Docker)

#### Запуск через Docker

```bash
# Клонирование репозитория
git clone https://github.com/koccyx/avito_assignment.git
cd avito_assignment

# Запуск через Docker Compose
docker-compose up --build    
```

#### Локальный запуск

1. Клонируйте репозиторий
```bash
git clone https://github.com/koccyx/avito_assignment.git
cd avito-shop
```

2. Установите зависимости
```bash
go mod download
```

3. Запустите PostgreSQL и примените миграции
```bash
make migrate
```

4. Запустите приложение
```bash
make run
```

#### Локальный запуск тестов

```bash
make test
```
