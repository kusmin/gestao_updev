# Implementação de PWA (Progressive Web App) no Gestão UpDev

## 1. Introdução

Progressive Web Apps (PWAs) são aplicações web que utilizam capacidades modernas do navegador para oferecer uma experiência de usuário semelhante a um aplicativo nativo. Eles são confiáveis, rápidos e envolventes. A implementação de PWAs no `backoffice` e `frontend` do Gestão UpDev trará benefícios como:

*   **Confiabilidade**: Carregamento instantâneo e desempenho consistente, mesmo em condições de rede instáveis ou offline.
*   **Velocidade**: Respostas rápidas às interações do usuário com animações suaves.
*   **Engajamento**: Experiência de usuário envolvente com recursos como instalação na tela inicial, notificações push e acesso a hardware do dispositivo.

## 2. Requisitos Essenciais para um PWA

Para que uma aplicação web seja considerada um PWA, ela deve atender a alguns requisitos fundamentais:

*   **HTTPS**: A aplicação deve ser servida via HTTPS para garantir a segurança e a integridade dos dados.
*   **Service Worker**: Um script JavaScript que o navegador executa em segundo plano, separado da página web. Ele permite o controle sobre as requisições de rede, cache de recursos e funcionalidades offline.
*   **Web App Manifest**: Um arquivo JSON que fornece informações sobre a aplicação (nome, ícones, cores, URL de início, etc.) ao navegador, permitindo que ela seja instalada na tela inicial do usuário e se comporte como um aplicativo nativo.

## 3. Análise da Estrutura Atual (Backoffice e Frontend)

Ambos os projetos `backoffice` e `frontend` são construídos com **React** e **Vite**. A estrutura atual é favorável à implementação de PWA, pois o Vite oferece um ecossistema moderno e flexível, e existem plugins que simplificam a configuração.

*   **`package.json`**: Ambos utilizam `vite` para desenvolvimento e build.
*   **`vite.config.ts`**: Configurações padrão do Vite.
*   **`index.html`**: HTML básico com um `div` para o React montar a aplicação, já contendo `meta name="viewport"`.

A abordagem para ambos os projetos será praticamente idêntica, com pequenas variações nos detalhes do `manifest.json` (nome, descrição, ícones).

## 4. Passo a Passo para Implementação de PWA

Este guia detalha as etapas para adicionar funcionalidades PWA aos projetos `backoffice` e `frontend`.

### 4.1. Instalação do Plugin `vite-plugin-pwa`

O `vite-plugin-pwa` é a ferramenta recomendada para integrar PWA com o Vite.

```bash
# No diretório raiz do projeto (backoffice ou frontend)
pnpm add -D vite-plugin-pwa
```

### 4.2. Configuração do `vite.config.ts`

Adicione e configure o plugin no seu arquivo `vite.config.ts`.

```typescript
// backoffice/vite.config.ts ou frontend/vite.config.ts
import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';
import { VitePWA } from 'vite-plugin-pwa'; // Importe o plugin

export default defineConfig({
  plugins: [
    react(),
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['favicon.ico', 'apple-touch-icon.png', 'masked-icon.svg'], // Inclua seus assets
      manifest: {
        name: 'Nome da Aplicação', // Ex: Gestão UpDev Backoffice ou Gestão UpDev Frontend
        short_name: 'ShortName', // Ex: Backoffice ou Frontend
        description: 'Descrição da sua aplicação',
        theme_color: '#ffffff',
        icons: [
          {
            src: 'pwa-192x192.png',
            sizes: '192x192',
            type: 'image/png',
          },
          {
            src: 'pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
          },
          {
            src: 'pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'any maskable',
          },
        ],
      },
      workbox: {
        // Configurações do Workbox para cache de assets
        globPatterns: ['**/*.{js,css,html,ico,png,svg}'],
      },
      devOptions: {
        enabled: true, // Habilita PWA em desenvolvimento para testes
      },
    }),
  ],
  // ... outras configurações
});
```

**Observações sobre a configuração do `manifest`:**

