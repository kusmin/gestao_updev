# Workflows do GitHub Actions

Este documento descreve os workflows de Integração Contínua (CI) e Entrega Contínua (CD) configurados no projeto `gestao_updev` usando GitHub Actions. Esses workflows garantem a qualidade do código, a segurança, a consistência da documentação e a automação de releases e builds.

## 1. `api-contract.yml` - Testes de Contrato da API

*   **Nome:** `API Contract Tests`
*   **Gatilhos:**
    *   `push` para alterações nos diretórios `backend/`, `docs/api.yaml`, `tests/dredd/`, `Makefile` ou no próprio arquivo de workflow.
    *   `pull_request` para as mesmas alterações.
*   **Objetivo:** Executar testes de contrato da API usando Dredd para garantir que as alterações no backend ou na especificação da API não quebrem o contrato estabelecido.
*   **Jobs Principais:**
    *   **`dredd`**: Configura o ambiente (Go, Node.js, PostgreSQL), instala dependências e executa os testes de contrato Dredd.

## 2. `api-spec.yml` - Qualidade da Especificação da API e Publicação

*   **Nome:** `API Spec Quality`
*   **Gatilhos:**
    *   `push` para alterações em `docs/api.yaml` ou no próprio arquivo de workflow.
    *   `pull_request` para as mesmas alterações.
*   **Objetivo:** Garantir a qualidade da especificação da API (`api.yaml`) usando Spectral, gerar a documentação HTML com Redoc e publicá-la no GitHub Pages.
*   **Jobs Principais:**
    *   **`spectral-lint`**: Linta a especificação da API com Spectral e gera o HTML da documentação com Redoc, fazendo upload como artefato.
    *   **`deploy-pages`**: Publica a documentação da API no GitHub Pages (apenas em `push` para a branch `main` e após o sucesso de `spectral-lint`).

## 3. `api-sync.yml` - Sincronização da Especificação da API

*   **Nome:** `API Spec Sync Check`
*   **Gatilhos:**
    *   `push` para a branch `main` com alterações em `backend/**` ou `docs/api.yaml`.
    *   `pull_request` para a branch `main` com alterações em `backend/**` ou `docs/api.yaml`.
*   **Objetivo:** Verificar se a especificação da API (`docs/api.yaml`) está sincronizada com a implementação do código Go, gerando a especificação a partir do código e comparando-a.
*   **Jobs Principais:**
    *   **`check-sync`**: Configura o ambiente (Go, `swag` CLI), gera a especificação Swagger a partir do código Go e compara com `docs/api.yaml`, falhando se houver diferenças.

## 4. `backend-ci.yml` - Integração Contínua do Backend

*   **Nome:** `Backend CI`
*   **Gatilhos:**
    *   `push` para a branch `main` com alterações em `backend/**`.
    *   `pull_request` para a branch `main` com alterações em `backend/**`.
*   **Objetivo:** Executar o pipeline de Integração Contínua (CI) para o backend, incluindo linting, testes, build e cobertura de código.
*   **Jobs Principais:**
    *   **`lint`**: Executa `golangci-lint` para verificar a qualidade do código Go.
    *   **`test`**: Executa os testes do backend usando Docker Compose.
    *   **`build`**: Compila a aplicação Go do backend.
    *   **`coverage`**: Gera um relatório de cobertura de testes e o envia para o Codecov.

## 5. `codeql.yml` - Análise de Segurança CodeQL

*   **Nome:** `CodeQL Security Analysis`
*   **Gatilhos:**
    *   `push` para a branch `main` com alterações em `backend/**`, `frontend/**` ou no próprio arquivo de workflow.
    *   `pull_request` para a branch `main` com alterações em `backend/**` ou `frontend/**`.
    *   `schedule`: Executa semanalmente às segundas-feiras às 02:30 UTC.
*   **Objetivo:** Realizar análises de segurança estática no código Go e JavaScript/TypeScript usando GitHub CodeQL para identificar vulnerabilidades.
*   **Jobs Principais:**
    *   **`go-analysis`**: Inicializa e executa a análise CodeQL para o código Go.
    *   **`javascript-analysis`**: Inicializa e executa a análise CodeQL para o código JavaScript/TypeScript.

## 6. `coverage.yml` - Cobertura Consolidada do Monorepo

