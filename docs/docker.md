# Guia Docker

Este documento explica como construir e executar o projeto usando Docker e Docker Compose.

## Pré-requisitos
- Docker 24+
- Docker Compose v2

## Estrutura
- `backend/Dockerfile`: build multi-stage para a API Go.
- `frontend/Dockerfile`: build da SPA React + Vite e entrega via Nginx.
- `docker-compose.yml`: orquestra PostgreSQL, backend e frontend.
- `.env.example`: variáveis consumidas pelo Compose (copie para `.env`).

## Passo a passo
1. Copie o arquivo de variáveis:
   ```bash
   cp .env.example .env
   ```
2. Ajuste valores conforme necessidade (senhas, portas, `VITE_API_BASE_URL`).
3. Construa e suba os serviços:
   ```bash
   docker compose up --build
   ```
4. Serviços disponíveis:
   - API: http://localhost:${HTTP_PORT:-8080}/v1/healthz
   - Frontend: http://localhost:${FRONTEND_PORT:-4173}
   - PostgreSQL: localhost:${POSTGRES_PORT:-5432}

## Variáveis Principais
- `APP_ENV`, `LOG_LEVEL`: repassadas ao backend.
- `DATABASE_URL`: string usada pela API (host `postgres` dentro da rede Compose).
- `VITE_API_BASE_URL`: endpoint usado na build do frontend para apontar para o backend (já default `http://backend:8080/v1`).

## Comandos Úteis
- Reconstruir apenas um serviço:
  ```bash
  docker compose build backend
  ```
- Ver logs:
  ```bash
  docker compose logs -f backend
  ```
- Encerrar tudo:
  ```bash
  docker compose down
  ```

## Próximos Passos
- Adicionar `docker compose -f docker-compose.yml -f docker-compose.dev.yml` para perfis locais (hot reload).
- Publicar imagens no Registry (GitHub Container Registry) para facilitar deploy.
