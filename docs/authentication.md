# Estratégia de Autenticação

Este documento descreve a estratégia para implementar as funcionalidades de autenticação (login, cadastro, recuperação de senha, login com Google) para as aplicações frontend e backoffice, integrando-as com o backend Go existente.

## 1. Visão Geral

O objetivo é fornecer um sistema de autenticação completo e seguro, aproveitando a infraestrutura JWT já presente no backend e adicionando suporte para autenticação via Google OAuth. As aplicações cliente (frontend e backoffice) serão responsáveis por apresentar a interface do usuário e interagir com os endpoints de autenticação do backend.

## 2. Estratégia de Backend (Go)

O backend já possui um sistema de autenticação robusto baseado em JWT (JSON Web Tokens), com suporte a tokens de acesso e refresh, e integração com multi-tenancy.

### 2.1. Mecanismos Existentes

*   **JWTManager (`backend/internal/auth/jwt.go`):** Responsável pela geração e validação de tokens de acesso (curta duração) e refresh (longa duração), utilizando HS256. Os claims do JWT incluem `UserID`, `TenantID` e `Role`.
*   **Endpoints de Autenticação (`backend/internal/http/handler/auth.go`):**
    *   `POST /v1/auth/signup`: Para registro de novos usuários e empresas.
    *   `POST /v1/auth/login`: Para autenticação de usuários existentes.
    *   `POST /v1/auth/refresh`: Para obter novos tokens de acesso usando um token de refresh válido.
*   **Middleware:**
    *   `middleware.TenantEnforcer`: Garante que todas as requisições protegidas estejam associadas a um `TenantID` válido.
    *   `middleware.Auth`: Valida o token de acesso JWT em requisições protegidas.
*   **Lógica de Serviço (`backend/internal/service/auth.go`):** Espera-se que este arquivo contenha a lógica de negócio para registro, login (incluindo hashing de senha) e renovação de tokens.

### 2.2. Integração Google OAuth

A integração com o Google OAuth será feita no backend para manter a lógica de segurança centralizada.

*   **Novos Endpoints:** Serão adicionados novos endpoints para iniciar e finalizar o fluxo OAuth, por exemplo:
    *   `GET /v1/auth/google/login`: Redireciona o usuário para a página de consentimento do Google.
    *   `GET /v1/auth/google/callback`: Recebe o callback do Google com o código de autorização.
*   **Fluxo OAuth 2.0:**
    1.  O frontend/backoffice redireciona o usuário para `/v1/auth/google/login`.
    2.  O backend inicia o fluxo OAuth com o Google, redirecionando o usuário para a página de consentimento do Google.
    3.  Após o consentimento, o Google redireciona o usuário de volta para `/v1/auth/google/callback` com um código de autorização.
    4.  O backend troca este código por tokens de acesso e ID do Google.
    5.  O backend valida o ID Token do Google para obter informações do usuário (email, nome, etc.).
    6.  Com base no email do Google, o backend verificará se o usuário já existe. Se não existir, um novo usuário será criado (ou associado a uma empresa existente, dependendo da lógica de negócio).
    7.  Finalmente, o backend emitirá seus próprios JWTs (acesso e refresh) para o cliente, seguindo o mesmo padrão do login tradicional.

## 3. Estratégia de Frontend e Backoffice (React/Next.js)

Ambas as aplicações cliente seguirão uma estratégia similar para interagir com o backend.

### 3.1. Componentes de UI Necessários

*   **Tela de Login:** Formulário para email/senha e botão "Login com Google".
*   **Tela de Cadastro:** Formulário para registro de novos usuários (email, senha, nome, etc.).
*   **Tela de Recuperação de Senha:** Formulário para solicitar redefinição de senha (envio de email).
*   **Redefinição de Senha:** Formulário para definir nova senha (após clique em link de email).

### 3.2. Lógica Cliente-Servidor

*   **Envio de Credenciais:** Enviar dados de formulário para os endpoints `/v1/auth/signup` ou `/v1/auth/login`.
*   **Armazenamento de Tokens:** Após o login/cadastro, os tokens de acesso e refresh recebidos do backend serão armazenados de forma segura (ex: `localStorage` para o token de acesso e `httpOnly` cookies para o token de refresh, se o backend for configurado para isso, ou ambos em `localStorage` com devidas precauções).
*   **Inclusão do Token de Acesso:** Todas as requisições subsequentes a endpoints protegidos do backend deverão incluir o token de acesso no cabeçalho `Authorization` (ex: `Bearer <access_token>`).
*   **Renovação de Token:** Implementar lógica para usar o token de refresh para obter um novo token de acesso quando o atual expirar.
*   **Fluxo Google Login:**
    1.  Ao clicar em "Login com Google", o cliente redireciona para o endpoint `/v1/auth/google/login` do backend.
    2.  O backend gerencia o fluxo OAuth com o Google.
    3.  Após a autenticação bem-sucedida no Google e no backend, o backend redirecionará o usuário de volta para o frontend/backoffice, possivelmente com os JWTs internos na URL ou em cookies.

### 3.3. Gerenciamento de Rotas Protegidas

*   Implementar um sistema de rotas protegidas que verifica a presença e validade do token de acesso antes de permitir o acesso a determinadas páginas. Usuários não autenticados serão redirecionados para a tela de login.

## 4. Considerações de Multi-tenancy

A autenticação deve respeitar a arquitetura multi-tenant. O `TenantID` presente nos claims do JWT será crucial para todas as operações subsequentes do usuário, garantindo que ele acesse apenas os dados de sua própria empresa.

## 5. Próximos Passos

1.  **Revisão e Aprovação:** Revisar esta estratégia e obter aprovação.
2.  **Implementação Backend:**
    *   Adicionar endpoints e lógica para Google OAuth.
    *   Implementar lógica de recuperação de senha (geração de token, envio de email, redefinição).
3.  **Implementação Frontend/Backoffice:**
    *   Desenvolver componentes de UI para login, cadastro, recuperação de senha.
    *   Integrar a lógica cliente-servidor com os endpoints do backend.
    *   Implementar o fluxo de login com Google.
    *   Gerenciar o estado de autenticação e rotas protegidas.
