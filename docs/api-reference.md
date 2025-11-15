# Referência de API – Plataforma de Gestão Local

## Objetivo
Documentar endpoints REST do MVP com foco em payloads, respostas e códigos de status para integração entre frontend e backend.

> **Nota**: A especificação OpenAPI (`docs/api.yaml`) é a fonte definitiva do contrato da API. Este documento serve como um guia complementar e de referência rápida.

## Autenticação
- **POST** `/v1/auth/signup`
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
- **POST** `/v1/auth/login`
  - Request: `{"email": "joao@x.com", "password": "Senha@123"}`
  - Response `200`: `{"data": {"access_token": "...", "refresh_token": "...", "expires_in": 3600}}`
- **POST** `/v1/auth/refresh`
  - Request: `{"refresh_token": "..."}`.
  - Response `200`: novos tokens.

## Saúde
- **GET** `/v1/healthz`
  - Sem autenticação.
  - Response `200`: `{"data": {"status": "ok", "env": "production"}, "error": null}`.

## Empresas e Usuários
- **GET** `/v1/companies/me`
  - Headers: `Authorization: Bearer`, `X-Tenant-ID`.
  - Response `200`: dados da empresa + configurações.
- **PUT** `/v1/companies/me`
  - Body: campos atualizáveis (nome, horários, timezone).
  - Response `200`: empresa atualizada.
- **POST** `/v1/users`
  - Body: `{"name": "...", "email": "...", "role": "manager", "phone": "...", "password": "..."}`
  - Response `201`: usuário criado.
- **GET** `/v1/users`
  - Query: `role`, `page`, `per_page`.
  - Response `200`: lista paginada.
- **PATCH** `/v1/users/{id}`
  - Body parcial (role, ativo, phone).
  - Response `200`.
- **DELETE** `/v1/users/{id}`
  - Soft delete (marca `deleted_at`).
  - Response `204`.

## Clientes
- **POST** `/v1/clients`
  - Body: `{"name": "...", "phone": "...", "email": "...", "notes": ""}`
  - Response `201`.
- **GET** `/v1/clients`
  - Query: `search`, `page`, `tags`.
  - Response `200`: lista + `meta.pagination`.
- **GET** `/v1/clients/{id}`
  - Response `200` com histórico resumido.
- **PUT** `/v1/clients/{id}`
  - Atualiza dados e preferências.
- **DELETE** `/v1/clients/{id}`
  - Soft delete (motivo opcional).

## Agenda
- **GET** `/v1/professionals`
  - Lista profissionais disponíveis (nome, especialidades, capacidade).
- **POST** `/v1/bookings`
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
- **GET** `/v1/bookings`
  - Query: `date`, `professional_id`, `status`.
  - Response `200`: lista ordenada por `start_at`.
- **PATCH** `/v1/bookings/{id}`
  - Campos: `status`, `notes`, `start_at`, `end_at`.
- **POST** `/v1/bookings/{id}/cancel`
  - Body: `{"reason": "Cliente não compareceu"}`; Response `200`.

## Serviços e Produtos
- **GET/POST/PUT/DELETE** `/v1/services`
  - Campos: `name`, `duration_minutes`, `price`, `category`, `description`.
- **GET/POST/PUT/DELETE** `/v1/products`
  - Campos: `name`, `sku`, `price`, `cost`, `stock_qty`, `min_stock`.
- **POST** `/v1/inventory/movements`
  - Body: `{"product_id": "uuid", "type": "in|out|adjustment", "quantity": 3, "reason": "Reposição"}`
  - Response `201`.
- **GET** `/v1/inventory/movements`
  - Filtros: `product_id`, `type`, `date_range`.

## Vendas
- **POST** `/v1/sales/orders`
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
- **GET** `/v1/sales/orders`
  - Query: `status`, `date`, `client_id`.
- **PATCH** `/v1/sales/orders/{id}`
  - Atualiza status (`confirmed`, `canceled`), notas, itens (restrito).
- **POST** `/v1/sales/orders/{id}/payments`
  - Body: `{"method": "pix", "amount": 120, "paid_at": "..." }`
  - Response `201`.
- **GET** `/v1/payments`
  - Filtros: `method`, `date_range`.

## Dashboard & Relatórios
- **GET** `/v1/dashboard/daily`
  - Query: `date`, `professional_id` opcional.
  - Response: KPIs (agendamentos, atendimentos, receita, top serviços).
- **GET** `/v1/reports/stock`
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
1. Publicar documentação interativa (Stoplight/Swagger UI) para o frontend.
2. Automatizar testes de contrato (Dredd/Postman) no pipeline CI.
3. Distribuir o spec 3.1 versionado (artefatos/SDK) para squads e integrações externas.
