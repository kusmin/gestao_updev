# Backoffice - Gerenciamento de Usuários

Este documento descreve a funcionalidade de gerenciamento de usuários no backoffice da plataforma Gestão UpDev.

## Visão Geral

A página de gerenciamento de usuários permite que os administradores da plataforma visualizem, criem, editem e removam usuários de todos os tenants.

## Funcionalidades

### Listagem de Usuários

A página principal exibe uma lista de todos os usuários cadastrados na plataforma. A lista inclui as seguintes informações:

-   **Nome:** O nome do usuário.
-   **Email:** O email do usuário.
-   **Telefone:** O telefone de contato do usuário.
-   **Role:** O papel do usuário (e.g., `admin`, `professional`).
-   **Ativo:** Um switch para ativar ou desativar o usuário.

### Criação e Edição de Usuários

É possível criar um novo usuário ou editar um existente através de um formulário modal. O formulário contém os seguintes campos:

-   Nome
-   Email
-   Telefone
-   Role
-   Senha (apenas para criação ou para alteração)
-   Tenant ID

### Remoção de Usuários

Cada usuário na lista pode ser removido através de um botão de exclusão.

## Endpoints da API

As seguintes rotas da API foram criadas no backend para dar suporte a esta funcionalidade:

-   `GET /admin/users`: Lista todos os usuários.
-   `POST /admin/users`: Cria um novo usuário para um tenant específico.
-   `PUT /admin/users/:id`: Atualiza um usuário existente.
-   `DELETE /admin/users/:id`: Remove um usuário.

**Nota:** Todas as rotas de `/admin` exigem que o usuário esteja autenticado e tenha a role de `admin`.
