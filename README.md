# Realty — Сайт продажи недвижимости

Микросервисное приложение на Go с React фронтендом, задеплоенное в Kubernetes (Minikube).

## Стек

**Backend:** Go, gRPC, REST  
**Frontend:** React, TypeScript, Vite, TailwindCSS  
**Инфраструктура:** Kubernetes (Minikube), Docker, PostgreSQL 15, Redis 7, Apache Kafka, Elasticsearch 8, MinIO  
**CI/CD:** GitHub Actions, self-hosted runner  

## Архитектура

                        ┌─────────────┐
                        │  api-gateway │ :8080
                        └──────┬──────┘
               ┌───────────────┼───────────────┐
               ▼               ▼               ▼
        user-service    listing-service   deal-service
           :50051           :50052           :50053
               
               ┌───────────────┐
               ▼               ▼
        search-service  notification-service
           :50054           :50055

Инфраструктура: PostgreSQL · Redis · Kafka · Elasticsearch · MinIO

## Требования

- Ubuntu 22.04+
- Docker 24+
- Minikube v1.35+
- kubectl
- Git

## Быстрый старт после установки

ansible-playbook ~/realty/ansible/playbook.yml


## Полная установка с нуля

### 1. Установи Docker

bash
sudo apt update && sudo apt install -y ca-certificates curl gnupg
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker.gpg] \
  https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list
sudo apt update && sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
sudo usermod -aG docker $USER && newgrp docker

### 2. Установи kubectl и Minikube

# kubectl
curl -LO "https://dl.k8s.io/release/$(curl -sL https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

### 3. Установи migrate

bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/



## Запусти локальный registry
docker run -d -p 5000:5000 --restart=always --name registry registry:2

## Узнай IP Minikube bridge (обычно 192.168.49.1)
ip addr show | grep "inet 192.168"

## Добавь insecure registry в Docker
sudo nano /etc/docker/daemon.json
## Вставь: { "insecure-registries": ["192.168.49.1:5000"] }
sudo systemctl restart docker
docker start registry

# Запусти Minikube
minikube start \
  --driver=docker \
  --cpus=3 \
  --memory=7g \
  --insecure-registry="192.168.49.1:5000"


### 5. Склонируй репозиторий

bash
git clone git@github.com:shmatkovapple-ctrl/realty.git
cd realty


### 6. Собери Docker образы

bash
REGISTRY=192.168.49.1:5000

for service in user-service listing-service deal-service search-service notification-service api-gateway; do
  docker build -t $REGISTRY/realty/$service:latest -f services/$service/Dockerfile .
  docker push $REGISTRY/realty/$service:latest
done

docker build -t $REGISTRY/realty/frontend:latest -f frontend/Dockerfile frontend/
docker push $REGISTRY/realty/frontend:latest

### 7. Задеплой в Kubernetes

bash
kubectl apply -f k8s/namespace.yml
kubectl apply -f k8s/secrets.yml
kubectl apply -f k8s/configmap.yml
kubectl apply -f k8s/infra/

# Ждём инфраструктуру
kubectl wait --for=condition=ready pod -l app=postgres -n realty --timeout=180s
kubectl wait --for=condition=ready pod -l app=redis -n realty --timeout=60s

kubectl apply -f k8s/services/

### 8. Применяй миграции БД

bash
kubectl port-forward -n realty service/postgres 5432:5432 &
sleep 2
migrate -path ./migrations \
  -database "postgres://usr:tr134sdfWE@localhost:5432/lets_goto_it?sslmode=disable" \
  up

### 9. Проверь что всё работает

bash
kubectl get pods -n realty
curl http://192.168.49.2:30080/health
# Ожидаемый ответ: {"status":"ok"}

## Доступ к приложению

### Локально (внутри VM)

# Фронтенд
http://192.168.49.2:30081

# API
http://192.168.49.2:30080

### Из интернета (Cloudflare Tunnel)

# Пробрось порты
kubectl port-forward -n realty service/frontend 8081:80 --address=0.0.0.0 &

# Запусти туннель
cloudflared tunnel --url http://localhost:8081
# Получишь URL вида: https://xxx.trycloudflare.com

## CI/CD

При каждом пуше в ветку `main`:
1. Собираются Docker образы всех сервисов
2. Образы пушатся в локальный registry (`192.168.49.1:5000`)
3. Применяются k8s манифесты
4. Все деплойменты перезапускаются с новым кодом

Требование: self-hosted runner должен быть запущен на Ubuntu VM.

# Проверь статус runner
sudo systemctl status actions.runner.shmatkovapple-ctrl-realty.nik-VirtualBox.service

# Запусти если не работает
cd ~/actions-runner && sudo ./svc.sh start

## Запуск после перезагрузки VM

ansible-playbook ~/realty/ansible/playbook.yml

Скрипт автоматически:
- Запускает Docker Registry
- Запускает Minikube
- Деплоит всю инфраструктуру и сервисы

### Структура проекта

realty/
├── api/                    # Proto файлы и сгенерированный код
├── frontend/               # React приложение
│   ├── Dockerfile
│   └── nginx.conf          # Проксирует /api/ на api-gateway
├── services/               # Go микросервисы
│   ├── user-service/
│   ├── listing-service/
│   ├── deal-service/
│   ├── search-service/
│   ├── notification-service/
│   └── api-gateway/
├── k8s/                    # Kubernetes манифесты
│   ├── namespace.yml
│   ├── secrets.yml
│   ├── configmap.yml
│   ├── infra/              # PostgreSQL, Redis, Kafka, Elasticsearch, MinIO
│   └── services/           # Все микросервисы + фронтенд
├── migrations/             # SQL миграции
├── docker-compose.yml      # Локальная разработка без k8s
└── Makefile                # Команды для разработки

## Полезные команды

bash
# Статус подов
kubectl get pods -n realty

# Логи сервиса
kubectl logs -n realty -l app=user-service --tail=50

# Применить миграции
migrate -path ./migrations \
  -database "postgres://usr:tr134sdfWE@localhost:5432/lets_goto_it?sslmode=disable" up

# Пересобрать и задеплоить один сервис
docker build -t 192.168.49.1:5000/realty/user-service:latest -f services/user-service/Dockerfile .
docker push 192.168.49.1:5000/realty/user-service:latest
kubectl rollout restart deployment/user-service -n realty

# Остановить всё
minikube stop
