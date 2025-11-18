# Guia de Instalação Manual do Storybook

Este guia detalha os passos para instalar e configurar o Storybook manualmente nos projetos `frontend` e `backoffice`. O Storybook será usado para desenvolver, documentar e testar componentes de UI de forma isolada.

## 1. Instalação das Dependências

Para cada projeto (`frontend` e `backoffice`), execute os seguintes comandos para adicionar as dependências de desenvolvimento do Storybook:

```bash
pnpm add -D storybook @storybook/addon-essentials @storybook/addon-interactions @storybook/addon-links @storybook/blocks @storybook/react @storybook/react-vite @storybook/test
```

## 2. Estrutura de Diretórios

Crie um diretório `.storybook` na raiz de cada projeto (`frontend/` e `backoffice/`).

```
frontend/
└── .storybook/
    ├── main.ts
    ├── preview.ts
    └── manager.ts
backoffice/
└── .storybook/
    ├── main.ts
    ├── preview.ts
    └── manager.ts
```

## 3. Arquivos de Configuração

### 3.1. `main.ts` (Configuração Principal)

Este arquivo configura o Storybook para localizar suas histórias e define os addons a serem usados.

**Conteúdo para `frontend/.storybook/main.ts`:**

```typescript
import type { StorybookConfig } from '@storybook/react-vite';

const config: StorybookConfig = {
  stories: ['../src/**/*.mdx', '../src/**/*.stories.@(js|jsx|mjs|ts|tsx)'],
  addons: [
    '@storybook/addon-links',
    '@storybook/addon-essentials',
    '@storybook/addon-interactions',
  ],
  framework: {
    name: '@storybook/react-vite',
    options: {},
  },
  docs: {
    autodocs: 'tag',
  },
};

export default config;
```

**Conteúdo para `backoffice/.storybook/main.ts`:**

```typescript
import type { StorybookConfig } from '@storybook/react-vite';

const config: StorybookConfig = {
  stories: ['../src/**/*.mdx', '../src/**/*.stories.@(js|jsx|mjs|ts|tsx)'],
  addons: [
    '@storybook/addon-links',
    '@storybook/addon-essentials',
    '@storybook/addon-interactions',
  ],
  framework: {
    name: '@storybook/react-vite',
    options: {},
  },
  docs: {
    autodocs: 'tag',
  },
};

export default config;
```

### 3.2. `preview.ts` (Configuração de Pré-visualização)

Este arquivo configura como as histórias são renderizadas no painel de pré-visualização. Pode ser usado para aplicar estilos globais, provedores de tema, etc.

**Conteúdo para `frontend/.storybook/preview.ts`:**

```typescript
import type { Preview } from '@storybook/react';

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
  },
};

export default preview;
```

**Conteúdo para `backoffice/.storybook/preview.ts`:**

```typescript
import type { Preview } from '@storybook/react';

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
  },
};

export default preview;
```

### 3.3. `manager.ts` (Configuração do Gerenciador de UI)

Este arquivo configura a UI do próprio Storybook (o painel lateral, cabeçalho, etc.). Para uma configuração básica, pode ser deixado vazio ou com configurações mínimas.

**Conteúdo para `frontend/.storybook/manager.ts` e `backoffice/.storybook/manager.ts`:**

```typescript
// Nenhuma configuração personalizada por padrão
// import { addons } from '@storybook/manager-api';
// import { create } from '@storybook/theming';

// addons.setConfig({
//   theme: create({
//     base: 'light',
//     brandTitle: 'Meu Projeto Storybook',
//   }),
// });
```

## 4. Scripts no `package.json`

Adicione os seguintes scripts ao `scripts` de cada `package.json` (`frontend/package.json` e `backoffice/package.json`):

```json
{
  "scripts": {
    // ... scripts existentes
    "storybook": "storybook dev -p 6006",
    "build-storybook": "storybook build"
  }
}
```

## 5. Exemplo de Componente e Story

Após a configuração, você pode criar suas histórias. Por exemplo, para um componente `Button`:

`src/components/Button.tsx`:

```typescript jsx
import React from 'react';

interface ButtonProps {
  label: string;
  onClick?: () => void;
  primary?: boolean;
}

export const Button: React.FC<ButtonProps> = ({ label, onClick, primary = false }) => {
  const mode = primary ? 'storybook-button--primary' : 'storybook-button--secondary';
  return (
    <button
      type="button"
      className={['storybook-button', mode].join(' ')}
      onClick={onClick}
    >
      {label}
    </button>
  );
};
```

`src/components/Button.stories.tsx`:

```typescript jsx
import type { Meta, StoryObj } from '@storybook/react';
import { Button } from './Button';

const meta = {
  title: 'Example/Button',
  component: Button,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
    primary: { control: 'boolean' },
    onClick: { action: 'clicked' },
  },
} satisfies Meta<typeof Button>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
  args: {
    primary: true,
    label: 'Button',
  },
};

export const Secondary: Story = {
  args: {
    label: 'Button',
  },
};
```

Com este guia, você pode configurar o Storybook em seus projetos e começar a desenvolver suas histórias de componentes.
