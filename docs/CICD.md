# Estratégia e Estado Atual de CI/CD

Este documento descreve a configuração de Integração Contínua (CI) e os planos para Entrega Contínua (CD) do projeto, utilizando GitHub Actions.

## Visão Geral

O projeto utiliza múltiplos workflows de GitHub Actions para garantir a qualidade e a consistência do código-fonte, da documentação da API e dos builds. Os workflows são separados por responsabilidade: backend, frontend e API.

---

## Workflows de CI Ativos

A seguir, uma descrição detalhada de cada workflow configurado no diretório `.github/workflows/`.

### 1. `backend-ci.yml`

- **Propósito:** Executa a integração contínua para o backend (Go).
- **Gatilhos:** Acionado em `push` ou `pull_request` para a branch `main` que modifiquem arquivos no diretório `backend/`.
- **Jobs:**
    - `lint`: Análise estática do código Go.
    - `test`: Testes de integração com Docker Compose.
    - `build`: Compilação do binário da aplicação.
    - `coverage`: Geração e envio do relatório de cobertura de testes unitários.

### 2. `frontend-ci.yml`

- **Propósito:** Responsável pela integração contínua do frontend.
- **Gatilhos:** Acionado em `push` ou `pull_request` para a branch `main` que modifiquem `frontend/`, `docs/api.yaml` ou o `Makefile`.
- **Jobs:**
    - `lint`: Executa o linter.
    - `test`: Executa os testes, gera e envia o relatório de cobertura.
    - `build`: Compila a aplicação.

### 3. `api-spec.yml`

- **Propósito:** Garante a qualidade da especificação OpenAPI e a publica como uma página de documentação.
- **Gatilhos:** Acionado por alterações no arquivo `docs/api.yaml`.
- **Jobs:**
    - `spectral-lint`: Valida o arquivo `docs/api.yaml` e gera um HTML da documentação.
    - `deploy-pages`: Faz o deploy da documentação da API no GitHub Pages.

### 4. `api-sync.yml`

- **Propósito:** Verifica se a especificação OpenAPI está sincronizada com o código-fonte do backend.
- **Gatilhos:** Acionado por alterações no `backend/` ou em `docs/api.yaml`.
- **Jobs:**
    - `check-sync`: Compara a especificação gerada pelo código com a documentação oficial.

### 5. `api-contract.yml`

- **Propósito:** Executa testes de contrato para validar se a implementação do backend corresponde à sua especificação.
- **Gatilhos:** Acionado por alterações no `backend/`, `docs/api.yaml`, ou nos próprios testes.
- **Jobs:**
    - `dredd`: Executa a ferramenta Dredd para testar o comportamento da API.

### 6. `codeql.yml`

- **Propósito:** Executa análise de segurança estática (SAST) no código-fonte.
- **Gatilhos:** Acionado em `push`/`pull_request` para `main` e também semanalmente.
- **Jobs:**
    - `go-analysis`: Analisa o código Go.
    - `javascript-analysis`: Analisa o código TypeScript/JavaScript.

### 7. `publish-docker.yml`

- **Propósito:** Constrói e publica a imagem Docker do backend no GitHub Container Registry (GHCR).
- **Gatilhos:** Acionado em `push` para a branch `main` que modifiquem arquivos no diretório `backend/`, o `backend/Dockerfile` ou o próprio workflow.
- **Jobs:**
    - `build-and-push`: Realiza o login no GHCR, extrai metadados da imagem e constrói/publica a imagem Docker do backend com tags `latest` e SHA do commit.

### 8. `pre-commit.yml`

- **Propósito:** Executa uma série de verificações de formatação e linting no código-fonte usando o framework `pre-commit`.
- **Gatilhos:** Acionado em `push` ou `pull_request` para a branch `main`.
- **Jobs:**
    - `pre-commit`:
        - **Configuração do Ambiente:** Instala todas as ferramentas necessárias para os hooks, como Go, Node.js, `golangci-lint` e dependências `npm`.
        - **Execução:** Roda `pre-commit run --all-files` para garantir que todo o código no repositório esteja em conformidade com os padrões definidos no arquivo `.pre-commit-config.yaml`. Isso ajuda a manter a consistência e a qualidade do código de forma automatizada.

---

## Melhorias Implementadas e Próximos Passos

Esta seção descreve as melhorias já realizadas e os próximos passos para otimizar os pipelines de CI/CD.

### Implementado

- **Cache de Dependências:** Adicionado cache para os módulos Go e `node_modules`, acelerando a execução dos jobs.
- **Análise de Segurança (SAST):** Implementado um workflow com GitHub CodeQL para escanear o código em busca de vulnerabilidades.
- **Relatórios de Cobertura de Testes (Code Coverage):** Os workflows de backend e frontend agora geram relatórios de cobertura de testes e os enviam para o Codecov, permitindo a análise da qualidade dos testes.
- **Construir e Publicar Imagens Docker:** Um novo workflow (`publish-docker.yml`) foi criado para automatizar a construção e publicação da imagem Docker do backend no GitHub Container Registry (GHCR).

### Próximos Passos

- **Implementar CD (Entrega Contínua):**
  - **Plataforma Escolhida (Staging):** **Railway**
    - **Justificativa:** A plataforma foi escolhida por sua simplicidade, facilidade de uso e por seu plano gratuito ser adequado para um ambiente de staging.
    - **Confirmação do Plano Gratuito:** O Railway oferece um plano gratuito que, após um período de trial com $5 em créditos, concede **$1 de uso de computação gratuito por mês**. Para um ambiente de staging com baixo tráfego, esse valor é suficiente para manter o serviço rodando sem custos, pois a cobrança é baseada apenas no uso ativo de CPU e memória.
  - **Estratégia de Deploy (Staging):**
    - **Gatilho:** Um novo workflow (`deploy-staging.yml`) será criado para ser acionado sempre que uma nova imagem do backend for publicada no GHCR.
    - **Ação:** O workflow usará a CLI do Railway para solicitar o deploy da nova imagem de Docker no ambiente de staging.
  - **Ações Necessárias:**
    1.  Criar um projeto no [Railway](https://railway.app/).
    2.  Dentro do projeto, criar um serviço que aponta para a imagem Docker do backend no GHCR: `ghcr.io/SEU_USUARIO/gestao_updev/backend`.
    3.  Gerar um [Token de API do Railway](https://docs.railway.app/reference/api-tokens) e adicioná-lo como um `secret` no repositório do GitHub com o nome `RAILWAY_TOKEN`.
  - **Deploy de Produção (Futuro):**
    - Após validar o fluxo de staging, o próximo passo será criar um workflow para produção, acionado pela criação de `tags` no Git, garantindo um processo de release mais controlado.

---
*Este documento foi atualizado para refletir o estado dos workflows após a conclusão das implementações.*
