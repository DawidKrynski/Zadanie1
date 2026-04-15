package com.example.app

import org.scalatra._
import org.scalatra.json._
import org.json4s.{DefaultFormats, Formats}
import org.scalatra.CorsSupport

trait ApiController extends ScalatraServlet with JacksonJsonSupport with CorsSupport {
  protected implicit lazy val jsonFormats: Formats = DefaultFormats

  options("/*") {
    response.setHeader("Access-Control-Allow-Headers", request.getHeader("Access-Control-Request-Headers"))
  }

  before() { contentType = formats("json") }

  def handleGet[T](item: Option[T]): ActionResult =
    item.map(Ok(_)).getOrElse(NotFound(Map("error" -> "Not found")))

  def handleDelete(condition: Boolean): ActionResult =
    if (condition) NoContent() else NotFound(Map("error" -> "Not found"))
}
