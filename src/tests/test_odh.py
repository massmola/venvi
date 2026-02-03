from datetime import UTC, datetime
from typing import Any, cast
from unittest.mock import AsyncMock, MagicMock, patch

import pytest
from sqlalchemy.ext.asyncio import AsyncSession
from venvi.models.event import Event
from venvi.services.providers.odh import ODHProvider

# Mock Response Data
MOCK_ODH_RESPONSE = {
    "TotalResults": 1,
    "Items": [
        {
            "Id": "test-event-1",
            "Detail": {
                "en": {"Title": "Test Event", "BaseText": "Description of Test Event"}
            },
            "DateBegin": "2024-01-01T10:00:00",
            "DateEnd": "2024-01-01T12:00:00",
            "ContactInfos": {"en": {"City": "Bozen"}},
            "ImageGallery": [{"ImageUrl": "http://example.com/image.jpg"}],
        }
    ],
}


@pytest.mark.asyncio
async def test_fetch_odh_events_mock() -> None:
    provider = ODHProvider()
    # Properly mock httpx.AsyncClient
    with patch("httpx.AsyncClient") as MockClient:
        mock_instance = MockClient.return_value
        # Mock __aenter__ to return the instance itself
        mock_instance.__aenter__.return_value = mock_instance

        # Setup the response object
        mock_response = MagicMock()
        mock_response.status_code = 200
        mock_response.json.return_value = MOCK_ODH_RESPONSE

        # Make .get() returns a coroutine that resolves to mock_response
        mock_instance.get = AsyncMock(return_value=mock_response)

        events = await provider.fetch_events()
        assert len(events) == 1
        assert events[0]["Id"] == "test-event-1"


@pytest.mark.asyncio
async def test_sync_odh_events() -> None:
    provider = ODHProvider()
    mock_session = AsyncMock(spec=AsyncSession)

    # Mock the execute result for 'existing' check
    mock_result = MagicMock()
    mock_result.scalar_one_or_none.return_value = None
    mock_session.execute.return_value = mock_result

    with patch("httpx.AsyncClient") as MockClient:
        mock_instance = MockClient.return_value
        mock_instance.__aenter__.return_value = mock_instance

        mock_response = MagicMock()
        mock_response.status_code = 200
        mock_response.json.return_value = MOCK_ODH_RESPONSE

        mock_instance.get = AsyncMock(return_value=mock_response)

        count = await provider.sync_events(mock_session)

        assert count == 1
        assert mock_session.add.called
        assert mock_session.commit.called

        added_obj = mock_session.add.call_args[0][0]
        assert isinstance(added_obj, Event)
        assert added_obj.title == "Test Event"
        assert added_obj.location == "Bozen"
        assert added_obj.source_name == "odh"


def test_map_odh_event_edge_cases() -> None:
    provider = ODHProvider()

    # Case 1: Minimal data
    raw: dict[str, Any] = {"Detail": {}}
    event = provider.map_event(raw)

    assert event.id is not None
    assert event.title == "Untitled Event"
    assert event.date_start is not None


@pytest.mark.asyncio
async def test_sync_odh_events_update_existing() -> None:
    provider = ODHProvider()
    mock_session = AsyncMock(spec=AsyncSession)

    # Simulate existing event found
    existing_event = Event(
        id="odh:test-event-1",
        title="Old Title",
        is_new=False,
        date_start=datetime.now(UTC),
        date_end=datetime.now(UTC),
        source_name="odh",
        source_id="test-event-1",
        url="http://test",
    )
    mock_result = MagicMock()
    mock_result.scalar_one_or_none.return_value = existing_event
    mock_session.execute.return_value = mock_result

    with patch("httpx.AsyncClient") as MockClient:
        mock_instance = MockClient.return_value
        mock_instance.__aenter__.return_value = mock_instance

        # Return same ID but new title
        new_data = cast(dict[str, Any], MOCK_ODH_RESPONSE.copy())
        new_data["Items"][0]["Detail"]["en"]["Title"] = "New Title"

        mock_response = MagicMock()
        mock_response.status_code = 200
        mock_response.json.return_value = new_data
        mock_instance.get = AsyncMock(return_value=mock_response)

        await provider.sync_events(mock_session)

        assert existing_event.title == "New Title"
        assert existing_event.is_new is False
        mock_session.add.assert_called_with(existing_event)
