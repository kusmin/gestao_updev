# Contexto do Gemini: `backend/cmd/api`

Este diretório contém o ponto de entrada principal (`main.go`) para a aplicação da API do backend.

## Responsabilidades

*   **Inicialização:** É responsável por inicializar todos os componentes essenciais da aplicação, como:
    *   Configuração (carregada de variáveis de ambiente).
    *   Logger (Zap).
    *   Conexão com o banco de dados (PostgreSQL com GORM).
    *   Repositório (`repository`).
    *   Gerenciador de JWT (`auth`).
    *   Camada de serviço (`service`).
    *   Servidor HTTP (Gin).
*   **Injeção de Dependência:** Orquestra a injeção de dependências entre os diferentes pacotes (ex: injeta o repositório no serviço, e o serviço nos handlers).
*   **Definição de Rotas:** Configura todas as rotas da API, associando os endpoints aos seus respectivos handlers e middlewares.
*   **Execução do Servidor:** Inicia o servidor HTTP para começar a ouvir as requisições.

Qualquer alteração relacionada à inicialização de novos serviços, configuração de middlewares globais ou adição de novas rotas de alto nível deve ser feita aqui.
