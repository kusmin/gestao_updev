# Testes de Contrato (API)

Objetivo: garantir que as respostas reais do backend estejam alinhadas com o spec OpenAPI (`docs/api.yaml`), evitando regressões entre squads.

## Ferramentas Recomendadas
- **Dredd**: executa o spec diretamente contra o serviço HTTP.
- **Newman/Postman**: alternativa para cenários com coleções já existentes (exportadas do OpenAPI).
- **Prism (Mock Server)**: útil para validar o frontend enquanto o backend não está pronto.

## Fluxo com Dredd
1. Instale globalmente:
   ```bash
   npm install -g dredd
   ```
2. O repositório já possui `tests/dredd/dredd.yml` apontando para o spec e para `make backend-run` (que sobe o servidor Go).  
3. Hooks ficam em `tests/dredd/hooks/`. Por padrão só o `GET /healthz` roda; demais endpoints são marcados como `skip` até serem implementados.
4. Execute localmente:
   ```bash
   make api-contract-test
   ```

## Integração no CI
- Adicionar job `api-contract-test` que:
  1. Sobe o backend (`make backend-run`, controlado automaticamente pelo Dredd).
  2. Roda migrations/fixtures.
  3. Executa `dredd`.
- Fails gate: PR só pode ser mergeado se todos os cenários estiverem verdes.

## Convenções
- Endpoints que dependem de autenticação devem usar tokens válidos gerados durante os hooks.
- Para rotas com efeitos colaterais (ex.: criação de venda), limpe ou isole os dados ao fim do teste.
- Documente cenários não testáveis (ex.: callbacks externos) em `README` dos testes para evitar falsos negativos.

## Próximos Passos
1. Configurar `tests/` com fixtures e hooks básicos (signup/login, CRUD simples).
2. Adicionar job de contrato ao pipeline principal após o backend estar disponível.
3. Expor resultados dos testes em badges ou comentários automáticos em PRs.
