# Auditoria de Segurança

## Resumo Executivo
A revisão identificou dois grupos principais de riscos: (1) parâmetros sensíveis com valores padrão fracos que podem ir para produção por engano e (2) falhas de validação de multi-tenancy que permitem correlacionar registros de outros locatários ao criar reservas, vendas, pagamentos e movimentos de estoque. Ambos comprometem a confidencialidade e a integridade dos dados dos clientes e exigem correções priorizadas.

## Vulnerabilidades

### 1. Segredos JWT e credenciais de banco com valores padrão fracos (Severidade: Alta)
A configuração carrega valores padrão para `DATABASE_URL`, `JWT_ACCESS_SECRET` e `JWT_REFRESH_SECRET`. Caso as variáveis de ambiente não estejam configuradas em produção, o serviço iniciará com esses valores triviais, expondo os tokens JWT e o banco de dados a ataques de força bruta ou acesso não autorizado. 【F:backend/internal/config/config.go†L10-L35】

**Mitigação sugerida:** Tornar obrigatória a definição explícita desses segredos em produção (por exemplo, removendo `envDefault` para segredos) e aplicar validações que derrubem o processo quando o valor permanecer nos padrões de desenvolvimento.

### 2. Falta de validação de pertencimento entre tenants ao criar registros relacionados (Severidade: Alta)
Diversos fluxos aceitam IDs de entidades relacionadas, mas não confirmam que esses IDs pertencem ao mesmo tenant do usuário autenticado. Isso permite que um atacante relacione, por tentativa ou vazamento de UUIDs, registros de outros locatários aos seus próprios dados, violando o isolamento multi-tenant:

- `CreateBooking` usa `ClientID` e `ProfessionalID` diretamente sem checar o tenant dos registros, garantindo apenas que o `ServiceID` exista para o tenant local. 【F:backend/internal/service/bookings.go†L66-L104】
- `CreateSalesOrder` aceita `ClientID`, `BookingID` e itens (`RefID`) sem confirmar a propriedade pelo tenant. O método `AddPayment` também permite vincular pagamentos a pedidos de outros tenants apenas conhecendo o UUID. 【F:backend/internal/service/sales.go†L84-L200】
- `CreateInventoryMovement` grava `ProductID` (e opcionalmente `OrderID`) sem validar a autoria, possibilitando alterar estoque de produtos de outro tenant. 【F:backend/internal/service/inventory.go†L56-L83】

O modelo de dados reforça o problema: as FKs exigem que a referência exista, mas não que pertença ao mesmo tenant, permitindo associações cruzadas válidas do ponto de vista do banco. 【F:backend/migrations/0002_agenda.up.sql†L1-L62】【F:backend/migrations/0003_sales.up.sql†L1-L61】

**Mitigação sugerida:** Em cada operação, buscar explicitamente os registros relacionados filtrando por `tenant_id` antes de prosseguir (e rejeitar quando não pertencerem), além de considerar restrições compostas ou triggers no banco para reforçar o vínculo entre `tenant_id` da entidade e de suas dependências.

**Mitigação implementada:** A migração `0005_multitenant_constraints` cria índices únicos em `(tenant_id, id)` e adiciona FKs compostas para `bookings`, `sales_orders`, `sales_items`, `payments` e `inventory_movements`, garantindo que o banco rejeite associações cruzadas mesmo que um UUID de outro tenant seja reutilizado.

## Próximos Passos Recomendados
1. Ajustar o carregamento de configuração para falhar quando segredos permanecerem com os defaults de desenvolvimento.
2. Revisar todos os casos de uso que recebem IDs de entidades relacionadas, garantindo que cada busca use filtros por `tenant_id` antes de criar ou atualizar registros.
3. (Concluído em `0005_multitenant_constraints`) Complementar as proteções na camada de banco (FKs compostas ou constraints) para impedir associações cruzadas mesmo em caso de regressões na aplicação.
4. Adicionar testes de integração cobrindo cenários de acesso cruzado entre tenants e validações de configuração obrigatória.
