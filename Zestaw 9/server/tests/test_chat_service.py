from fastapi.testclient import TestClient

from main import app, extract_response_text, get_openai_service


class FakeOpenAIService:
    model = "test-model"

    def ask(self, message: str) -> str:
        return f"answer: {message}"


def test_health_endpoint() -> None:
    client = TestClient(app)

    response = client.get("/health")

    assert response.status_code == 200
    assert response.json() == {"status": "ok"}


def test_chat_endpoint_uses_openai_service() -> None:
    app.dependency_overrides[get_openai_service] = lambda: FakeOpenAIService()
    client = TestClient(app)

    response = client.post("/chat", json={"message": "hello"})

    app.dependency_overrides.clear()
    assert response.status_code == 200
    assert response.json() == {"response": "answer: hello", "model": "test-model"}


def test_chat_rejects_empty_message() -> None:
    client = TestClient(app)

    response = client.post("/chat", json={"message": ""})

    assert response.status_code == 422


def test_extract_response_text_from_output_text() -> None:
    data = {"output_text": " hello "}

    assert extract_response_text(data) == "hello"


def test_extract_response_text_from_output_content() -> None:
    data = {
        "output": [
            {"content": [{"text": "first"}, {"text": "second"}]},
            {"content": [{"text": "third"}]},
        ]
    }

    assert extract_response_text(data) == "first\nsecond\nthird"
