# Backoffice - Gerenciamento de Agendamentos

Este documento descreve a funcionalidade de gerenciamento de agendamentos no backoffice da plataforma Gestão UpDev.

## Visão Geral

A página de gerenciamento de agendamentos permite que os administradores da plataforma visualizem, criem, editem e removam agendamentos de todos os tenants.

## Funcionalidades

### Listagem de Agendamentos

A página principal exibe uma lista de todos os agendamentos cadastrados na plataforma. A lista inclui as seguintes informações:

-   **Cliente:** O ID do cliente.
-   **Profissional:** O ID do profissional.
-   **Serviço:** O ID do serviço.
-   **Início:** A data e hora de início do agendamento.
-   **Fim:** A data e hora de fim do agendamento.
-   **Status:** O status do agendamento (e.g., `pending`, `confirmed`, `done`, `canceled`).

### Criação e Edição de Agendamentos

É possível criar um novo agendamento ou editar um existente através de um formulário modal. O formulário contém os seguintes campos:

-   Client ID
-   Professional ID
-   Service ID
-   Início
-   Fim
-   Status
-   Tenant ID

### Remoção de Agendamentos

Cada agendamento na lista pode ser removido através de um botão de exclusão.

## Endpoints da API

As seguintes rotas da API foram criadas no backend para dar suporte a esta funcionalidade:

-   `GET /admin/bookings`: Lista todos os agendamentos.
-   `POST /admin/bookings`: Cria um novo agendamento para um tenant específico.
-   `PUT /admin/bookings/:id`: Atualiza um agendamento existente.
-   `DELETE /admin/bookings/:id`: Remove um agendamento.

**Nota:** Todas as rotas de `/admin` exigem que o usuário esteja autenticado e tenha a role de `admin`.
