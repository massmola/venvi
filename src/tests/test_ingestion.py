from typing import Any
from unittest.mock import patch

import pytest
from sqlalchemy.ext.asyncio import AsyncSession
from sqlmodel import select
from venvi.models.event import Event
from venvi.services.ingestion import sync_all_events

SAMPLE_DATA: list[dict[str, Any]] = [
    {
        "id": "209e1cfb-5040-4086-bfec-f67bdc3380ff",
        "name": "Test Event",
        "city": "Test City",
        "country_code": "TC",
        "date_start": "2026-02-03T00:00:00+00:00",
        "date_end": "2026-02-04T00:00:00+00:00",
        "topics": [],
        "notes": None,
        "url": "https://example.com",
    }
]


@pytest.mark.asyncio
async def test_sync_all_events(session: AsyncSession) -> None:
    # Build a mock response for httpx in EuroHackathonsProvider
    mock_response = {"data": SAMPLE_DATA, "success": True}

    with (
        patch("httpx.AsyncClient.get") as mock_get,
        patch("venvi.services.providers.odh.ODHProvider.fetch_events", return_value=[]),
    ):
        from unittest.mock import MagicMock

        response_mock = MagicMock()
        response_mock.status_code = 200
        response_mock.json.return_value = mock_response
        mock_get.return_value = response_mock

        results = await sync_all_events(session)
        assert results["euro_hackathons"] == 1

        result = await session.execute(select(Event))
        events = result.scalars().all()
        assert len(events) == 1
        assert events[0].title == "Test Event"


@pytest.mark.asyncio
async def test_sync_events_update(session: AsyncSession) -> None:
    # Initial sync
    mock_response_initial = {"data": SAMPLE_DATA, "success": True}

    from unittest.mock import MagicMock

    with (
        patch("httpx.AsyncClient.get") as mock_get,
        patch("venvi.services.providers.odh.ODHProvider.fetch_events", return_value=[]),
    ):
        response_mock_initial = MagicMock()
        response_mock_initial.status_code = 200
        response_mock_initial.json.return_value = mock_response_initial
        mock_get.return_value = response_mock_initial

        await sync_all_events(session)

    # Update data
    updated_data = list(SAMPLE_DATA)
    updated_data[0]["name"] = "Updated Name"
    mock_response_updated = {"data": updated_data, "success": True}

    with (
        patch("httpx.AsyncClient.get") as mock_get,
        patch("venvi.services.providers.odh.ODHProvider.fetch_events", return_value=[]),
    ):
        response_mock_updated = MagicMock()
        response_mock_updated.status_code = 200
        response_mock_updated.json.return_value = mock_response_updated
        mock_get.return_value = response_mock_updated

        results = await sync_all_events(session)
        assert results["euro_hackathons"] == 0  # No new items

        result = await session.execute(select(Event))
        event = result.scalars().first()
        assert event is not None
        assert event.title == "Updated Name"
