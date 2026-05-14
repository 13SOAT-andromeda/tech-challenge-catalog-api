# Especificação: Estratégia de Infraestrutura e IaC (tech-challenge-catalog-api)

Esta especificação define como a infraestrutura da **tech-challenge-catalog-api** será gerenciada, seguindo o modelo de centralização de Platform Engineering adotado pelo projeto.

## 1. Visão Geral da Estratégia
A Catalog API adotará um modelo **Cloud-Free**, onde o repositório da aplicação não gerencia recursos de infraestrutura da AWS diretamente. A responsabilidade é dividida entre os repositórios de IaC centrais e os manifestos Kubernetes locais.

## 2. Divisão de Responsabilidades (IaC Central)

### 2.1. iac-tech-challenge-infra (Computação e Rede)
Responsável por provisionar o ambiente onde a aplicação será executada.
*   **ECR:** Provisionamento do registro de container `tech-challenge-catalog-api-repo`.
*   **EKS:** Utilização do cluster compartilhado `eks-tech-challenge`.
*   **Networking:** Configuração de Ingress/Gateway para roteamento do path `/api/v1/catalog`.
*   **IAM:** Criação de IAM Roles for Service Accounts (IRSA) se necessário.

### 2.2. iac-tech-challenge-data (Persistência e Mensageria)
Responsável por gerenciar o estado e a comunicação entre serviços.
*   **Banco de Dados (RDS):** Criação de um database isolado `catalog_db` na instância PostgreSQL compartilhada.
*   **Credenciais:** Gerenciamento de usuários e senhas exclusivos para a Catalog API (via Secrets Manager).
- **Broker:** Configuração das filas e permissões de acesso ao tópico `orders.approved` (para consumo) e provisionamento do novo tópico `catalog.events` (para publicação).

## 3. Manifesto Kubernetes (Local)
O repositório `tech-challenge-catalog-api` manterá apenas os arquivos necessários para a orquestração no cluster:
*   **Deployment:** Configuração de replicas, recursos (CPU/Memória) e variáveis de ambiente.
*   **Service:** Exposição interna para o Ingress.
*   **HPA:** Escalonamento automático baseado em demanda de processamento assíncrono.
*   **ConfigMap/Secrets:** Mapeamento de configurações e credenciais de banco.

## 4. Fluxo de CI/CD
1.  **Build:** Geração da imagem Docker.
2.  **Push:** Envio para o ECR (gerenciado pelo `iac-tech-challenge-infra`).
3.  **Deploy:** Aplicação dos manifestos K8s no EKS (gerenciado pelo `iac-tech-challenge-infra`).

## 5. Monitoramento
*   **Datadog:** Integração via Agent já existente no cluster EKS.
*   **Logs:** Coleta centralizada via CloudWatch/Datadog Logs.
