import { useContext } from 'react';
import { AuthContext } from '../contexts/AuthContext';
import { AuthContextValue } from '../contexts/AuthContext'; // Importar o tipo AuthContextValue

export const useAuth = (): AuthContextValue => {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error('useAuth deve ser usado dentro de AuthProvider');
  }
  return ctx;
};
