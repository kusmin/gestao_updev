# Plano de Implementação de Monitoramento (Sentry e Google Analytics 4)

## 1. Visão Geral

Este documento detalha o plano para integrar ferramentas de monitoramento de erros e análise de comportamento (Sentry e Google Analytics 4 - GA4) nas aplicações frontend (`frontend/` e `backoffice/`) e backend (`backend/`) do projeto `gestao_updev`. O objetivo é obter insights sobre a performance, identificar e resolver erros proativamente, e entender o comportamento do usuário para futuras melhorias.

## 2. Ferramentas Escolhidas

*   **Sentry:** Para monitoramento de erros e performance (Real User Monitoring - RUM para frontend, APM para backend).
*   **Google Analytics 4 (GA4):** Para análise de comportamento e métricas de uso da aplicação (apenas para frontends).

## 3. Plano de Implementação Detalhado

A implementação será realizada separadamente para cada aplicação (`backend/`, `frontend/` e `backoffice/`), garantindo que cada uma tenha sua própria configuração e dados de monitoramento.

### 3.1. Implementação do Sentry

O Sentry será configurado para capturar erros, rastrear transações e monitorar a performance das aplicações.

#### Configuração Prévia:

*   Criar projetos separados no Sentry para "gestao-frontend-cliente", "gestao-backoffice" e "gestao-backend".
*   Obter os DSNs (Data Source Names) para cada projeto Sentry.

#### 3.1.1. Implementação do Sentry nos Frontends (`frontend/` e `backoffice/`)

#### Passos para `frontend/` e `backoffice/` (repetir para ambos):

1.  **Instalação das Dependências:**
    ```bash
    cd frontend # ou cd backoffice
    npm install @sentry/react @sentry/tracing
    ```

2.  **Inicialização do Sentry:**
    *   No ponto de entrada da aplicação (geralmente `src/main.tsx` ou `src/App.tsx`), adicionar o código de inicialização do Sentry.

    ```typescript
    // Exemplo para src/main.tsx
    import React from 'react';
    import ReactDOM from 'react-dom/client';
    import * as Sentry from '@sentry/react';
    import { BrowserTracing } from '@sentry/tracing';
    import App from './App';
    import './index.css'; // ou o arquivo CSS principal

    // Configurar o DSN do Sentry (substitua pelo DSN real do seu projeto)
    const SENTRY_DSN = import.meta.env.VITE_SENTRY_DSN; // Exemplo de variável de ambiente

    if (SENTRY_DSN) {
      Sentry.init({
        dsn: SENTRY_DSN,
        integrations: [
          new BrowserTracing({
            tracePropagationTargets: ["localhost", /^\//], // Ajuste conforme seus domínios de API
          }),
          // Adicione outras integrações conforme necessário
        ],
        // Taxa de amostragem para rastreamento de transações (0 a 1)
        tracesSampleRate: 1.0, // Ajuste para um valor menor em produção (ex: 0.2)

        // Taxa de amostragem para replays de sessão (0 a 1)
        replaysSessionSampleRate: 0.1, // Amostra 10% das sessões
        replaysOnErrorSampleRate: 1.0, // Amostra 100% das sessões com erro
      });
    }

    ReactDOM.createRoot(document.getElementById('root')!).render(
      <React.StrictMode>
        <App />
      </ReactStrictMode>,
    );
    ```

3.  **Configuração de Variáveis de Ambiente:**
    *   Adicionar `VITE_SENTRY_DSN` (ou similar) ao arquivo `.env.example` e `.env` de cada aplicação frontend.

4.  **Identificação de Usuário (Opcional, mas Recomendado):**
    *   Após o login do usuário, chamar `Sentry.setUser()` para associar erros e transações a usuários específicos.

    ```typescript
    // Exemplo após o login
    Sentry.setUser({
      id: user.id,
      email: user.email,
      username: user.name,
      // Outros dados do usuário que sejam úteis para depuração
    });
    ```

5.  **Adição de Contexto Adicional (Opcional):**
    *   Usar `Sentry.setContext()` para adicionar informações relevantes sobre a aplicação ou o estado atual.

    ```typescript
    Sentry.setContext("app_state", {
      theme: "dark",
      currentRoute: "/dashboard",
    });
    ```

6.  **Verificação:**
    *   Forçar um erro em desenvolvimento (ex: `throw new Error("Teste de erro Sentry");`) e verificar se ele aparece no painel do Sentry.

