# Backoffice - Gerenciamento de Produtos

Este documento descreve a funcionalidade de gerenciamento de produtos no backoffice da plataforma Gestão UpDev.

## Visão Geral

A página de gerenciamento de produtos permite que os administradores da plataforma visualizem, criem, editem e removam produtos de todos os tenants.

## Funcionalidades

### Listagem de Produtos

A página principal exibe uma lista de todos os produtos cadastrados na plataforma. A lista inclui as seguintes informações:

-   **Nome:** O nome do produto.
-   **SKU:** O SKU do produto.
-   **Preço:** O preço do produto.
-   **Estoque:** A quantidade em estoque do produto.

### Criação e Edição de Produtos

É possível criar um novo produto ou editar um existente através de um formulário modal. O formulário contém os seguintes campos:

-   Nome
-   SKU
-   Preço
-   Custo
-   Estoque
-   Estoque Mínimo
-   Descrição
-   Tenant ID

### Remoção de Produtos

Cada produto na lista pode ser removido através de um botão de exclusão.

## Endpoints da API

As seguintes rotas da API foram criadas no backend para dar suporte a esta funcionalidade:

-   `GET /admin/products`: Lista todos os produtos.
-   `POST /admin/products`: Cria um novo produto para um tenant específico.
-   `PUT /admin/products/:id`: Atualiza um produto existente.
-   `DELETE /admin/products/:id`: Remove um produto.

**Nota:** Todas as rotas de `/admin` exigem que o usuário esteja autenticado e tenha a role de `admin`.
