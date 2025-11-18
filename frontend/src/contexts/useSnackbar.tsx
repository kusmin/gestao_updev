import { createContext, useContext } from 'react';
import { AlertColor } from '@mui/material';

export interface SnackbarContextData {
  showSnackbar: (message: string, severity: AlertColor) => void;
}

export const SnackbarContext = createContext<SnackbarContextData | undefined>(undefined);

export const useSnackbar = () => {
  const context = useContext(SnackbarContext);
  if (context === undefined) {
    throw new Error('useSnackbar must be used within a SnackbarProvider');
  }
  return context;
};
