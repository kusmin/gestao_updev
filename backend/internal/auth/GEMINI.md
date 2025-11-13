# Contexto do Gemini: `backend/internal/auth`

Este pacote lida com a lógica de autenticação da aplicação, especificamente a geração e validação de JSON Web Tokens (JWT).

## Responsabilidades

*   **Geração de Token:** Cria tokens de acesso e de atualização (refresh tokens) para usuários autenticados. Os tokens contêm `claims` essenciais, como ID do usuário, ID do tenant e roles.
*   **Validação de Token:** Verifica a validade de um token JWT, incluindo a assinatura, o tempo de expiração e os `claims` esperados.
*   **Gerenciamento de Segredos:** Utiliza os segredos de JWT (JWT secrets) fornecidos pela configuração para assinar e verificar os tokens.

## Convenções

*   A lógica é encapsulada na struct `JWTManager`.
*   As funções devem ser puras e testáveis, dependendo apenas das chaves secretas e dos dados do usuário para operar.
