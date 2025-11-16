import AppRouter from './routes/AppRouter';
import { ColorModeProvider } from './theme/ColorModeProvider';

const App: React.FC = () => {
  return (
    <ColorModeProvider>
      <AppRouter />
    </ColorModeProvider>
  );
};

export default App;
