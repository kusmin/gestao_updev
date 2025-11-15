# Multi-Tenancy na Plataforma Gestão UpDev

Este documento detalha a estratégia e a implementação da arquitetura multi-tenant na plataforma Gestão UpDev, focando no backend Go e na interação com o banco de dados PostgreSQL.

## Estratégia Geral

A plataforma adota uma abordagem de multi-tenancy baseada em **coluna (`tenant_id`)** para a maioria das entidades de negócio. Isso significa que cada registro em tabelas compartilhadas entre tenants inclui uma coluna `tenant_id` que identifica a qual tenant aquele dado pertence.

## Implementação no Backend (Go)

### 1. Identificação do Tenant

*   **Header `X-Tenant-ID`**: Todas as requisições autenticadas para endpoints protegidos devem incluir o header `X-Tenant-ID`. Este header é utilizado para identificar o tenant da requisição.
*   **Token JWT**: Após o login, o token JWT emitido contém o `tenant_id` do usuário. Este `tenant_id` é extraído do token e injetado no contexto da requisição.
*   **Middleware**: Um middleware específico (`middleware.TenantEnforcer`) é responsável por:
    *   Validar a presença do `X-Tenant-ID` para rotas protegidas.
    *   Extrair o `tenant_id` do token JWT e/ou do header `X-Tenant-ID`.
    *   Injetar o `tenant_id` no `context.Context` da requisição, tornando-o acessível para as camadas de serviço e repositório.
    *   Para rotas públicas (ex: `/auth/signup`, `/auth/login`), o `tenant_id` não é obrigatório.

### 2. Propagação do `tenant_id`

O `tenant_id` é propagado através do `context.Context` do Go. Isso garante que todas as operações subsequentes na cadeia de execução (serviços, repositórios) tenham acesso ao `tenant_id` da requisição atual.

### 3. Camada de Repositório

A camada de repositório é crucial para garantir o isolamento dos dados entre tenants.

*   **Queries Parametrizadas**: Todas as queries de banco de dados que acessam tabelas multi-tenant devem incluir uma cláusula `WHERE tenant_id = $1` (ou equivalente, dependendo do driver/ORM) para filtrar os dados pertencentes ao tenant da requisição.
*   **SQLC**: Ao utilizar SQLC, as queries são escritas com placeholders para o `tenant_id`, que é passado como um parâmetro para a função gerada.
*   **GORM**: Se GORM for utilizado, o `tenant_id` deve ser adicionado como um filtro global ou em cada operação de consulta/modificação.

### 4. Criação de Novas Entidades

Ao criar novas entidades em tabelas multi-tenant, o `tenant_id` deve ser obrigatoriamente preenchido com o `tenant_id` da requisição atual.

## Banco de Dados (PostgreSQL)

### 1. Design de Schema

*   **Coluna `tenant_id`**: Todas as tabelas que armazenam dados específicos de tenants devem incluir uma coluna `tenant_id` do tipo `UUID` (ou `VARCHAR` se for o caso).
*   **Chaves Estrangeiras**: Chaves estrangeiras para tabelas multi-tenant devem incluir o `tenant_id` como parte da chave composta para garantir a integridade referencial dentro do escopo do tenant.
*   **Índices**: Índices devem ser criados nas colunas `tenant_id` (e em combinação com outras colunas frequentemente consultadas) para otimizar o desempenho das queries multi-tenant.

### 2. Tabelas Globais

Algumas tabelas podem ser globais e não pertencer a nenhum tenant específico (ex: planos de assinatura, configurações gerais da plataforma). Essas tabelas não possuem a coluna `tenant_id` e são acessíveis por todos os tenants ou pelo sistema.

### 3. Migrações

As migrações de banco de dados (`backend/migrations`) devem garantir que a coluna `tenant_id` seja adicionada às tabelas apropriadas e que as chaves estrangeiras e índices sejam configurados corretamente para suportar a estratégia multi-tenant.

## Considerações de Segurança

*   **Isolamento de Dados**: A implementação rigorosa do filtro `tenant_id` em todas as queries é fundamental para garantir que um tenant não possa acessar os dados de outro.
*   **Validação de Entrada**: Além da validação do `tenant_id` no middleware, a lógica de negócio deve sempre considerar o `tenant_id` ao realizar operações sensíveis.

## Futuras Evoluções

*   **Schema por Tenant**: Se a necessidade de isolamento de dados for mais rigorosa ou se houver requisitos de conformidade específicos, uma migração para uma abordagem de "schema por tenant" pode ser considerada. Isso envolveria a criação de um schema de banco de dados separado para cada tenant.
*   **Sharding**: Para escalar a plataforma para um grande número de tenants, estratégias de sharding (distribuição de tenants em diferentes instâncias de banco de dados) podem ser exploradas.
