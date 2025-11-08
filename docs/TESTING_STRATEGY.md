# Estratégia de Testes

Este documento descreve as diferentes estratégias de teste adotadas no projeto `gestao_updev`.

## Testes End-to-End (E2E)

Para garantir que os fluxos de usuário funcionem de ponta a ponta, utilizamos testes E2E que simulam interações reais em um navegador.

### Ferramenta: Playwright

- **O que é?** Playwright é um framework de automação de navegador moderno da Microsoft.
- **Por que foi escolhido?**
  - **Cross-Browser:** Permite executar os mesmos testes em Chromium (Chrome, Edge), Firefox e WebKit (Safari).
  - **Rapidez e Robustez:** Possui mecanismos de espera automática que reduzem a instabilidade ("flakiness") dos testes.
  - **Excelente Developer Experience:** Oferece ferramentas poderosas como o Trace Viewer, que grava um traço completo da execução do teste (ações, snapshots, logs, requisições de rede), e o Codegen, que grava interações e gera o código do teste automaticamente.
  - **Integração:** O `webServer` em sua configuração permite que ele inicie automaticamente o servidor de desenvolvimento do frontend antes de executar os testes.

### Localização dos Testes

- Os testes E2E, sua configuração e dependências estão localizados no diretório `tests/e2e`.

### Como Executar os Testes

1.  **Pré-requisito:** Certifique-se de que o backend e outras dependências (como o banco de dados) estejam em execução, caso o teste dependa deles. O servidor de desenvolvimento do frontend será iniciado automaticamente pelo Playwright.

2.  **Instalar Dependências (primeira vez):**
    ```bash
    npm install --prefix tests/e2e
    npx playwright install --prefix tests/e2e
    ```

3.  **Executar os Testes:**
    A partir da raiz do projeto, execute o seguinte comando:
    ```bash
    npm run test:e2e
    ```
    Este comando irá iniciar o Playwright, que por sua vez iniciará o servidor do frontend e executará todos os testes no diretório `tests/e2e/tests`.

4.  **Ver Relatórios:**
    Após a execução, um relatório HTML será gerado na pasta `playwright-report`. Para visualizá-lo, execute:
    ```bash
    npx playwright show-report tests/e2e/playwright-report
    ```

## Outros Tipos de Teste (TODO)

- **Testes de Unidade (Backend):** A serem implementados em Go.
- **Testes de Unidade e Componente (Frontend/Backoffice):** A serem implementados com Vitest e React Testing Library.
- **Testes de Contrato:** A estrutura inicial com Dredd está em `tests/dredd`.
