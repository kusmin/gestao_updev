# Backoffice - Gerenciamento de Vendas

Este documento descreve a funcionalidade de gerenciamento de vendas no backoffice da plataforma Gestão UpDev.

## Visão Geral

A página de gerenciamento de vendas permite que os administradores da plataforma visualizem, criem, editem e removam ordens de venda de todos os tenants.

## Funcionalidades

### Listagem de Vendas

A página principal exibe uma lista de todas as ordens de venda cadastradas na plataforma. A lista inclui as seguintes informações:

-   **Cliente:** O ID do cliente associado à venda.
-   **Status:** O status da ordem de venda (e.g., `draft`, `pending`, `paid`, `canceled`).
-   **Total:** O valor total da ordem de venda.

### Criação e Edição de Vendas

É possível criar uma nova ordem de venda ou editar uma existente através de um formulário modal. O formulário contém os seguintes campos:

-   Client ID
-   Booking ID (opcional)
-   Desconto
-   Notas
-   Itens (JSON - para produtos e serviços)
-   Tenant ID

### Remoção de Vendas

Cada ordem de venda na lista pode ser removida através de um botão de exclusão.

## Endpoints da API

As seguintes rotas da API foram criadas no backend para dar suporte a esta funcionalidade:

-   `GET /admin/sales/orders`: Lista todas as ordens de venda.
-   `POST /admin/sales/orders`: Cria uma nova ordem de venda para um tenant específico.
-   `PUT /admin/sales/orders/:id`: Atualiza uma ordem de venda existente.
-   `DELETE /admin/sales/orders/:id`: Remove uma ordem de venda.

**Nota:** Todas as rotas de `/admin` exigem que o usuário esteja autenticado e tenha a role de `admin`.
