from datetime import UTC, datetime
from typing import Any, cast
from unittest.mock import AsyncMock, MagicMock, patch

import pytest
from sqlalchemy.ext.asyncio import AsyncSession
from venvi.models.odh_event import ODHEvent
from venvi.services.odh import fetch_odh_events, sync_odh_events

# Mock Response Data
MOCK_ODH_RESPONSE = {
    "TotalResults": 1,
    "TotalPages": 1,
    "CurrentPage": 1,
    "Seed": None,
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

        events = await fetch_odh_events()
        assert len(events) == 1
        assert events[0]["Id"] == "test-event-1"


@pytest.mark.asyncio
async def test_sync_odh_events() -> None:
    # Mock the session completely to avoid DB interaction/Greenlet issues
    mock_session = AsyncMock(spec=AsyncSession)

    # Mock the execute result for 'existing' check
    # We want to simulate that the event does NOT exist first
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

        count = await sync_odh_events(mock_session)

        assert count == 1

        # Verify session.add was called
        assert mock_session.add.called
        # Verify commit was called
        assert mock_session.commit.called

        added_obj = mock_session.add.call_args[0][0]
        assert added_obj.title == "Test Event"
        assert added_obj.location == "Bozen"
        assert added_obj.image_url == "http://example.com/image.jpg"


def test_map_odh_event_edge_cases() -> None:
    from venvi.services.odh import map_odh_event

    # Case 1: Minimal data to trigger fallbacks
    raw: dict[str, Any] = {"Detail": {}}  # No Id, No DateBegin, No Title in detail
    event = map_odh_event(raw)

    # Should generate ID from title (which defaults to "Untitled Event")
    # + DateBegin (defaults to None->None/Now)
    # Actually code says: raw_id = str(hash(title + str(raw.get("DateBegin"))))
    assert event.id is not None
    assert event.title == "Untitled Event"
    # Date fallback is now
    assert event.date_start is not None

    # Case 2: ID from Mapping
    raw_mapping = {
        "Detail": {"en": {"Title": "Mapped ID"}},
        "Mapping": {"some_source": {"rid": "mapped-id-123"}},
    }
    event_mapped = map_odh_event(raw_mapping)
    assert event_mapped.id == "mapped-id-123"

    # Case 3: Missing dates
    raw_dates = {
        "Id": "date-test",
        "Detail": {"en": {"Title": "Date Test"}},
        # No DateBegin/End
    }
    event_dates = map_odh_event(raw_dates)
    assert event_dates.date_start is not None  # Should be now()


@pytest.mark.asyncio
async def test_sync_odh_events_update_existing() -> None:
    mock_session = AsyncMock(spec=AsyncSession)

    # Simulate existing event found
    existing_event = ODHEvent(
        id="test-event-1",
        title="Old Title",
        is_new=False,
        date_start=datetime.now(UTC),
        date_end=datetime.now(UTC),
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

        await sync_odh_events(mock_session)

        # Verify title was updated on existing object
        assert existing_event.title == "New Title"
        # Verify is_new was NOT updated (should stay False)
        assert existing_event.is_new is False

        # Verify session.add was called with existing object
        mock_session.add.assert_called_with(existing_event)
