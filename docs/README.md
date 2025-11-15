# Documentação Inicial – Plataforma de Gestão Local

## Visão Geral
- Plataforma SaaS para barbearias, lojas de roupa e outros comércios locais.
- Foco em centralizar clientes, agenda, estoque e vendas em um painel único.
- Cada empresa possui um espaço dedicado (multi-tenant leve).

## Objetivos do MVP
- Criar cadastro de empresas, usuários e perfis (admin, gerente, colaborador).
- Disponibilizar agenda de serviços e registro de atendimento/vendas.
- Controlar catálogo simples de produtos/serviços com estoque básico.
- Emitir relatórios essenciais (faturamento, agenda do dia, estoque crítico).

## Público-Alvo e Persona
- Pequenos empresários com baixa maturidade digital.
- Precisam de ferramenta rápida, responsiva e acessível em desktop/mobile.
- Valor decisivo: redução de custos operacionais e visão consolidada do negócio.

## Módulos Principais
1. **Empresas & Usuários** – onboarding guiado, gestão de permissões simples.
2. **Clientes** – histórico de visitas/compras, preferências e contatos.
3. **Agenda** – agendamentos por profissional, bloqueios de horário, lembretes.
4. **Produtos & Serviços** – catálogo, combos e pacotes, controle de estoque.
5. **Vendas & Pagamentos** – pedidos rápidos, recibos, suporte a Pix manual.
6. **Relatórios** – visão diária/semanal, ranking de serviços/profissionais.

## Arquitetura Inicial
- **Backend**: Go 1.25+ com Gin ou Fiber; camada de serviço com validações e regras de negócio.
- **Persistência**: PostgreSQL (multi-schema ou coluna `tenant_id`); migrations com Atlas ou golang-migrate; uso de SQLC para queries seguras.
- **Autenticação**: JWT + refresh tokens, armazenamento seguro de senhas (bcrypt/argon2).
- **Frontend**: React + Vite com UI responsiva (SPA); consumo via REST/JSON e suporte futuro a WebSockets para agenda em tempo real.
- **Infraestrutura**: Containers Docker, deploy inicial em plataforma PaaS econômica (Railway/Render/Fly.io); monitoramento com Prometheus + Grafana leves.
- **Observabilidade**: logs estruturados (zap/logrus), tracing via OpenTelemetry, health checks expostos.

## Escopo MVP (Sprint 0–2)
- Fluxo de cadastro/login + criação de empresa.
- CRUD de clientes, serviços e profissionais.
- Agenda diária/semanal com filtros básicos e notificações por e-mail.
- Registro de vendas simples (seleciona cliente, itens, forma de pagamento).
- Dashboard inicial com KPIs diários e alertas de estoque.

## Roadmap de Evolução
1. **Fase 1 (MVP – 5/6 semanas)**
   - Core backend Go, autenticação, empresas/usuários.
   - Agenda, clientes, serviços e vendas básicas.
   - Dashboard com métricas essenciais.
2. **Fase 2 (6/8 semanas)**
   - Estoque avançado, relatórios customizáveis, exportações.
   - Melhoria UX, atalhos mobile, templates de promoções.
   - Integrações de calendário externo (Google) e webhooks.
3. **Fase 3 (8+ semanas)**
   - Pagamentos Pix/Boleto via parceiros, automações (WhatsApp/SMS).
   - App mobile híbrido para colaboradores.
   - BI leve com indicadores preditivos e segmentação de clientes.

## Próximos Passos
- Validar tecnologia frontend e componentes UI.
- Definir modelo de dados físico (diagramas) e convenções de multi-tenancy.
- Preparar repositório backend (Go module, configs, migrações iniciais).
- Padrões de código: `docs/padroes-codigo.md`.
- Multi-tenancy: `docs/multi-tenancy.md`.
