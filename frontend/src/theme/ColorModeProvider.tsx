import {
  CssBaseline,
  PaletteMode,
  ThemeProvider as MuiThemeProvider,
  createTheme,
} from '@mui/material';
import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';

type ColorMode = PaletteMode;

type ColorModeContextValue = {
  mode: ColorMode;
  toggleColorMode: () => void;
};

const STORAGE_KEY = 'gestao-frontend-theme';

const ColorModeContext = createContext<ColorModeContextValue | undefined>(undefined);

const getDesignTokens = (mode: PaletteMode) => ({
  palette: {
    mode,
    primary: {
      main: mode === 'light' ? '#0044ff' : '#86a9ff',
    },
    secondary: {
      main: mode === 'light' ? '#00b0b9' : '#66d8de',
    },
    background: {
      default: mode === 'light' ? '#f7f9fc' : '#0f172a',
      paper: mode === 'light' ? '#ffffff' : '#1e293b',
    },
  },
  shape: {
    borderRadius: 16,
  },
  typography: {
    fontFamily: [
      'Roboto',
      'Inter',
      'system-ui',
      '-apple-system',
      'BlinkMacSystemFont',
      'Segoe UI',
      'sans-serif',
    ].join(','),
    h3: {
      fontWeight: 600,
    },
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        ':root': {
          fontFamily:
            "'Roboto', 'Inter', 'system-ui', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', sans-serif",
          lineHeight: 1.5,
          fontWeight: 400,
        },
        body: {
          margin: 0,
          minHeight: '100vh',
        },
        '#root': {
          minHeight: '100vh',
        },
        '*, *::before, *::after': {
          boxSizing: 'border-box',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          borderRadius: 999,
        },
      },
    },
  },
});

const resolveInitialMode = (): PaletteMode => {
  if (typeof window === 'undefined') {
    return 'light';
  }
  const stored = window.localStorage.getItem(STORAGE_KEY) as PaletteMode | null;
  if (stored === 'light' || stored === 'dark') {
    return stored;
  }
  const prefersDark = window.matchMedia?.('(prefers-color-scheme: dark)').matches;
  return prefersDark ? 'dark' : 'light';
};

export const ColorModeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [mode, setMode] = useState<PaletteMode>(resolveInitialMode);

  useEffect(() => {
    window.localStorage.setItem(STORAGE_KEY, mode);
  }, [mode]);

  const toggleColorMode = useCallback(() => {
    setMode((prev) => (prev === 'light' ? 'dark' : 'light'));
  }, []);

  const theme = useMemo(() => createTheme(getDesignTokens(mode)), [mode]);

  const value = useMemo(
    () => ({
      mode,
      toggleColorMode,
    }),
    [mode, toggleColorMode],
  );

  return (
    <ColorModeContext.Provider value={value}>
      <MuiThemeProvider theme={theme}>
        <CssBaseline />
        {children}
      </MuiThemeProvider>
    </ColorModeContext.Provider>
  );
};

export const useColorMode = (): ColorModeContextValue => {
  const ctx = useContext(ColorModeContext);
  if (!ctx) {
    throw new Error('useColorMode must be used within ColorModeProvider');
  }
  return ctx;
};
