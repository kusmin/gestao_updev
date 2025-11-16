# Plano de Implementação do Frontend

Este documento resume os requisitos identificados no `docs/api.yaml`, a arquitetura proposta e o plano de ação para criar as páginas do frontend.

## 1. Levantamento de requisitos

- **Autenticação (`/auth/login`, `/auth/signup`, `/auth/refresh`)**: precisamos de telas públicas para cadastro inicial da empresa/usuário administrador, login e renovação automática de tokens.
- **Perfil da empresa (`/companies/me`)**: página autenticada para exibir/editar nome, documento, telefone e timezone do tenant.
- **CRUDs principais**: clientes, usuários, serviços, produtos, profissionais, bookings, sales orders, inventory e payments seguem o padrão listagem + criação + edição + exclusão. Cada página deve consumir os endpoints correspondentes e respeitar o contrato `response.APIResponse`.
- **Dashboard/KPIs (`/dashboard/*`)**: fornecer cards na landing autenticada com métricas diárias e agregadas.
- **Fluxos correlatos**:
  - Bookings dependem de clientes/profissionais/serviços.
  - Sales orders e payments exigem seleção de itens e registro de pagamentos.
  - Inventory movements impactam produtos.
- **Resposta padrão**: quase todas as rotas devolvem `response.APIResponse { data, meta, error }`. O cliente deve normalizar o payload e padronizar erros.

## 2. Arquitetura proposta

1. **Roteamento**: adicionar React Router, separando rotas públicas (`/login`, `/signup`) e privadas (`/app/*`) com um `ProtectedRoute`.
2. **Contextos**:
   - `AuthContext`: armazena tokens, tenant e dados do usuário; executa `/auth/refresh` automaticamente e injeta o header `Authorization`.
   - `TenantContext` (opcional): centraliza o `X-Tenant-ID`.
3. **Camada de API**:
   - Manter `scripts/generate-types.mjs` para gerar tipos e padronizar a camada `lib/apiClient`.
   - Considerar TanStack Query para cache/carregamento.
4. **Componentes compartilhados**:
   - `Layout` (AppBar + Drawer), formulários (React Hook Form + MUI), tabelas reutilizáveis e componentes de feedback (snackbars, modais).
   - Hooks utilitários (`useApi`, `useToast`, `useConfirmDialog`).
5. **Configuração**:
   - Usar `VITE_API_BASE_URL` e normalizar URLs no `apiClient`.
   - Incluir loaders, tratamento uniforme de erros (401 → logout) e logs para debugging.

## 3. Plano de ação

### 1. Infraestrutura inicial
- Adicionar React Router e definir rotas públicas/privadas.
- Implementar `AuthContext`, persistência segura (Storage) e fluxo de refresh.
- Atualizar `App.tsx` para montar `Layout`, `Sidebar` e menus dos módulos previstos.
- Garantir que `pnpm run generate:api-types` rode antes de build/lint/test (script `prebuild/prelint/predev` já existe; validar em pre-commit/CI).

### 2. Páginas prioritárias
- **Auth**: telas de signup/login com validações, mensagens de erro e redirecionamento pós-login.
- **Dashboard**: cards com dados de `/dashboard/daily` e KPIs úteis.
- **Clientes/Usuários**: consolidar a página existente de clientes (remover IDs fixos) e criar módulo equivalente para usuários (filtros, paginação, CRUD completo).

### 3. Cadastros restantes
- **Serviços/Produtos/Profissionais**: replicar o padrão de CRUD com formulários reutilizáveis.
- **Bookings**: lista/calendário com ações (confirmar/cancelar) e formulário com seleção de cliente/serviço/profissional.
- **Sales orders & Payments**: wizard para seleção de itens, status e registro de pagamentos (`/sales/orders/{id}/payments`).
- **Inventory movements**: tela simples para entradas/saídas (formulário + lista histórica).

### 4. Testes e validações
- Para cada página, adicionar testes com Vitest + Testing Library cobrindo carregamento, interações e mensagens de erro.
- Atualizar pipelines: `pnpm --filter frontend run test -- --coverage`, `pnpm --filter backoffice test -- --coverage`, `make api-contract-test` e `make coverage-backend` antes de PRs.

### 5. Automação & documentação
- Documentar o fluxo “docs → gerar tipos → implementar → rodar testes” no README/Contributing.
- Expandir o pre-commit (se necessário) para rodar `pnpm run generate:api-types` quando `docs/api.yaml` muda.
- Manter o `api-reference.html` (Redoc) atualizado no pipeline “API Spec Quality”.

## 4. Dependências e riscos

- **Back-end**: cobertura ainda falha (`backend/internal/config` e `bookings_test`); precisamos do fix antes de depender totalmente dos testes de contrato.
- **CI/CD**: `make coverage-backend` e `make api-contract-test` devem rodar limpos para aceitar merges; planejar janelas para ajustes na spec.
- **Autenticação/tenant**: é fundamental definir como o front obtém o `tenant_id` (armazenado após signup/login) para eliminar IDs hardcoded.

Com essa estrutura, podemos evoluir o frontend de maneira incremental, garantindo consistência com a spec OpenAPI e cobrindo os principais fluxos do produto.
