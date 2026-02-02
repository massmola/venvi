from datetime import datetime

from sqlmodel import Column, DateTime, Field, SQLModel


class ODHEvent(SQLModel, table=True):
    """
    Represents an event from the South Tyrol Open Data Hub.
    """

    id: str = Field(primary_key=True)
    title: str
    description: str | None = None
    date_start: datetime = Field(sa_column=Column(DateTime(timezone=True)))
    date_end: datetime = Field(sa_column=Column(DateTime(timezone=True)))
    location: str | None = None
    image_url: str | None = None
    source_url: str | None = None
    is_new: bool = True
    taken: bool = False
