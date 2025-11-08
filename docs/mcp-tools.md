# Ferramentas do Servidor MCP (Model Context Protocol)

Este documento descreve as ferramentas disponíveis através do servidor MCP local configurado para o projeto `gestao_updev`. O servidor MCP estende as capacidades do Gemini CLI, permitindo a execução de tarefas específicas do projeto diretamente da linha de comando.

## Como Iniciar o Servidor MCP

Para que as ferramentas estejam disponíveis, o servidor MCP deve estar em execução.

1.  **Certifique-se de ter o Node.js instalado.** Se não tiver, você pode baixá-lo em [nodejs.org](https://nodejs.org/).
2.  **Inicie o servidor MCP:**
    Abra um terminal na raiz do projeto `gestao_updev` e execute:
    ```bash
    node scripts/mcp-server.js &
    ```
    (O `&` no final fará com que o servidor rode em segundo plano, liberando seu terminal.)
    Você deverá ver uma mensagem como `MCP Server running on port 8081`.

## Como Usar as Ferramentas

Uma vez que o servidor MCP esteja em execução, você pode interagir com as ferramentas usando o Gemini CLI.

### Listar Ferramentas Disponíveis

Para ver todas as ferramentas expostas pelo servidor MCP, execute:

```bash
gemini tool list
```

Você deverá ver as ferramentas listadas sob o nome do servidor `gestao_updev_local_mcp/`.

### Executar uma Ferramenta

Para executar uma ferramenta específica, use o comando `gemini tool execute` seguido do nome completo da ferramenta (incluindo o prefixo do servidor).

## Ferramentas Disponíveis

### `gestao_updev_local_mcp/hello-world`

*   **Descrição:** Uma ferramenta simples que retorna uma mensagem de saudação.
*   **Uso:**
    ```bash
    gemini tool execute gestao_updev_local_mcp/hello-world
    ```
*   **Saída Esperada:** `Hello from gestao_updev MCP Server!`

### `gestao_updev_local_mcp/generate-api-types`

*   **Descrição:** Gera os tipos TypeScript para o frontend a partir da especificação OpenAPI do backend. Isso garante que o frontend esteja sempre sincronizado com as definições da API.
*   **Uso:**
    ```bash
    gemini tool execute gestao_updev_local_mcp/generate-api-types
    ```
*   **Comando Interno:** Executa `npm run generate:api-types` no diretório `frontend`.

### `gestao_updev_local_mcp/run-all-tests`

*   **Descrição:** Executa todos os testes do backend (Go) e do frontend (React/Next.js) em sequência.
*   **Uso:**
    ```bash
    gemini tool execute gestao_updev_local_mcp/run-all-tests
    ```
*   **Comandos Internos:**
    *   Backend: `docker compose -f docker-compose.test.yml up --build --abort-on-container-exit`
    *   Frontend: `npm run test` (no diretório `frontend`)

### `gestao_updev_local_mcp/deploy-backend-staging`

*   **Descrição:** Inicia o processo de implantação do backend Go para o ambiente de staging. Atualmente, esta é uma simulação e precisa ser implementada com a lógica de deploy real.
*   **Uso:**
    ```bash
    gemini tool execute gestao_updev_local_mcp/deploy-backend-staging
    ```
*   **Observação:** A lógica de implantação real (construção da aplicação Go, imagem Docker, push para registro, atualização de deployment) precisa ser adicionada ao script `scripts/mcp-server.js`.

---

Este documento será atualizado à medida que mais ferramentas forem adicionadas ao servidor MCP.
