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
3. Hooks ficam em `tests/dredd/hooks/`. O arquivo `basic-flow.js` valida `signup -> login -> refresh -> companies/me` e continua o fluxo cobrindo o CRUD completo de usuários e clientes usando apenas dados criados durante o teste. Operações ainda não tratadas seguem marcadas com `skip` para evitar falsos negativos.
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
- Antes de rodar o contrato, execute `DATABASE_URL=postgres://gestao:gestao@localhost:5432/gestao_updev?sslmode=disable make -C backend migrate` para aplicar as migrations e, se necessário, `make -C backend seed` para garantir dados mínimos em ambientes de desenvolvimento compartilhados (o fluxo padrão do hook já cria uma empresa própria).
- Endpoints autenticados devem depender apenas de dados criados pelos próprios testes (hooks) ou pelo seed oficial (`cmd/seed`); evite dependências externas.
- Para rotas com efeitos colaterais (ex.: vendas, estoque), prefira criar registros específicos e, quando necessário, limpá-los ao fim do teste.
- Documente cenários não testáveis (ex.: callbacks externos) no README dos testes para evitar falsos negativos.

## Próximos Passos
1. Expandir os hooks para cobrir catálogos de serviços/profissionais e fluxos financeiros (estoque, vendas e pagamentos).
2. Adicionar job de contrato ao pipeline principal após o backend estar disponível.
3. Expor resultados dos testes em badges ou comentários automáticos em PRs.
