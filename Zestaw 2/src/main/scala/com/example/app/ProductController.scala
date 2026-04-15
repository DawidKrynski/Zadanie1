package com.example.app

import org.scalatra._

case class Product(id: Int, name: String, price: Double)

class ProductController extends ApiController {
  var data = List(Product(1, "Laptop", 3000.0), Product(2, "Myszka", 150.0))

  get("/") { data }
  get("/:id") { handleGet(data.find(_.id == params("id").toInt)) }
  post("/") { val p = parsedBody.extract[Product]; data = data :+ p; Created(p) }
  put("/:id") { val p = parsedBody.extract[Product]; data = data.map(x => if (x.id == params("id").toInt) p else x); Ok(p) }
  delete("/:id") {
    val id = params("id").toInt
    val exists = data.exists(_.id == id)
    data = data.filterNot(_.id == id)
    handleDelete(exists)
  }
}
