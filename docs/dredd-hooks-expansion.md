# Plano de Expansão para Dredd Hooks (`basic-flow.js`)

Este documento detalha o plano para expandir a cobertura dos testes de contrato da API utilizando o Dredd, especificamente o arquivo de hooks `basic-flow.js`. O objetivo é aumentar a abrangência dos testes, incluindo mais endpoints e fluxos de recursos, garantindo que o contrato da API seja validado de forma mais completa.

## Contexto Atual (`basic-flow.js`)

O arquivo `basic-flow.js` é responsável por preparar o ambiente e os dados para cada requisição de teste do Dredd. Ele gerencia:
*   **Autenticação:** Realiza o signup e login de um usuário administrativo, armazenando `access_token`, `refresh_token` e `tenantId`.
*   **Criação de Recursos:** Cria um usuário e um cliente de teste, armazenando seus IDs.
*   **Manipulação de Recursos:** Testa operações básicas (GET, PATCH, PUT, DELETE) para usuários e clientes.
*   **Controle de Fluxo:** Pula testes para endpoints não suportados ou quando recursos necessários não estão disponíveis.

## Recursos a Serem Expandidos

O plano de expansão visa incluir a cobertura para os seguintes recursos da API, que são fundamentais para a plataforma `gestao_updev`:

*   **Profissionais (`/professionals`)**:
    *   `GET /professionals`
    *   `POST /professionals`
    *   `GET /professionals/{id}`
    *   `PUT /professionals/{id}`
    *   `DELETE /professionals/{id}`
*   **Serviços (`/services`)**:
    *   `GET /services`
    *   `POST /services`
    *   `GET /services/{id}`
    *   `PUT /services/{id}`
    *   `DELETE /services/{id}`
*   **Produtos (`/products`)**:
    *   `GET /products`
    *   `POST /products`
    *   `GET /products/{id}`
    *   `PUT /products/{id}`
    *   `DELETE /products/{id}`
*   **Agendamentos (`/bookings`)**:
    *   `GET /bookings`
    *   `POST /bookings`
    *   `GET /bookings/{id}`
    *   `PUT /bookings/{id}`
    *   `DELETE /bookings/{id}`
*   **Vendas (`/sales`)**:
    *   `GET /sales`
    *   `POST /sales`
    *   `GET /sales/{id}`
    *   `PUT /sales/{id}`
    *   `DELETE /sales/{id}`

## Detalhamento do Plano de Ação

### 1. Atualizar `SUPPORTED` e `PUBLIC_ENDPOINTS`

*   Adicionar todos os novos endpoints listados acima ao `Set` `SUPPORTED`.
*   Verificar se algum desses novos endpoints deve ser adicionado ao `Set` `PUBLIC_ENDPOINTS` (provavelmente não, pois a maioria das operações CRUD exige autenticação).

### 2. Adicionar Variáveis de Contexto (`ctx`)

Expandir o objeto `ctx` para armazenar os IDs dos recursos criados durante os testes, permitindo que testes subsequentes os referenciem.

Exemplos de novas variáveis:
*   `ctx.createdProfessionalId`
*   `ctx.createdServiceId`
*   `ctx.createdProductId`
*   `ctx.createdBookingId`
*   `ctx.createdSalesOrderId`

### 3. Expandir `hooks.beforeEach()`

Esta função será modificada para incluir a lógica de preparação para os novos endpoints:

*   **Para `POST /professionals`**:
    *   Gerar dados de teste para um novo profissional (nome, email, telefone, etc.).
    *   Usar `setJSONBody(transaction, payload)` para definir o corpo da requisição.
*   **Para `PATCH /professionals/{id}`, `PUT /professionals/{id}`, `DELETE /professionals/{id}`**:
    *   Usar `ensureResource(transaction, 'professionals', ctx.createdProfessionalId, 'Profissional de teste indisponível...')` para garantir que um profissional tenha sido criado e definir o ID na URL.
    *   Para `PUT`/`PATCH`, gerar dados de atualização e usar `setJSONBody()`.
*   **Repetir o padrão** para `services`, `products`, `bookings` e `sales`.
    *   Para `POST /bookings` e `POST /sales`, será necessário garantir que `clientId`, `professionalId`, `serviceId`, `productId` (conforme o caso) estejam disponíveis no `ctx` antes de criar o recurso. Isso pode exigir a criação desses recursos em testes anteriores ou a adição de lógica para pular o teste se eles não existirem.
*   **Para `GET /professionals`, `GET /services`, `GET /products`, `GET /bookings`, `GET /sales`**:
    *   Apenas garantir que os cabeçalhos de autenticação estejam definidos (já tratado pela lógica `!PUBLIC_ENDPOINTS.has(key)`).

### 4. Expandir `hooks.afterEach()`

Esta função será modificada para extrair os IDs dos recursos recém-criados das respostas da API:

*   **Para `POST /professionals`**:
    *   Extrair `data.id` da resposta e armazenar em `ctx.createdProfessionalId`.
*   **Repetir o padrão** para `services`, `products`, `bookings` e `sales`.

## Considerações Adicionais

*   **Ordem dos Testes:** A ordem em que os testes são executados pelo Dredd é crucial. Garantir que os recursos sejam criados antes de serem manipulados (GET, PUT, PATCH, DELETE). O Dredd geralmente segue a ordem definida no arquivo OpenAPI.
*   **Limpeza de Dados:** Embora não seja estritamente parte dos hooks `beforeEach`/`afterEach` para criação/manipulação, é importante considerar uma estratégia de limpeza de dados para evitar que os testes de contrato deixem dados sujos no ambiente de teste. Isso pode ser feito com hooks `afterAll` ou scripts externos.
*   **Dados de Teste:** Utilizar `uniqueSuffix()` para garantir que os dados de teste sejam únicos a cada execução, evitando colisões.
*   **Mensagens de Skip:** Manter as mensagens de `skipReason` claras para facilitar a depuração.

Este plano fornecerá uma cobertura de teste de contrato muito mais robusta para a API `gestao_updev`.