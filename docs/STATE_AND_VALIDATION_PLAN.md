# Plano de Implementação: Zod e Zustand para Frontend e Backoffice

## 1. Visão Geral

Este documento descreve a proposta de adoção da biblioteca `zod` para validação de dados e da `zustand` para gerenciamento de estado global nos projetos `frontend` e `backoffice`.

A análise atual revela que:
- Não há validação de dados em tempo de execução para as respostas da API, o que pode levar a erros inesperados caso a API retorne um formato de dados diferente do esperado.
- Não há uma solução padronizada para gerenciamento de estado global do cliente, como informações de autenticação do usuário, tenant ativo ou estado da UI (ex: tema claro/escuro).

## 2. Ferramentas Propostas

### 2.1. Zod (Validação de Dados)

**Por quê?**
- **Segurança de Tipos em Tempo de Execução:** Enquanto o `openapi-typescript` nos dá segurança de tipos em tempo de compilação, o Zod garante que os dados recebidos da API em tempo de execução correspondam ao schema esperado.
- **Mensagens de Erro Claras:** Facilita a depuração de inconsistências de dados.
- **Leve e Simples:** Possui uma API intuitiva e é fácil de integrar.

### 2.2. Zustand (Gerenciamento de Estado)

**Por quê?**
- **Simplicidade e Baixo Boilerplate:** `Zustand` oferece uma API minimalista baseada em hooks, tornando a criação e o uso de stores globais muito mais simples em comparação com o Redux.
- **Complementa o React Query:** O `@tanstack/react-query` (já em uso no `frontend`) é excelente para o estado do servidor (dados da API). O Zustand é ideal para o estado do cliente (estado da UI, autenticação), preenchendo a lacuna existente.
- **Performance:** Evita renderizações desnecessárias, pois os componentes só são atualizados quando a parte específica do estado que eles consomem muda.

## 3. Plano de Implementação

A implementação será idêntica para os projetos `frontend` e `backoffice`.

### 3.1. Instalação das Dependências

Execute o seguinte comando na raiz de cada projeto (`frontend/` e `backoffice/`):

```bash
pnpm add zod zustand
```

### 3.2. Estrutura de Diretórios

Sugerimos a seguinte estrutura de diretórios em `src/` para ambos os projetos:

```
src/
├── ...
├── lib/
│   └── apiClient.ts  // Onde a validação com Zod pode ser integrada
├── schemas/
│   └── index.ts      // Arquivo para exportar todos os schemas Zod
│   └── user.ts       // Exemplo: schemas de usuário
├── store/
│   └── index.ts      // Arquivo para exportar todos os stores Zustand
│   └── auth.ts       // Exemplo: store de autenticação
└── ...
```

## 4. Exemplos de Uso

### 4.1. Validando Dados da API com Zod

**Passo 1: Definir o Schema**

Crie um schema Zod que corresponda à estrutura de dados esperada da API.

`src/schemas/user.ts`:
```typescript
import { z } from 'zod';

// Supondo que a API retorna um usuário com esta estrutura
export const userSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  email: z.string().email(),
  createdAt: z.string().datetime(),
});

export const userListSchema = z.array(userSchema);
```

**Passo 2: Validar na Chamada da API**

Modifique o `apiClient` ou a função de busca de dados para validar a resposta antes de retorná-la.

Exemplo de uma função de busca de dados:
```typescript
import { userSchema } from '../schemas/user';
import { components } from '../types/api'; // Tipos gerados pelo openapi-typescript

type User = components['schemas']['User'];

async function fetchUser(userId: string): Promise<User> {
  const response = await fetch(`/api/users/${userId}`);
  const data = await response.json();

  // Valida os dados recebidos com o schema Zod
  const validatedData = userSchema.parse(data);

  // Se a validação passar, retorna os dados.
  // O tipo de `validatedData` será inferido como o tipo do schema,
  // garantindo consistência com o tipo `User`.
  return validatedData;
}
```

### 4.2. Gerenciando Estado Global com Zustand

**Passo 1: Criar o Store (Slice)**

Vamos criar um store para gerenciar o estado de autenticação.

`src/store/auth.ts`:
```typescript
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthState {
  token: string | null;
  user: {
    id: string;
    name: string;
  } | null;
  setToken: (token: string | null) => void;
  setUser: (user: AuthState['user']) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  // O middleware `persist` salva o estado no localStorage automaticamente
  persist(
    (set) => ({
      token: null,
      user: null,
      setToken: (token) => set({ token }),
      setUser: (user) => set({ user }),
      logout: () => set({ token: null, user: null }),
    }),
    {
      name: 'auth-storage', // Nome da chave no localStorage
    }
  )
);
```
*Nota: Para usar o middleware `persist`, talvez seja necessário instalar uma dependência adicional, dependendo da configuração do projeto (`pnpm add zustand-middleware-persist`). A documentação oficial do Zustand deve ser consultada.*

**Passo 2: Usar o Store em um Componente**

Agora você pode usar o hook `useAuthStore` em qualquer componente para acessar e modificar o estado.

```tsx
import React from 'react';
import { useAuthStore } from '../store/auth';

const UserProfile: React.FC = () => {
  const { user, logout } = useAuthStore();

  if (!user) {
    return <div>Você não está logado.</div>;
  }

  return (
    <div>
      <h1>Bem-vindo, {user.name}!</h1>
      <button onClick={logout}>Sair</button>
    </div>
  );
};

const LoginButton: React.FC = () => {
  const { setToken, setUser } = useAuthStore();

  const handleLogin = () => {
    // Simula um login
    const fakeToken = 'abc-123';
    const fakeUser = { id: 'uuid-456', name: 'Usuário Exemplo' };
    
    setToken(fakeToken);
    setUser(fakeUser);
  };

  return <button onClick={handleLogin}>Entrar</button>;
};
```

## 5. Próximos Passos

1.  Aprovar este plano.
2.  Criar as tarefas no gerenciador de projetos para instalar as dependências e criar a estrutura de diretórios inicial.
3.  Refatorar gradualmente as chamadas de API existentes para usar a validação com Zod.
4.  Refatorar o gerenciamento de estado de autenticação e outros estados globais para usar o store do Zustand.
