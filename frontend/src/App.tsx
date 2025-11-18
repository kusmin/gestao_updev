import AppRouter from './routes/AppRouter';
import { ColorModeProvider } from './theme/ColorModeProvider';
import { AuthProvider } from './contexts/AuthContext';
import { SnackbarProvider } from './contexts/SnackbarProvider';

const App: React.FC = () => {
  return (
    <ColorModeProvider>
      <AuthProvider>
        <SnackbarProvider>
          <AppRouter />
        </SnackbarProvider>
      </AuthProvider>
    </ColorModeProvider>
  );
};

export default App;
