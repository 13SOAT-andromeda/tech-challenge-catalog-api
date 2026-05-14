# Spec: Replenishment Strategy

**Author(s):** Tech Challenge Team
**Status:** DRAFT
**Date:** 2026-05-14

## 1. Visão Geral

Esta especificação detalha as estratégias de reabastecimento de estoque (`replenishment`) para a `tech-challenge-catalog-api`. O objetivo é garantir a disponibilidade de produtos, otimizando custos e respondendo de forma eficaz à demanda. A estratégia é definida no nível de cada produto e acionada automaticamente quando o sistema detecta uma necessidade de reabastecimento após a aprovação de uma ordem de serviço.

## 2. Estratégias Suportadas

O sistema suporta duas estratégias distintas: `IMMEDIATE` e `BATCH`.

### 2.1. Estratégia: `IMMEDIATE`

**Descrição:** Reabastecimento just-in-time para itens de alta prioridade.

*   **Gatilho:**
    1.  Um evento `OrderApproved` é consumido.
    2.  O sistema verifica o estoque para um produto da ordem.
    3.  O estoque é insuficiente para atender à demanda.
    4.  A estratégia do produto é `IMMEDIATE`.

*   **Fluxo de Execução:**
    1.  Um registro de `Backorder` é criado com o status `PENDING`.
    2.  Imediatamente, o sistema monta uma requisição de compra.
    3.  Uma chamada de API é feita para o endpoint do fornecedor correspondente (`POST /purchase-order`).
    4.  **Se a chamada for bem-sucedida:** O status do `Backorder` é atualizado para `ORDERED`.
    5.  **Se a chamada falhar:** O status do `Backorder` é atualizado para `FAILED`, e um alerta é enviado para a equipe de backoffice.

*   **Caso de Uso:**
    *   Produtos de alto valor agregado.
    *   Itens críticos que não podem faltar no estoque.
    *   Produtos com fornecedores que oferecem APIs para pedidos em tempo real.

### 2.2. Estratégia: `BATCH`

**Descrição:** Reabastecimento em lote para otimização de custos.

*   **Gatilho:**
    1.  Um evento `OrderApproved` é consumido.
    2.  O sistema verifica o estoque e constata a insuficiência.
    3.  A estratégia do produto é `BATCH`.

*   **Fluxo de Execução:**
    1.  Um registro de `Backorder` é criado com o status `PENDING`.
    2.  Nenhuma ação imediata é tomada. O sistema aguarda o processo em lote.

*   **Processo em Lote (Cron Job):**
    *   **Frequência:** O cron job será executado em intervalos configuráveis (e.g., a cada 24 horas, durante a madrugada).
    *   **Lógica:**
        1.  O job faz uma varredura na tabela `Backorder` e coleta todos os registros com status `PENDING`.
        2.  Os itens são agrupados por `supplier_sku` (ou seja, por fornecedor).
        3.  Para cada fornecedor, o sistema consolida todos os itens necessários em um único pedido de compra.
        4.  Uma chamada de API é feita para o fornecedor com o pedido em lote.
        5.  **Se a chamada for bem-sucedida:** Todos os `Backorder`s incluídos no lote têm seu status atualizado para `ORDERED`.
        6.  **Se a chamada falhar:** Os `Backorder`s permanecem como `PENDING`, e um alerta é gerado para a equipe de backoffice.

*   **Caso de Uso:**
    *   Produtos de baixo custo e alta rotatividade.
    *   Itens não críticos, onde um pequeno atraso no reabastecimento é aceitável.
    *   Otimização de custos de frete e obtenção de melhores preços em compras por volume.

## 3. Considerações de Implementação

*   **Rastreabilidade:** O `order_id` na tabela `Backorder` é crucial para rastrear qual pedido de cliente originou a necessidade de reabastecimento.
*   **Configurabilidade:** A frequência do cron job para a estratégia `BATCH` deve ser configurável via variáveis de ambiente.
*   **Alertas:** Falhas na comunicação com a API do fornecedor ou rejeições de pedidos devem gerar alertas visíveis no dashboard de backoffice.
