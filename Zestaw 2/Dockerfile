FROM kprzystalski/ebiznes-scala:latest

USER root
WORKDIR /app
COPY . /app
RUN chown -R student /app
USER student

EXPOSE 8080

CMD ["sbt", "-Djetty.host=0.0.0.0", "~Jetty/start"]
