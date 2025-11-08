# Modelo de Dados – Plataforma de Gestão Local

## Visão Geral
- Banco primário: PostgreSQL 14+.
- Multi-tenant por coluna (`tenant_id`) nas tabelas compartilhadas.
- Convenções:
  - Chaves primárias `id` (UUID v4).
  - Timestamps auditáveis: `created_at`, `updated_at`, `deleted_at` (soft delete opcional).
  - Índices compostos `tenant_id + coluna_critica` para consultas frequentes.

## Entidades Principais
| Tabela | Descrição | Campos-chave |
| --- | --- | --- |
| `companies` | Empresas cadastradas. | `id`, `name`, `document`, `settings (jsonb)` |
| `users` | Usuários vinculados à empresa. | `tenant_id`, `role`, `email`, `password_hash`, `profile` |
| `clients` | Clientes finais. | `tenant_id`, `name`, `contact`, `notes`, `tags (jsonb)` |
| `professionals` | Colaboradores que executam serviços (barbeiros, vendedores). | `tenant_id`, `user_id` (opcional), `specialties` |
| `services` | Serviços ofertados (corte, coloração, consultoria). | `tenant_id`, `name`, `duration`, `price`, `category` |
| `products` | Produtos físicos (shampoos, roupas). | `tenant_id`, `sku`, `stock_qty`, `price`, `cost` |
| `inventory_movements` | Movimentações de estoque. | `tenant_id`, `product_id`, `type (+/-)`, `quantity`, `reason` |
| `bookings` | Agendamentos. | `tenant_id`, `client_id`, `professional_id`, `service_id`, `status`, `start_at`, `end_at`, `notes` |
| `sales_orders` | Pedidos/vendas. | `tenant_id`, `client_id`, `status`, `payment_method`, `total`, `discount` |
| `sales_items` | Itens da venda. | `tenant_id`, `order_id`, `item_type (service/product)`, `item_ref_id`, `quantity`, `unit_price` |
| `payments` | Pagamentos efetivados. | `tenant_id`, `order_id`, `method`, `amount`, `paid_at`, `pix_payload` |
| `audit_logs` | Eventos relevantes (login, alteração de permissões). | `tenant_id`, `entity`, `action`, `actor_id`, `metadata` |

## Relacionamentos
- `companies 1:N users`, `companies 1:N clients`, `companies 1:N bookings`, etc. via `tenant_id`.
- `users` podem referenciar `professionals` quando o colaborador precisa de acesso ao painel.
- `bookings` vinculam `clients`, `professionals` e `services`. Exclusão lógica mantém histórico.
- `sales_orders` podem referenciar `bookings` (campo `booking_id` opcional) para conciliar atendimento → venda.
- `sales_items` usam `item_type` + `item_ref_id` para permitir mix entre serviços e produtos.
- `inventory_movements` relacionam `sales_items` (via `order_id`) para baixar estoque automaticamente.

## Regras e Restrições
- `tenant_id` obrigatório em todas as tabelas exceto `companies` e `audit_logs` globais.
- Unique indexes:
  - `users(tenant_id, email)` (case insensitive).
  - `services(tenant_id, name)`, `products(tenant_id, sku)`.
  - `bookings(tenant_id, professional_id, start_at)` para evitar conflitos.
- `bookings.status`: `pending`, `confirmed`, `done`, `canceled`. Transições validadas na camada de serviço.
- `sales_orders.status`: `draft`, `confirmed`, `paid`, `canceled`.
- `payments.method`: `cash`, `debit`, `credit`, `pix`, `transfer`.
- `inventory_movements.type`: `in`, `out`, `adjustment`.

## Diagrama ER (texto)
```
companies (id, name, document, settings)
  ├─< users (id, tenant_id, role, email, password_hash)
  ├─< clients (id, tenant_id, ...)
  ├─< professionals (id, tenant_id, user_id)
  ├─< services (id, tenant_id, ...)
  ├─< products (id, tenant_id, ...)
  ├─< bookings (id, tenant_id, client_id, professional_id, service_id)
        └─ sales_orders (id, tenant_id, client_id, booking_id)
              └─< sales_items (id, tenant_id, order_id, item_ref_id)
                    └─ products/services (referência indireta)
  └─< inventory_movements (id, tenant_id, product_id, order_id?)
```

## Migrations & Convenções
- Ferramenta: golang-migrate (scripts versionados em `backend/migrations`).
- Estrutura:
  - `0001_init.sql`: companies, users, clients, services, products.
  - `0002_agenda.sql`: professionals, bookings, regras de disponibilidade (tabela `availability_rules`).
  - `0003_sales.sql`: sales_orders, sales_items, payments.
  - `0004_inventory.sql`: inventory_movements, triggers de atualização de estoque.
- Naming:
  - Colunas snake_case.
  - FKs `fk_<tabela>_<coluna>`.
  - Índices `idx_<tabela>_<coluna(s)>`.
- Triggers úteis:
  - Atualizar `updated_at` em mudanças.
  - Reduzir estoque após `sales_orders` confirmados, reverter em cancelamentos.

## Dados Derivados & Relatórios
- Views materializadas:
  - `vw_daily_sales` (total por dia/profissional).
  - `vw_inventory_status` (estoque atual + nível crítico).
- Campos agregados guardados em `companies.settings` (jsonb) para thresholds.
- Possível uso de `timescaledb` ou partições para `audit_logs` se volume crescer.

## Próximos Passos
1. Detalhar atributos opcionais/obrigatórios em um dicionário de dados (coluna, tipo, descrição).
2. Criar diagramas visuais (dbdiagram.io/DrawSQL) e anexar imagem/PDF na pasta `docs/`.
3. Especificar regras de consistência para disponibilidade de profissionais (intervalos, folgas).
4. Definir seeds iniciais para ambientes de desenvolvimento (empresa demo, usuários, serviços).
