from datetime import UTC, datetime

import pytest
from httpx import AsyncClient
from sqlalchemy.ext.asyncio import AsyncSession
from venvi.models.hackathon import Hackathon


@pytest.mark.asyncio
async def test_read_hackathons_empty(client: AsyncClient) -> None:
    response = await client.get("/hackathons/")
    assert response.status_code == 200
    assert response.json() == []


@pytest.mark.asyncio
async def test_sync_endpoint_mocked(client: AsyncClient, session: AsyncSession) -> None:
    from unittest.mock import patch

    with patch("venvi.services.ingestion.fetch_euro_hackathons", return_value=[]):
        response = await client.post("/hackathons/sync")
        assert response.status_code == 200
        assert response.json() == {"message": "Sync complete", "new_items": 0}


@pytest.mark.asyncio
async def test_create_and_read_hackathon(
    client: AsyncClient, session: AsyncSession
) -> None:
    hackathon = Hackathon(
        id="209e1cfb-5040-4086-bfec-f67bdc3380ff",
        name="API Test Hackathon",
        city="API City",
        country_code="AC",
        date_start=datetime(2026, 2, 3, 0, 0, 0, tzinfo=UTC),
        date_end=datetime(2026, 2, 4, 0, 0, 0, tzinfo=UTC),
        topics=[],
        url="https://api.test",
        status="upcoming",
    )
    session.add(hackathon)
    await session.commit()

    response = await client.get("/hackathons/")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["name"] == "API Test Hackathon"


@pytest.mark.asyncio
async def test_read_hackathons_filtering(
    client: AsyncClient, session: AsyncSession
) -> None:
    """Test filtering hackathons by status."""
    h1 = Hackathon(
        id="a09e1cfb-5040-4086-bfec-f67bdc3380ff",
        name="Upcoming Hack",
        city="City A",
        country_code="AA",
        date_start=datetime(2026, 2, 3, 0, 0, 0, tzinfo=UTC),
        date_end=datetime(2026, 2, 4, 0, 0, 0, tzinfo=UTC),
        topics=[],
        url="https://a.test",
        status="upcoming",
    )
    h2 = Hackathon(
        id="b09e1cfb-5040-4086-bfec-f67bdc3380ff",
        name="Past Hack",
        city="City B",
        country_code="BB",
        date_start=datetime(2025, 2, 3, 0, 0, 0, tzinfo=UTC),
        date_end=datetime(2025, 2, 4, 0, 0, 0, tzinfo=UTC),
        topics=[],
        url="https://b.test",
        status="past",
    )
    session.add(h1)
    session.add(h2)
    await session.commit()

    # Test filtering for upcoming
    response = await client.get("/hackathons/", params={"status": "upcoming"})
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["name"] == "Upcoming Hack"

    # Test filtering for past
    response = await client.get("/hackathons/", params={"status": "past"})
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["name"] == "Past Hack"


@pytest.mark.asyncio
async def test_sync_endpoint_failure(
    client: AsyncClient, session: AsyncSession
) -> None:
    """Test the sync endpoint when it encounters an error."""
    from unittest.mock import patch

    with patch(
        "venvi.api.routers.hackathons.sync_hackathons",
        side_effect=Exception("Database Connection Error"),
    ):
        response = await client.post("/hackathons/sync")
        assert response.status_code == 500
        assert "Database Connection Error" in response.json()["detail"]
