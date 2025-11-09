# Observabilidade – Plano de Implementação

## Contexto Atual
- **Logs**: o backend já inicia um `zap.Logger` em `backend/cmd/api/main.go:23` e adiciona `middleware.RequestID` no router, mas a stack ainda usa o `gin.Logger` padrão (`backend/internal/server/server.go:39`), que imprime em texto simples para `stdout` sem correlação automática com `tenant_id` ou usuário.
- **Métricas**: não existe endpoint `/metrics` nem coleta de métricas de aplicação, GC ou banco.
- **Tracing**: não há instrumentação OpenTelemetry (`rg` sem ocorrências), logo não é possível rastrear chamadas entre serviços ou medir spans de GORM/Gin.
- **Alertas/Dashboards**: apenas diretrizes genéricas constam em `docs/operacao-devops.md:44`, sem detalhes de painéis, consultas ou integrações com Grafana/Alertmanager.

## Objetivos
1. **Logs estruturados consistentes**: middleware próprio que envia todas as requisições para o `zap.Logger` com `request_id`, `tenant_id`, `user_id` (quando autenticado), status HTTP, tempo de resposta e tamanho das cargas.
2. **Métricas exportadas**: expor `/metrics` em Prometheus com counters de requisições, histogramas de latência, métricas de filas (agenda, vendas), além de observar GORM (`db_conns`, `queries_failed_total`).
3. **Tracing distribuído**: habilitar OpenTelemetry (OTLP/HTTP) para Gin, GORM e chamadas HTTP a terceiros, permitindo reconstruir fluxos de reservas e vendas.
4. **Painéis + Alertas**: definir dashboards base na stack Grafana/Loki/Tempo + alertas (latência p95, taxa de erro 5xx, saturação do DB e filas de inventário) com playbooks associados.

## Arquitetura Proposta
| Camada | Ferramenta | Detalhes |
|--------|------------|----------|
| Logs   | `zap` + `middleware.Logger` → Collector Loki | Middleware escreve JSON consolidado, collector envia para Grafana Cloud (Loki). |
| Métricas | OpenTelemetry SDK → `otel-collector` → Prometheus/Grafana | Exportador `prometheus` embutido para `/metrics` local + pipeline OTLP para Grafana Cloud. |
| Traces | `otelgin`, `otelgorm`, `otelhttp` → OTLP → Grafana Tempo | Headers W3C (`traceparent`) propagados; contexto herdado pelo middleware de autenticação. |
| Alertas | Grafana Alerting | Regras baseadas nas métricas acima; notificação para Slack/Teams. |

### Componentes
1. **SDK OpenTelemetry** (Go): inicializado no `cmd/api/main.go` antes de criar o `zapLogger`, carregando config via env (`OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_SERVICE_NAME=gestao-api`).
2. **Middleware de logging**: substitui `gin.Logger()` e usa `zap.Logger` com os campos `request_id`, `tenant_id`, `user_id`, `path`, `method`, `status`, `latency_ms`, `bytes_in`, `bytes_out`.
3. **Collector (docker-compose)**: adicionar serviço `otel-collector` com pipelines `otlp -> loki` e `otlp -> tempo -> grafana` (para ambientes locais/staging).
4. **Dashboards**: JSON exportado para `docs/grafana/*.json` com painéis “Visão Geral API”, “Banco & Fila de Agenda”.

### Execução local
- `.env.example` agora inclui os envs `SERVICE_NAME`, `OTEL_ENABLED`, `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_EXPORTER_OTLP_HEADERS`, `OTEL_EXPORTER_OTLP_INSECURE` e `METRICS_ROUTE`.
- `docker-compose.yml` sobe o `otel-collector` (imagem `otel/opentelemetry-collector-contrib`) usando `docker/otel-collector-config.yaml`, que por padrão exporta traces/métricas para o próprio log da ferramenta.
- Para inspecionar métricas localmente, basta acessar `http://localhost:8080/metrics` (ou o caminho configurado) ou apontar um Prometheus externo com `job_name: gestao-api`.

## Plano de Implementação
### Fase 1 – Fundamentos de Logging ✅
- `internal/middleware/logging.go` substitui o `gin.Logger`, emitindo eventos estruturados no `zap`.
- `RequestID` continua sendo a fonte de correlação e é propagado para cada log.

### Fase 2 – Métricas Prometheus ✅
- Pacote `pkg/telemetry` inicializa OpenTelemetry (MeterProvider + exporter Prometheus) com atributos `service.name`, `deployment.environment`.
- `/metrics` é exposto no caminho configurável `METRICS_ROUTE` (default `/metrics`) e pronto para scrape Prometheus/Loki.

### Fase 3 – Tracing (em andamento)
- `cmd/api/main.go` inicializa o tracer via OTLP HTTP (configurado pelos envs `OTEL_*`) e injeta o middleware `otelgin`.
- Próximos passos: propagar `trace_id` nas respostas HTTP, instrumentar GORM (plugin dedicado), clientes externos e adicionar spans customizados nos casos de uso críticos (CreateBooking, CreateSalesOrder).

### Fase 4 – Stack Operacional
1. Estender `docker-compose.yml` com `grafana`, `loki`, `tempo`, `promtail` (perfil `observability`).
2. Versionar dashboards (`docs/grafana/*.json`) e definir pipeline “Observability smoke-tests” em CI (verificar `/metrics` e `/healthz`).
3. Configurar alertas iniciais e documentar no runbook (links nos docs devops).

## Backlog / Pendências
- 🧩 Definir política de retenção para logs/traces (mínimo 7 dias).
- 🔐 Garantir redaction de dados sensíveis (token JWT, dados de pagamento) antes de enviar ao collector.
- 🧪 Adicionar testes para middleware de logging (verificar presença do `request_id`) e endpoint `/metrics`.
- 📊 Criar painel “Financeiro” com métricas derivadas (faturamento diário) usando PromQL/Recording Rules.
- 📥 Avaliar integração com Sentry ou Honeycomb para alertas de erros aplicacionais antes da fase 3.
- 🧱 Adicionar exporters reais (Grafana Cloud / Tempo / Loki) no `docker/otel-collector-config.yaml`.

Com este plano, conseguimos iniciar a implementação incremental focando primeiro em visibilidade básica (logs e métricas), seguido por tracing e a stack completa de observabilidade/alertas.
