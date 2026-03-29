#ifndef WIFI_MANAGER_H
#define WIFI_MANAGER_H

#include "esp_err.h"

#define WIFI_SSID "Anya"
#define WIFI_PASS "cutiepie"

#define WIFI_STATIC_IP "10.0.0.2"
#define WIFI_GATEWAY   "10.0.0.1"
#define WIFI_NETMASK  "255.255.255.0"

esp_err_t wifi_init_sta(void);

#endif // WIFI_MANAGER_H
