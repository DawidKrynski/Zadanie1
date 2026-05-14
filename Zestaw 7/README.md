# Zadanie 7 - SonarCloud

Aplikacja z `Zestaw 5` zostala przygotowana do analizy w SonarCloud jako dwa
osobne projekty: serwer Go/Echo oraz klient React. Kod serwera jest dodatkowo
sprawdzany lokalnie przez `golangci-lint` w hooku gita z katalogu tego zadania.

## Zakres

- 3.0 linter dla aplikacji serwerowej w hooku gita: `Zestaw 7/githooks/pre-commit`
- 3.5 poprawki bledow w kodzie serwera
- 4.0 poprawki code smells w kodzie serwera
- 4.5 poprawki podatnosci i security hotspots w kodzie serwera
- 5.0 poprawki bledow i code smells w kodzie klienta

## SonarCloud

### Server

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=DawidKrynski_Zadanie1_server&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=DawidKrynski_Zadanie1_server)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=DawidKrynski_Zadanie1_server&metric=bugs)](https://sonarcloud.io/summary/new_code?id=DawidKrynski_Zadanie1_server)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=DawidKrynski_Zadanie1_server&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=DawidKrynski_Zadanie1_server)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=DawidKrynski_Zadanie1_server&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=DawidKrynski_Zadanie1_server)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=DawidKrynski_Zadanie1_server&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=DawidKrynski_Zadanie1_server)
[![Security Review Rating](https://sonarcloud.io/api/project_badges/measure?project=DawidKrynski_Zadanie1_server&metric=security_review_rating)](https://sonarcloud.io/summary/new_code?id=DawidKrynski_Zadanie1_server)

### Client

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=DawidKrynski_Zadanie1_client&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=DawidKrynski_Zadanie1_client)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=DawidKrynski_Zadanie1_client&metric=bugs)](https://sonarcloud.io/summary/new_code?id=DawidKrynski_Zadanie1_client)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=DawidKrynski_Zadanie1_client&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=DawidKrynski_Zadanie1_client)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=DawidKrynski_Zadanie1_client&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=DawidKrynski_Zadanie1_client)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=DawidKrynski_Zadanie1_client&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=DawidKrynski_Zadanie1_client)
[![Security Review Rating](https://sonarcloud.io/api/project_badges/measure?project=DawidKrynski_Zadanie1_client&metric=security_review_rating)](https://sonarcloud.io/summary/new_code?id=DawidKrynski_Zadanie1_client)

## Konfiguracja

Projekty w SonarCloud:

```text
DawidKrynski_Zadanie1_server
DawidKrynski_Zadanie1_client
```

Do uruchomienia analizy w SonarCloud trzeba ustawic token:

```text
SONAR_TOKEN
```

Konfiguracja skanowania znajduje sie w plikach `sonar-project.properties`:

```text
Zestaw 5/server/sonar-project.properties
Zestaw 5/client/sonar-project.properties
```

## Uruchomienie lokalne

Hook gita:

```bash
git config core.hooksPath "Zestaw 7/githooks"
```

Serwer:

```bash
cd "Zestaw 5/server"
go test ./...
golangci-lint run --config .golangci.yml ./...
```

Klient:

```bash
cd "Zestaw 5/client"
npm ci
npm run build
npm audit --audit-level=moderate
```
