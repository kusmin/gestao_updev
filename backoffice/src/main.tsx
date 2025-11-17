import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.tsx';
import { AuthProvider } from './contexts/AuthContext.tsx';
import { SnackbarProvider } from './contexts/SnackbarContext.tsx';

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
    <AuthProvider>
      <SnackbarProvider>
        <App />
      </SnackbarProvider>
    </AuthProvider>
  </React.StrictMode>,
);
