# Auction App

## Descrição

Este é um sistema de leilões desenvolvido em Go. O sistema permite criar leilões, fazer lances e consultar leilões ativos. A aplicação pode ser executada localmente usando Docker e Docker Compose.


**Acesse a aplicação:**

   A aplicação estará disponível em `http://localhost:8080`.

### Exemplos de Requisição

- **Para criar leilão:**

   ```sh
   curl -X POST http://localhost:8080/auction \
   -H "Content-Type: application/json" \
   -d '{
        "product_name": "Gol",
        "category": "Car",
        "description": "Volkswagen brand car",
        "condition": 1
       }'
   ```

- **Para buscar leilão por ID:**

   ```sh
   curl -X GET "http://localhost:8080/auction/{auctionId}" \
   -H "Content-Type: application/json"
   ```

- **Para buscar leilões por parâmetros de consulta:**

   ```sh
   curl -X GET "http://localhost:8080/auction?category=Pet%20Supplies&condition=1&status=0" \
   -H "Content-Type: application/json"
   ```


   ** O tempo de encerramento automático pode ser encerrado usando a seguinte variável no .env

```bash
   AUCTION_DURATION=30s
```
