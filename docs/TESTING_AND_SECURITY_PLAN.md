# Plano de Melhoria de Testes e Segurança

## 1. Introdução e Objetivos

Este documento descreve um plano estratégico para aprimorar a qualidade, a confiabilidade e a segurança da aplicação `gestao_updev`. O objetivo é expandir a cobertura de testes em áreas críticas e implementar práticas de segurança robustas em todo o ciclo de desenvolvimento.

**Objetivos Principais:**

*   **Aumentar a Cobertura de Testes:** Atingir uma cobertura de testes de no mínimo 80% para os pacotes críticos do backend.
*   **Fortalecer a Segurança da Aplicação:** Identificar e mitigar vulnerabilidades de segurança, adotando uma abordagem de "segurança por padrão" (security by default).
*   **Garantir a Qualidade do Código:** Estabelecer um processo de desenvolvimento que previna a introdução de bugs e falhas de segurança.
*   **Melhorar a Manutenibilidade:** Criar uma base de código bem testada e segura, facilitando futuras manutenções e evoluções.

## 2. Análise do Estado Atual

### 2.1. Cobertura de Testes

Uma análise inicial revelou a seguinte situação:

*   **Testes Existentes:** O projeto já possui uma base de testes de integração para o pacote `service`, que interagem com o banco de dados. Existem também testes para middlewares, configuração e respostas de API.
*   **Baixa Cobertura:** Os logs de CI indicam uma cobertura de testes muito baixa para pacotes críticos, como `internal/service` (em torno de 2-10%).
*   **Áreas Críticas Sem Testes:**
    *   **`internal/http/handler`:** Não há testes para os manipuladores de API. Isso representa um risco significativo, pois a validação de entrada, o binding de payloads e a orquestração das chamadas de serviço não estão sendo validados automaticamente.
    *   **`internal/auth`:** A lógica de geração e validação de tokens JWT não possui testes de unidade, o que é crítico para a segurança da autenticação.

### 2.2. Segurança da Aplicação

*   **SAST (Static Application Security Testing):** O projeto já utiliza o **CodeQL**, o que é um excelente ponto de partida.
*   **SCA (Software Composition Analysis):** Não há uma verificação automatizada de dependências vulneráveis.
    *   **Backend:** A ferramenta `govulncheck` não está integrada ao CI.
    *   **Frontend/Backoffice:** O comando `npm audit` não é executado como parte do processo de CI.
*   **Gerenciamento de Segredos:** O vazamento acidental de um token do Codecov nos logs indica a necessidade de reforçar a cultura e as ferramentas para o gerenciamento seguro de segredos.
*   **Controles de Acesso:** A lógica de autorização (ex: garantir que um usuário só possa ver/editar os próprios dados ou os dados de seu tenant) precisa ser rigorosamente testada.

## 3. Plano de Expansão de Testes

A estratégia será focada em criar uma pirâmide de testes saudável, com uma base sólida de testes de unidade, seguida por testes de integração e E2E.

### 3.1. Testes de Unidade

O foco será em testar a lógica pura, sem dependências externas como banco de dados ou APIs.

*   **Pacote `internal/auth`:**
    *   **Objetivo:** Validar a geração e verificação de tokens JWT.
    *   **Ações:**
        1.  Criar `jwt_test.go`.
        2.  Testar a função `NewJWTManager`: garantir que inicializa corretamente.
        3.  Testar a função `Generate`: verificar se o token gerado contém os claims corretos (ID do usuário, ID do tenant, roles, etc.) e tem o tempo de expiração esperado.
        4.  Testar a função `Verify`: testar com tokens válidos, inválidos (assinatura incorreta), expirados e com claims ausentes.

*   **Pacote `internal/service`:**
    *   **Objetivo:** Testar a lógica de negócio de forma isolada.
    *   **Ações:**
        1.  Refatorar os serviços para dependerem de interfaces (`interface`) do repositório, em vez de implementações concretas (`*repository.Repository`).
        2.  Usar mocks (como os gerados pela biblioteca `testify/mock`) para simular o comportamento do banco de dados.
        3.  Escrever testes de unidade para cada função de serviço, cobrindo os cenários de sucesso e de erro (ex: o que acontece se o repositório retornar um erro?).

