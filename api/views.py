from rest_framework.response import Response
from rest_framework.decorators import api_view
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework import status
import requests
from rest_framework import viewsets, permissions
from rest_framework.decorators import action
from django.contrib.auth.models import User
from .services import GoServerService


class GetKeyFromESP(APIView):
    def get(self, request, esp_id):
        try:
            url = f"http://192.168.1.104:3333/key?ESP_ID={esp_id}"
            resp = requests.get(url, timeout=5)

            if resp.status_code == 200:
                if resp.text == "ERR":
                    return Response({"error": "ESP повернув помилку"}, status=502)
                return Response({"uid": resp.text.strip()}, status=200)
            else:
                return Response({"error": "Помилка з'єднання з ESP"}, status=resp.status_code)

        except requests.exceptions.RequestException:
            return Response({"error": "ESP не відповідає"}, status=504)


@api_view(['GET'])
def test_get(request):
    person = {'Name': 'Artyom', 'Age': 28}
    return Response(person)


class RoomViewSet(viewsets.ViewSet):
    permission_classes = [permissions.IsAuthenticated]
    go_service = GoServerService()

    def list(self, request):
        """Get list of all rooms."""
        rooms = self.go_service.get_rooms()
        return Response(rooms)

    @action(detail=True, methods=['post'])
    def check_access(self, request, pk=None):
        """Check if a user has access to a room."""
        card_uid = request.data.get('card_uid')
        if not card_uid:
            return Response(
                {"error": "card_uid is required"},
                status=status.HTTP_400_BAD_REQUEST
            )

        has_access = self.go_service.check_access(card_uid, int(pk))
        return Response({"has_access": has_access})


class BookingViewSet(viewsets.ViewSet):
    permission_classes = [permissions.IsAuthenticated]
    go_service = GoServerService()

    def create(self, request):
        """Create a new booking."""
        booking_data = {
            "room_id": request.data.get('room_id'),
            "user_id": request.user.id,
            "start_time": request.data.get('start_time'),
            "end_time": request.data.get('end_time')
        }

        booking = self.go_service.create_booking(booking_data)
        if booking:
            return Response(booking, status=status.HTTP_201_CREATED)
        return Response(
            {"error": "Failed to create booking"},
            status=status.HTTP_400_BAD_REQUEST
        )


