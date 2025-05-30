from django.db import models

# Create your models here.
class UserProfile(models.Model):
    phone = models.CharField(max_length=20, unique=True)
    role = models.IntegerField(choices=((0, 'User'), (1, 'Admin')), default=0)
    uid = models.CharField(max_length=100, unique=True)


class Audience(models.Model):
    room_number = models.CharField(max_length=10, unique=True)
    is_reserved = models.BooleanField(default=False)
    reserved_by = models.ForeignKey(UserProfile, null=True, blank=True, on_delete=models.SET_NULL)

class KeyCard(models.Model):
    uid_keycard = models.CharField(max_length=100)
    user = models.ForeignKey(UserProfile, on_delete=models.CASCADE)



