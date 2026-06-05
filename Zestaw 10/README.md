# Zadanie 10 - Docker w chmurze

Minimalna realizacja punktu 3.0: przygotowano instancje aplikacji serwerowej
uruchamiane w Dockerze na maszynie w chmurze.

## Zakres

- 3.0 instancja VM w chmurze z Dockerem
- 3.0 kontener `server` z API z `Zestaw 8/server`
- 3.0 publiczny port API: `8080`

Nie dodano GitHub Actions, Sonara, klienta ani BrowserStacka, bo to progi
3.5-5.0.

## Pliki

- `cloud-init.yml` - automatyczna konfiguracja VM: Docker + clone repo + start kontenera
- `docker-compose.cloud.yml` - uruchomienie kontenera serwera na VM

## Uruchomienie na VM

Przy tworzeniu VM w chmurze wkleic `cloud-init.yml` jako user data.
W firewallu/security group otworzyc TCP `8080`.

Po starcie VM sprawdzenie:

```bash
docker ps
curl http://PUBLICZNY_ADRES_VM:8080/products
```

## Reczne uruchomienie na VM

```bash
git clone https://github.com/DawidKrynski/Zadanie1.git
cd Zadanie1
docker compose -f "Zestaw 10/docker-compose.cloud.yml" up -d --build
```

API:

```text
http://PUBLICZNY_ADRES_VM:8080/products
```
