# Contexto do Gemini: `backend/internal/service`

Este pacote contém a lógica de negócio central da aplicação. Ele é responsável por orquestrar as operações, aplicar regras de negócio e coordenar o acesso aos dados através da camada de repositório.

## Responsabilidades

*   **Lógica de Negócio:** Implementa todos os casos de uso da aplicação (ex: criar um cliente, processar uma venda, agendar um horário).
*   **Coordenação de Repositórios:** Utiliza os métodos do pacote `repository` para ler e escrever dados no banco de dados. Pode coordenar operações que envolvem múltiplos repositórios dentro de uma transação.
*   **Validação de Negócio:** Realiza validações que dependem do estado da aplicação (ex: verificar se um produto tem estoque suficiente antes de uma venda), em contraste com a validação de formato feita na camada de handler.
*   **Autorização:** Aplica regras de autorização refinadas (ex: um usuário só pode ver agendamentos do seu próprio tenant).

## Convenções

*   A lógica de serviço é encapsulada na struct `Service`, que recebe o `repository` e outras dependências.
*   Os métodos de serviço devem ser agnósticos em relação ao protocolo HTTP. Eles recebem e retornam structs de domínio, não DTOs de HTTP.
*   Para facilitar os testes de unidade, os serviços devem depender de **interfaces** do repositório, não de implementações concretas, permitindo o uso de mocks.
