from datetime import datetime
from typing import List, Optional
from uuid import UUID

from sqlmodel import Field, SQLModel, JSON, Column, DateTime


class Hackathon(SQLModel, table=True):
    id: UUID = Field(primary_key=True)
    name: str
    city: str
    country_code: str
    date_start: datetime = Field(sa_column=Column(DateTime(timezone=True)))
    date_end: datetime = Field(sa_column=Column(DateTime(timezone=True)))
    topics: List[str] = Field(sa_column=Column(JSON))
    notes: Optional[str] = None
    url: str
    status: str
    is_new: bool = False
