from datetime import UTC, datetime
from unittest.mock import patch

import pytest
from httpx import AsyncClient
from sqlalchemy.ext.asyncio import AsyncSession
from venvi.models.event import Event

# Mock Data for Sync
MOCK_ODH_API_DATA = {
    "Items": [
        {
            "Id": "integration-test-event-1",
            "Detail": {
                "en": {
                    "Title": "Integration Event",
                    "BaseText": "Integration Description",
                }
            },
            "DateBegin": "2026-05-01T10:00:00",
            "DateEnd": "2026-05-01T12:00:00",
            "ContactInfos": {"en": {"City": "Integration City"}},
        }
    ],
}


@pytest.mark.asyncio
async def test_read_events_empty(client: AsyncClient, session: AsyncSession) -> None:
    """Test reading events when DB is empty."""
    response = await client.get("/events/")
    assert response.status_code == 200
    assert response.json() == []


@pytest.mark.asyncio
async def test_sync_events_integration(
    client: AsyncClient, session: AsyncSession
) -> None:
    """Test full sync flow with mocked providers."""

    with (
        patch("venvi.services.providers.odh.ODHProvider.fetch_events") as mock_fetch,
        patch(
            "venvi.services.providers.euro_hackathons.EuroHackathonsProvider.fetch_events"
        ) as mock_fetch_hack,
    ):
        mock_fetch.return_value = MOCK_ODH_API_DATA["Items"]
        mock_fetch_hack.return_value = []

        # Trigger Sync
        response = await client.post("/events/sync")
        assert response.status_code == 200
        data = response.json()
        assert data["message"] == "Sync complete"
        assert data["total"] == 1

        # Verify persistence via API
        list_response = await client.get("/events/")
        assert list_response.status_code == 200
        events = list_response.json()
        assert len(events) == 1
        assert events[0]["title"] == "Integration Event"
        assert events[0]["location"] == "Integration City"
        assert events[0]["source_name"] == "odh"


@pytest.mark.asyncio
async def test_read_events_existing(client: AsyncClient, session: AsyncSession) -> None:
    """Test reading pre-inserted events."""
    event = Event(
        id="pre-existing-1",
        title="Existing Event",
        description="Description",
        date_start=datetime(2026, 1, 1, 10, 0, 0, tzinfo=UTC),
        date_end=datetime(2026, 1, 1, 12, 0, 0, tzinfo=UTC),
        location="Existing City",
        url="https://test.com",
        source_name="test",
        source_id="1",
        is_new=True,
    )
    session.add(event)
    await session.commit()

    response = await client.get("/events/")
    assert response.status_code == 200
    events = response.json()
    assert len(events) == 1
    assert events[0]["title"] == "Existing Event"


@pytest.mark.asyncio
async def test_sync_events_error(client: AsyncClient, session: AsyncSession) -> None:
    """Test error handling in sync endpoint."""
    with patch(
        "venvi.api.routers.events.sync_all_events", side_effect=Exception("Sync Failed")
    ):
        response = await client.post("/events/sync")
        assert response.status_code == 500
        assert "Sync Failed" in response.json()["detail"]
