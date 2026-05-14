# Design Técnico: tech-challenge-catalog-api

Este documento consolida o design do microserviço de Catálogo, unificando a gestão de Produtos (Estoque) e Serviços (Maintenances).

## 1. Visão Geral
A Catalog API é a única fonte da verdade para itens de Ordem de Serviço.
- **Produtos:** Controle de saldo, gatilhos de segurança (Min Threshold) e reabastecimento dual.
- **Serviços (Maintenances):** Padronização de mão de obra e preços tabelados.

## 2. Fluxo Híbrido (Ponta-a-Ponta)
### Fase 1: Validação (Síncrona)
Durante a montagem do orçamento, o serviço de Ordens consulta `POST /api/v1/products/search` para produtos e `POST /api/v1/maintenances/search` para serviços. As APIs retornam a disponibilidade e os preços base para o snapshot financeiro.

### Fase 2: Execução (Assíncrona)
Após a aprovação da OS, a API consome o evento `OrderApproved`:
1. **Para Peças:** Realiza o `SELECT FOR UPDATE` para garantir atomicidade. Se faltar saldo, entra no fluxo de Backorder.
2. **Para Serviços:** Registra o início da prestação para auditoria.

## 3. Estratégia de Reabastecimento (Replenishment)
- **IMMEDIATE:** Pedido urgente via API REST ao fornecedor.
- **BATCH:** Consolidação de itens de baixo valor para pedido em lote.
- **Reconciliação:** Caso o fornecedor rejeite um pedido imediato, o status do Backorder muda para `REJECTED_BY_SUPPLIER`, disparando um alerta administrativo no dashboard de Backoffice.

## 4. Segurança e Performance
- **RBAC:** Headers `X-User-Role` (Admin/Mechanic/Attendant).
- **Integridade:** Lock pessimista no banco para gestão de unidades físicas.
