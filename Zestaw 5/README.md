Zadanie 5 – Aplikacja kliencka React.js + Go REST API
Aplikacja kliencka napisana w React.js z wykorzystaniem biblioteki axios oraz react-router-dom. Serwer REST API z poprzedniego zadania (Go + Echo + GORM). Dane między komponentami przesyłane za pomocą React hooks.
Docker Hub
Obrazy dostępne na:

https://hub.docker.com/r/dawidk12/zadanie5-server
https://hub.docker.com/r/dawidk12/zadanie5-client

bashdocker-compose up
Realizacja zadania
✅ 3.0 Dwa komponenty: Produkty (GET /products) oraz Płatności (POST); dane pobierane i wysyłane do serwera Go
✅ 3.5 Komponent Koszyk z widokiem; routing za pomocą react-router-dom
✅ 4.0 Dane między wszystkimi komponentami przesyłane za pomocą React hooks (useContext, useReducer, useLocation, useNavigate, useEffect)
✅ 4.5 docker-compose uruchamiający aplikację serwerową (Go) oraz kliencką (React + nginx)
✅ 5.0 Axios z nagłówkami Content-Type i Accept; CORS middleware w Echo (AllowOrigins, AllowHeaders, AllowMethods)


Uruchomienie lokalne
bash# Serwer
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
bashdocker-compose up --build

Klient: http://localhost
Serwer: http://localhost:8080