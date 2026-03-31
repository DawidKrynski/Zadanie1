# Zadanie 2 - Scalatra REST API & Docker

Autor: Dawid Kryński (dawidk12)

Aplikacja oparta na frameworku **Scalatra**, implementująca w pamięci operacje CRUD dla modułów e-commerce oraz skonteneryzowana w oparciu o obraz zdefiniowany w wymaganiach.

## Obraz Docker Hub

**https://hub.docker.com/r/dawidk12/scalatra-ebiznes**

## Uruchomienie

```bash
docker pull dawidk12/scalatra-ebiznes:latest
docker run -p 8080:8080 dawidk12/scalatra-ebiznes:latest
```

Ngrok:
```bash
bash ngrok.sh
```

## Endpointy

| Metoda | Produkty | Kategorie | Koszyk |
|--------|----------|-----------|--------|
| GET all | `/products/` | `/categories/` | `/cart/` |
| GET by id | `/products/:id` | `/categories/:id` | `/cart/:id` |
| POST | `/products/` | `/categories/` | `/cart/` |
| PUT | `/products/:id` | `/categories/:id` | `/cart/:id` |
| DELETE | `/products/:id` | `/categories/:id` | `/cart/:id` |

## CORS

Skonfigurowany dla dwóch hostów:
- `http://localhost:3000`
- `http://127.0.0.1:5500`

Dozwolone metody: GET, POST, PUT, DELETE, OPTIONS
