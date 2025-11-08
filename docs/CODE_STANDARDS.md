# Padrões de Código

Este documento define os padrões de código e estilo para o projeto.

## Geral

- **Idioma:** Todo o código, comentários e documentação devem ser escritos em português.
- **Formatação:** Use `prettier` para o frontend e `gofmt` para o backend para garantir um estilo consistente.

## Backend (Go)

- **Estrutura do projeto:** Siga o layout de projeto padrão da comunidade Go.
- **Nomenclatura:** Use `camelCase` para variáveis e `PascalCase` para tipos e funções exportadas.
- **Tratamento de erros:** Verifique sempre os erros e retorne-os quando apropriado. Não use `panic` para erros esperados.
- **Logs:** Use logs estruturados (ex: `zap` ou `logrus`).

## Frontend (React + Vite)

- **Componentes:** Sempre funcionais, com hooks e React Query para operações assíncronas.
- **Tipagem:** TypeScript estrito; tipos de API gerados automaticamente via `openapi-typescript`.
- **Estilo:** Tailwind ou CSS Modules; evite estilos inline complexos.
- **Estado:** Estado remoto com React Query; use `Zustand`/Context apenas quando necessário.

## Banco de Dados

- **Nomenclatura:** Use `snake_case` para tabelas e colunas.
- **Migrations:** Todas as alterações de esquema devem ser feitas através de arquivos de migração.
