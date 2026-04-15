# Obraz Ubuntu 22.04 z Pythonem 3.10
# Autor: Dawid Kryński (dawidk12)

FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

# instalacja pythona 3.10 - jest domyslnie w repo Ubuntu 22.04
RUN apt-get update && \
    apt-get install -y python3.10 python3-pip && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# sprawdzenie wersji
RUN python3.10 --version

CMD ["python3.10", "-c", "print('Hello from Python 3.10 on Ubuntu 22.04!')"]
