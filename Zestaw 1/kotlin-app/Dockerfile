# Obraz Ubuntu 24.04 z Java 8, Kotlinem, Gradle + SQLite JDBC
# Autor: Dawid Kryński (dawidk12)

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive
ENV GRADLE_VERSION=8.13
ENV KOTLIN_VERSION=2.1.0

# potrzebne narzedzia do pobierania paczek
RUN apt-get update && apt-get install -y \
    wget unzip gnupg ca-certificates apt-transport-https \
    && rm -rf /var/lib/apt/lists/*

# Java 8 - trzeba dodac repo Adoptium bo Ubuntu 24.04 nie ma Javy 8
RUN wget -qO - https://packages.adoptium.net/artifactory/api/gpg/key/public \
      | gpg --dearmor -o /usr/share/keyrings/adoptium.gpg && \
    echo "deb [signed-by=/usr/share/keyrings/adoptium.gpg] \
      https://packages.adoptium.net/artifactory/deb noble main" \
      > /etc/apt/sources.list.d/adoptium.list && \
    apt-get update && apt-get install -y temurin-8-jdk && \
    rm -rf /var/lib/apt/lists/*

ENV JAVA_HOME=/usr/lib/jvm/temurin-8-amd64
ENV PATH=${JAVA_HOME}/bin:${PATH}

# Kotlin - sciagamy kompilator
RUN wget -q https://github.com/JetBrains/kotlin/releases/download/v${KOTLIN_VERSION}/kotlin-compiler-${KOTLIN_VERSION}.zip \
    && unzip -q kotlin-compiler-${KOTLIN_VERSION}.zip -d /opt \
    && rm kotlin-compiler-${KOTLIN_VERSION}.zip

ENV PATH=/opt/kotlinc/bin:${PATH}

# Gradle - najnowszy
RUN wget -q https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip \
    && unzip -q gradle-${GRADLE_VERSION}-bin.zip -d /opt \
    && rm gradle-${GRADLE_VERSION}-bin.zip

ENV GRADLE_HOME=/opt/gradle-${GRADLE_VERSION}
ENV PATH=${GRADLE_HOME}/bin:${PATH}

# sprawdzenie czy wszystko sie zainstalowalo
RUN java -version && kotlinc -version && gradle --version

WORKDIR /app

# kopiujemy pliki projektu
COPY build.gradle settings.gradle ./
COPY src ./src

# budujemy JAR-a (fat jar ze wszystkimi zaleznoscia)
RUN gradle jar --no-daemon

# domyslnie odpala przez java -jar
# mozna tez: gradle run --no-daemon
CMD ["java", "-jar", "build/libs/app.jar"]
