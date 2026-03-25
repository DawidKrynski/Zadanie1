# Zadanie 1 Docker

Autor: Dawid Kryński (dawidk12)

✅ 3.0 Obraz Ubuntu z Pythonem w wersji 3.10
[Link do commita 1](#)

✅ 3.5 Obraz ubuntu:24.04 z Javą 8 oraz Kotlinem
[Link do commita 2](#)

✅ 4.0 Dodanie najnowszego Gradle'a oraz paczki JDBC SQLite w projekcie Gradle (build.gradle)
[Link do commita 3](#)

✅ 4.5 Przykład HelloWorld + uruchomienie przez CMD (`java -jar`) oraz Gradle (`gradle run`)
[Link do commita 4](#)

✅ 5.0 Konfiguracja docker-compose
[Link do commita 5](#)

## Obrazy na Docker Hub

- [dawidk12/python-ubuntu](https://hub.docker.com/r/dawidk12/python-ubuntu)
- [dawidk12/kotlin-hello-world](https://hub.docker.com/r/dawidk12/kotlin-hello-world)

## Uruchomienie

```bash
# budowanie i uruchomienie obu kontenerow
docker-compose up --build

# albo osobno:
docker build -t dawidk12/python-ubuntu ./python
docker build -t dawidk12/kotlin-hello-world ./kotlin-app

docker run dawidk12/python-ubuntu
docker run dawidk12/kotlin-hello-world

# uruchomienie HelloWorld przez gradle (wewnatrz kontenera)
docker run dawidk12/kotlin-hello-world gradle run --no-daemon
```

Kod: [Link do repozytorium](#)
