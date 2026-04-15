package com.example.app

import org.scalatra._

case class Category(id: Int, name: String)

class CategoryController extends ApiController {
  var data = List(Category(1, "Elektronika"), Category(2, "Akcesoria"))

  get("/") { data }
  get("/:id") { handleGet(data.find(_.id == params("id").toInt)) }
  post("/") { val c = parsedBody.extract[Category]; data = data :+ c; Created(c) }
  put("/:id") { val c = parsedBody.extract[Category]; data = data.map(x => if (x.id == params("id").toInt) c else x); Ok(c) }
  delete("/:id") {
    val id = params("id").toInt
    val exists = data.exists(_.id == id)
    data = data.filterNot(_.id == id)
    handleDelete(exists)
  }
}
