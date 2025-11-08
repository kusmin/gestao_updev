# Referência de API – Plataforma de Gestão Local

## Objetivo
Documentar endpoints REST do MVP com foco em payloads, respostas e códigos de status para integração entre frontend e backend.

> **Nota**: versão inicial textual; futura migração para OpenAPI (`docs/api.yaml`).

## Autenticação
- **POST** `/auth/signup`
  - Request:
    ```json
    {
      "company": {"name": "Barbearia X", "document": "12.345.678/0001-99"},
      "user": {"name": "João", "email": "joao@x.com", "password": "Senha@123"},
      "phone": "+55 11 99999-0000"
    }
    ```
  - Response `201`:
    ```json
    {"data": {"tenant_id": "uuid", "user_id": "uuid"}, "error": null}
    ```
- **POST** `/auth/login`
  - Request: `{"email": "joao@x.com", "password": "Senha@123"}`
  - Response `200`: `{"data": {"access_token": "...", "refresh_token": "...", "expires_in": 3600}}`
- **POST** `/auth/refresh`
  - Request: `{"refresh_token": "..."}`.
  - Response `200`: novos tokens.

## Saúde
- **GET** `/healthz`
  - Sem autenticação.
  - Response `200`: `{"data": {"status": "ok", "env": "production"}, "error": null}`.

## Empresas e Usuários
- **GET** `/companies/me`
  - Headers: `Authorization: Bearer`, `X-Tenant-ID`.
  - Response `200`: dados da empresa + configurações.
- **PUT** `/companies/me`
  - Body: campos atualizáveis (nome, horários, timezone).
  - Response `200`: empresa atualizada.
- **POST** `/users`
  - Body: `{"name": "...", "email": "...", "role": "manager", "phone": "...", "password": "..."}`
  - Response `201`: usuário criado.
- **GET** `/users`
  - Query: `role`, `page`, `per_page`.
  - Response `200`: lista paginada.
- **PATCH** `/users/{id}`
  - Body parcial (role, ativo, phone).
  - Response `200`.
- **DELETE** `/users/{id}`
  - Soft delete (marca `deleted_at`).
  - Response `204`.

## Clientes
- **POST** `/clients`
  - Body: `{"name": "...", "phone": "...", "email": "...", "notes": ""}`
  - Response `201`.
- **GET** `/clients`
  - Query: `search`, `page`, `tags`.
  - Response `200`: lista + `meta.pagination`.
- **GET** `/clients/{id}`
  - Response `200` com histórico resumido.
- **PUT** `/clients/{id}`
  - Atualiza dados e preferências.
- **DELETE** `/clients/{id}`
  - Soft delete (motivo opcional).

## Agenda
- **GET** `/professionals`
  - Lista profissionais disponíveis (nome, especialidades, capacidade).
- **POST** `/bookings`
  - Body:
    ```json
    {
      "client_id": "uuid",
      "professional_id": "uuid",
      "service_id": "uuid",
      "start_at": "2024-04-01T10:00:00-03:00",
      "status": "confirmed",
      "notes": ""
    }
    ```
  - Response `201`: booking criado.
- **GET** `/bookings`
  - Query: `date`, `professional_id`, `status`.
  - Response `200`: lista ordenada por `start_at`.
- **PATCH** `/bookings/{id}`
  - Campos: `status`, `notes`, `start_at`, `end_at`.
- **POST** `/bookings/{id}/cancel`
  - Body: `{"reason": "Cliente não compareceu"}`; Response `200`.

## Serviços e Produtos
- **GET/POST/PUT/DELETE** `/services`
  - Campos: `name`, `duration_minutes`, `price`, `category`, `description`.
- **GET/POST/PUT/DELETE** `/products`
  - Campos: `name`, `sku`, `price`, `cost`, `stock_qty`, `min_stock`.
- **POST** `/inventory/movements`
  - Body: `{"product_id": "uuid", "type": "in|out|adjustment", "quantity": 3, "reason": "Reposição"}`
  - Response `201`.
- **GET** `/inventory/movements`
  - Filtros: `product_id`, `type`, `date_range`.

## Vendas
- **POST** `/sales/orders`
  - Body:
    ```json
    {
      "client_id": "uuid",
      "booking_id": "uuid?",
      "items": [
        {"type": "service", "ref_id": "uuid", "quantity": 1, "unit_price": 50},
        {"type": "product", "ref_id": "uuid", "quantity": 2, "unit_price": 30}
      ],
      "discount": 5,
      "notes": ""
    }
    ```
  - Response `201`: `order_id`.
- **GET** `/sales/orders`
  - Query: `status`, `date`, `client_id`.
- **PATCH** `/sales/orders/{id}`
  - Atualiza status (`confirmed`, `canceled`), notas, itens (restrito).
- **POST** `/sales/orders/{id}/payments`
  - Body: `{"method": "pix", "amount": 120, "paid_at": "..." }`
  - Response `201`.
- **GET** `/payments`
  - Filtros: `method`, `date_range`.

## Dashboard & Relatórios
- **GET** `/dashboard/daily`
  - Query: `date`, `professional_id` opcional.
  - Response: KPIs (agendamentos, atendimentos, receita, top serviços).
- **GET** `/reports/stock`
  - Retorna produtos abaixo do mínimo + export CSV (header `Accept: text/csv`).

## Convenções de Resposta
```json
{
  "data": {...},
  "meta": {"pagination": {"page": 1, "per_page": 20, "total": 53}},
  "error": null
}
```
- Erros seguem:
```json
{
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Campo email é obrigatório",
    "details": [{"field": "email", "reason": "required"}]
  }
}
```

## Status Codes
- `200` sucesso padrão.
- `201` recurso criado.
- `204` sem conteúdo (delete).
- `400` validação inválida.
- `401` token inválido/expirado.
- `403` permissão insuficiente.
- `404` recurso inexistente.
- `409` conflito (agendamento duplicado).
- `422` regra de negócio (estoque insuficiente).
- `500` erro interno.

## Próximos Passos
1. Migrar para OpenAPI 3.1 em `docs/api.yaml` com schemas reutilizáveis.
2. Publicar documentação interativa (Stoplight/Swagger UI) para o frontend.
3. Automatizar testes de contrato (Dredd/Postman) no pipeline CI.
