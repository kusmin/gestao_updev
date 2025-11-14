# Status Atual do Projeto `gestao_updev`

## Visão Geral

O projeto `gestao_updev` é uma plataforma SaaS projetada para negócios locais, como barbearias e lojas de roupas. O objetivo é fornecer um painel centralizado para gerenciar clientes, agendamentos, estoque e vendas. A arquitetura é baseada em um backend Go e um frontend React/Next.js, com um banco de dados PostgreSQL. O sistema é projetado para ser multi-tenant.

O projeto encontra-se em uma fase inicial de desenvolvimento, mas já possui uma estrutura bem definida e um planejamento robusto, evidenciado pela vasta documentação existente.

## Status do Backend (Go)

O backend é construído em Go e segue uma arquitetura em camadas, promovendo modularidade e manutenibilidade.

### Módulos e Estrutura:

*   **`backend/cmd`**: Contém os pontos de entrada da aplicação:
    *   `api`: O ponto de entrada principal para a API.
    *   `migrate`: Para gerenciar migrações de banco de dados.
    *   `seed`: Para preenchimento inicial de dados no banco de dados.
*   **`backend/internal/domain`**: Define as entidades e a lógica de negócio.
    *   `models.go`: Contém modelos de dados abrangentes e bem definidos para as principais entidades do negócio, incluindo:
        *   `BaseModel`: Para campos de auditoria (ID, CreatedAt, UpdatedAt, DeletedAt).
        *   `TenantModel`: Para suporte a multi-tenancy, associando entidades a um `TenantID`.
        *   Entidades Principais: `Company`, `User`, `Client`, `Professional`, `Service`, `Product`, `InventoryMovement`, `Booking`, `SalesOrder`, `SalesItem`, `Payment`, `AuditLog`.
        *   Constantes: Para status de agendamento (`BookingStatus`), status de ordem de venda (`SalesOrderStatus`) e tipos de movimento de inventário (`InventoryMovement`).
*   **`backend/internal/repository`**: Responsável pela persistência de dados.
    *   Repositórios para `booking`, `company`, `product`, `sales` e `service` já estão presentes, indicando o início da implementação da camada de persistência para essas entidades.
    *   `repository.go`: Provavelmente define interfaces genéricas para repositórios.
*   **`backend/internal/service`**: Contém a lógica de negócio.
    *   Arquivos de serviço para diversas funcionalidades: `auth`, `bookings`, `catalog`, `clients`, `company`, `dashboard`, `inventory`, `professionals`, `sales`, `service` e `users`.
    *   Presença de arquivos de teste (`_test.go`) associados a alguns serviços, indicando preocupação com a qualidade do código.
*   **`backend/internal/http`**: Define os manipuladores de requisição HTTP.
    *   `contextutil`: Utilitários para manipulação de contexto HTTP.
    *   `handler`: Contém os manipuladores de requisição (controllers) que interagem com a camada de serviço. Inclui handlers para administração de diversas entidades, autenticação, agendamentos, catálogo, clientes, empresa, dashboard, inventário, profissionais, vendas e usuários.
    *   `response`: Define as estruturas de resposta padronizadas da API.
*   **`backend/internal/middleware`**: Contém middlewares HTTP.
*   **`backend/internal/config`**: Gerencia a configuração da aplicação.
*   **`backend/internal/server`**: Configuração e inicialização do servidor HTTP.

### Tecnologias e Ferramentas do Backend:

*   **Linguagem:** Go.
*   **Framework Web:** Gin (conforme `backend/GEMINI.md`).
*   **ORM:** GORM (`gorm.io/gorm`).
*   **UUIDs:** `github.com/google/uuid` para identificadores únicos.
*   **JSON:** `gorm.io/datatypes.JSONMap` para campos JSON flexíveis.
*   **Logging:** `zap` (conforme `backend/GEMINI.md`).
*   **Documentação API:** Swagger/OpenAPI (`swaggo/swag`).
*   **Linting:** `golangci-lint`.

## Status do Frontend (Aplicação Cliente)

O diretório `frontend/` representa uma Single Page Application (SPA) construída com React, TypeScript e Vite, que pode ser a aplicação voltada para o cliente final.

### Módulos e Estrutura (`frontend/src`):

*   **`App.tsx`**: O componente raiz da aplicação.
*   **`main.tsx`**: O ponto de entrada da aplicação React.
*   **`vite-env.d.ts`**: Definições de tipo para o ambiente Vite.
*   **`App.test.tsx`**: Testes para o componente `App`.
*   **`frontend/src/pages`**: Contém os componentes de página da aplicação.
    *   `clients`: Atualmente, apenas o diretório `clients` está presente, sugerindo que a página de gerenciamento de clientes é uma das primeiras a ser desenvolvida.
*   **`frontend/src/lib`**: Provavelmente contém utilitários, hooks personalizados ou outras lógicas reutilizáveis.
*   **`frontend/src/types`**: Contém as definições de tipos TypeScript. O arquivo `api.d.ts`, gerado a partir da especificação OpenAPI, foi criado com sucesso.
*   **`frontend/scripts/generate-types.mjs`**: Script responsável por gerar os tipos TypeScript da API a partir de `docs/api.yaml` para `frontend/src/types/api.d.ts`.

### Tecnologias e Ferramentas do Frontend (Aplicação Cliente):

