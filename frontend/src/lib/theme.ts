import { createTheme } from '@mui/material/styles';

const theme = createTheme({
  palette: {
    primary: {
      main: '#0044ff',
    },
    secondary: {
      main: '#00b0b9',
    },
    background: {
      default: '#f7f9fc',
      paper: '#ffffff',
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
    ].join(', '),
    h3: {
      fontWeight: 600,
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          borderRadius: 999,
        },
      },
    },
    MuiCssBaseline: {
      styleOverrides: {
        ':root': {
          fontFamily:
            "'Roboto', 'Inter', 'system-ui', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', sans-serif",
          lineHeight: 1.5,
          fontWeight: 400,
          color: '#0f172a',
          backgroundColor: '#f7f9fc',
        },
        body: {
          margin: 0,
          minHeight: '100vh',
          backgroundColor: '#f7f9fc',
        },
        '#root': {
          minHeight: '100vh',
        },
        '*, *::before, *::after': {
          boxSizing: 'border-box',
        },
      },
    },
  },
});

export default theme;
