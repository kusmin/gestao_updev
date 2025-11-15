# Padrões de Código – Plataforma de Gestão Local

## Objetivo
Estabelecer convenções para backend Go e frontend React + Vite, garantindo consistência entre times e facilitando revisões.

## Geral
- **Idioma:** Todo o código, comentários e documentação devem ser escritos em português.
- **Formatação:** Use `prettier` para o frontend e `gofmt` para o backend para garantir um estilo consistente.

## Backend (Go)
- **Estilo**: seguir `gofmt`/`goimports` obrigatoriamente; usar `golangci-lint` com linters `govet`, `staticcheck`, `errcheck`, `gocyclo`.
- **Estrutura**: módulos dentro de `internal/` por domínio (ex.: `internal/service/agenda`); código compartilhado em `pkg/`.
- **Nomenclatura**:
  - Use `camelCase` para variáveis e `PascalCase` para tipos e funções exportadas.
  - Interfaces com sufixo `Service`, `Repository`.
  - Structs de request terminam com `Request`, responses com `Response`.
  - Use `ErrAlgo` para erros exportados.
- **Erros**:
  - Use `errors.Join`/`fmt.Errorf("%w", err)` para wrap.
  - Traduzir erros para respostas HTTP em middleware central.
- **Contexto**:
  - Todas funções públicas que tocam I/O recebem `context.Context`.
  - Deadlines/timeouts definidos no handler.
- **Tests**:
  - `*_test.go` com tabela de casos; usar `stretchr/testify` (require/assert).
  - Cobrir services e handlers críticos (meta 70%+).
- **Dependências**:
  - Gerenciar com `go mod tidy`.
  - Evitar libs pesadas; priorizar stdlib + pacotes bem mantidos.

## Frontend (React + Vite)
- **Estilo**: ESLint + Prettier; TypeScript obrigatório (`strict: true`). Regra de lint principal via `eslint.config.js`.
- **Arquitetura**:
  - Entrypoint único `src/main.tsx`, componentes em `src/components`, hooks em `src/hooks` e abstrações HTTP em `src/lib`.
  - Separar estilos globais (`src/index.css`) de estilos específicos (`*.module.css` ou Tailwind).
- **Estado**:
  - React Query para dados remotos; Zustand ou Context para estado global simples.
  - Nunca mutar cache diretamente; usar imutabilidade.
- **Naming**:
  - Componentes em PascalCase, hooks `useAlgo`.
  - CSS modules ou styled-components com nomes descritivos.
- **Acessibilidade**:
  - Semântica HTML, aria-labels em botões icônicos, foco visível.
- **Tests**:
  - Vitest + Testing Library para unidades/componentes; snapshots apenas para UI estática.
  - E2E (Playwright) focando fluxos críticos (signup, agenda, venda).

## Git & Pull Requests
- Commits no formato Conventional Commits (`feat:`, `fix:`, `docs:` etc.).
- PR deve incluir checklist: testes rodados, screenshots (quando aplicável), descrição do impacto.
- Branches: `feature/<slug>`, `fix/<slug>`, `chore/<slug>`.
- Revisões obrigatórias (2 revisores ou 1 revisor sênior para mudanças críticas).

## Quality Gates
- CI executa: lint, testes unitários, build backend/frontend.
- Codecov ou similar para monitorar cobertura (mínimo 60% backend, 50% frontend inicialmente).
- Nenhuma dependência com licenças restritivas sem aprovação.

## Segurança e Dados
- Nunca commitar secrets; usar `.env.example`.
- Sanitizar inputs no backend e escapar outputs HTML no frontend.
- Logs não devem conter dados sensíveis (CPF, tokens).

## Documentação
- Cada módulo deve ter README curto (dependências, comandos).
- Comentários no código só quando necessário para explicar decisões não triviais.
- ADRs (Architecture Decision Records) para mudanças significativas na arquitetura.

## Banco de Dados
- **Nomenclatura:** Use `snake_case` para tabelas e colunas.
- **Migrations:** Todas as alterações de esquema devem ser feitas através de arquivos de migração.

## Próximos Passos
1. Configurar linters e scripts (`make lint`, `npm run lint`) no repositório.
2. Criar templates de PR/issue com as convenções descritas.
3. Publicar `.editorconfig` para unificar indentação e line endings.
