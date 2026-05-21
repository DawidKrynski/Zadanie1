Zadanie 9 - osobny serwis ChatGPT

Pythonowy serwis HTTP realizujacy punkt 3.0: osobny komponent po stronie serwerowej do laczenia z ChatGPT przez OpenAI API.

Realizacja zadania

3.0 Stworzono osobny serwis serwerowy w Pythonie:

- FastAPI jako warstwa HTTP
- POST /chat przyjmuje tekst uzytkownika i wysyla go do OpenAI API
- GET /health sluzy do sprawdzenia dzialania serwisu
- konfiguracja przez zmienne srodowiskowe

Zmienne srodowiskowe

OPENAI_API_KEY=...
OPENAI_MODEL=gpt-4o-mini
OPENAI_API_URL=https://api.openai.com/v1/responses

Uruchomienie lokalne

```bash
cd server
pip install -r requirements.txt
$env:OPENAI_API_KEY="twoj_klucz"
uvicorn main:app --reload --host 0.0.0.0 --port 8000
```

Test zapytania

```bash
curl -X POST http://localhost:8000/chat \
  -H "Content-Type: application/json" \
  -d "{\"message\":\"Napisz krotka odpowiedz testowa\"}"
```

Uruchomienie w Dockerze

```bash
cd server
docker build -t zadanie9-chatgpt-service .
docker run --rm -p 8000:8000 -e OPENAI_API_KEY=twoj_klucz zadanie9-chatgpt-service
```

Testy

```bash
cd server
pip install -r requirements-dev.txt
pytest
```
