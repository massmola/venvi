from unittest.mock import AsyncMock, MagicMock, patch

import pytest
from venvi.api.routers.hackathons import read_hackathons, sync_data
from venvi.main import lifespan
from venvi.web.router import get_hackathons_partial


@pytest.mark.asyncio
async def test_read_hackathons_direct() -> None:
    session = AsyncMock()
    mock_result = MagicMock()
    mock_result.scalars.return_value.all.return_value = []
    session.execute.return_value = mock_result

    res = await read_hackathons(status=None, session=session)
    assert res == []
    session.execute.assert_called()


@pytest.mark.asyncio
async def test_sync_data_direct() -> None:
    session = AsyncMock()
    with patch("venvi.api.routers.hackathons.sync_hackathons", return_value=5):
        res = await sync_data(session=session)
        assert res == {"message": "Sync complete", "new_items": 5}


@pytest.mark.asyncio
async def test_get_hackathons_partial_direct() -> None:
    session = AsyncMock()
    mock_result = MagicMock()
    mock_result.scalars.return_value.all.return_value = []
    session.execute.return_value = mock_result
    request = MagicMock()

    with patch("venvi.web.router.templates.TemplateResponse") as mock_template:
        await get_hackathons_partial(request=request, session=session)
        mock_template.assert_called()


@pytest.mark.asyncio
async def test_lifespan_direct() -> None:
    app = MagicMock()
    with patch("venvi.main.init_db", new_callable=AsyncMock) as mock_init:
        async with lifespan(app):
            pass
        mock_init.assert_called_once()
