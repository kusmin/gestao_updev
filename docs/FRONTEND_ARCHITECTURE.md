# Arquitetura do Frontend

Este documento detalha a arquitetura do frontend da plataforma, incluindo decisões de tecnologia, estrutura e padrões.

## Visão Geral

A arquitetura do frontend foi projetada para ser modular, escalável e de fácil manutenção, com foco na experiência do usuário.

## Frontend (React + Vite)

- **Framework/Bundler:** React 18 + Vite para entregar um SPA leve, com build rápido e deploy simplificado (assets estáticos).
- **UI Kit:** A ser definido. Opções incluem `Material-UI`, `Chakra UI` ou `Tailwind CSS` com `Headless UI`. Critérios de escolha: acessibilidade, velocidade de customização e cobertura de componentes.
- **Gerenciamento de Estado:** `React Query` para dados remotos + `Zustand`/Context para estados globais simples.
- **Comunicação com o Backend:** `fetch` (com wrappers tipados) ou `axios`. `WebSockets` (Socket.IO/Native WS) planejados para agenda em tempo real.
- **Estrutura sugerida:**
  ```
  /src
    /components
      /common
      /layout
    /hooks
    /lib        # apiClient, armazenamento local, helpers
    /pages      # rotas SPA (React Router futuramente)
    /styles
    /types      # tipos derivados do OpenAPI
  ```

## Backoffice (React + Vite com React-Admin)

O backoffice será desenvolvido utilizando o **React-Admin**, um framework para construção de painéis administrativos, que se integra com a mesma stack tecnológica do frontend principal para garantir consistência, reuso de conhecimento e, potencialmente, de componentes.

- **Framework:** React-Admin (construído sobre React 18 + Vite).
- **UI Kit:** **Material-UI (MUI)** (integrado ao React-Admin).
- **Gerenciamento de Estado e Data Fetching:** React Query (integrado ao React-Admin).
- **Comunicação com o Backend:** Através de Data Providers do React-Admin (ex: `ra-data-json-server` para prototipagem, ou customizado para a API Go).
- **Estrutura sugerida:** Seguirá a estrutura recomendada pelo React-Admin, com recursos definidos e componentes de lista, edição, criação, etc.

## Validação de Tecnologia e Componentes UI

- **Prova de Conceito (PoC):** Uma PoC será realizada para avaliar os UI Kits candidatos.
- **Critérios de Avaliação:**
  - Facilidade de customização.
  - Acessibilidade.
  - Performance.
  - Qualidade da documentação.
- **Componentes Prioritários para a PoC:**
  - Tabela de dados com paginação e filtros.
  - Formulário complexo com validação.
  - Componente de calendário/agenda.
  - Gráficos para o dashboard.
