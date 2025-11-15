# Análise de Segurança Open Source

Este documento detalha as ferramentas de análise de segurança open source integradas ao pipeline de CI/CD do projeto `gestao_updev`, substituindo a necessidade de ferramentas proprietárias como o Snyk. As ferramentas foram escolhidas com base na sua eficácia para as tecnologias utilizadas (Go para o backend, Node.js/React para o frontend e backoffice) e facilidade de integração em ambientes de CI/CD.

## Ferramentas Implementadas

### Para o Backend (Go)

Para o backend desenvolvido em Go, utilizamos uma combinação de ferramentas para cobrir tanto a análise de dependências quanto a análise estática do código-fonte.

#### 1. `govulncheck` (Go Vulnerability Management)

*   **Tipo:** Análise de Composição de Software (SCA)
*   **Propósito:** Identifica vulnerabilidades conhecidas em módulos Go utilizados no projeto. É a ferramenta oficial do ecossistema Go para este fim, garantindo alta precisão e relevância.
*   **Integração:** Executado como parte do workflow de CI do backend.

#### 2. `GoSec` (Go Security Checker)

*   **Tipo:** Análise Estática de Código (SAST)
*   **Propósito:** Escaneia o código-fonte Go em busca de padrões de código que podem levar a vulnerabilidades de segurança (por exemplo, injeção de SQL, uso inseguro de criptografia, problemas de concorrência).
*   **Integração:** Executado como parte do workflow de CI do backend.

### Para o Frontend e Backoffice (Node.js/React)

Para as aplicações frontend e backoffice baseadas em Node.js/React, focamos na análise de dependências e na manutenção da qualidade do código através do ESLint.

#### 1. `Trivy`

*   **Tipo:** Análise de Composição de Software (SCA)
*   **Propósito:** Embora seja amplamente utilizado para varredura de imagens de contêiner, o Trivy também é muito eficaz para escanear sistemas de arquivos e repositórios em busca de vulnerabilidades em dependências de linguagens como Node.js. Ele oferece uma varredura rápida e abrangente.
*   **Integração:** Executado como parte dos workflows de CI do frontend e backoffice.

#### 2. `ESLint` (com regras de segurança)

*   **Tipo:** Análise Estática de Código (SAST)
*   **Propósito:** O ESLint já é utilizado para manter a qualidade e o estilo do código. Com a configuração adequada, ele também pode ser usado para identificar padrões de código que podem introduzir vulnerabilidades de segurança em JavaScript/TypeScript.
*   **Integração:** Já faz parte dos workflows de CI do frontend e backoffice. As regras de segurança são gerenciadas através da configuração do ESLint.

## Integração no CI/CD (GitHub Actions)

As ferramentas acima são integradas nos respectivos workflows de CI/CD do GitHub Actions (`backend-ci.yml`, `frontend-ci.yml`, `backoffice-ci.yml`). Cada ferramenta é executada em um passo dedicado, e a falha em qualquer uma dessas verificações de segurança resultará na falha do build, garantindo que apenas código seguro seja mesclado.

### Exemplo de Configuração (Passos no Workflow)

```yaml
# Exemplo para GoSec no backend-ci.yml
- name: Run GoSec Security Scanner
  uses: securego/gosec@master
  with:
    args: './...' # Escaneia todos os pacotes Go no diretório atual

# Exemplo para govulncheck no backend-ci.yml
- name: Run govulncheck
  run: |
    go install golang.org/x/vuln/cmd/govulncheck@latest
    govulncheck ./...

# Exemplo para Trivy no frontend-ci.yml ou backoffice-ci.yml
- name: Run Trivy vulnerability scan
  uses: aquasecurity/trivy-action@master
  with:
    scan-type: 'fs' # Escaneia o sistema de arquivos
    severity: 'HIGH,CRITICAL' # Define a severidade mínima para falha
    format: 'table'
    output: 'trivy-results.txt'
    # Adicione um token se necessário para limites de API ou varreduras privadas
    # token: ${{ secrets.TRIVY_TOKEN }}
```

## Próximos Passos

*   Monitorar os resultados das varreduras de segurança e ajustar as configurações ou regras conforme necessário.
*   Explorar a adição de mais ferramentas SAST específicas para JavaScript/TypeScript, se as regras do ESLint não forem suficientes.
*   Considerar a integração de ferramentas de análise de segredos para evitar o vazamento de credenciais.
