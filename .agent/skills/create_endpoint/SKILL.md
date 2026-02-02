---
name: create_endpoint
description: Guide for adding a new API endpoint to the FastAPI application.
---
# Create Endpoint Skill

Follow these steps to add a new endpoint to the application.

## 1. Plan the Endpoint
Define the following:
- **Path**: e.g., `/users/{id}`
- **Method**: GET, POST, PUT, DELETE
- **Request Body**: What data is coming in? (Pydantic Model)
- **Response Model**: What data is returning? (Pydantic Model)

## 2. Create Data Models (`src/app/models/`)
If this endpoint involves a new entity, create a model in `src/app/models/`.
- Use `SQLModel` for database tables.
- Use `Pydantic` models for schemas (DTOs).

## 3. Create the Route (`src/app/api/`)
Create a new file or edit an existing one in `src/app/api/`.

```python
from fastapi import APIRouter
from app.models.item import Item

router = APIRouter()

@router.get("/items/", response_model=list[Item])
async def read_items():
    return []
```

## 4. Register the Router
Ensure the router is included in `src/app/main.py` or the parent router.

## 5. Write Tests (`src/tests/`)
Create a test file `src/tests/test_api_<name>.py`.
- Use `TestClient`.
- Verify status codes and response bodies.

## 6. Verify
Run `./.agent/workflows/verify_changes.md` to ensure everything is correct.
