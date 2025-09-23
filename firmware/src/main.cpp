#include <Arduino.h>
#include <MFRC522.h>
#include <WiFi.h>
#include <HTTPClient.h>

#define SDA_PIN 22
#define RST_PIN 5
#define LED_PIN 2   // status LED pin

MFRC522 rfid = MFRC522(SDA_PIN, RST_PIN);

String ssid = "Pixel_3504";
String password = "pixel123";
String serverURL = "http://10.154.5.203:3333/key?ESP_ID=";
int ESP_ID = 315;

String adminUID = "99A39359";
String responseUID = "";
int updateTime = 2000; // millis

String getKeyUID();
void doorOpen();
void doorClose();

void setup() {
  Serial.begin(115200);
  SPI.begin();
  rfid.PCD_Init();
  pinMode(LED_PIN, OUTPUT);

  // Set up Wi-Fi connection
  WiFi.begin(ssid, password);
  Serial.print("Connecting to WiFi.");
  while (WiFi.status() != WL_CONNECTED) {
    Serial.print('.');
    delay(500);
  }
  Serial.println("Connected to WiFi!");
}

int lastMillis = 0;
bool fl = false;

void loop() {
  // updating key from server every updateTime milliseconds
  if (millis() > lastMillis + updateTime) {
    lastMillis = millis();

    if (WiFi.status() != WL_CONNECTED ) {
      responseUID = "FF";
      Serial.println("WiFi connection lost");
    } else {
      responseUID = getKeyUID();
    }
    digitalWrite(LED_PIN, (responseUID == "FF"));
  }

  if (rfid.PICC_IsNewCardPresent()) {
    rfid.PICC_ReadCardSerial();

    String cardUID;
    for (int i = 0; i < rfid.uid.size; i++) {
      cardUID += rfid.uid.uidByte[i] < 0x10 ? "0" + String(rfid.uid.uidByte[i], HEX) : String(rfid.uid.uidByte[i], HEX);
    }

    cardUID.toUpperCase();
    if (cardUID == responseUID || cardUID == adminUID) {
      if (!fl) doorOpen();
      else doorClose();
      fl = !fl;
    }
    
    rfid.PICC_HaltA();
  }
}

String getKeyUID() {
  HTTPClient http;
  http.begin(serverURL + String(ESP_ID));

  int httpResponseCode = http.GET();
  if (httpResponseCode == HTTP_CODE_OK) {
    String payload = http.getString();
    if (payload == "ERR") {
      Serial.printf("Device with ESP_ID = %u is not registered in server DB\n", ESP_ID);
      return "";
    }

    Serial.println("Response: " + payload);
    http.end();
    return payload;

  } else {
    Serial.print("HTTP GET request failed with error code: ");
    Serial.println(httpResponseCode);
    http.end();
    return "";
  }
}

void doorOpen() {
  Serial.println("Door opened!");
}
void doorClose() {
  Serial.println("Door closed!");
}