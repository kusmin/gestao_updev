# Usando Imagens Docker do GHCR

Este documento descreve como usar as imagens Docker do projeto que são publicadas no GitHub Container Registry (GHCR).

## Visão Geral

As imagens Docker para os serviços da aplicação (como o backend) são construídas e publicadas automaticamente no GHCR sempre que há uma alteração na branch `main`.

Isso permite que você execute a aplicação usando imagens pré-construídas, sem a necessidade de compilar o código-fonte em sua máquina local.

## Imagem do Backend

A imagem para o serviço de backend está disponível no seguinte repositório do GHCR:

- **Nome da Imagem:** `ghcr.io/kusmin/gestao_updev/backend`

As seguintes tags estão disponíveis:
- `latest`: A versão mais recente da branch `main`.
- `<sha>`: A imagem construída a partir de um commit específico (ex: `sha-a1b2c3d`).

## Como Usar

### 1. Autenticação no GHCR

Para baixar imagens de um repositório privado no GHCR, você precisa se autenticar.

**a. Crie um Personal Access Token (PAT)**

- Vá para as [configurações de desenvolvedor](https://github.com/settings/tokens) na sua conta do GitHub.
- Clique em "Generate new token".
- Dê um nome para o token (ex: `docker-ghcr`).
- Selecione o escopo (permissão) `read:packages`.
- Clique em "Generate token" e copie o valor gerado.

**b. Faça o login no Docker**

Abra seu terminal e execute o comando a seguir, substituindo `SEU_USUARIO` pelo seu nome de usuário do GitHub e `SEU_PAT` pelo token que você acabou de criar.

```bash
export CR_PAT="SEU_PAT"
echo $CR_PAT | docker login ghcr.io -u SEU_USUARIO --password-stdin
```

### 2. Executando com Docker Compose

O arquivo `docker-compose.yml` na raiz do projeto já está configurado para usar a imagem do backend do GHCR.

Para iniciar a aplicação, basta executar:

```bash
docker-compose up -d
```

O Docker Compose irá baixar (pull) a imagem `ghcr.io/kusmin/gestao_updev/backend:latest` e iniciar o container.

## Imagem do Frontend

**Nota:** Atualmente, a imagem Docker para o serviço de frontend ainda não está sendo publicada no GHCR. O serviço de frontend ainda é construído a partir do código-fonte localmente.

Quando a imagem do frontend estiver disponível, a seção `frontend` no `docker-compose.yml` será semelhante a esta:

```yaml
  frontend:
    image: ghcr.io/kusmin/gestao_updev/frontend:latest
    depends_on:
      - backend
    ports:
      - "4173:80"
    # A variável de ambiente para a URL da API seria passada aqui
    # environment:
    #   VITE_API_BASE_URL: http://backend:8080/v1
```

## Usando com Portainer

Você também pode implantar a aplicação em um ambiente Portainer.

### 1. Adicionando o Registry no Portainer

Primeiro, você precisa configurar o Portainer para que ele possa se autenticar no GHCR.

1.  No menu do Portainer, vá para **Registries**.
2.  Clique em **Add registry**.
3.  Selecione **Custom registry**.
4.  Preencha os seguintes campos:
    - **Name**: `GitHub GHCR`
    - **Registry URL**: `ghcr.io`
    - **Authentication**: Ative a opção.
    - **Username**: Seu nome de usuário do GitHub.
    - **Password**: O seu Personal Access Token (PAT) que você criou anteriormente.
5.  Clique em **Add registry**.

### 2. Deployando a Stack

Com o registry configurado, você pode implantar a aplicação.

1.  No menu do Portainer, vá para **Stacks**.
2.  Clique em **Add stack**.
3.  Dê um nome para a stack (ex: `gestao-updev`).
4.  No **Web editor**, cole o conteúdo do arquivo `docker-compose.yml` do projeto.
5.  Certifique-se de que o Portainer está usando o registry que você configurou.
6.  Clique em **Deploy the stack**.

O Portainer irá baixar as imagens necessárias (incluindo a do backend do GHCR) e iniciar os containers.
