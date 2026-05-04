# Virtual Chat

Локальный чат с ИИ на базе Ollama.

---

## Быстрый запуск

### 1. Клонировать репозиторий
```bash
git clone https://github.com/Vlad-Ali/Virtual-chat.git
cd virtual-chat
```

### 2. Настроить `.env`
```bash
cp .env.example .env
```

Минимальный `.env` для примера:
```env
HTTP_ADDRESS=:8080
OLLAMA_URL=http://ollama:11434
OLLAMA_MODEL_NAME=qwen2.5:3b
FRONTEND_URL=http://localhost
```

### 3. Скачать модель
```bash
docker-compose up -d ollama

docker exec ollama ollama pull qwen2.5:3b
```

### 4. Запустить всё
```bash
docker-compose up -d
```

### 5. Открыть в браузере
- [http://localhost](http://localhost)

- `/` — описание проекта
- `/chat.html` — чат с ИИ

---

## Полезные команды

```bash
docker-compose logs -f

docker-compose down

docker-compose build virtual-chat && docker-compose up -d

docker exec ollama ollama list
```

---


## HTTPS / WSS (для продакшена)

1. Положите сертификаты в `nginx/ssl/cert.pem` и `key.pem`
2. В `nginx.conf` замените `listen 80` на:
```nginx
listen 443 ssl;
ssl_certificate /etc/nginx/ssl/cert.pem;
ssl_certificate_key /etc/nginx/ssl/key.pem;
```
3. В `.env` укажите: `FRONTEND_URL=https://yourdomain.com`


---

## Структура

```
├── frontend/          # Статика (HTML/CSS/JS)
├── cmd/server/        # Точка входа Go
├── internal/          # Логика приложения
├── docker-compose.yml # Оркестрация
├── nginx.conf         # Конфиг прокси
└── .env               # Переменные окружения
```
