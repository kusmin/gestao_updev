import { createContext, useContext } from 'react';

// Define a interface para o contexto
interface SnackbarContextType {
  showSnackbar: (message: string, severity: 'success' | 'error' | 'info' | 'warning') => void;
}

// Cria o contexto com um valor padrão undefined
export const SnackbarContext = createContext<SnackbarContextType | undefined>(undefined);

export const useSnackbar = () => {
  const context = useContext(SnackbarContext);
  if (!context) {
    throw new Error('useSnackbar must be used within a SnackbarProvider');
  }
  return context;
};