#### 3.1.2. Implementação do Sentry no Backend (Go)

O Sentry será configurado no backend Go para capturar erros, panics e rastrear transações HTTP.

#### Configuração Prévia:

*   Certificar-se de que um projeto Sentry para "gestao-backend" foi criado e o DSN obtido.

#### Passos para `backend/`:

1.  **Instalação da Dependência:**
    ```bash
    cd backend
    go get github.com/getsentry/sentry-go
    go get github.com/getsentry/sentry-go/gin
    ```

2.  **Inicialização do Sentry:**
    *   No ponto de entrada da aplicação (ex: `cmd/api/main.go`), inicializar o Sentry.

    ```go
    package main

    import (
        "log"
        "os"
        "time" // Importar time

        "github.com/getsentry/sentry-go"
        sentrygin "github.com/getsentry/sentry-go/gin"
        "github.com/gin-gonic/gin"
        // ... outras importações
    )

    func main() {
        // Carregar variáveis de ambiente (se ainda não estiverem carregadas)
        // ...

        if err := sentry.Init(sentry.ClientOptions{
            Dsn: os.Getenv("SENTRY_DSN"), // Obter DSN de variável de ambiente
            // Set TracesSampleRate to 1.0 to capture 100%
            // of transactions for performance monitoring.
            // We recommend adjusting this value in production.
            TracesSampleRate: 1.0,
            EnableTracing: true,
            // Debug: true, // Ativar para depuração do Sentry
        }); err != nil {
            log.Fatalf("sentry.Init: %s", err)
        }
        // Flush buffered events before exiting
        defer sentry.Flush(2 * time.Second)

        router := gin.New()
        // ... outros middlewares

        // Adicionar middleware do Sentry para Gin
        router.Use(sentrygin.New(sentrygin.Options{
            Repanic: true, // Re-panic para que outros middlewares possam capturar
        }))

        // Opcional: Capturar panics que ocorrem fora do contexto do Gin
        defer func() {
            if r := recover(); r != nil {
                sentry.CurrentHub().Recover(r)
                sentry.Flush(2 * time.Second)
                panic(r) // Re-panic para que o programa termine como esperado
            }
        }()

        // ... suas rotas e handlers
        // router.GET("/ping", func(c *gin.Context) {
        //     c.String(http.StatusOK, "pong")
        // })

        // Exemplo de rota que pode causar um erro
        // router.GET("/error", func(c *gin.Context) {
        //     panic("Um erro inesperado ocorreu!")
        // })

        // router.Run(":8080")
    }
    ```

3.  **Configuração de Variáveis de Ambiente:**
    *   Adicionar `SENTRY_DSN` ao arquivo `.env.example` e `.env` do backend.

4.  **Integração com Logger (Zap):**
    *   Configurar o logger `zap` para enviar logs de nível `Error` ou superior para o Sentry. Isso geralmente envolve a criação de um `Core` personalizado para o Zap que encaminha os logs para o Sentry.

    ```go
    // Exemplo de integração com Zap (simplificado)
    // No seu pacote de logger (ex: pkg/logger/logger.go)
    import (
        "github.com/getsentry/sentry-go"
        "go.uber.org/zap"
        "go.uber.org/zap/zapcore"
        sentryzap "github.com/TheZeroSlave/zapsentry" // Importar zapsentry
        "os" // Importar os
    )

    func NewLogger(sentryDsn string) *zap.Logger {
        // ... configuração normal do Zap (assumindo que 'core' já está definido)
        // Exemplo de core básico, ajuste conforme a implementação real do projeto
        encoderConfig := zap.NewProductionEncoderConfig()
        encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
        core := zapcore.NewCore(
            zapcore.NewJSONEncoder(encoderConfig),
            zapcore.AddSync(os.Stdout), // Adicionado para evitar erro de 'os'
            zap.InfoLevel,
        )

        // Se o Sentry estiver inicializado, adicione um core para ele
        if sentry.CurrentHub().Client() != nil {
            sentryCore := sentryzap.New(zapcore.ErrorLevel) // Envia logs de erro para o Sentry
            return zap.New(zapcore.NewTee(core, sentryCore))
        }

        return zap.New(core)
    }
    ```
    *   Será necessário instalar `go get github.com/TheZeroSlave/zapsentry`.

