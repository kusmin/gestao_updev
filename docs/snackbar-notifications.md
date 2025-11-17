# Sistema de Notificações (Snackbar/Toast)

Este documento descreve o plano e o guia de implementação para um sistema de notificações (também conhecido como snackbar ou toast) para as aplicações `frontend` e `backoffice`.

## Objetivo

Prover um sistema centralizado e fácil de usar para exibir mensagens de feedback para o usuário, como confirmações de sucesso, erros ou alertas.

## Estratégia

A estratégia adotada é a criação de um sistema de notificações reutilizável, baseado no componente `Snackbar` do Material-UI, que já é uma dependência em ambos os projetos.

O sistema será composto por:
-   Um `SnackbarProvider`: Um provedor de contexto que gerencia o estado e a exibição das notificações.
-   Um hook `useSnackbar`: Um hook customizado para disparar notificações de qualquer componente da aplicação.

Esta abordagem centraliza a lógica e a aparência das notificações, garantindo consistência e facilitando a manutenção.

## Plano de Implementação

### Passo 1: Criar o Contexto e o Provedor do Snackbar

1.  **Criar o arquivo**: `frontend/src/contexts/SnackbarContext.tsx`
2.  **Definir o contexto**: Usar `createContext` para criar um `SnackbarContext`.
3.  **Criar o `SnackbarProvider`**:
    -   Gerenciar o estado da notificação (mensagem, severidade, visibilidade).
    -   Renderizar o componente `Snackbar` e `Alert` do Material-UI.
    -   Expor uma função `showSnackbar` através do contexto.

### Passo 2: Criar o Hook `useSnackbar`

1.  **Criar o arquivo**: `frontend/src/hooks/useSnackbar.ts`
2.  **Implementar o hook**: O hook `useSnackbar` simplesmente consumirá o `SnackbarContext` e retornará a função `showSnackbar`.

### Passo 3: Integrar o `SnackbarProvider` na Aplicação

1.  **Envolver a aplicação**: No arquivo `frontend/src/App.tsx`, envolver o `AppRouter` com o `SnackbarProvider`. Isso garantirá que o hook `useSnackbar` possa ser usado em qualquer página.

### Passo 4: Exemplo de Uso

Após a implementação, para exibir uma notificação em qualquer componente, o uso será simples:

```tsx
import { useSnackbar } from '../hooks/useSnackbar';

const MyComponent = () => {
  const { showSnackbar } = useSnackbar();

  const handleClick = () => {
    showSnackbar('Operação realizada com sucesso!', 'success');
  };

  return <button onClick={handleClick}>Realizar Operação</button>;
};
```

### Passo 5: Replicar para o Backoffice

A mesma estrutura (`SnackbarContext`, `useSnackbar`) será replicada para o projeto `backoffice`, seguindo os mesmos passos.

## Documentação e Boas Práticas

-   A severidade da notificação (`error`, `warning`, `info`, `success`) deve ser usada corretamente para fornecer o feedback visual adequado ao usuário.
-   As mensagens devem ser curtas e informativas.
-   Evitar o uso excessivo de notificações para não sobrecarregar o usuário.
