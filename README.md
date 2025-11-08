# Gestão UpDev

Plataforma de gestão para negócios locais (barbearias, lojas de roupa e comércios de bairro), composta por um backend Go, frontend web em React + Vite (a ser implementado) e documentação centralizada em `docs/`.

## Estrutura do Projeto
- `backend/`: futuro serviço Go (APIs REST multi-tenant).
- `frontend/`: aplicação web em React + Vite (scripts já configurados para gerar tipos a partir do OpenAPI).
- `backoffice/`: aplicação web em React + Vite para gestão administrativa.
- `docs/`: documentação funcional, técnica e operacional.
- `docker/`: espaço reservado para composições e imagens.

## Documentação-chave
- Visão geral e arquitetura: `docs/README.md`, `docs/arquitetura-backend.md`.
- Modelo de dados: `docs/modelo-dados.md`.
- Fluxos de uso: `docs/fluxos-uso.md`.
- Padrões de código: `docs/padroes-codigo.md`.
- Operação & DevOps: `docs/operacao-devops.md`.
- Guia Docker/Compose: `docs/docker.md`.
- API (referência + spec + changelog):
  - Resumo textual: `docs/api-reference.md`.
  - OpenAPI 3.1: `docs/api.yaml`.
  - Guia de uso/CI/pipelines: `docs/api-usage.md`.
  - Histórico/versionamento: `docs/api-changelog.md`.

## API pública (GitHub Pages)
A pipeline `API Spec Quality` publica automaticamente o HTML do Redoc na branch `main`. A URL pública padrão é:

```
https://kusmin.github.io/gestao_updev/
```

(A página é atualizada sempre que `docs/api.yaml` muda e passa na validação do CI.)

## Comandos Úteis
- `make api-lint`: valida o OpenAPI com Spectral.
- `make api-preview`: abre o Redoc localmente.
- `make api-types`: gera definições TypeScript em `frontend/src/types/api.d.ts`.
- `npm run build` (frontend): executa `generate:api-types` automaticamente antes do build real (placeholder por enquanto).
- `docker compose up --build`: sobe PostgreSQL + API + frontend (ver `docs/docker.md`).

## CI/CD
- `API Spec Quality`: lint + publicação do spec.
- `Frontend Build`: instala dependências e roda o build (com geração de tipos) a cada alteração relevante.
- `API Contract Tests`: em breve executará Dredd contra o backend (estrutura inicial já criada).

## Próximos Passos
- Implementar o backend Go conforme docs.
- Criar o frontend React + Vite consumindo os endpoints definidos.
- Adicionar testes de contrato (Dredd/Postman) para garantir conformidade entre implementação e OpenAPI (ver `docs/tests-contrato.md`).