*   **`name` e `short_name`**: Devem ser específicos para cada aplicação (Backoffice/Frontend).
*   **`theme_color`**: Cor da barra de status do navegador.
*   **`icons`**: É crucial fornecer ícones em vários tamanhos e com `purpose: 'maskable'` para garantir que a aplicação seja exibida corretamente em diferentes dispositivos e modos de exibição. Você precisará criar esses arquivos de imagem e colocá-los na pasta `public/`.
*   **`workbox`**: O `vite-plugin-pwa` usa o Workbox para gerar o service worker. `globPatterns` define quais arquivos serão cacheados.

### 4.3. Criação dos Ícones

Crie os ícones necessários e coloque-os na pasta `public/` de cada projeto.

*   `public/favicon.ico`
*   `public/apple-touch-icon.png` (ex: 180x180px)
*   `public/masked-icon.svg`
*   `public/pwa-192x192.png`
*   `public/pwa-512x512.png`

### 4.4. Registro do Service Worker

O `vite-plugin-pwa` geralmente injeta o registro do service worker automaticamente. No entanto, você pode adicionar um registro explícito no seu `main.tsx` (ou arquivo de entrada principal) para ter mais controle ou para exibir prompts de instalação/atualização.

```typescript
// backoffice/src/main.tsx ou frontend/src/main.tsx
import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.tsx';
import './index.css';

// Adicione este bloco para registrar o service worker
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js') // O nome do arquivo gerado pelo vite-plugin-pwa
      .then(registration => {
        console.log('SW registered: ', registration);
      })
      .catch(registrationError => {
        console.log('SW registration failed: ', registrationError);
      });
  });
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
```

**Nota**: O `vite-plugin-pwa` gera o arquivo `sw.js` (ou o nome configurado) automaticamente no diretório `dist` durante o build.

### 4.5. Atualização do `index.html`

O `vite-plugin-pwa` injeta automaticamente a tag `<link rel="manifest" ...>` no `index.html` durante o build. No entanto, você pode adicionar manualmente algumas tags para melhor compatibilidade e experiência:

```html
<!-- backoffice/index.html ou frontend/index.html -->
<!doctype html>
<html lang="pt-BR">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Gestão UpDev</title>
    <!-- Adicione as seguintes tags para PWA -->
    <meta name="theme-color" content="#ffffff" />
    <link rel="apple-touch-icon" href="/apple-touch-icon.png" sizes="180x180" />
    <link rel="mask-icon" href="/masked-icon.svg" color="#FFFFFF" />
    <!-- O manifest será injetado automaticamente pelo vite-plugin-pwa -->
    <!-- <link rel="manifest" href="/manifest.webmanifest" /> -->
    <!-- ... outras tags ... -->
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
```

### 4.6. Testando o PWA

Após a implementação, você pode testar seu PWA:

1.  **Build da Aplicação**: `pnpm run build`
2.  **Servir a Aplicação**: Use um servidor estático (ex: `npx serve dist`) ou o `pnpm run preview` do Vite.
3.  **Chrome DevTools**: Abra a aba "Application" -> "Manifest" e "Service Workers" para verificar se o manifest está sendo carregado corretamente e se o service worker está ativo.
4.  **Lighthouse**: Execute uma auditoria Lighthouse (na aba "Lighthouse" do DevTools) e verifique a categoria "Progressive Web App" para ver se todos os requisitos estão sendo atendidos.
5.  **Instalação**: Tente instalar a aplicação na tela inicial do seu dispositivo (desktop ou mobile).

## 5. Considerações Adicionais

*   **Estratégias de Cache**: O Workbox permite configurar diferentes estratégias de cache (Cache First, Network First, Stale While Revalidate) para diferentes tipos de assets. Ajuste `workbox.globPatterns` e adicione `runtimeCaching` conforme necessário.
*   **Notificações Push**: Para implementar notificações push, você precisará de um servidor para enviar as notificações e de mais configurações no service worker para recebê-las e exibi-las.
*   **Atualizações do PWA**: O `vite-plugin-pwa` com `registerType: 'autoUpdate'` já lida com a atualização do service worker. Você pode adicionar lógica na sua aplicação para notificar o usuário sobre uma nova versão disponível.
*   **Experiência Offline**: Certifique-se de que as rotas críticas da sua aplicação funcionem offline, exibindo uma página de fallback ou dados cacheados.

Este documento servirá como base para a implementação do PWA em ambos os projetos.
