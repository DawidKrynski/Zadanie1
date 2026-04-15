import com.example.app._
import org.scalatra._
import javax.servlet.ServletContext

class ScalatraBootstrap extends LifeCycle {
  override def init(context: ServletContext): Unit = {
    context.setInitParameter("org.scalatra.cors.allowedOrigins", "http://localhost:3000 http://127.0.0.1:5500")
    context.setInitParameter("org.scalatra.cors.allowedMethods", "GET, POST, PUT, DELETE, OPTIONS")

    context.mount(new ProductController, "/products/*")
    context.mount(new CategoryController, "/categories/*")
    context.mount(new CartController, "/cart/*")
  }
}
