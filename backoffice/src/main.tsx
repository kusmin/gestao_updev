import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.tsx';
import { AuthProvider } from './contexts/AuthContext.tsx';
import * as Sentry from '@sentry/react';
import { BrowserTracing } from '@sentry/tracing';

const SENTRY_DSN = import.meta.env.VITE_SENTRY_DSN_BACKOFFICE;

if (SENTRY_DSN) {
  Sentry.init({
    dsn: SENTRY_DSN,
    integrations: [
      new BrowserTracing({
        tracePropagationTargets: ["localhost", /^\//],
      }),
    ],
    tracesSampleRate: 1.0,
    replaysSessionSampleRate: 0.1,
    replaysOnErrorSampleRate: 1.0,
  });
}

const GA4_MEASUREMENT_ID = import.meta.env.VITE_GA4_MEASUREMENT_ID_BACKOFFICE;

if (GA4_MEASUREMENT_ID && typeof window.gtag === 'function') {
  window.gtag('config', GA4_MEASUREMENT_ID, {
    send_page_view: false,
  });
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <AuthProvider>
      <App />
    </AuthProvider>
  </React.StrictMode>,
);
