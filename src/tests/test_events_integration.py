from datetime import UTC, datetime
from unittest.mock import patch

import pytest
from httpx import AsyncClient
from sqlalchemy.ext.asyncio import AsyncSession
from venvi.models.odh_event import ODHEvent

# Mock Data for Sync
MOCK_ODH_API_DATA = {
    "TotalResults": 1,
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
async def test_read_odh_events_empty(
    client: AsyncClient, session: AsyncSession
) -> None:
    """Test reading ODH events when DB is empty."""
    # Ensure DB is empty (fixture creates fresh DB but session might persist
    # if not careful, but based on conftest it should be isolated per test
    # if using in-memory)
    # Actually conftest uses session fixture.

    response = await client.get("/events/odh")
    assert response.status_code == 200
    assert response.json() == []


@pytest.mark.asyncio
async def test_sync_odh_events_integration(
    client: AsyncClient, session: AsyncSession
) -> None:
    """Test full sync flow with mocked external API but real DB persistence."""

    # Mock the external API call in services.odh.fetch_odh_events
    # We mock 'httpx.AsyncClient.get' essentially, or simpler: mock fetch_odh_events

    with patch("venvi.services.odh.fetch_odh_events") as mock_fetch:
        # fetch_odh_events returns list[dict]
        mock_fetch.return_value = MOCK_ODH_API_DATA["Items"]

        # Trigger Sync
        response = await client.post("/events/odh/sync")
        assert response.status_code == 200
        data = response.json()
        assert data["message"] == "Sync complete"
        assert data["new_items"] == 1

        # Verify persistence via API
        list_response = await client.get("/events/odh")
        assert list_response.status_code == 200
        events = list_response.json()
        assert len(events) == 1
        assert events[0]["title"] == "Integration Event"
        assert events[0]["location"] == "Integration City"
        assert events[0]["id"] == "integration-test-event-1"


@pytest.mark.asyncio
async def test_read_odh_events_existing(
    client: AsyncClient, session: AsyncSession
) -> None:
    """Test reading pre-inserted ODH events."""
    event = ODHEvent(
        id="pre-existing-1",
        title="Existing Event",
        description="Description",
        date_start=datetime(2026, 1, 1, 10, 0, 0, tzinfo=UTC),
        date_end=datetime(2026, 1, 1, 12, 0, 0, tzinfo=UTC),
        location="Existing City",
        is_new=True,
    )
    session.add(event)
    await session.commit()

    response = await client.get("/events/odh")
    assert response.status_code == 200
    events = response.json()
    assert len(events) == 1
    assert events[0]["title"] == "Existing Event"


@pytest.mark.asyncio
async def test_read_odh_events_filter_taken(
    client: AsyncClient, session: AsyncSession
) -> None:
    """Test filtering ODH events by taken status."""
    e1 = ODHEvent(
        id="e1",
        title="Taken",
        taken=True,
        is_new=True,
        date_start=datetime.now(UTC),
        date_end=datetime.now(UTC),
    )
    e2 = ODHEvent(
        id="e2",
        title="Free",
        taken=False,
        is_new=True,
        date_start=datetime.now(UTC),
        date_end=datetime.now(UTC),
    )
    session.add_all([e1, e2])
    await session.commit()

    # Filter taken=True
    response = await client.get("/events/odh", params={"taken": "true"})
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["id"] == "e1"

    # Filter taken=False
    response = await client.get("/events/odh", params={"taken": "false"})
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["id"] == "e2"


@pytest.mark.asyncio
async def test_toggle_odh_event_taken(
    client: AsyncClient, session: AsyncSession
) -> None:
    """Test toggling the taken status of an ODH event."""
    event = ODHEvent(
        id="toggle-test",
        title="Toggle",
        taken=False,
        is_new=True,
        date_start=datetime.now(UTC),
        date_end=datetime.now(UTC),
    )
    session.add(event)
    await session.commit()

    # Toggle to True
    response = await client.patch(
        f"/events/odh/{event.id}/taken", params={"taken": "true"}
    )
    assert response.status_code == 200
    assert response.json()["taken"] is True

    # Verify persistence
    await session.refresh(event)
    assert event.taken is True

    # Toggle to False
    response = await client.patch(
        f"/events/odh/{event.id}/taken", params={"taken": "false"}
    )
    assert response.status_code == 200
    assert response.json()["taken"] is False


@pytest.mark.asyncio
async def test_toggle_odh_event_not_found(
    client: AsyncClient, session: AsyncSession
) -> None:
    """Test toggling a non-existent event."""
    response = await client.patch(
        "/events/odh/non-existent/taken", params={"taken": "true"}
    )
    assert response.status_code == 404
    assert response.json()["detail"] == "Event not found"


@pytest.mark.asyncio
async def test_sync_odh_events_error(
    client: AsyncClient, session: AsyncSession
) -> None:
    """Test error handling in sync endpoint."""
    with patch(
        "venvi.api.routers.events.sync_odh_events", side_effect=Exception("Sync Failed")
    ):
        response = await client.post("/events/odh/sync")
        assert response.status_code == 500
        assert "Sync Failed" in response.json()["detail"]
