# Contratos JSON Finais: Catalog API

## 1. Produtos (`/products`)
- `POST /search`: Busca por IDs.
- `GET /`: Lista com paginação e filtros.
- `GET /{id}`: Detalhe público.
- `POST /`: (Admin) Cria produto.
- `PUT /{id}`: (Admin) Edita dados de catálogo.
- `DELETE /{id}`: (Admin) Desativa produto.
- `GET /{id}/inventory`: (Admin) Consulta inventário.
- `PUT /{id}/inventory`: (Admin) Edita inventário.

## 2. Serviços (`/maintenances`)
- `POST /search`: Busca por IDs.
- `GET /`: Lista com paginação e filtros.
- `GET /{id}`: Detalhe.
- `POST /`: (Admin) Cria serviço.
- `PUT /{id}`: (Admin) Edita serviço.
- `DELETE /{id}`: (Admin) Desativa serviço.

## 3. Mensageria (`catalog.events`)
### Evento: `BackorderCreated`
```json
{
  "event_type": "BackorderCreated",
  "order_id": "uuid-da-os",
  "product_id": "uuid-do-produto",
  "quantity_missing": 2,
  "estimated_arrival": "2026-05-15T18:00:00Z"
}
```