*   **Framework:** React.
*   **Build Tool:** Vite.
*   **Linguagem:** TypeScript.
*   **Gerenciamento de Estado/Dados:** `@tanstack/react-query` para busca e gerenciamento de dados.
*   **UI Kit:** `@mui/material` (Material UI) e `@emotion/react`, `@emotion/styled` para componentes de UI e estilização.
*   **Linting:** ESLint.
*   **Testes:** Vitest, `@testing-library/react`.
*   **Tipagem da API:** `openapi-typescript` para geração de tipos.

## Status do Backoffice (Painel Administrativo)

O diretório `backoffice/` representa uma Single Page Application (SPA) separada, construída com React, TypeScript e Vite, que serve como o painel administrativo da plataforma.

### Módulos e Estrutura (`backoffice/src`):

*   **`App.tsx`**: O componente raiz da aplicação.
*   **`main.tsx`**: O ponto de entrada da aplicação React.
*   **`vite-env.d.ts`**: Definições de tipo para o ambiente Vite.
*   **`backoffice/src/pages`**: Contém os componentes de página da aplicação.
    *   Similar ao `frontend/`, espera-se que contenha páginas para dashboard, agendamentos, produtos, serviços, usuários, etc.
*   **`backoffice/src/lib`**: Provavelmente contém utilitários, hooks personalizados ou outras lógicas reutilizáveis.
*   **`backoffice/src/types`**: Contém as definições de tipos TypeScript. O arquivo `api.d.ts`, gerado a partir da especificação OpenAPI, também é esperado aqui.
*   **`backoffice/scripts/generate-types.mjs`**: Script responsável por gerar os tipos TypeScript da API a partir de `docs/api.yaml` para `backoffice/src/types/api.d.ts`.

### Tecnologias e Ferramentas do Backoffice:

*   **Framework:** React.
*   **Build Tool:** Vite.
*   **Linguagem:** TypeScript.
*   **Gerenciamento de Estado/Dados:** `@tanstack/react-query` (presumido, similar ao frontend).
*   **UI Kit:** `@mui/material` (Material UI) e `@emotion/react`, `@emotion/styled` para componentes de UI e estilização.
*   **Linting:** ESLint.
*   **Testes:** Vitest, `@testing-library/react`.
*   **Tipagem da API:** `openapi-typescript` para geração de tipos.

## Status da Documentação (`docs`)

O diretório `docs` é um ponto forte do projeto, sendo extremamente completo e bem organizado, cobrindo diversos aspectos cruciais para o desenvolvimento e a compreensão do projeto.

### Conteúdo Principal:

*   **Especificação da API:**
    *   `api.yaml`: A especificação OpenAPI (Swagger) que define o contrato da API REST, servindo como a "fonte da verdade".
    *   `api-reference.md`, `api-usage.md`, `ANALISE_API.md`, `api-changelog.md`: Documentos complementares sobre a API.
*   **Arquitetura:**
    *   `arquitetura-backend.md`: Detalhes sobre a arquitetura do backend.
    *   `FRONTEND_ARCHITECTURE.md`: Detalhes sobre a arquitetura do frontend.
*   **Padrões e Convenções:**
    *   `CODE_STANDARDS.md`, `padroes-codigo.md`: Diretrizes claras para padrões de código.
*   **Guias de Desenvolvimento:**
    *   `DEVELOPMENT.md`, `backend-setup.md`, `docker.md`, `USANDO_IMAGENS_DOCKER.md`: Informações para configurar e desenvolver o projeto, incluindo o uso de Docker.
*   **Funcionalidades do Backoffice:**
    *   Documentos detalhados para várias seções do backoffice: `backoffice-appointments.md`, `backoffice-clients.md`, `backoffice-dashboard.md`, `backoffice-products.md`, `backoffice-sales.md`, `backoffice-services.md`, `backoffice-tenants.md`, `backoffice-users.md`.
*   **Outros Tópicos Relevantes:**
    *   `authentication.md`, `CICD.md`, `CONTRIBUTING.md`, `fluxos-uso.md`, `modelo-dados.md`, `observabilidade.md`, `operacao-devops.md`, `plano-cobertura-codigo.md`, `security-review.md`, `TESTING_AND_SECURITY_PLAN.md`, `TESTING_STRATEGY.md`, `tests-contrato.md`.

## Conclusão

O projeto `gestao_updev` está em um estágio inicial promissor. O backend possui uma arquitetura sólida e um domínio de negócio bem modelado, com muitas funcionalidades já em desenvolvimento ou com sua estrutura definida. Existem dois projetos frontend distintos: uma **aplicação cliente (`frontend/`)** e um **painel administrativo (`backoffice/`)**, ambos começando a ser estruturados com a integração de ferramentas modernas. A documentação é um ponto forte, fornecendo um guia abrangente para o desenvolvimento e a compreensão do projeto.

**Próximos Passos Sugeridos:**

1.  **Gerar Tipos da API nos Frontends:** **CONCLUÍDO (frontend).** O comando `npm run generate:api-types` foi executado no diretório `frontend`, e o arquivo `frontend/src/types/api.d.ts` foi criado com sucesso. **PRÓXIMO PASSO:** Executar o mesmo comando no diretório `backoffice` para gerar `backoffice/src/types/api.d.ts`, garantindo a comunicação tipada para ambos os frontends.
2.  **Implementação Contínua:** Continuar o desenvolvimento das funcionalidades do backend, da aplicação cliente (`frontend/`) e do painel administrativo (`backoffice/`), seguindo as diretrizes e a arquitetura estabelecidas.
3.  **Testes Abrangentes:** Expandir a cobertura de testes em todas as camadas do backend, da aplicação cliente (`frontend/`) e do painel administrativo (`backoffice/`).
4.  **Refinamento da Documentação:** Manter a documentação atualizada conforme o projeto evolui.
