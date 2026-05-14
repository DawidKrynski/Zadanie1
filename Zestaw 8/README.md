Zadanie 8 - OAuth2 + logowanie po stronie serwera

Aplikacja z zadania 5 rozszerzona o logowanie lokalne oraz OAuth2.
Klient React nie tworzy klienta OAuth bezposrednio po swojej stronie.
Przeplyw logowania idzie przez aplikacje serwerowa Go:

React -> serwer -> dostawca OAuth2 -> serwer callback -> React

Serwer zapisuje dane uzytkownika, konto OAuth oraz token dostawcy w SQLite.
Do klienta wysylany jest osobny token sesji wygenerowany przez serwer.

Konfiguracja OAuth2

Zmienne srodowiskowe dla serwera:

GOOGLE_CLIENT_ID=...
GOOGLE_CLIENT_SECRET=...
GITHUB_CLIENT_ID=...
GITHUB_CLIENT_SECRET=...
SERVER_BASE_URL=http://localhost:8080
CLIENT_BASE_URL=http://localhost:5173
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

Callback URL do ustawienia u dostawcow:

http://localhost:8080/auth/google/callback
http://localhost:8080/auth/github/callback

Realizacja zadania

3.0 Logowanie przez aplikacje serwerowa: POST /auth/login
3.5 Rejestracja przez aplikacje serwerowa: POST /auth/register
4.0 Logowanie przez Google OAuth2: GET /auth/google/login
4.5 Logowanie przez GitHub OAuth2: GET /auth/github/login
5.0 Dane logowania OAuth2 i token dostawcy zapisywane sa po stronie serwera

Najwazniejsze endpointy

POST /auth/register - rejestracja lokalna
POST /auth/login - logowanie lokalne
POST /auth/logout - wylogowanie i usuniecie sesji
GET /auth/me - dane zalogowanego uzytkownika
GET /auth/google/login - start logowania Google
GET /auth/google/callback - callback Google
GET /auth/github/login - start logowania GitHub
GET /auth/github/callback - callback GitHub

Uruchomienie lokalne

# Serwer
cd server
go mod tidy
go run ./seed
go run .

# Klient (nowy terminal)
cd client
npm install
npm run dev

Klient: http://localhost:5173
Serwer: http://localhost:8080

Uruchomienie przez Docker

docker-compose up --build

Klient: http://localhost
Serwer: http://localhost:8080
