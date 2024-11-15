#include <WiFi.h>
#include <HTTPClient.h>
#include <esp_wifi.h>


//const int led = 9; // Led positive terminal to the digital pin 9.              
const  int pirSensor = 5; // signal pin of sensor to digital pin 5.               
int state = LOW;            
int pirValue = 0;     


const char* ssid = "..."; // TODO
const char* password = "..."; // TODO
const String api_endpoint_notify = "https://callboxes.etincelle-coworking.com/api/motion";
const String api_endpoint_alive = "https://callboxes.etincelle-coworking.com/api/alive";

unsigned long lastTime = 0;

const int DELAY_EVERY_MINUTE = 60000;
// Timer set to 10 minutes (600000)
//unsigned long timerDelay = 600000;
// Set timer to 5 seconds (5000)
unsigned long aliveDelay = DELAY_EVERY_MINUTE;


void setup() {
  Serial.begin(9600);   
  
  pinMode (LED_BUILTIN, OUTPUT);  // set the LED pin mode
  pinMode(pirSensor, INPUT); // PIR motion sensor is determined is an input here.  

  delay(1000);

  WiFi.mode(WIFI_MODE_STA);
  //WiFi.config(INADDR_NONE, INADDR_NONE, INADDR_NONE, INADDR_NONE);

  //WiFi.STA.begin();

  Serial.print("[DEFAULT] ESP32 Board MAC Address: " + WiFi.macAddress());
  
  
  // We start by connecting to a WiFi network
  Serial.println();
  Serial.println();
  Serial.print ("Connecting to ");
  Serial.println (ssid);

  WiFi.begin (ssid, password);

  while (WiFi.status() != WL_CONNECTED) {
    delay (500);
    Serial.print (".");
  }

  Serial.println ("");
  Serial.println ("WiFi connected.");
  Serial.println ("IP address: ");
  Serial.println (WiFi.localIP());

  lastTime = millis();
}

void loop(){ 
  delay (500);
  pirValue = digitalRead(pirSensor);

  if  (pirValue == HIGH) {           
    if (state == LOW) {
      //Serial.println("  Motion detected "); 
      digitalWrite(LED_BUILTIN, HIGH);
      state = HIGH;       
      notifymMotion();
    }
  } else {
      if  (state == HIGH){
       // Serial.println("The action/ motion has stopped");
        digitalWrite(LED_BUILTIN, LOW);
        state = LOW;       
    }
  }

  if ( (millis() - lastTime) > aliveDelay) {
      notifyWatchDog();
      lastTime = millis();
  }
  
}

void notifyWatchDog(){
  HTTPClient http;
  String uri = api_endpoint_alive + "?mac=" + WiFi.macAddress();
  Serial.println("Calling " + uri);
  http.begin(uri.c_str());

  int httpResponseCode = http.GET();
  if (httpResponseCode>0) {
    Serial.print("HTTP Response code: ");
    Serial.println(httpResponseCode);
    String payload = http.getString();
    Serial.println(payload);
  }else {
    Serial.print("Error code: ");
    Serial.println(httpResponseCode);
  }
  http.end();
}

void notifymMotion(){
  HTTPClient http;
  String uri = api_endpoint_notify + "?mac=" + WiFi.macAddress();
  Serial.println("Calling "+uri);
  http.begin(uri.c_str());

  int httpResponseCode = http.GET();
  if (httpResponseCode>0) {
    Serial.print("HTTP Response code: ");
    Serial.println(httpResponseCode);
    String payload = http.getString();
    Serial.println(payload);
  }else {
    Serial.print(" -> Error code: ");
    Serial.println(httpResponseCode);
  }
  http.end();
}