### 3.2. Testes de Integração

O foco será em testar a interação entre os componentes, especialmente com o banco de dados. Os testes existentes já seguem este padrão e devem ser expandidos.

*   **Pacote `internal/http/handler`:**
    *   **Objetivo:** Testar os handlers da API de ponta a ponta (sem o servidor HTTP real), validando o ciclo completo de request/response.
    *   **Ações:**
        1.  Criar arquivos `*_test.go` para cada handler (ex: `users_handler_test.go`).
        2.  Usar a biblioteca `net/http/httptest` para criar requests HTTP simulados e um `ResponseRecorder` para capturar as respostas.
        3.  Para cada endpoint, testar:
            *   **Caminho Feliz:** Request válido, status code 200/201 e corpo da resposta correto.
            *   **Validação de Entrada:** Requests com payloads inválidos (campos faltando, tipos errados) e garantir que um status code 400 (Bad Request) seja retornado com uma mensagem de erro clara.
            *   **Autorização:** Tentar acessar um endpoint protegido sem token ou com um token de outro usuário/tenant e garantir que um status code 401/403 seja retornado.
            *   **Cenários de Erro:** Simular erros no `service` layer (usando mocks) e garantir que o handler retorne o status code de erro apropriado (ex: 500).

### 3.3. Aumentar a Meta de Cobertura

*   **Ação:** No arquivo `codecov.yml`, aumentar gradualmente o `target` de cobertura para `80%` à medida que os novos testes forem sendo adicionados.

## 4. Plano de Melhoria de Segurança

### 4.1. Análise de Dependências Vulneráveis (SCA)

*   **Backend:**
    *   **Ação:** Adicionar um novo job no workflow `backend-ci.yml` para executar o `govulncheck`.
    *   **Exemplo de Job:**
        ```yaml
          security-scan:
            name: Security Scan
            runs-on: ubuntu-latest
            steps:
              - uses: actions/checkout@v5
              - uses: actions/setup-go@v6
                with:
                  go-version: '1.25'
              - name: Run govulncheck
                run: go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./...
        ```

*   **Frontend/Backoffice:**
    *   **Ação:** Adicionar uma etapa de `npm audit` nos workflows de CI correspondentes (`frontend-ci.yml`, etc.).
    *   **Exemplo de Etapa:**
        ```yaml
          - name: Run npm audit
            run: npm audit --audit-level=high
        ```

### 4.2. Revisão de Código e Práticas Seguras

*   **Validação de Entrada:** Padronizar o uso de `struct tags` de validação (ex: `binding:"required"`) nos DTOs (Data Transfer Objects) do Gin e garantir que todos os handlers validem os payloads de entrada.
*   **Gerenciamento de Segredos:** Realizar uma auditoria para garantir que nenhum segredo (chaves de API, senhas, segredos de JWT) esteja hardcoded no código. Utilizar exclusivamente variáveis de ambiente e os segredos do GitHub Actions.
*   **Princípio do Menor Privilégio:** Revisar os controles de acesso para garantir que as rotas da API e as consultas ao banco de dados apliquem rigorosamente o filtro de `tenant_id`, prevenindo vazamento de dados entre tenants. Os testes de integração para os handlers devem validar exaustivamente este cenário.

## 5. Documentação e Próximos Passos

1.  **Criar este Documento:** Salvar este plano como `docs/TESTING_AND_SECURITY_PLAN.md`.
2.  **Priorizar Tarefas:** Criar issues no GitHub para cada uma das ações descritas neste plano.
3.  **Implementação Incremental:** Começar implementando os testes de unidade para o pacote `auth`, seguido pelos testes de handler para as funcionalidades mais críticas. Em paralelo, adicionar os jobs de SCA aos workflows de CI.
4.  **Monitorar Métricas:** Acompanhar a evolução da cobertura de testes no Codecov e os resultados das varreduras de segurança a cada pull request.
