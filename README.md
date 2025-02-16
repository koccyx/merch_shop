# Avito Trainee assignment 2025

### Технологический стек

- **Go 1.23.3+**
- **PostgreSQL**
- **Docker** и **Docker Compose**

### Установка и запуск

#### Предварительные требования

- Go 1.23.3 или выше
- Docker и Docker Compose
- PostgreSQL (если запуск без Docker)

#### Запуск через Docker

```bash
# Клонирование репозитория
git clone https://github.com/your-username/avito-shop.git
cd avito_assignment

# Запуск через Docker Compose
docker-compose up --build    
```

#### Локальный запуск

1. Клонируйте репозиторий
```bash
git clone https://github.com/your-username/avito-shop.git
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