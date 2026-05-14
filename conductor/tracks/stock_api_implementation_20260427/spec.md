# Especificação Final: tech-challenge-catalog-api

## 1. Visão Geral
A **Catalog API** é o microserviço autoritativo para o catálogo unificado de **Produtos** (com gestão de estoque) e **Serviços/Maintenances** (mão de obra).

## 2. Requisitos Técnicos
- **Stack:** Go, Gin, GORM, PostgreSQL.
- **Arquitetura:** Hexagonal (Ports & Adapters).
- **Mensageria:**
  - **Consome:** Tópico `orders.approved`.
  - **Publica:** Tópico `catalog.events` (eventos de feedback como `StockReserved`, `BackorderCreated`).

## 3. Arquitetura de Fluxo
- **Fase 1 (Validação Síncrona):** O serviço de Ordens realiza chamadas paralelas `POST` aos endpoints `/products/search` e `/maintenances/search` para obter dados de disponibilidade e preço para o orçamento.
- **Fase 2 (Execução Assíncrona):** A API consome o evento `OrderApproved` para iniciar a baixa de estoque (com `SELECT FOR UPDATE` para evitar race conditions) ou criar um `Backorder`.
- **Fase 3 (Feedback Assíncrono):** A API publica eventos como `StockReserved` ou `BackorderCreated` no tópico `catalog.events`, permitindo que o serviço de Ordens atualize o status para o cliente final.

## 4. Definições de API
### Produtos
- `POST /products/search`: Consulta em lote.
- `GET /products`: Listagem paginada.
- `GET /products/{id}`: Detalhe público.
- `POST /products`: (Admin) Criação de produto.
- `PUT /products/{id}`: (Admin) Edição de dados de catálogo.
- `GET /products/{id}/inventory`: (Admin) Consulta de dados de inventário.
- `PUT /products/{id}/inventory`: (Admin) Edição de dados de inventário.
- `DELETE /products/{id}`: (Admin) Deleção lógica.

### Serviços (Maintenances)
- `POST /maintenances/search`: Consulta em lote.
- `GET /maintenances`: Listagem paginada.
- `GET /maintenances/{id}`: Detalhe.
- `POST /maintenances`: (Admin) Criação de serviço.
- `PUT /maintenances/{id}`: (Admin) Edição de serviço.
- `DELETE /maintenances/{id}`: (Admin) Deleção lógica.

## 5. Modelo de Dados
- **Product:** `id`, `name`, `description`, `price`, `stock_quantity`, `min_threshold`, `max_stock_level`, `replenishment_strategy`, `supplier_sku`, `default_lead_time_days`.
- **Maintenance:** `id`, `description`, `base_price`, `estimated_duration`.
- **Backorder:** `id`, `product_id`, `order_id`, `quantity`, `status`, `estimated_arrival`.
