from datetime import datetime
from uuid import UUID

from sqlmodel import JSON, Column, DateTime, Field, SQLModel


class Hackathon(SQLModel, table=True):
    """
    Represents a hackathon event in the EU.

    This model stores core information about hackathons aggregated from various sources,
    including their location, dates, topics, and status.
    """

    id: UUID = Field(primary_key=True)
    name: str
    city: str
    country_code: str
    date_start: datetime = Field(sa_column=Column(DateTime(timezone=True)))
    date_end: datetime = Field(sa_column=Column(DateTime(timezone=True)))
    topics: list[str] = Field(sa_column=Column(JSON))
    notes: str | None = None
    url: str
    status: str
    is_new: bool = False
