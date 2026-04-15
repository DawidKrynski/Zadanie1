# Zadanie 3 Kotlin – Discord + Slack Bot

✅ 3.0 Aplikacja kliencka w Kotlinie (Ktor) pozwalająca na przesyłanie wiadomości na platformę Discord

✅ 3.5 Aplikacja odbiera wiadomości użytkowników z platformy Discord skierowane do bota

✅ 4.0 Zwraca listę kategorii na żądanie użytkownika (`!kategorie`)

✅ 4.5 Zwraca listę produktów wg żądanej kategorii (`!produkty <kategoria>`)

✅ 5.0 Aplikacja obsługuje dodatkowo platformę Slack (Socket Mode)

✅ Docker – aplikacja uruchomiona na Dockerze

Obraz Docker Hub: https://hub.docker.com/r/dawidk12/discord-slack-bot

```
src/main/kotlin/com/example/bot/
├── Main.kt                     # Punkt wejścia – uruchamia Discord + Slack + Ktor
├── catalog/
│   └── ProductCatalog.kt       # Katalog produktów (in-memory)
├── discord/
│   └── DiscordBot.kt           # Obsługa Discord (Kord Gateway)
└── slack/
    └── SlackBot.kt             # Obsługa Slack (Bolt SDK, Socket Mode)
```

## Health check

```bash
curl http://localhost:8080/health
# OK
```
