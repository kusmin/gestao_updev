# Changelog da API

Histórico das mudanças do contrato OpenAPI (`docs/api.yaml`). Cada versão deve ser publicada junto com uma tag Git correspondente e referências no README.

## [0.2.0] - 2025-11-08
- Expansão completa do spec cobrindo empresas, usuários, clientes, profissionais, serviços, produtos, estoque, agenda, vendas, pagamentos e dashboard.
- Adicionados schemas reutilizáveis (Company, User, Client, Booking, SalesOrder, etc.) e suportes para multi-tenant (`X-Tenant-ID`).
- Pipeline passa a gerar artefatos HTML e publicar no GitHub Pages.
- Adicionado endpoint público `/healthz` (schema `HealthStatus`) utilizado pelos clients para checar disponibilidade do backend.

## [0.1.0] - 2025-11-08
- Criação inicial do spec com autenticação e CRUD básico de clientes.

---

### Processo de Versionamento
1. Atualize `info.version` em `docs/api.yaml`.
2. Adicione nova seção no changelog com data e descrição do impacto (breaking changes / novos recursos).
3. Abra PR com as mudanças e, após merge, crie uma tag `api-vX.Y.Z` apontando para `main`.
4. Publique release no GitHub anexando o HTML gerado pelo workflow (download do artefato ou link do GitHub Pages).
