from datetime import UTC, datetime

import pytest
from httpx import AsyncClient
from sqlalchemy.ext.asyncio import AsyncSession
from venvi.models.hackathon import Hackathon


@pytest.mark.asyncio
async def test_index_page(client: AsyncClient) -> None:
    """Test that the index page loads and contains HTMX triggers."""
    response = await client.get("/")
    assert response.status_code == 200
    assert "Venvi" in response.text
    # Check for HTMX trigger
    assert 'hx-get="/partials/hackathons"' in response.text
    assert 'hx-trigger="load, reload from:body"' in response.text


@pytest.mark.asyncio
async def test_hackathons_partial_empty(client: AsyncClient) -> None:
    """Test the hackathons partial when no data is present."""
    response = await client.get("/partials/hackathons")
    assert response.status_code == 200
    # Should render an empty list or appropriate partial content
    # (Check for lack of 500 status and presence of expected HTML structure)
    assert "text-brand-400" not in response.text  # No cards should be rendered


@pytest.mark.asyncio
async def test_hackathons_partial_with_data(
    client: AsyncClient, session: AsyncSession
) -> None:
    """Test the hackathons partial with actual data, catching rendering/async errors."""
    hackathon = Hackathon(
        id="309e1cfb-5040-4086-bfec-f67bdc3380ff",
        name="Integration Test Hackathon",
        city="Test City",
        country_code="TC",
        date_start=datetime(2026, 2, 3, 0, 0, 0, tzinfo=UTC),
        date_end=datetime(2026, 2, 4, 0, 0, 0, tzinfo=UTC),
        topics=["test", "integration"],
        url="https://integration.test",
        status="upcoming",
    )
    session.add(hackathon)
    await session.commit()

    response = await client.get("/partials/hackathons")
    assert response.status_code == 200
    assert "Integration Test Hackathon" in response.text
    assert "Test City, TC" in response.text
    assert "#test" in response.text
