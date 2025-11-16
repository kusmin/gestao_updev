import {
  CssBaseline,
  ThemeProvider as MuiThemeProvider,
  createTheme,
} from '@mui/material';
import { useEffect, useMemo, useState } from 'react';
import { getDesignTokens, resolveInitialMode, STORAGE_KEY } from './utils';
import { ColorModeContext } from './ColorModeContextDefinition';

export const ColorModeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [mode, setMode] = useState<PaletteMode>(resolveInitialMode);

  useEffect(() => {
    window.localStorage.setItem(STORAGE_KEY, mode);
  }, [mode]);

  const value = useMemo(
    () => {
      const toggleColorMode = () => {
        setMode((prev) => (prev === 'light' ? 'dark' : 'light'));
      };
      return {
        mode,
        toggleColorMode,
      };
    },
    [mode],
  );

  const theme = useMemo(() => createTheme(getDesignTokens(mode)), [mode]);

  return (
    <ColorModeContext.Provider value={value}>
      <MuiThemeProvider theme={theme}>
        <CssBaseline />
        {children}
      </MuiThemeProvider>
    </ColorModeContext.Provider>
  );
};
