# Plano de Ação – Aumento de Cobertura de Código

Documento vivo que descreve como mensurar, priorizar e executar o aumento de cobertura de testes no `gestao_updev`, abrangendo backend (Go), frontend e backoffice (React + Vite), além dos testes end-to-end existentes.

## 1. Objetivos e Escopo
- **Qualidade preventiva:** reduzir regressões em fluxos críticos antes de chegarem ao ambiente de produção.
- **Confiança em deploys contínuos:** permitir releases menores com maior frequência e rollback menos frequente.
- **Métricas rastreáveis:** expor cobertura por área no pipeline (GitHub Actions) e no Codecov para acompanhar evolução sprint a sprint.
- **Critério de aceite:** nenhum PR pode reduzir a cobertura global do módulo afetado sem justificativa aprovada pela engenharia.

## 2. Linha de Base e Métricas
Registrar cobertura inicial antes de iniciar o plano. Salvar relatórios HTML e resumos no diretório `docs/coverage-history/<YYYY-MM-DD>` ou em anexos do Codecov. Utilize `make coverage` na raiz para executar todos os relatórios de uma vez antes de preencher a tabela.

### Backend (Go)
```bash
cd backend
go test ./... -covermode=atomic -coverpkg=./... -coverprofile=coverage.out
go tool cover -func=coverage.out              # Resumo por função/pacote
go tool cover -html=coverage.out -o coverage.html
```
- Exportar `coverage.out` para o Codecov (já há `codecov.yml` na raiz).
- Identificar pacotes com `<60%` de cobertura para priorização imediata.

### Frontend
```bash
cd frontend
vitest run --coverage.enabled true --coverage.reporter=text-summary --coverage.reporter=lcov \
  --coverage.include="src/**/*.{ts,tsx}" --passWithNoTests
```
- Publicar `coverage/lcov.info` no pipeline e usar `npm run test -- --watch=false` localmente para feedback rápido.
- Correlacionar gaps com fluxos críticos (onboarding, dashboards, cadastros).

### Backoffice
```bash
cd backoffice
vitest run --coverage.enabled true --coverage.reporter=text-summary --coverage.reporter=lcov \
  --coverage.include="src/**/*.{ts,tsx}" --passWithNoTests
```
- Mapear telas com lógica de permissão, aprovações e conciliações manualmente para garantir casos positivos/negativos.

### Testes End-to-End
- Executar `npm run test:e2e` para validar fluxos completos.
- Registrar taxa de sucesso/flake e ligar resultados a estórias que dependem de integrações reais.

## 3. Metas Quantitativas (sugeridas)
| Área       | Linha de base (2025-11-14) | Meta curto prazo (4 sprints) | Meta médio prazo | Observações |
|------------|---------------------------|------------------------------|------------------|-------------|
| Backend    | 17% statements            | 70%                          | 80%+             | Cobrir services críticos e regras de negócio. |
| Frontend   | 31.6% statements          | 65%                          | 75%+             | Componentes de página + hooks compartilhados. |
| Backoffice | 0% (sem testes ainda)     | 60%                          | 75%+             | Telas administrativas e validações de permissão. |
| E2E        | A medir                   | 12 fluxos                    | 15 fluxos        | Focar fluxos receita, faturamento e suporte. |

> Ajustar metas após capturar a linha de base real; manter gráficos no Codecov e no dashboard do time.

## 4. Backlog Tático por Área

### Backend (Go)
1. **Mapeamento de serviços e pacotes sem testes** usando `go tool cover` + planilha compartilhada.
2. **Testes unitários** para regras de negócio em `internal/<domínio>` utilizando `testing` + `testify`.
3. **Testes de integração leves** com banco em memória (ou SQLite) para repositórios; usar `docker-compose.test.yml` para cenários completos.
4. **Testes de contrato** reutilizando `tests/dredd` para validar compatibilidade com `api.yaml`.
5. **Factories/mocks reutilizáveis** em `internal/testhelpers` para reduzir custo de criação de casos.

### Frontend
1. **Cobrir hooks e utilitários** com Vitest focando em estados edge (erros de API, loading múltiplo).
2. **React Testing Library** para componentes de página: renderização condicional, modais, tabelas.
3. **Mock de React Query** com `QueryClientProvider` específico para cada suíte, garantindo cenários success/error.
4. **Snapshot de contratos** usando MSW ou mocks locais para assegurar alinhamento com backend.
5. **Smoke e2e**: alinhar com Playwright para fluxos mais críticos do app principal (login, dashboard, solicitações).

### Backoffice
1. **Testes de permissão** (admin vs. operador) garantindo ocultação de botões/permissões.
2. **Validações de formulários complexos** (cadastro, auditoria) com `user-event`.
3. **Cobertura de hooks de API** garantindo caching e sincronização com React Query.
4. **Integração com relatórios**: validar exportações e filtros, usando mocks de `FileReader` ou `URL.createObjectURL`.
5. **Fluxos críticos com Playwright** compartilhando fixtures do frontend principal quando possível.

## 5. Processo e Governança
- **Pipeline:** adicionar etapas `go test` e `vitest --coverage` em CI; falhar build se cobertura cair abaixo do alvo configurado no `codecov.yml`.
- **Template de PR:** exigir link para relatório de cobertura ou justificativa formal quando houver queda.
- **Rotina semanal:** quinta-feira revisar métricas no Codecov, atualizar tendências e decidir próximos focos.
- **Definition of Done:** histórias só podem ser concluídas com testes cobrindo novos caminhos de código.
- **Pairing/Dojo:** sessões quinzenais para destravar áreas com testes mais complexos (ex.: autenticação, billing).

## 6. Cronograma Sugerido (8 semanas)
| Semana | Objetivo principal | Resultado esperado |
|--------|-------------------|--------------------|
| 1-2    | Medir baseline, configurar pipelines e Codecov | Relatórios publicados e alertas no CI. |
| 3-4    | Quick wins em módulos críticos (backend serviços + hooks frontend) | +10-15 p.p. em áreas prioritárias. |
| 5-6    | Testes de integração/contrato e cobertura de telas sensíveis no backoffice | Fluxos administrativos cobertos. |
| 7-8    | Endurecer PR gate e revisar metas | Cobertura estabilizada + rotinas documentadas. |

## 7. Papéis e Responsabilidades
| Papel                | Responsabilidades chave |
|----------------------|-------------------------|
| Tech Lead Backend    | Definir prioridades de pacotes, revisar testes de integração, manter helpers. |
| Tech Lead Frontend   | Garantir padrões de RTL/Vitest, apoiar mocks de API e práticas de MSW. |
| QA / Engenheiro(a)   | Consolidar métricas, coordenar testes E2E e monitorar flakes. |
| Devs do squad        | Criar/atualizar testes em cada PR seguindo o plano e metas. |

## 8. Riscos e Mitigações
- **Falta de baseline confiável:** priorizar instrumentação no início e automatizar coleta.
- **Builds lentos:** usar `vitest --runInBand=false` e cache de Go; dividir suites longas em jobs paralelos.
- **Flakiness E2E:** investir em fixtures determinísticas e usar `retry` apenas como último recurso.
- **Rotatividade de time:** registrar exemplos nas pastas `docs/tests-*` e manter templates de casos.

## 9. Próximos Passos Imediatos
1. Executar comandos da Seção 2 e preencher a tabela de metas com valores reais.
2. Criar tarefas no board (um item por iniciativa listada na Seção 4).
3. Atualizar o pipeline CI com upload automático de cobertura (Go + Vitest + Playwright).
4. Revisitar o plano trimestralmente e ajustar metas conforme maturidade e roadmap.
