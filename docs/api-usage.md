# Guia de Uso da API (OpenAPI)

## Objetivo
Explicar como validar, visualizar e compartilhar o arquivo `docs/api.yaml` durante o desenvolvimento.

## Visualização Interativa
### Redocly (recomendado)
1. Instale a CLI (global ou local):
   ```bash
   npm install -g @redocly/cli
   ```
2. Rode o preview:
   ```bash
   redocly preview-docs docs/api.yaml
   ```
3. Acesse `http://127.0.0.1:8080` para navegar pela documentação.

### Swagger UI Watcher (alternativa)
```bash
npx swagger-ui-watcher docs/api.yaml --port 8090
```

## Validação
- Execute o lint com Spectral (Inclua em CI futuramente):
  ```bash
  npx @stoplight/spectral-cli lint docs/api.yaml
  ```
- O pipeline deve falhar se houver erros de esquema ou referências quebradas.

## Comandos Make (atalhos)
- `make api-lint`: roda o Spectral sobre `docs/api.yaml`.
- `make api-preview`: abre o Redocly em modo live preview.
- `make api-types`: gera `frontend/src/types/api.d.ts` com `openapi-typescript` (garante pasta automaticamente).

## Atualizações de Esquema
- Adicione novos endpoints e schemas sempre em `docs/api.yaml`.
- Garanta que os mesmos campos estejam refletidos em `docs/api-reference.md` e nas responses reais.
- Use `$ref` para reutilizar estruturas e evitar divergências.
- Registre toda mudança no `docs/api-changelog.md` e sincronize a versão (`info.version`) com a tag publicada.

## Integração com o Frontend
- Gere clientes automaticamente (opcional):
  ```bash
  npx openapi-typescript docs/api.yaml -o frontend/src/types/api.d.ts
  ```
- Combine com React Query para tipagem forte das requisições/respostas.
- `frontend/package.json` executa `npm run generate:api-types` automaticamente antes de `npm run dev`, `npm run build` e `npm run preview`, garantindo que os tipos estejam atualizados.
- Workflow `.github/workflows/frontend-build.yml` roda `npm run build` em cada push/PR que altera frontend/Makefile/spec, forçando a geração dos tipos durante o CI.

## Publicação
- Hospede o arquivo como artefato do build (ex.: `api-spec` no GitHub Actions).
- Opcional: publicar no Stoplight ou Redocly Workflows para compartilhamento externo.
- Workflow `.github/workflows/api-spec.yml` valida o spec, gera um HTML do Redoc como artefato e, em pushes na `main`, publica automaticamente no GitHub Pages (environment `github-pages`).

## Próximos Passos
1. Adicionar testes de contrato automatizados (ver `docs/tests-contrato.md`) ao pipeline do backend.
2. Automatizar o versionamento/publicação criando tags `api-vX.Y.Z` diretamente no workflow (release drafter).
3. Consumir o spec versionado em ambientes externos (SDK ou Portal do Cliente).
