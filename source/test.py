# example_pydantic_bytes.py
from datetime import datetime
from pydantic import BaseModel
from typing import List, Optional, Dict
import json

class Item(BaseModel):
    id: int
    name: str
    price: float
    tags: List[str] = []

class Order(BaseModel):
    order_id: str
    created_at: datetime
    items: List[Item]
    shipped: bool = False
    metadata: Optional[Dict[str, str]] = None

if __name__ == "__main__":
    # create example data using pydantic models (this is your "json dict")
    order = Order(
        order_id="ORD-1001",
        created_at=datetime.utcnow(),
        items=[
            Item(id=1, name="Widget", price=9.99, tags=["blue", "small"]),
            Item(id=2, name="Gadget", price=19.95, tags=["red"]),
        ],
        shipped=False,
        metadata={"customer": "ACME Corp", "priority": "high"},
    )

    # serialize to a JSON string, then to bytes
    json_str = order.model_dump_json(exclude_none=True)  # returns str
    json_bytes = json_str.encode("utf-8")     # serialize as bytes

    # print outputs
    print("JSON string:")
    print(json_str)
    print("\nBytes (raw):")
    print(json_bytes)
    print("\nBytes (repr):")
    print(repr(json_bytes))