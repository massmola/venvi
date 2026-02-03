from datetime import UTC, datetime

import pytest
from httpx import AsyncClient
from sqlalchemy.ext.asyncio import AsyncSession
from venvi.models.event import Event


@pytest.mark.asyncio
async def test_index_page(client: AsyncClient) -> None:
    """Test that the index page loads and contains HTMX triggers."""
    response = await client.get("/")
    assert response.status_code == 200
    assert "Venvi" in response.text
    # Check for HTMX trigger
    assert 'hx-get="/partials/events"' in response.text
    assert 'hx-trigger="load, reload from:body"' in response.text


@pytest.mark.asyncio
async def test_events_partial_empty(client: AsyncClient) -> None:
    """Test the events partial when no data is present."""
    response = await client.get("/partials/events")
    assert response.status_code == 200
    assert "event-list" not in response.text  # Adjust if needed based on partial logic


@pytest.mark.asyncio
async def test_events_partial_with_data(
    client: AsyncClient, session: AsyncSession
) -> None:
    """Test the events partial with actual data."""
    event = Event(
        id="test-id-123",
        title="Web Test Event",
        date_start=datetime(2026, 2, 3, 0, 0, 0, tzinfo=UTC),
        date_end=datetime(2026, 2, 4, 0, 0, 0, tzinfo=UTC),
        location="Web City",
        url="https://web.test",
        source_name="web_test",
        source_id="123",
        topics=["web", "test"],
        category="general",
    )
    session.add(event)
    await session.commit()

    response = await client.get("/partials/events")
    assert response.status_code == 200
    assert "Web Test Event" in response.text
    assert "Web City" in response.text
    assert "#web" in response.text