*   **Nome:** `Monorepo Coverage`
*   **Gatilhos:**
    *   `push` para a branch `main` com alterações em `backend/**`, `frontend/**`, `backoffice/**`, `Makefile`, `docker-compose.test.yml`, `package.json`, `package-lock.json` ou no próprio arquivo de workflow.
    *   `pull_request` para a branch `main` com as mesmas alterações.
*   **Objetivo:** Executar testes de cobertura consolidados para todo o monorepo (backend, frontend e backoffice) e enviar os resultados para o Codecov.
*   **Jobs Principais:**
    *   **`coverage`**: Configura o ambiente (Node.js, Go), instala dependências, executa `make coverage` (que orquestra os testes de cobertura em cada subprojeto) e faz upload dos relatórios para o Codecov.

## 7. `e2e.yml` - Testes End-to-End (E2E) com Playwright

*   **Nome:** `Playwright E2E`
*   **Gatilhos:**
    *   `push` para a branch `main` com alterações em `frontend/**`, `backoffice/**`, `tests/e2e/**` ou no próprio arquivo de workflow.
    *   `pull_request` para a branch `main` com as mesmas alterações.
*   **Objetivo:** Executar testes End-to-End (E2E) usando Playwright para verificar a funcionalidade completa da aplicação, simulando interações do usuário.
*   **Jobs Principais:**
    *   **`playwright`**: Configura o ambiente (Node.js, Playwright), instala dependências e navegadores, executa os testes E2E e faz upload do relatório.

## 8. `frontend-build.yml` - Build do Frontend

*   **Nome:** `Frontend Build`
*   **Gatilhos:**
    *   `push` para alterações em `frontend/**`, `docs/api.yaml` ou `Makefile`.
    *   `pull_request` para as mesmas alterações.
*   **Objetivo:** Construir a aplicação frontend, incluindo a geração automática dos tipos da API.
*   **Jobs Principais:**
    *   **`build`**: Configura o ambiente (Node.js), instala dependências e executa o comando `npm run build` no diretório `frontend`.

## 9. `frontend-ci.yml` - Integração Contínua do Frontend

*   **Nome:** `Frontend CI`
*   **Gatilhos:**
    *   `push` para a branch `main` com alterações em `frontend/**`, `docs/api.yaml` ou `Makefile`.
    *   `pull_request` para a branch `main` com as mesmas alterações.
*   **Objetivo:** Executar o pipeline de Integração Contínua (CI) para o frontend, incluindo linting, testes e build.
*   **Jobs Principais:**
    *   **`lint`**: Executa ESLint para verificar a qualidade do código frontend.
    *   **`test`**: Executa os testes do frontend com cobertura e envia o relatório para o Codecov.
    *   **`build`**: Compila a aplicação frontend.

## 10. `pre-commit.yml` - Verificações de Pre-commit

*   **Nome:** `Pre-commit Checks`
*   **Gatilhos:**
    *   `push` para a branch `main`.
    *   `pull_request` para a branch `main`.
*   **Objetivo:** Executar os hooks de pre-commit configurados no projeto para garantir a qualidade e consistência do código antes que as alterações sejam commitadas ou mescladas na branch principal.
*   **Jobs Principais:**
    *   **`pre-commit`**: Configura o ambiente (Python, Go, Node.js, `golangci-lint`, `pre-commit`), instala dependências e executa os checks de pre-commit em todos os arquivos.

## 11. `publish-docker.yml` - Publicação da Imagem Docker do Backend

*   **Nome:** `Publish Backend Docker Image`
*   **Gatilhos:**
    *   `push` para a branch `main` com alterações em `backend/**`, no próprio arquivo de workflow ou em `backend/Dockerfile`.
*   **Objetivo:** Construir e publicar a imagem Docker do backend no GitHub Container Registry (GHCR).
*   **Jobs Principais:**
    *   **`build-and-push`**: Configura o ambiente Docker, faz login no GHCR, extrai metadados Docker e, em seguida, constrói e envia a imagem Docker do backend para o registro.

## 12. `release-please.yml` - Automação de Release

*   **Nome:** `Release Automation`
*   **Gatilhos:**
    *   `push` para a branch `main`.
*   **Objetivo:** Automatizar o processo de release do projeto usando `release-please`, incluindo a criação de pull requests para novas versões, atualização de changelogs e criação de tags de release.
*   **Jobs Principais:**
    *   **`release-please`**: Executa a ação `googleapis/release-please-action` para gerenciar o processo de release.
