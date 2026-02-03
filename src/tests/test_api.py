from datetime import UTC, datetime

import pytest
from httpx import AsyncClient
from sqlalchemy.ext.asyncio import AsyncSession
from venvi.models.event import Event


@pytest.mark.asyncio
async def test_read_events_empty(client: AsyncClient) -> None:
    response = await client.get("/events/")
    assert response.status_code == 200
    assert response.json() == []


@pytest.mark.asyncio
async def test_sync_endpoint_mocked(client: AsyncClient, session: AsyncSession) -> None:
    from unittest.mock import patch

    with patch("venvi.services.ingestion.sync_all_events", return_value={"test": 0}):
        response = await client.post("/events/sync")
        assert response.status_code == 200
        assert response.json()["message"] == "Sync complete"


@pytest.mark.asyncio
async def test_create_and_read_event(
    client: AsyncClient, session: AsyncSession
) -> None:
    event = Event(
        id="test-event-id",
        title="API Test Event",
        date_start=datetime(2026, 2, 3, 0, 0, 0, tzinfo=UTC),
        date_end=datetime(2026, 2, 4, 0, 0, 0, tzinfo=UTC),
        location="API City",
        url="https://api.test",
        source_name="test_source",
        source_id="original_id",
        category="hackathon",
    )
    session.add(event)
    await session.commit()

    response = await client.get("/events/")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["title"] == "API Test Event"


@pytest.mark.asyncio
async def test_read_events_filtering(
    client: AsyncClient, session: AsyncSession
) -> None:
    """Test filtering events by category and source."""
    e1 = Event(
        id="event-1",
        title="Hackathon Event",
        date_start=datetime(2026, 2, 3, 0, 0, 0, tzinfo=UTC),
        date_end=datetime(2026, 2, 4, 0, 0, 0, tzinfo=UTC),
        location="City A",
        url="https://a.test",
        category="hackathon",
        source_name="source_a",
        source_id="1",
    )
    e2 = Event(
        id="event-2",
        title="Other Event",
        date_start=datetime(2026, 2, 3, 0, 0, 0, tzinfo=UTC),
        date_end=datetime(2026, 2, 4, 0, 0, 0, tzinfo=UTC),
        location="City B",
        url="https://b.test",
        category="general",
        source_name="source_b",
        source_id="2",
    )
    session.add(e1)
    session.add(e2)
    await session.commit()

    # Test filtering for hackathon
    response = await client.get("/events/", params={"category": "hackathon"})
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["title"] == "Hackathon Event"

    # Test filtering for source_b
    response = await client.get("/events/", params={"source": "source_b"})
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["title"] == "Other Event"


@pytest.mark.asyncio
async def test_sync_endpoint_failure(
    client: AsyncClient, session: AsyncSession
) -> None:
    """Test the sync endpoint when it encounters an error."""
    from unittest.mock import patch

    with patch(
        "venvi.api.routers.events.sync_all_events",
        side_effect=Exception("Database Connection Error"),
    ):
        response = await client.post("/events/sync")
        assert response.status_code == 500
        assert "Database Connection Error" in response.json()["detail"]
