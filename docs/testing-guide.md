# Guia Detalhado de Estratégia de Testes

Este documento detalha a estratégia de testes para a plataforma Gestão UpDev, abrangendo testes unitários, de integração e End-to-End (E2E) para o backend Go e as aplicações frontend/backoffice React. Ele visa fornecer diretrizes, exemplos e melhores práticas para garantir a qualidade e a confiabilidade do software.

## Filosofia de Testes

Adotamos uma abordagem de "pirâmide de testes", com a maioria dos testes sendo unitários, seguidos por testes de integração e uma menor quantidade de testes E2E.

*   **Testes Unitários:** Foco em componentes isolados, funções e métodos. Rápidos e fáceis de escrever.
*   **Testes de Integração:** Verificam a interação entre diferentes módulos ou serviços (ex: serviço e repositório, API e banco de dados).
*   **Testes End-to-End (E2E):** Simulam o fluxo completo do usuário através da interface gráfica, validando a funcionalidade do sistema como um todo.

## Backend (Go)

### Ferramentas

*   **Go Testing Package:** `testing` (stdlib) para testes unitários e de integração.
*   **Testify:** `stretchr/testify` para asserções e mocks.
*   **Docker Compose:** Para ambientes de teste de integração (banco de dados, serviços externos).

### Testes Unitários

*   **Foco:** Funções puras, lógica de negócio isolada, validações.
*   **Localização:** Arquivos `_test.go` no mesmo pacote do código testado.
*   **Convenções:**
    *   Nomes de funções de teste: `Test[NomeDaFuncao]` ou `Test[NomeDoMetodo]_[Cenario]`.
    *   Uso de `t.Run()` para subtestes.
    *   Mocks/Stubs para dependências externas (ex: repositórios, serviços HTTP).
*   **Exemplo:**
    ```go
    // internal/service/auth_test.go
    func TestAuthService_Login(t *testing.T) {
        // ... setup mocks ...
        svc := NewAuthService(mockRepo, mockJWT)

        t.Run("should return tokens for valid credentials", func(t *testing.T) {
            // ... test logic ...
        })

        t.Run("should return error for invalid password", func(t *testing.T) {
            // ... test logic ...
        })
    }
    ```

### Testes de Integração

*   **Foco:** Interação entre camadas (handler-service-repository), acesso ao banco de dados, middlewares.
*   **Localização:** Arquivos `_test.go` no mesmo pacote, ou em um pacote `[nome_do_pacote]_test` separado para testes de integração mais pesados.
*   **Convenções:**
    *   Utilizar um banco de dados de teste (ex: PostgreSQL via Docker Compose).
    *   Limpar o estado do banco de dados antes de cada teste ou suíte de testes.
    *   Testar endpoints HTTP completos (ex: usando `httptest.NewRecorder`).
*   **Exemplo:**
    ```go
    // internal/http/user_handler_test.go
    func TestUserHandler_CreateUser(t *testing.T) {
        // ... setup test DB, API server ...
        req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(jsonPayload))
        // ... perform request ...
        // ... assert response ...
    }
    ```

### Testes de Contrato (Dredd)

*   **Foco:** Garantir que a implementação da API esteja em conformidade com a especificação OpenAPI (`docs/api.yaml`).
*   **Ferramenta:** Dredd.
*   **Localização:** `tests/dredd/dredd.yml` e arquivos de blueprint/spec.
*   **Execução:** Via workflow `api-contract.yml`.

## Frontend (React + Vite) e Backoffice (React + Vite)

### Ferramentas

*   **Vitest:** Framework de testes rápido.
*   **React Testing Library:** Para testar componentes React de forma que simule o uso real do usuário.
*   **MSW (Mock Service Worker):** Para mockar requisições de rede em testes de integração.
*   **Playwright:** Para testes End-to-End (E2E).

### Testes Unitários (Componentes e Hooks)

*   **Foco:** Lógica de componentes isolados, hooks personalizados, utilitários.
*   **Localização:** Arquivos `*.test.tsx` ou `*.test.ts` próximos ao código testado.
*   **Convenções:**
    *   Testar a renderização e interação do usuário com os componentes.
    *   Mockar dependências externas (ex: chamadas de API, contexto React).
*   **Exemplo (React Testing Library):**
    ```typescript jsx
    // src/components/Button.test.tsx
    import { render, screen } from '@testing-library/react';
    import Button from './Button';

    test('renders button with correct text', () => {
      render(<Button>Click me</Button>);
      expect(screen.getByText(/click me/i)).toBeInTheDocument();
    });
    ```

### Testes de Integração (Comunicação com API)

*   **Foco:** Interação entre componentes, comunicação com a API (mockada), gerenciamento de estado.
*   **Ferramenta:** MSW para interceptar requisições HTTP.
*   **Exemplo (MSW + React Query):**
    ```typescript jsx
    // src/pages/Clients.test.tsx
    import { render, screen, waitFor } from '@testing-library/react';
    import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
    import { setupServer } from 'msw/node';
    import { rest } from 'msw';
    import ClientsPage from './ClientsPage';

    const server = setupServer(
      rest.get('/v1/clients', (req, res, ctx) => {
        return res(ctx.json({ data: [{ id: '1', name: 'Test Client' }] }));
      })
    );

    beforeAll(() => server.listen());
    afterEach(() => server.resetHandlers());
    afterAll(() => server.close());

    test('fetches and displays clients', async () => {
      const queryClient = new QueryClient();
      render(
        <QueryClientProvider client={queryClient}>
          <ClientsPage />
        </QueryClientProvider>
      );
      await waitFor(() => expect(screen.getByText('Test Client')).toBeInTheDocument());
    });
    ```

### Testes End-to-End (E2E)

*   **Foco:** Validar fluxos completos do usuário através da interface, interagindo com o backend real.
*   **Ferramenta:** Playwright.
*   **Localização:** `tests/e2e/`.
*   **Execução:** Via workflow `e2e.yml`.
*   **Convenções:**
    *   Testar cenários críticos (ex: signup, login, criação de agendamento, registro de venda).
    *   Utilizar seletores robustos (ex: `data-testid`, `role`).
    *   Limpar o estado do sistema (ex: banco de dados) antes de cada suíte de testes E2E.

## Cobertura de Código

*   **Ferramenta:** `go tool cover` (Go), Vitest (Frontend/Backoffice).
*   **Relatórios:** Codecov para agregação e visualização.
*   **Metas:**
    *   Backend: 70% de cobertura de linhas.
    *   Frontend/Backoffice: 50% de cobertura de linhas.
*   **Workflow:** `coverage.yml` para execução e upload.

## Boas Práticas Gerais

*   **Testes Rápidos e Confiáveis:** Testes devem ser rápidos para não atrasar o ciclo de desenvolvimento e confiáveis (não devem falhar intermitentemente).
*   **Testes Legíveis:** O código de teste deve ser tão legível quanto o código de produção.
*   **Testar Comportamento, Não Implementação:** Focar no que o código faz, não em como ele faz.
*   **Testes de Borda e Casos de Erro:** Incluir testes para cenários de erro, entradas inválidas e condições de borda.
*   **Manter Testes Atualizados:** Atualizar os testes sempre que o código de produção for alterado.
