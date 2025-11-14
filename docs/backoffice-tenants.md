# Backoffice - Gerenciamento de Tenants

Este documento descreve a funcionalidade de gerenciamento de tenants no backoffice da plataforma Gestão UpDev.

## Visão Geral

A página de gerenciamento de tenants permite que os administradores da plataforma visualizem, criem, editem e removam tenants (empresas).

## Funcionalidades

### Listagem de Tenants

A página principal exibe uma lista de todos os tenants cadastrados na plataforma. A lista inclui as seguintes informações:

-   **Nome:** O nome do tenant.
-   **Documento:** O documento (CNPJ/CPF) do tenant.
-   **Telefone:** O telefone de contato do tenant.
-   **Email:** O email de contato do tenant.

### Criação e Edição de Tenants

É possível criar um novo tenant ou editar um existente através de um formulário modal. O formulário contém os seguintes campos:

-   Nome
-   Documento
-   Telefone
-   Email

### Remoção de Tenants

Cada tenant na lista pode ser removido através de um botão de exclusão.

## Endpoints da API

As seguintes rotas da API foram criadas no backend para dar suporte a esta funcionalidade:

-   `GET /admin/tenants`: Lista todos os tenants.
-   `GET /admin/tenants/:id`: Busca um tenant pelo ID.
-   `POST /admin/tenants`: Cria um novo tenant.
-   `PUT /admin/tenants/:id`: Atualiza um tenant existente.
-   `DELETE /admin/tenants/:id`: Remove um tenant.

**Nota:** Todas as rotas de `/admin` exigem que o usuário esteja autenticado e tenha a role de `admin`.
