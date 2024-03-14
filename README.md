# soa-coursework

Это репозиторий для проекта в рамках курса "Сервис-ориентированные архитектуры".

**Тема**: трекер задач.

**Студент**: Сафарова Элина Ильнуровна, БПМИ215.

## Запуск сервиса

Перед запуском необходимо сгенерировать пару RSA ключей:
```
openssl genrsa -out signature.pem 1024
openssl rsa -in signature.pem -out signature.pub -pubout -outform PEM 
```

Пути до файлов с ключами нужно положить в env переменные:
```
export PRIVATE_KEY_PATH=./signature.pem
export PUBLIC_KEY_PATH=./signature.pub
```

Запуск сервиса в контейнере с базой данных:
```bash
docker-compose up --build
```
