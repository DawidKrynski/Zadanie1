# Zadanie 4 – E-commerce REST API (Go + Echo + GORM)

Aplikacja REST API napisana w Go z użyciem frameworka Echo oraz ORM GORM. Wykorzystuje SQLite jako bazę danych.

## Docker Hub

Obraz dostępny na: https://hub.docker.com/r/dawidk12/zadanie4

```bash
docker pull dawidk12/zadanie4
docker run -p 8080:8080 dawidk12/zadanie4
```

## Realizacja zadania

✅ 3.0 Aplikacja we frameworku Echo w Go z kontrolerem Produktów zgodnym z CRUD  
✅ 3.5 Model Produktów z GORM, obsługa CRUD w kontrolerze zamiast listy, baza SQLite  
✅ 4.0 Model Koszyka (Cart + CartItem) z endpointami do zarządzania koszykiem  
✅ 4.5 Model Kategorii z relacją Category → Product (foreignKey CategoryID)  
✅ 5.0 Zapytania pogrupowane w GORM-owe scope'y (ByCategory, PriceRange, InStock, Paginate, WithCategory)  

## Modele

- **Category** – kategoria produktów
- **Product** – produkt (należy do kategorii, relacja `CategoryID`)
- **Cart** – koszyk zakupowy
- **CartItem** – pozycja w koszyku (należy do koszyka i produktu)

## Endpointy

### Produkty
| Metoda | Ścieżka | Opis |
|--------|---------|------|
| GET | /products | Lista produktów (filtry: category_id, min_price, max_price, in_stock, page, page_size) |
| GET | /products/:id | Pojedynczy produkt |
| POST | /products | Utwórz produkt |
| PUT | /products/:id | Aktualizuj produkt |
| DELETE | /products/:id | Usuń produkt |

### Kategorie
| Metoda | Ścieżka | Opis |
|--------|---------|------|
| GET | /categories | Lista kategorii |
| GET | /categories/:id | Kategoria z produktami |
| POST | /categories | Utwórz kategorię |
| PUT | /categories/:id | Aktualizuj kategorię |
| DELETE | /categories/:id | Usuń kategorię |

### Koszyk
| Metoda | Ścieżka | Opis |
|--------|---------|------|
| POST | /carts | Utwórz koszyk |
| GET | /carts/:id | Pobierz koszyk z itemami |
| POST | /carts/:id/items | Dodaj produkt do koszyka |
| PUT | /carts/:id/items/:itemId | Aktualizuj ilość |
| DELETE | /carts/:id/items/:itemId | Usuń item |
| DELETE | /carts/:id | Usuń koszyk |

## Przykłady

```bash
# Utwórz kategorię
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"Elektronika","description":"Sprzet elektroniczny"}'

# Utwórz produkt
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","price":2999.99,"quantity":10,"category_id":1}'

# Produkty z filtrami
curl "http://localhost:8080/products?category_id=1&min_price=100&in_stock=true&page=1&page_size=5"
```