# Baseline E2E – 14/11/2025

- **Comando:** `npm run test:e2e`
- **Navegadores:** Chromium e Firefox para frontend/backoffice por padrão. WebKit pode ser habilitado exportando `PLAYWRIGHT_INCLUDE_WEBKIT=1` quando as dependências do sistema estiverem instaladas.
- **Servidores:** `npm run dev` de `frontend` (porta 5173) e `backoffice` (porta 5174) sob demanda via `webServer` do Playwright.
- **Spec files executados:**
  - `tests/e2e/tests/frontend/example.spec.ts`: valida a renderização da página de clientes e os cabeçalhos da tabela.
  - `tests/e2e/tests/backoffice/example.spec.ts`: garante redirecionamento automático para `/login` e comportamento básico do formulário.
- **Resultado:** 8 testes aprovados em ~12 segundos (2 fluxos × 2 apps × Chromium/Firefox). Relatório HTML disponível em `tests/e2e/playwright-report`.

Use este snapshot como linha de base. Ao adicionar novos fluxos críticos, atualize as specs correspondentes e anexos neste diretório.
