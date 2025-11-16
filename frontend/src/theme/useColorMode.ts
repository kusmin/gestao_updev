import { useContext } from 'react';
import { ColorModeContext, ColorModeContextValue } from './ColorModeProvider';

export const useColorMode = (): ColorModeContextValue => {
  const ctx = useContext(ColorModeContext);
  if (!ctx) {
    throw new Error('useColorMode must be used within ColorModeProvider');
  }
  return ctx;
};
