<div align="center">

  # Retail Vision Catalog Backend

  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go"></a>
  <a href="https://www.postgresql.org/"><img src="https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL"></a>
  <a href="https://redis.io/"><img src="https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white" alt="Redis"></a>
  <a href="https://kafka.apache.org/"><img src="https://img.shields.io/badge/Apache_Kafka-231F20?style=for-the-badge&logo=apache-kafka&logoColor=white" alt="Kafka"></a>
  <a href="https://qdrant.tech/"><img src="https://img.shields.io/badge/Qdrant-f82329?style=for-the-badge&logo=qdrant&logoColor=white" alt="Qdrant"></a>
  <a href="https://min.io/"><img src="https://img.shields.io/badge/MinIO-C72E49?style=for-the-badge&logo=minio&logoColor=white" alt="MinIO"></a>
  <a href="https://grpc.io/"><img src="https://img.shields.io/badge/gRPC-4285F4?style=for-the-badge&logo=grpc&logoColor=white" alt="gRPC"></a>
  <a href="https://swagger.io/"><img src="https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=swagger&logoColor=black" alt="Swagger"></a>
  <a href="https://www.docker.com/"><img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker"></a>
</div>

**Retail Vision Catalog** – это высокопроизводительный бэкенд-сервис для систем умного ритейла. Мы избавляем покупателей от необходимости искать штрихкоды, используя технологии компьютерного зрения для автоматического распознавания товаров на кассе.

## 🚀 Быстрый старт
1. Клонируйте репозитории (backend и ml-service)
    ```bash
    git clone git@github.com:DRSN-tech/go-backend.git
    gir clone git@github.com:DRSN-tech/ml-service.git
    ```
2. Перейдите в корневую папку проекта и добавьте подмодули (submodules)
    ```bash
    cd go-backend/
    git submodule update --init --recursive --remote
    ```
3. Создайте `.env` файл. Образец переменных окружения находится в файле `.env.example`
4. Выполните `sudo docker-compose up --build -d` в корневой папке проекта `go-backend`

## ⚡️ Преимущества
- **Computer Vision First**: Распознавание товаров в реальном времени без использования штрихкодов.
- **Reliability**: Надежная идентификация товаров даже при частичном перекрытии или плохом освещении.
- **Speed**: Максимальное ускорение пути покупателя «от полки до оплаты».
- **Seamless Integration**: Быстрая передача данных в кассовое ПО.

## 📖 API Documentation
- **Swagger**: `http://localhost:8080/swagger/index.html` (после запуска)
- **gRPC**: Описание сервисов в папке `api/proto/`

## 📊 Диаграммы последовательности
Регистрация продукта и ML-обработки товара
![register_product](images/register_product.svg)

Получение списка продуктов
![get_products](images/get_products.svg)

Схема работы фонового воркера и доставки событий
![worker](images/worker.svg)

## ⚙️ Стек технологий

| Категория | Используемые инструменты |
| :--- | :--- |
| **Основной язык** | [![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/) **1.24** |
| **Базы данных** | [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=flat-square&logo=postgresql&logoColor=white)](https://www.postgresql.org/) + [![Redis](https://img.shields.io/badge/Redis-DC382D?style=flat-square&logo=redis&logoColor=white)](https://redis.io/) (кэширование) |
| **ML & Search** | [![gRPC](https://img.shields.io/badge/gRPC-4285F4?style=flat-square&logo=grpc&logoColor=white)](https://grpc.io/) + [![Qdrant](https://img.shields.io/badge/Qdrant-f82329?style=flat-square&logo=qdrant&logoColor=white)](https://qdrant.tech/) (векторный поиск) |
| **Messaging** | [![Kafka](https://img.shields.io/badge/Apache_Kafka-231F20?style=flat-square&logo=apache-kafka&logoColor=white)](https://kafka.apache.org/) (событийная архитектура) |
| **Хранилище** | [![MinIO](https://img.shields.io/badge/MinIO-C72E49?style=flat-square&logo=minio&logoColor=white)](https://min.io/) (S3 совместимое) |
| **Миграции** | [![Go Migrate](https://img.shields.io/badge/Migrate-00ADD8?style=flat-square&logo=go&logoColor=white)](https://github.com/golang-migrate/migrate) (versioning) |
| **API & Docs** | [![Chi](https://img.shields.io/badge/Chi-00ADD8?style=flat-square&logo=go&logoColor=white)](https://github.com/go-chi/chi) + [![Swagger](https://img.shields.io/badge/Swagger-85EA2D?style=flat-square&logo=swagger&logoColor=black)](https://swagger.io/) + [![gRPC](https://img.shields.io/badge/gRPC-4285F4?style=flat-square&logo=grpc&logoColor=white)](https://grpc.io/) |
| **DevOps** | [![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white)](https://www.docker.com/) |

## ⌛️ Будущие изменения 
- Рефакторинг `RegisterNewProduct` — Переход на надёжную асинхронную регистрацию через **Outbox** паттерн. Обеспечить максимальную надёжность и восстанавливаемость при регистрации нового продукта с изображениями и векторами. Устранить риск орфаненных файлов в MinIO и несогласованности с Qdrant/Kafka при любых сбоях (включая SIGKILL, сетевые ошибки, перезапуски). Единственным источником истины будет основная база данных – в нашем случае PostgreSQL.
- Удалить `pkg/logger` и перейти на `slog/logger`. Реализовать информативное логирование, продумать метрики.
