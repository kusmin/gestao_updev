import AppRouter from './routes/AppRouter';
import { ColorModeProvider } from './theme/ColorModeProvider';
import { AuthProvider } from './contexts/AuthContext';

const App: React.FC = () => {
  return (
    <ColorModeProvider>
      <AuthProvider>
        <AppRouter />
      </AuthProvider>
    </ColorModeProvider>
  );
};

export default App;
