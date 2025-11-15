# Plano para Início das Rotas no Frontend

## Contexto Atual
- O frontend (`frontend/`) usa React + Vite + TypeScript, Material UI e React Query, porém ainda não possui uma camada de roteamento configurada.
- Não há dependência como `react-router-dom` declarada; o app inicial renderiza um único componente.

## Objetivos
1. Estruturar a árvore de rotas baseada nos fluxos descritos em `docs/backoffice-*.md` e `docs/fluxos-uso.md`.
2. Garantir navegação segura (public vs. autenticado) e suporte a layouts compartilhados.
3. Facilitar o carregamento de dados da API (`/v1/...`) com React Query e pré-busca de metadados (empresa/usuário logado).

## Dependências e Setup
1. Adicionar `react-router-dom@^7` e tipos correspondentes:
   ```bash
   pnpm --dir frontend add react-router-dom
   pnpm --dir frontend add -D @types/react-router-dom # apenas se necessário
   ```
2. Criar utilitário para mapear constantes de rota (ex.: `src/routes/paths.ts`) para evitar strings soltas.
3. Atualizar `src/main.tsx` para envolver a aplicação com `<RouterProvider>` e `<QueryClientProvider>`.

## Arquitetura de Rotas Proposta
```
/
├── public
│   ├── /login
│   └── /forgot-password (futuro)
└── app (rota protegida)
    ├── /dashboard
    ├── /clients
    │   ├── /clients/new
    │   └── /clients/:clientId
    ├── /services
    ├── /products
    ├── /appointments
    ├── /sales
    └── /settings (empresa/usuário)
```
- Usar `<createBrowserRouter>` com loaders opcionais para verificar sessão.
- Implementar `ProtectedRoute` reutilizando contexto de autenticação ou hook (`useAuth()`).
- Layout principal (`AppLayout`) deve conter sidebar, topbar e `<Outlet />`.

## Plano Passo a Passo

### Fase 0 – Preparação e Boilerplate
1. Instalar `react-router-dom` (e tipos, se necessário).
2. Criar pasta `src/routes/` com:
   - `paths.ts`: constantes para rotas (`ROOT`, `LOGIN`, `CLIENTS`, etc.).
   - `router.tsx`: função `createRouter()` que exporta a configuração do `createBrowserRouter`.
3. Atualizar `src/main.tsx`:
   - Instanciar `QueryClient`.
   - Renderizar `<RouterProvider router={router} />` dentro de `<QueryClientProvider>`.
4. Criar página placeholder para `/login` (ex.: `LoginPage.tsx`) e para `/app` (`DashboardPlaceholder`).
5. Validar rotas básicas via `pnpm --dir frontend run dev`.

### Fase 1 – Contexto de Autenticação e Guards
1. Implementar `src/contexts/AuthContext.tsx` com estados:
   - usuário logado, tokens, métodos `login/logout`.
   - Persistir session em `localStorage` (ou cookies).
2. Criar hook `useAuth()` retornando contexto.
3. Implementar `ProtectedRoute` (componente ou loader) que:
   - verifica `auth.user`;
   - redireciona para `/login` se ausente;
   - injeta tenant header/id para loaders futuros.
4. Atualizar rotas para usar `ProtectedRoute` na árvore `/app`.

### Fase 2 – Layout Principal e Navegação
1. Criar `src/components/layout/AppLayout.tsx` contendo:
   - `<Sidebar>` com links (clientes, serviços, produtos...);
   - `<Topbar>` com nome do usuário/empresa e botão de logout;
   - `<Outlet />` para telas internas.
2. Adicionar `src/components/layout/PublicLayout.tsx` para telas de auth (sem sidebar).
3. Criar `src/components/navigation/AppMenu.ts` com definição de itens (título + path + ícone).
4. Atualizar `router.tsx` para usar `AppLayout` na rota protegida e `PublicLayout` na rota `/login`.

### Fase 3 – Páginas Prioritárias
1. **Dashboard**:
   - Página com cards/resumo (`src/pages/dashboard/DashboardPage.tsx`).
   - Loader usando React Query para `/v1/dashboard/daily`.
2. **Clientes**:
   - `ClientListPage` com tabela/pesquisa.
   - Hook `useClientsQuery` em `src/features/clients/api.ts`.
3. **Serviços/Produtos**:
   - Listagens básicas reaproveitando componentes de tabela.
4. Adicionar rotas correspondentes em `paths.ts` e navegação lateral.

### Fase 4 – Formulários e Rotas Aninhadas
1. Rotas aninhadas para `/clients/new`, `/clients/:clientId`.
2. Componentes de formulário isolados (`src/features/clients/components/ClientForm.tsx`).
3. Utilizar `react-hook-form` (se aprovado) e validações compartilhadas.
4. Loader + action (opcional) para carregar dados ao entrar em `/clients/:clientId`.

### Fase 5 – Estados Avançados e UX
1. Rotas com query params (filtros de agenda, vendas).
2. Tabs internas (ex.: `/clients/:clientId/history`).
3. Breadcrumbs dinâmicos e indicadores de carregamento global.
4. Code splitting com `React.lazy` para rotas pesadas.

## Boas Práticas
- **Code-splitting:** usar `lazy()` e `Suspense` a partir da Fase 3 para telas pesadas.
- **Definição de tipos:** cada rota deve exportar tipos de loader/params; criar `RouteParams` em `src/routes/types.ts`.
- **Testes:** adicionar testes com React Testing Library simulando navegação (`MemoryRouter`).
- **Acessibilidade:** assegurar foco automático e breadcrumbs onde aplicável.

## Próximos Passos
1. Revisar e aprovar este plano com o time.
2. Registrar épicos/tarefas por fase (0 a 5) no board.
3. Implementar Fase 0 + Fase 1 em um único PR inicial (roteador + guard + layouts básicos).
4. Sequenciar Fases 2–5 em sprints subsequentes com validações de design/UX.
