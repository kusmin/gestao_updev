# Backoffice - Gerenciamento de Serviços

Este documento descreve a funcionalidade de gerenciamento de serviços no backoffice da plataforma Gestão UpDev.

## Visão Geral

A página de gerenciamento de serviços permite que os administradores da plataforma visualizem, criem, editem e removam serviços de todos os tenants.

## Funcionalidades

### Listagem de Serviços

A página principal exibe uma lista de todos os serviços cadastrados na plataforma. A lista inclui as seguintes informações:

-   **Nome:** O nome do serviço.
-   **Duração (min):** A duração do serviço em minutos.
-   **Preço:** O preço do serviço.

### Criação e Edição de Serviços

É possível criar um novo serviço ou editar um existente através de um formulário modal. O formulário contém os seguintes campos:

-   Nome
-   Categoria
-   Descrição
-   Duração (minutos)
-   Preço
-   Cor
-   Tenant ID

### Remoção de Serviços

Cada serviço na lista pode ser removido através de um botão de exclusão.

## Endpoints da API

As seguintes rotas da API foram criadas no backend para dar suporte a esta funcionalidade:

-   `GET /admin/services`: Lista todos os serviços.
-   `POST /admin/services`: Cria um novo serviço para um tenant específico.
-   `PUT /admin/services/:id`: Atualiza um serviço existente.
-   `DELETE /admin/services/:id`: Remove um serviço.

**Nota:** Todas as rotas de `/admin` exigem que o usuário esteja autenticado e tenha a role de `admin`.
