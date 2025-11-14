import React, { createContext, useState, useContext, ReactNode, useEffect } from 'react';
import { setAuthToken } from '../lib/apiClient';

interface AuthContextType {
  token: string | null;
  login: (token: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    const storedToken = localStorage.getItem('backoffice_token');
    if (storedToken) {
      setToken(storedToken);
      setAuthToken(storedToken);
    }
  }, []);

  const login = (newToken: string) => {
    setToken(newToken);
    setAuthToken(newToken);
    localStorage.setItem('backoffice_token', newToken);
  };

  const logout = () => {
    setToken(null);
    setAuthToken(null);
    localStorage.removeItem('backoffice_token');
  };

  return <AuthContext.Provider value={{ token, login, logout }}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