5.  **Identificação de Usuário e Contexto:**
    *   Em seus handlers, você pode usar `sentry.CurrentHub().Scope().SetUser()` e `sentry.CurrentHub().Scope().SetContext()` para adicionar informações relevantes aos eventos de erro.

    ```go
    // Exemplo em um handler Gin
    func MyHandler(c *gin.Context) {
        // ...
        sentry.CurrentHub().Scope().SetUser(sentry.User{
            ID:    "user-id-from-auth",
            Email: "user@example.com",
        })
        sentry.CurrentHub().Scope().SetContext("request_info", map[string]interface{}{
            "ip_address": c.ClientIP(),
            "user_agent": c.Request.UserAgent(),
        })
        // ...
    }
    ```

6.  **Verificação:**
    *   Forçar um erro ou panic no backend (ex: em uma rota de teste) e verificar se ele aparece no painel do Sentry.
    *   Verificar se os logs de erro estão sendo enviados corretamente.

### 3.2. Implementação do Google Analytics 4 (GA4)

O GA4 será configurado para coletar dados sobre o comportamento do usuário, visualizações de página e eventos personalizados.

#### Configuração Prévia:

*   Criar propriedades GA4 separadas para "gestao-frontend-cliente" e "gestao-backoffice" no Google Analytics.
*   Obter os IDs de Medição (Measurement IDs) para cada propriedade (ex: `G-XXXXXXXXXX`).

#### Passos para `frontend/` e `backoffice/` (repetir para ambos):

1.  **Adição da Tag do Google (gtag.js):
    *   No arquivo `index.html` de cada aplicação, adicionar o snippet da tag do Google dentro da seção `<head>`.

    ```html
    <!-- Google tag (gtag.js) -->
    <script async src="https://www.googletagmanager.com/gtag/js?id=G-XXXXXXXXXX"></script>
    <script>
      window.dataLayer = window.dataLayer || [];
      function gtag(){dataLayer.push(arguments);}
      gtag('js', new Date());

      gtag('config', 'G-XXXXXXXXXX'); // Substitua pelo seu Measurement ID
    </script>
    ```
    *   Substituir `G-XXXXXXXXXX` pelo ID de Medição real de cada propriedade GA4.

2.  **Rastreamento de Páginas:**
    *   Se estiver usando `react-router-dom`, configurar um listener para rastrear as visualizações de página a cada mudança de rota.

    ```typescript
    // Exemplo de hook personalizado para rastreamento de página
    import { useEffect } from 'react';
    import { useLocation } from 'react-router-dom';

    const usePageTracking = (measurementId: string) => {
      const location = useLocation();

      useEffect(() => {
        if (window.gtag) {
          window.gtag('config', measurementId, {
            'page_path': location.pathname + location.search,
          });
        }
      }, [location, measurementId]);
    };

    // No seu componente App.tsx ou similar
    // usePageTracking(import.meta.env.VITE_GA4_MEASUREMENT_ID);
    ```

3.  **Rastreamento de Eventos Personalizados:**
    *   Criar uma função utilitária para enviar eventos personalizados ao GA4.

    ```typescript
    // utils/ga4.ts
    export const trackEvent = (eventName: string, eventParams: Record<string, any> = {}) => {
      if (window.gtag) {
        window.gtag('event', eventName, eventParams);
      }
    };

    // Exemplo de uso em um componente
    // import { trackEvent } from '../utils/ga4';
    // <button onClick={() => trackEvent('button_click', { button_name: 'submit_form' })}>Enviar</button>
    ```

4.  **Identificação de Usuário (Opcional, mas Recomendado):**
    *   Após o login do usuário, enviar o `user_id` para o GA4.

    ```typescript
    // Exemplo após o login
    if (window.gtag) {
      window.gtag('set', 'user_properties', {
        user_id: user.id,
        user_email: user.email,
      });
    }
    ```

5.  **Verificação:**
    *   Utilizar o "DebugView" no painel do Google Analytics para verificar se os eventos e visualizações de página estão sendo coletados corretamente em tempo real.
    *   Verificar o console do navegador para erros relacionados ao GA4.

## 4. Próximos Passos

*   Implementar as configurações e o código conforme detalhado acima no `backend/`, `frontend/` e `backoffice/`.
*   Testar exaustivamente em ambientes de desenvolvimento e staging.
*   Monitorar os painéis do Sentry e GA4 após o deploy em produção.
*   Refinar as configurações de amostragem e os eventos personalizados conforme a necessidade do negócio.
