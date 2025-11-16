import React from 'react';
import ReactDOM from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
// import * as Sentry from '@sentry/react';
// import { BrowserTracing } from '@sentry/tracing'; // Reverter para importação direta

import App from './App';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';

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

// const SENTRY_DSN = import.meta.env.VITE_SENTRY_DSN_FRONTEND;

// if (SENTRY_DSN) {
//   Sentry.init({
//     dsn: SENTRY_DSN,
//     integrations: [
//       new BrowserTracing({ // Reverter para instanciação direta
//         tracingOrigins: ['localhost', /^\//], // Manter tracingOrigins
//       }),
//       Sentry.replayIntegration({
//         maskAllInputs: true,
//         blockAllMedia: true,
//       }),
//     ],
//     tracesSampleRate: 1.0,
//     replaysSessionSampleRate: 0.1,
//     replaysOnErrorSampleRate: 1.0,
//   });
// }

const GA4_MEASUREMENT_ID = import.meta.env.VITE_GA4_MEASUREMENT_ID_FRONTEND;

if (GA4_MEASUREMENT_ID && typeof window.gtag === 'function') {
  window.gtag('config', GA4_MEASUREMENT_ID, {
    send_page_view: false,
  });
}

const queryClient = new QueryClient();

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>
  </React.StrictMode>,
);
