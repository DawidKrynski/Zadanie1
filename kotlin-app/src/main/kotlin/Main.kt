import java.sql.DriverManager

// HelloWorld z przykladem uzycia SQLite JDBC
// Autor: Dawid Kryński (dawidk12)

fun main() {
    println("Hello, World!")

    // laczenie z baza SQLite w pamieci (nie tworzy pliku)
    val connection = DriverManager.getConnection("jdbc:sqlite::memory:")

    connection.createStatement().use { stmt ->
        stmt.execute("CREATE TABLE greetings (id INTEGER PRIMARY KEY, message TEXT)")
        stmt.execute("INSERT INTO greetings (message) VALUES ('Hello from SQLite!')")
    }

    // odczyt danych z bazy
    connection.createStatement().use { stmt ->
        val rs = stmt.executeQuery("SELECT message FROM greetings")
        while (rs.next()) {
            println("SQLite: ${rs.getString("message")}")
        }
    }

    connection.close()
    println("Polaczenie z SQLite zamkniete.")
}
