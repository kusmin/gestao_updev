# Plano de Expansão do Frontend

Este documento descreve o estado atual do frontend e o plano para adicionar novas funcionalidades, com base nos documentos de arquitetura e planejamento existentes.

## Estado Atual

O frontend está estruturado com React, Vite, TypeScript e Material-UI, seguindo as diretrizes de `FRONTEND_ARCHITECTURE.md`. O roteamento é gerenciado pelo `react-router-dom`, conforme o `frontend-routing-plan.md`.

Atualmente, as seguintes seções estão implementadas (ou em fase de implementação):

- **Autenticação**:
  - `LoginPage.tsx`
  - `SignupPage.tsx`
- **Dashboard**:
  - `DashboardPage.tsx`
- **Clientes**:
  - `ClientListPage.tsx`
  - `ClientForm.tsx`
- **Agendamentos**:
  - `AppointmentListPage.tsx` (recém-adicionada)
  - `AppointmentForm.tsx` (recém-adicionada)

## Plano de Expansão

O objetivo é completar as funcionalidades principais da plataforma, conforme definido no `frontend-pages-plan.md`. As próximas seções a serem desenvolvidas são **Inventário** e **Vendas**.

### 1. Módulo de Inventário (Produtos)

Este módulo permitirá o gerenciamento de produtos e controle de estoque.

- **Páginas a serem criadas**:
  - `frontend/src/pages/products/ProductListPage.tsx`: Listagem de produtos com busca e filtros.
  - `frontend/src/pages/products/ProductForm.tsx`: Formulário para criar e editar produtos (nome, preço, quantidade em estoque).
  - `frontend/src/pages/products/InventoryMovementPage.tsx`: Uma página para registrar entradas e saídas de estoque.

- **Passos**:
  1. Criar o diretório `frontend/src/pages/products`.
  2. Implementar `ProductListPage.tsx` e `ProductForm.tsx`, seguindo o padrão de `ClientListPage.tsx`.
  3. Adicionar a rota `/products` em `AppRouter.tsx`.
  4. Adicionar o link "Produtos" na navegação principal em `AppLayout.tsx`.
  5. (Opcional) Implementar a página de movimentação de estoque.

### 2. Módulo de Vendas

Este módulo permitirá o registro de vendas, associando produtos e clientes.

- **Páginas a serem criadas**:
  - `frontend/src/pages/sales/SaleListPage.tsx`: Listagem de vendas realizadas.
  - `frontend/src/pages/sales/SaleForm.tsx`: Formulário para registrar uma nova venda, permitindo selecionar cliente e adicionar produtos.

- **Passos**:
  1. Criar o diretório `frontend/src/pages/sales`.
  2. Implementar `SaleListPage.tsx` e `SaleForm.tsx`. O formulário de vendas será mais complexo, envolvendo a seleção de múltiplos produtos.
  3. Adicionar a rota `/sales` em `AppRouter.tsx`.
  4. Adicionar o link "Vendas" na navegação principal em `AppLayout.tsx`.

## Integração com o Backend

Todas as novas páginas serão inicialmente implementadas com dados mocados (placeholders), assim como foi feito para a página de agendamentos.

A integração com o backend exigirá:
1.  A definição dos modelos e endpoints da API no `docs/api.yaml`.
2.  A implementação dos serviços e repositórios no backend (Go).
3.  A implementação das chamadas de API no frontend, substituindo os dados mocados.

Este plano permite um desenvolvimento incremental e paralelo do frontend, enquanto o backend é desenvolvido.
