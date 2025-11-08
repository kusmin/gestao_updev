# Arquitetura Backend – Plataforma de Gestão Local

## Visão Geral
- Linguagem: Go 1.25+.
- Framework HTTP: Gin (preferência) ou Fiber (alternativo).
- Estilo: RESTful stateless com autenticação bearer (JWT).
- Persistência: PostgreSQL 14+ com abordagem multi-tenant por coluna (`tenant_id`) em todas as entidades compartilhadas.

## Stack Técnica
- **Camada HTTP**: Gin, middlewares para logging, auth, rate limit leve.
- **Validação/DTO**: Pacote `go-playground/validator` com structs dedicados a requests/responses.
- **Service Layer**: regras de negócio e orquestração transacional; isolada em interfaces para facilitar testes.
- **Repositórios**: SQL gerado via SQLC (preferido) ou GORM para cenários dinâmicos.
- **Infra**: `pkg/db` (pool + migrations), `pkg/config` (Viper/env), `pkg/log` (Zap/zerolog), `pkg/auth`.

## Organização do Código
```
backend/
  cmd/api/main.go        # bootstrap HTTP
  internal/
    http/                # controllers + routes
    service/             # casos de uso (agenda, vendas, etc.)
    repository/          # interfaces + impl SQLC
    domain/              # entidades e valores
    auth/                # JWT, RBAC
    middleware/          # middlewares Gin
  pkg/                   # utilidades reutilizáveis (config, logger)
  migrations/            # arquivos SQL (golang-migrate)
```

## Fluxo de Requisição (alto nível)
1. **HTTP Handler** valida entrada, injeta `tenant_id` e `user_id` a partir do token.
2. **Service** aplica regras (ex.: checar permissões, consistência de agenda).
3. **Repository** executa queries parametrizadas sempre filtrando por `tenant_id`.
4. **Resposta** padronizada (`data`, `error`, `meta`) retornada ao cliente.

## Estratégia Multi-Tenant
- Cada requisição inclui `X-Tenant-ID` (após login) e o token carrega `tenant_id`.
- Middlewares garantem:
  - `tenant_id` obrigatório salvo para rotas públicas (signup/login).
  - Queries sempre adicionam `WHERE tenant_id = $1`.
- Tabelas globais (planos, configurações) ficam em schema separado (`public`).
- Futuro: migrar para `schema por tenant` se necessário para isolamento maior.

## Contratos REST (MVP)
- `/auth/signup` (POST): cria empresa + usuário admin.
- `/auth/login` (POST): retorna `access_token` + `refresh_token`.
- `/companies/me` (GET/PUT): dados da empresa.
- `/users` (CRUD): gestão de colaboradores com níveis (admin, gerente, operador).
- `/clients`, `/services`, `/professionals`: CRUDs básicos.
- `/agenda/bookings`: criar/listar agendamentos por profissional/data.
- `/sales/orders`: registrar vendas e formas de pagamento (cash, Pix manual).
- Respostas padronizadas:
  ```json
  {
    "data": {...},
    "meta": {"pagination": {...}},
    "error": null
  }
  ```

## Observabilidade & Operação
- Logging estruturado (JSON) com correlação (`request_id`, `tenant_id`).
- Métricas Prometheus via `/metrics`: latência, throughput, erros, uso DB.
- Health checks (`/healthz` liveness / `/readyz` readiness).
- Feature flags simples via config (ex.: habilitar agenda avançada).

## Configuração
- Arquivo `.env` + variáveis de ambiente (com suporte Viper).
- Principais variáveis:
  - `APP_ENV`, `HTTP_PORT`.
  - `DATABASE_URL`.
  - `JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET`, `JWT_ACCESS_TTL`, `JWT_REFRESH_TTL`.
  - `RATE_LIMIT_RPS`.
- Docker Compose para levantar API + PostgreSQL localmente.

## Segurança
- Hash de senha com bcrypt (cost 12) ou argon2id.
- Refresh tokens armazenados com hash + expiração; endpoint `/auth/refresh`.
- Rate limiting por IP e por `tenant_id` em rotas sensíveis.
- Validação de payload rigorosa + sanitização de inputs.
- Politica RBAC:
  - `admin`: tudo.
  - `manager`: agenda, vendas, relatórios, gestão básica de usuários.
  - `operator`: agenda/vendas próprias.

## Próximos Passos
1. Definir o framework HTTP (fixar Gin) e gerar boilerplate `cmd/api`.
2. Documentar DTOs iniciais (OpenAPI/Swagger) e publicar em `docs/api.yaml`.
3. Modelar migrations base (empresas, usuários, clientes, serviços, agenda, vendas).
4. Escrever guia de contribuição e padrões de commit/tests para o backend.
