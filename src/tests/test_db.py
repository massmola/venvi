from typing import Any
from unittest.mock import patch

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine
from sqlmodel.pool import StaticPool
from venvi.core.db import get_session, init_db


@pytest.mark.asyncio
async def test_init_db() -> None:
    """Test init_db logic (using a test engine to avoid real DB)."""
    test_engine = create_async_engine(
        "sqlite+aiosqlite://",
        connect_args={"check_same_thread": False},
        poolclass=StaticPool,
    )

    with patch("venvi.core.db.engine", test_engine):
        await init_db()
        # If it runs without error, it's covered.
        # We verify it by checking if metadata was created.
        async with test_engine.begin() as conn:
            # Simple check if any table exists
            from sqlalchemy import inspect

            def check_tables(connection: Any) -> Any:
                return inspect(connection).get_table_names()

            tables = await conn.run_sync(check_tables)
            assert "hackathon" in tables


@pytest.mark.asyncio
async def test_get_session_generator() -> None:
    """Test the get_session generator explicitly."""
    test_engine = create_async_engine(
        "sqlite+aiosqlite://",
        connect_args={"check_same_thread": False},
        poolclass=StaticPool,
    )

    with patch("venvi.core.db.engine", test_engine):
        generator = get_session()
        session = await anext(generator)
        assert isinstance(session, AsyncSession)
        await session.close()
        try:
            await anext(generator)
        except StopAsyncIteration:
            pass
