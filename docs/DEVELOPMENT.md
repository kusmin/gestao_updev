# Guia de Desenvolvimento Backend

Este documento fornece instruções para o desenvolvimento do backend.

## Pré-requisitos

- [Go](https://go.dev/dl/) (versão 1.24 ou superior)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [golangci-lint](https://golangci-lint.run/usage/install/)

## Linter

Usamos o `golangci-lint` para garantir a qualidade e a consistência do código.

Para executar o linter localmente, a partir da pasta `backend`:

```bash
golangci-lint run
```

O linter também é executado automaticamente em cada pull request.

## Documentação da API (Swagger)

A documentação da API é gerada automaticamente a partir dos comentários do código usando `swag`.

### Visualizando a Documentação

1.  **Inicie a aplicação:** (instruções a serem adicionadas, ex: `docker-compose up`)
2.  **Acesse no navegador:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### Atualizando a Documentação

Após adicionar ou modificar os comentários de anotação da API, regenere os arquivos da documentação com o seguinte comando, a partir da pasta `backend`:

```bash
go run github.com/swaggo/swag/cmd/swag init -g cmd/api/main.go
```
