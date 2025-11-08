# Operação & DevOps – Plataforma de Gestão Local

## Objetivo
Definir práticas para desenvolvimento local, pipelines CI/CD e operação em ambientes de staging/produção com custo controlado.

## Desenvolvimento Local
- **Pré-requisitos**: Go 1.25+, Node 20+, Docker/Docker Compose, Make.
- **Setup manual**:
  1. `cp .env.example .env` na raiz e exporte variáveis conforme necessidade.
  2. Banco: `docker compose up postgres -d` ou `docker compose up` para subir tudo.
  3. Backend: `make backend-run` (ou `air` para hot reload).
  4. Frontend: `cd frontend && npm run dev` (gera tipos automaticamente).
- **Setup full via Docker**:
  1. Copiar `.env.example -> .env`.
  2. `docker compose up --build`.
  3. Frontend acessível em `http://localhost:4173`, API em `http://localhost:8080/v1/healthz`.
- **Ferramentas**:
  - `air` para hot reload Go.
  - `golang-migrate` via `make migrate-up/down`.
  - `task` opcional para orquestrar scripts multi-stack.

## Pipelines CI/CD
- **CI (GitHub Actions)**:
  - Jobs paralelos: `backend-lint-test`, `frontend-lint-test`.
  - Passos backend: `go test ./...`, `golangci-lint run`, build binário.
  - Passos frontend: `npm ci`, `npm run lint`, `npm run test -- --runInBand`.
  - Upload coverage para Codecov.
- **CD**:
  - Deploy automático para staging em cada merge na `main`.
  - Produção via tag (`v*`) com aprovação manual (environments GitHub).
  - Artefatos Docker publicados no GHCR (`ghcr.io/org/gestao-api`).

## Infraestrutura
- **Ambiente**: Railway/Render/Fly.io (escolha baseada em custo e latência).
- **Componentes**:
  - API Go em container (256-512MB RAM).
  - PostgreSQL gerenciado (mín. 1 vCPU, 1GB RAM, storage 20GB).
  - Redis opcional para cache/session (futuro).
  - Bucket S3 compatível (DigitalOcean Spaces) para anexos (nota fiscal, imagens).
- **Networking**:
  - HTTPS obrigatório (Let’s Encrypt integrado na plataforma).
  - WAF básico/limitação de IP via plataforma quando disponível.

## Observabilidade
- Logs centralizados (Grafana Loki ou Logtail) com correlação `request_id`.
- Métricas Prometheus/OTEL exportadas para Grafana Cloud free tier.
- Alertas:
  - Latência p95 > 500ms.
  - Taxa de erro 5xx > 2% em 5 min.
  - Espaço em disco > 80%.
- SLO inicial: disponibilidade 99.5% mensal.

## Gestão de Configurações
- Segredos via ambiente seguro da plataforma (Railway variables, Render envs).
- Configs versionadas em `.env.example` e `config/*.yaml`.
- Rotacionar chaves JWT e DB passwords a cada 90 dias.

## Backup & Recuperação
- Snapshots automáticos do banco 1x/dia com retenção 7 dias.
- Exportação semanal para armazenamento externo (S3) criptografado.
- Testar restore trimestralmente em ambiente isolado.

## Segurança Operacional
- Acesso aos painéis restrito via SSO (Google/Microsoft) ou 2FA.
- Logs de auditoria para ações admin (deploys, rotação de chaves).
- Dependabot habilitado para backend/frontend.

## Plano de Incidentes
- Canal dedicado (Slack/Teams) com rota de escalonamento.
- Playbook mínimo:
  1. Detectar e classificar (P1/P2/P3).
  2. Mitigar (rollback, escala horizontal, feature flag).
  3. Comunicar status aos stakeholders.
  4. Post-mortem em até 48h (conter ações corretivas).

## Próximos Passos
1. Criar perfil `docker-compose.dev.yml` com hot reload e montar volume do código.
2. Prototipar pipeline GitHub Actions com jobs descritos (lint/test/build).
3. Especificar variáveis e segredos necessários por ambiente (local/staging/prod).
4. Escrever guia de release (checklist, rollback) e anexar ao `docs/`.
