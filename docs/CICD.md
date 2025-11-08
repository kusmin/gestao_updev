# Estratégia de CI/CD

Este documento descreve a estratégia de Integração Contínua (CI) e Entrega Contínua (CD) para o projeto.

## Visão Geral

- **Ferramenta:** GitHub Actions.
- **Ambientes:**
  - `development`: Branches de feature.
  - `staging`: Ambiente de homologação, deploy automático a partir da branch `main`.
  - `production`: Ambiente de produção, deploy manual ou automático a partir de tags/releases.

## Pipeline de CI (Integração Contínua)

A pipeline de CI será executada em cada `push` para a branch `main` e em cada `pull request` aberto.

### Backend (Go)

1.  **Trigger:** `push` em `main`, `pull_request` para `main`.
2.  **Jobs:**
    - `lint`: Executa `golangci-lint` para análise estática.
    - `test`: Executa os testes unitários e de integração (`go test ./...`).
    - `build`: Compila o binário da aplicação (`go build`).
    - `dockerize`: (Opcional) Constroi e publica a imagem Docker no Docker Hub ou GitHub Container Registry.

### Frontend (React + Vite)

1.  **Trigger:** `push` em `main`, `pull_request` para `main`.
2.  **Jobs:**
    - `lint`: Executa `eslint`/`prettier`.
    - `test`: Executa os testes unitários/componente com `vitest` + Testing Library (futuro).
    - `build`: Gera o bundle com `vite build` (pré-step `generate:api-types`).
    - `dockerize`: (Opcional) Constroi e publica a imagem Docker.

## Pipeline de CD (Entrega Contínua)

A pipeline de CD será responsável por automatizar o deploy da aplicação.

1.  **Trigger:** `push` na branch `main`.
2.  **Steps:**
    - **Deploy para Staging:** A versão mais recente da branch `main` será automaticamente enviada para o ambiente de `staging` na PaaS escolhida (Railway, Render, etc.).
    - **Deploy para Produção:** O deploy para produção será um passo manual, acionado pela criação de uma `tag` ou `release` no GitHub. Isso garante que apenas versões estáveis e validadas cheguem aos usuários finais.

## Próximos Passos

1.  **Configurar Workflows Básicos:**
    - Criar o arquivo `.github/workflows/backend-ci.yml` com os jobs de `lint`, `test` e `build` para o backend.
    - Criar o arquivo `.github/workflows/frontend-ci.yml` com os jobs de `lint`, `test` e `build` para o frontend.
2.  **Gerenciar Segredos:**
    - Configurar os `secrets` no repositório do GitHub para armazenar tokens de acesso (ex: `DOCKER_HUB_TOKEN`, `RAILWAY_API_TOKEN`).
3.  **Definir Estratégia de Deploy:**
    - Detalhar a estratégia de deploy (ex: blue-green) para evitar downtime durante as atualizações em produção.
4.  **Criar Workflow de CD:**
    - Implementar o workflow de deploy para o ambiente de `staging`.
    - Implementar o workflow de deploy para o ambiente de `production`, com o gatilho manual.

## Workflows Ativos
- `API Spec Quality`: lint do OpenAPI e publicação no GitHub Pages.
- `Frontend Build`: garante build do app React + geração de tipos.
- `API Contract Tests`: executa Dredd contra o backend (atualmente apenas health check liberado).
