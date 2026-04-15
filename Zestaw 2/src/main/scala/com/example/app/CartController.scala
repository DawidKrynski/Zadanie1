package com.example.app

import org.scalatra._

case class CartItem(id: Int, productId: Int, quantity: Int)

class CartController extends ApiController {
  var data = List(CartItem(1, 1, 2))

  get("/") { data }
  get("/:id") { handleGet(data.find(_.id == params("id").toInt)) }
  post("/") { val c = parsedBody.extract[CartItem]; data = data :+ c; Created(c) }
  put("/:id") { val c = parsedBody.extract[CartItem]; data = data.map(x => if (x.id == params("id").toInt) c else x); Ok(c) }
  delete("/:id") {
    val id = params("id").toInt
    val exists = data.exists(_.id == id)
    data = data.filterNot(_.id == id)
    handleDelete(exists)
  }
}
