# Plano Detalhado de Integração do Redis

Este documento descreve um plano detalhado e os passos necessários para integrar o Redis à plataforma Gestão UpDev, visando otimizar o desempenho, a escalabilidade e habilitar novas funcionalidades.

## 1. Objetivo

Integrar o Redis como uma camada de cache em memória e/ou para outras funcionalidades avançadas (como rate limiting ou filas de mensagens) no backend Go da plataforma Gestão UpDev.

## 2. Justificativa

A adição do Redis trará os seguintes benefícios, conforme analisado em `docs/cache-analysis.md`:

*   **Melhora de Desempenho:** Redução significativa da latência de resposta da API ao servir dados da memória.
*   **Redução da Carga no PostgreSQL:** Diminuição do número de requisições ao banco de dados principal, liberando recursos para operações mais complexas.
*   **Aumento da Escalabilidade:** Capacidade de suportar um maior volume de requisições e usuários com os mesmos recursos.
*   **Habilitação de Novas Funcionalidades:** Suporte a rate limiting, filas de mensagens e outras estruturas de dados avançadas.

## 3. Escopo da Integração (Fase 1 - Prova de Conceito)

Para a fase inicial (PoC), o foco será em cenários de leitura intensiva e de alto impacto:

*   **Cache de Configurações da Empresa (Tenant):** As configurações da `Company` são frequentemente acessadas após a autenticação.
*   **Cache de Dados de Catálogo:** `Service` e `Product` são listados com frequência.
*   **Rate Limiting (Opcional):** Implementação básica de rate limiting por IP ou por tenant para endpoints críticos.

## 4. Impacto na Arquitetura

### Backend (Go)

*   **Nova Camada/Módulo:** Criação de um módulo `pkg/cache` para encapsular a lógica de conexão e operações com o Redis.
*   **Modificação de Serviços/Repositórios:** A camada de serviço ou repositório será modificada para primeiro verificar o cache antes de consultar o banco de dados (estratégia Cache-Aside).
*   **Invalidação de Cache:** Implementação de lógica para invalidar o cache quando os dados correspondentes forem alterados no banco de dados.

### Infraestrutura

*   **Docker Compose:** Adição de um serviço Redis ao `docker-compose.yml` para o ambiente de desenvolvimento local e testes.
*   **Ambiente de Produção:** Provisionamento de uma instância gerenciada de Redis (ex: AWS ElastiCache, Google Cloud Memorystore, Azure Cache for Redis) para garantir alta disponibilidade, escalabilidade e segurança.

## 5. Plano de Implementação Detalhado

### Passo 1: Configuração do Ambiente de Desenvolvimento

1.  **Adicionar Serviço Redis ao `docker-compose.yml`:**
    *   Incluir um serviço Redis básico no arquivo `docker-compose.yml` principal do projeto.
    *   Exemplo de configuração:
        ```yaml
        # docker-compose.yml
        services:
          redis:
            image: redis:7-alpine
            ports:
              - "6379:6379"
            command: redis-server --appendonly yes
            volumes:
              - redis_data:/data
        volumes:
          redis_data:
        ```
2.  **Atualizar Variáveis de Ambiente:**
    *   Adicionar `REDIS_ADDR`, `REDIS_PASSWORD` (se aplicável), `REDIS_DB` ao `.env.example` e `.env.test`.
    *   Atualizar o `config.Config` em `backend/internal/config/config.go` para carregar essas variáveis.

### Passo 2: Integração no Backend Go

1.  **Instalar Biblioteca `go-redis`:**
    *   No diretório `backend/`: `go get github.com/go-redis/redis/v8`
2.  **Criar Cliente Redis (`pkg/cache`):**
    *   Criar um novo pacote `backend/pkg/cache`.
    *   Implementar uma função `NewRedisClient(cfg *config.Config) *redis.Client` que inicializa e retorna um cliente Redis.
    *   Implementar um método `Close()` para fechar a conexão.
    *   Adicionar o cliente Redis ao `server.Server` struct em `backend/internal/server/server.go`.
3.  **Implementar Wrapper de Cache:**
    *   No `backend/pkg/cache`, criar uma interface `Cache` com métodos como `Get(key string, dest interface{}) error`, `Set(key string, value interface{}, expiration time.Duration) error`, `Delete(key string) error`.
    *   Implementar `RedisCache` que satisfaça essa interface.
    *   Considerar um mecanismo de serialização/deserialização (ex: `json.Marshal`/`json.Unmarshal`) para armazenar objetos Go no Redis.

### Passo 3: Implementação de Cenários de Cache (PoC)

1.  **Cache de Configurações da Empresa:**
    *   **Modificar `CompanyService`:** No `backend/internal/service/company.go`, antes de buscar a `Company` no repositório, verificar o cache. Se encontrado, retornar do cache. Se não, buscar no DB, armazenar no cache e retornar.
    *   **Invalidação:** No método `UpdateCompany` do `CompanyService`, após a atualização bem-sucedida no DB, invalidar a entrada correspondente no cache.
2.  **Cache de Dados de Catálogo (Serviços/Produtos):**
    *   **Modificar `CatalogService`:** Implementar lógica de cache para os métodos `GetServices` e `GetProducts`.
    *   **Invalidação:** Invalidar o cache de um serviço/produto específico após operações de `Create`, `Update` ou `Delete`.

### Passo 4: Testes

1.  **Testes Unitários:**
    *   Escrever testes unitários para o pacote `backend/pkg/cache` para garantir que as operações de cache funcionam corretamente.
2.  **Testes de Integração:**
    *   Modificar testes de integração existentes (ex: `CompanyService_test.go`, `CatalogService_test.go`) para incluir o Redis.
    *   Garantir que o ambiente de teste (via Docker Compose) inclua o serviço Redis.
    *   Verificar que o cache é populado e invalidado conforme esperado.
3.  **Testes de Carga/Desempenho:**
    *   Utilizar ferramentas como `k6` ou `JMeter` para simular carga e medir os ganhos de desempenho nos endpoints cacheados.

### Passo 5: Monitoramento

1.  **Métricas Redis:**
    *   Integrar métricas do Redis (hits, misses, latência de comandos) ao sistema de observabilidade (Prometheus/Grafana). A biblioteca `go-redis` pode expor métricas ou um exporter Redis pode ser usado.
2.  **Logs:**
    *   Adicionar logs informativos sobre operações de cache (cache hit/miss) para facilitar a depuração e o monitoramento.

## 6. Considerações de Segurança

*   **Autenticação:** Configurar o Redis com uma senha (`requirepass`) em ambientes de produção.
*   **Criptografia (TLS):** Se o Redis for acessível externamente, configurar TLS para criptografar o tráfego.
*   **Acesso Restrito:** Limitar o acesso ao Redis apenas para a aplicação backend.

## 7. Próximos Passos (Pós-PoC)

*   **Expansão de Cenários:**
    *   Implementar cache para resultados de relatórios ou KPIs do dashboard.
    *   Utilizar Redis para rate limiting em endpoints críticos.
    *   Explorar o uso de Pub/Sub para comunicação entre serviços ou para invalidação de cache distribuída.
    *   Utilizar Redis como broker de mensagens para filas de tarefas assíncronas.
*   **Configuração de Produção:**
    *   Configurar o Redis em um ambiente de produção gerenciado, com replicação e alta disponibilidade.
    *   Ajustar configurações de persistência do Redis (RDB/AOF) conforme a necessidade.
