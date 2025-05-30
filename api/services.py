import requests
from django.conf import settings
from typing import Dict, Any, Optional

class GoServerService:
    def __init__(self):
        self.base_url = settings.GO_SERVER_URL

    def check_access(self, card_uid: str, esp_id: int) -> bool:
        """Check if a user has access to a room."""
        response = requests.get(
            f"{self.base_url}/access",
            params={"card_uid": card_uid, "esp_id": esp_id}
        )
        return response.text == "GRANTED"

    def get_rooms(self) -> list:
        """Get list of all rooms."""
        response = requests.get(f"{self.base_url}/rooms")
        return response.json()

    def create_booking(self, booking_data: Dict[str, Any]) -> Optional[Dict[str, Any]]:
        """Create a new booking."""
        response = requests.post(
            f"{self.base_url}/bookings",
            json=booking_data
        )
        if response.status_code == 200:
            return response.json()
        return None 