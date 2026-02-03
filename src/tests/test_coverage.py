from unittest.mock import AsyncMock, MagicMock, patch

import pytest
from venvi.api.routers.events import read_events, sync_events
from venvi.main import lifespan
from venvi.web.router import get_events_partial


@pytest.mark.asyncio
async def test_read_events_direct() -> None:
    session = AsyncMock()
    mock_result = MagicMock()
    mock_result.scalars.return_value.all.return_value = []
    session.execute.return_value = mock_result

    res = await read_events(category=None, source=None, session=session)
    assert res == []
    session.execute.assert_called()


@pytest.mark.asyncio
async def test_sync_events_direct() -> None:
    session = AsyncMock()
    with patch("venvi.api.routers.events.sync_all_events", return_value={"test": 5}):
        res = await sync_events(session=session)
        assert res["message"] == "Sync complete"
        assert res["total"] == 5


@pytest.mark.asyncio
async def test_get_events_partial_direct() -> None:
    session = AsyncMock()
    mock_result = MagicMock()
    mock_result.scalars.return_value.all.return_value = []
    session.execute.return_value = mock_result
    request = MagicMock()

    with patch("venvi.web.router.templates.TemplateResponse") as mock_template:
        await get_events_partial(request=request, session=session)
        mock_template.assert_called()


@pytest.mark.asyncio
async def test_lifespan_direct() -> None:
    app = MagicMock()
    with patch("venvi.main.init_db", new_callable=AsyncMock) as mock_init:
        async with lifespan(app):
            pass
        mock_init.assert_called_once()
