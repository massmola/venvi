from datetime import datetime

from sqlmodel import JSON, Column, DateTime, Field, SQLModel


class Event(SQLModel, table=True):
    """
    A unified model representing an event from any source.
    """

    id: str = Field(primary_key=True)
    title: str
    description: str | None = None
    date_start: datetime = Field(sa_column=Column(DateTime(timezone=True)))
    date_end: datetime = Field(sa_column=Column(DateTime(timezone=True)))
    location: str | None = None
    url: str
    image_url: str | None = None
    source_name: str
    source_id: str
    topics: list[str] = Field(default=[], sa_column=Column(JSON))
    category: str = "general"
    is_new: bool = True

    # Note: 'taken' field was removed globally from the project.
