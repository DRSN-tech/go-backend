# Retail Vision Catalog Backend

---
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
4. Выполните `sudo docker-compose up --build -d` в корневой папке проекта