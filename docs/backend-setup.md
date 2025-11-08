# Guia de Configuração do Backend

Este documento complementa a visão geral em `docs/arquitetura-backend.md` e o contrato em `docs/api-reference.md`, descrevendo como inicializar o serviço Go, executar as migrações e validar os endpoints principais.

## Visão Geral Rápida
- **Stack:** Go 1.25, Gin, GORM, PostgreSQL 14+, JWT.
- **Camadas:** `internal/http` (handlers + middlewares), `internal/service` (regras), `internal/domain` (modelos), `pkg/database` (conexão), `migrations` (SQL versionado).
- **Multi-tenant:** coluna `tenant_id` em todas as tabelas compartilhadas. O header `X-Tenant-ID` precisa acompanhar cada requisição autenticada.

## Modelos e Migrações
Os schemas SQL residem em `backend/migrations` e seguem a sequência descrita em `docs/modelo-dados.md`:

| Arquivo | Conteúdo |
| --- | --- |
| `0001_init` | empresas, usuários, clientes, serviços, produtos, audit logs. |
| `0002_agenda` | profissionais, regras de disponibilidade, bookings. |
| `0003_sales` | pedidos, itens, pagamentos. |
| `0004_inventory` | movimentações de estoque + trigger de ajuste. |

### Executando com golang-migrate
```bash
cd backend
DATABASE_URL="postgres://postgres:postgres@localhost:5432/gestao_updev?sslmode=disable"
migrate -path migrations -database "$DATABASE_URL" up
```

> Consulte `docs/operacao-devops.md` para detalhes de pipeline e opções Docker.

## Variáveis de Ambiente
Todas as chaves lidas em `internal/config/config.go`:

| Variável | Descrição | Default |
| --- | --- | --- |
| `APP_ENV` | Ambiente (`development`, `staging`, `production`). | `development` |
| `HTTP_PORT` | Porta HTTP. | `8080` |
| `DATABASE_URL` | DSN PostgreSQL. | `postgres://postgres:postgres@localhost:5432/gestao_updev?sslmode=disable` |
| `JWT_ACCESS_SECRET` / `JWT_REFRESH_SECRET` | Segredos HMAC. | `dev-access-secret` / `dev-refresh-secret` |
| `JWT_ACCESS_TTL` / `JWT_REFRESH_TTL` | Expirações (`15m`, `720h`, ...). | `15m` / `720h` |
| `TENANT_HEADER` | Nome do header multi-tenant. | `X-Tenant-ID` |
| `LOG_LEVEL` | `debug`, `info`, `warn`, `error`. | `info` |

Para desenvolvimento local crie um `.env` (ou exporte no shell) antes de usar `make run`.

## Passo a Passo para Subir o Backend

1. **Instalar dependências:**
   ```bash
   cd backend
   go mod download
   ```
2. **Provisionar PostgreSQL:**
   - Rode `docker compose up db` (ver `docs/docker.md`) ou use uma instância existente.
3. **Aplicar migrações:**
   - Execute `migrate ... up` conforme seção anterior.
4. **Rodar a API:**
   ```bash
   make run
   # ou
   APP_ENV=development HTTP_PORT=8080 go run ./cmd/api
   ```
5. **Popular dados de referência (opcional para dev):**
   ```bash
   make -C backend seed
   ```
6. **Rodar testes com cobertura (opcional):**
   ```bash
   make -C backend test
   # gera backend/coverage.out e imprime o resumo via go tool cover -func
   ```

7. **Validar endpoints:**
   - Health check: `curl http://localhost:8080/v1/healthz`.
   - Swagger: `http://localhost:8080/swagger/index.html`.
   - Fluxo típico:
     1. `POST /v1/auth/signup` → cria empresa e admin.
     2. `POST /v1/auth/login` → obtém tokens.
     3. Envia `Authorization: Bearer <token>` + `X-Tenant-ID` nas rotas protegidas.

> Sempre que alterar handlers, rode `go test ./...` e, se necessário, `make swagger` (que executa `swag init`) para atualizar `backend/docs/swagger.*`.

## Endpoints Disponíveis
Todos os handlers listados em `docs/api-reference.md` agora estão implementados e expostos via `/v1`:

- **Autenticação:** signup, login, refresh.
- **Empresa/Usuários:** `/companies/me`, `/users`.
- **Clientes & Profissionais:** `/clients`, `/professionals`.
- **Catálogo:** `/services`, `/products`.
- **Agenda:** `/bookings`.
- **Estoque:** `/inventory/movements`.
- **Vendas & Pagamentos:** `/sales/orders`, `/payments`.
- **Dashboard:** `/dashboard/daily`.

As respostas seguem o padrão `internal/http/response.APIResponse`, garantindo campos `data`, `meta` (quando aplicável) e `error`.

## Boas Práticas
- Utilize o middleware de tenant (`middleware.TenantEnforcer`) e o middleware de auth (`middleware.Auth`) para qualquer novo grupo de rotas.
- Adicione validações nos DTOs via tags `binding:"..."` ao criar novos handlers.
- Para novas entidades, crie primeiro o modelo em `internal/domain`, adicione migrações SQL, exponha métodos no `service` e, por fim, escreva o handler correspondente.

Referências adicionais:
- `docs/api-usage.md` – exemplos de chamadas e CI.
- `docs/tests-contrato.md` – plano para Dredd/Postman.
- `docs/padroes-codigo.md` – convenções de Go e REST no projeto.
