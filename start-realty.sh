#!/bin/bash
set -e

echo "=== [1/6] Запускаем Docker Registry ==="
docker start registry 2>/dev/null || \
  docker run -d -p 5000:5000 --restart=always --name registry registry:2
echo "Registry OK"

echo "=== [2/6] Запускаем Minikube ==="
minikube start \
  --driver=docker \
  --cpus=3 \
  --memory=7g \
  --insecure-registry="192.168.49.1:5000"
echo "Minikube OK"

echo "=== [3/6] Деплоим конфиги ==="
cd ~/realty
kubectl apply -f k8s/namespace.yml
kubectl apply -f k8s/secrets.yml
kubectl apply -f k8s/configmap.yml

echo "=== [4/6] Деплоим инфраструктуру ==="
kubectl apply -f k8s/infra/
kubectl wait --for=condition=ready pod \
  -l app=postgres -n realty --timeout=180s
kubectl wait --for=condition=ready pod \
  -l app=redis -n realty --timeout=60s

echo "=== [5/6] Деплоим сервисы ==="
kubectl apply -f k8s/services/

echo "=== [6/6] Пробрасываем порты ==="
# Убиваем старые port-forward если есть
pkill -f "kubectl port-forward" 2>/dev/null || true
sleep 2

kubectl port-forward -n realty service/frontend 8081:80 --address=0.0.0.0 &
kubectl port-forward -n realty service/api-gateway 8080:8080 --address=0.0.0.0 &
kubectl port-forward -n monitoring service/monitoring-grafana 3000:80 --address=0.0.0.0 &

sleep 3
echo ""
echo "=== Готово! ==="
kubectl get pods -n realty
echo ""
echo "Фронтенд: http://192.168.49.2:30081"
echo "API:      http://192.168.49.2:30080"
echo "Grafana:  http://localhost:3000 (admin/admin123)"
echo ""
echo "Для доступа из интернета:"
echo "cloudflared tunnel --url http://localhost:8081"
