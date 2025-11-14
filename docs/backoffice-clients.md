# Backoffice - Gerenciamento de Clientes

Este documento descreve a funcionalidade de gerenciamento de clientes no backoffice da plataforma Gestão UpDev.

## Visão Geral

A página de gerenciamento de clientes permite que os administradores da plataforma visualizem, criem, editem e removam clientes de todos os tenants.

## Funcionalidades

### Listagem de Clientes

A página principal exibe uma lista de todos os clientes cadastrados na plataforma. A lista inclui as seguintes informações:

-   **Nome:** O nome do cliente.
-   **Email:** O email do cliente.
-   **Telefone:** O telefone de contato do cliente.
-   **Tenant ID:** O ID do tenant ao qual o cliente pertence.

### Criação e Edição de Clientes

É possível criar um novo cliente ou editar um existente através de um formulário modal. O formulário contém os seguintes campos:

-   Nome
-   Email
-   Telefone
-   Tenant ID

### Remoção de Clientes

Cada cliente na lista pode ser removido através de um botão de exclusão.

## Endpoints da API

As seguintes rotas da API foram criadas no backend para dar suporte a esta funcionalidade:

-   `GET /admin/clients`: Lista todos os clientes.
-   `POST /admin/clients`: Cria um novo cliente para um tenant específico.
-   `PUT /admin/clients/:id`: Atualiza um cliente existente.
-   `DELETE /admin/clients/:id`: Remove um cliente.

**Nota:** Todas as rotas de `/admin` exigem que o usuário esteja autenticado e tenha a role de `admin`.
