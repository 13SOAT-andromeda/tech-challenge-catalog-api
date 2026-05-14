# Plano de Implementação: Catalog API

## Fase 1: Setup e Infra
- [ ] Criar repositório `tech-challenge-catalog-api`.
- [ ] Configurar DB isolado `catalog_db`.

## Fase 2: Domínio e Persistência
- [ ] Implementar entidade `Product` e `Maintenance`.
- [ ] Criar Migrations para o catálogo completo.
- [ ] Implementar Repositório com Locks para produtos.

## Fase 3: API e Backoffice
- [ ] Implementar CRUD de Produtos.
- [ ] Implementar CRUD de Maintenances (Serviços).
- [ ] Implementar `POST /products/search` (valida peças).
- [ ] Implementar `POST /maintenances/search` (valida serviços).

## Fase 4: Integração Assíncrona
- [ ] Consumir `OrderApproved`.
- [ ] Processar baixa de peças e log de serviços realizados.

## Fase 5: Replenishment e Finalização
- [ ] Integrar com fornecedor (Produtos).
- [ ] Job de Backorder.
