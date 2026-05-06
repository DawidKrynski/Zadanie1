# Zadanie 6 - testy Cypress + Go

Projekt zawiera testy funkcjonalne Cypress dla aplikacji `saucedemo.com` oraz testy jednostkowe/API Go dla projektu z katalogu `Zestaw4`.

## Realizacja wymagan

- 3.0: 20 przypadkow testowych w CypressJS.
- 3.5: testy funkcjonalne zawieraja 53 asercje.
- 4.0: testy jednostkowe Go do kontrolerow z projektu `Zestaw4`.
- 4.5: testy API pokrywaja endpointy produktow, kategorii i koszyka, razem ze scenariuszami negatywnymi.
- 5.0: konfiguracja BrowserStack znajduje sie w `browserstack.json`.

## Struktura

```text
Zadanie 6/
|-- cypress/
Zadanie 6 — Podsumowanie wymagań i status

✅ 3.0: Stworzono 20 przypadków testowych w CypressJS (cypress/e2e). Link do commita 1

✅ 3.5: Testy funkcjonalne zawierają co najmniej 50 asercji (użycie `should` i `expect` w cypress/e2e) — spełnione. Link do commita 2

✅ 4.0: Stworzono testy jednostkowe dla projektu w `Zestaw4` (kontrolery i logika) z łącznie >=50 asercjami. Link do commita 3

✅ 4.5: Dodano testy API pokrywające endpointy kategorii, produktów i koszyka wraz z scenariuszami negatywnymi (BadRequest / NotFound). Link do commita 4

✅ 5.0: Przygotowano konfigurację BrowserStack: plik `browserstack.json` i skrypt `npm run cypress:browserstack` w `package.json`. Aby uruchomić na BrowserStack uzupełnij `auth.username` i `auth.access_key` (konto Education/GitHub Pack). Link do commita 5


Struktura projektu (istotne pliki):

 - [cypress/e2e](cypress/e2e)
 - [cypress/support/e2e.js](cypress/support/e2e.js)
 - [browserstack.json](browserstack.json)
 - [cypress.config.js](cypress.config.js)
 - [Zestaw4 (Go tests)](Zestaw4)


Jak uruchomić lokalnie

```bash
npm install
npm run cypress:run
```

Testy Go (unit/API)

```bash
cd Zestaw4
go test ./...
```
