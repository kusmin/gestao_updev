# Backoffice - Dashboard

Este documento descreve a funcionalidade do Dashboard no backoffice da plataforma Gestão UpDev.

## Visão Geral

O Dashboard oferece uma visão geral das principais métricas da plataforma, permitindo que os administradores monitorem o estado geral e o desempenho.

## Métricas Exibidas

O Dashboard exibe as seguintes métricas agregadas de todos os tenants:

-   **Total de Tenants:** Número total de empresas/tenants registrados na plataforma.
-   **Total de Usuários:** Número total de usuários cadastrados em todos os tenants.
-   **Total de Clientes:** Número total de clientes registrados em todos os tenants.
-   **Total de Produtos:** Número total de produtos cadastrados em todos os tenants.
-   **Total de Serviços:** Número total de serviços cadastrados em todos os tenants.
-   **Total de Agendamentos:** Número total de agendamentos realizados em todos os tenants.
-   **Receita Total:** A receita total gerada por todas as vendas na plataforma.

## Endpoints da API

A seguinte rota da API foi criada no backend para dar suporte a esta funcionalidade:

-   `GET /admin/dashboard/metrics`: Retorna as métricas gerais da plataforma.

**Nota:** Todas as rotas de `/admin` exigem que o usuário esteja autenticado e tenha a role de `admin`.
