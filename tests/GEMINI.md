# GEMINI Project Context: gestao_updev - Módulo de Testes

## Visão Geral do Projeto

Este diretório (`tests`) contém as configurações e scripts para os diferentes tipos de testes da plataforma SaaS `gestao_updev`. O objetivo é garantir a qualidade e o correto funcionamento das APIs e das aplicações frontend/backoffice através de testes de contrato, end-to-end e funcionais.

## Ferramentas de Teste Utilizadas

*   **Dredd (Testes de Contrato de API):** Utilizado para validar se a implementação da API está em conformidade com sua especificação (`docs/api.yaml`).
*   **Playwright (Testes End-to-End - E2E):** Utilizado para testar a funcionalidade completa das aplicações frontend e backoffice, simulando interações de usuário em navegadores reais.
*   **Postman/Newman (Testes Funcionais de API):** Utilizado para executar coleções de testes de API, verificando o comportamento dos endpoints.

## Como Executar os Testes

### Testes de Contrato (Dredd)

Para executar os testes de contrato, é necessário que o servidor da API esteja em execução. O Dredd utilizará o script `scripts/run_dredd_server.sh` para iniciar o servidor, se configurado.

```bash
# Navegue até o diretório 'tests'
cd tests

# Execute os testes Dredd
# Certifique-se de que o blueprint da API (docs/api.yaml) e os hooks estejam configurados corretamente.
dredd
```

### Testes End-to-End (Playwright)

Os testes E2E com Playwright exigem que as aplicações `frontend` e `backoffice` estejam em execução. O Playwright é configurado para iniciar esses servidores automaticamente se não estiverem rodando.

```bash
# Navegue até o diretório 'tests/e2e'
cd tests/e2e

# Instale as dependências (se ainda não o fez)
npm install

# Instale os navegadores necessários para o Playwright
npm run install:browsers

# Execute os testes E2E
npm run test
```

### Testes Funcionais de API (Postman/Newman)

Os testes funcionais de API são executados usando Newman, a ferramenta de linha de comando do Postman.

```bash
# Navegue até o diretório 'tests/postman'
cd tests/postman

# Instale as dependências (se ainda não o fez)
npm install

# Execute a coleção de testes do Postman
npm run test
```

## Convenções de Desenvolvimento

*   **Estrutura de Diretórios:** Os testes são organizados por tipo (`dredd`, `e2e`, `postman`).
*   **Configuração:** As configurações específicas de cada ferramenta de teste são mantidas em seus respectivos diretórios (ex: `dredd.yml`, `playwright.config.ts`, `collection.json`).
*   **Scripts:** Scripts de execução e auxiliares são definidos nos arquivos `package.json` ou em scripts shell dedicados.
*   **Idioma:** Comentários e descrições nos arquivos de teste devem ser em português, seguindo a convenção geral do projeto.
