name := "Zadanie2-Scalatra"
version := "0.1.0"
scalaVersion := "2.13.12"

val ScalatraVersion = "2.8.2"

libraryDependencies ++= Seq(
  "org.scalatra" %% "scalatra"      % ScalatraVersion,
  "org.scalatra" %% "scalatra-json" % ScalatraVersion,
  "org.json4s"   %% "json4s-jackson" % "4.0.6",
  "javax.servlet" % "javax.servlet-api" % "3.1.0" % "provided",
  "org.eclipse.jetty" % "jetty-webapp" % "9.4.43.v20210629" % "container"
)

enablePlugins(JettyPlugin)
Jetty / containerLibs := Seq("org.eclipse.jetty" % "jetty-runner" % "9.4.43.v20210629")
Jetty / containerPort := 8080
