# Train Departures MQTT

A service to poll the Realtime Trains API for train departures between two stations. Departures are published as MQTT messages.

## Usage

| Argument           | Env var           | Description                                   | Default       |
|--------------------|-------------------|-----------------------------------------------|---------------|
| -rtt-username      | RTT_USERNAME      | Realtime Trains API Auth Credentials username | _required_    |
| -rtt-password      | RTT_PASSWORD      | Realtime Trains API Auth Credentials password | _required_    |
| -rtt-from          | RTT_FROM          | 3 char code of the departure station          | _required_    |
| -rtt-to            | RTT_TO            | 3 char code of the arrival station            | _required_    |
| -mqtt-broker       | MQTT_BROKER       | MQTT broker URL                               | _required_    |
| -mqtt-port         | MQTT_PORT         | MQTT broker port                              | 1883          |
| -mqtt-username     | MQTT_USERNAME     | MQTT broker username                          | <none>        |
| -mqtt-password     | MQTT_PASSWORD     | MQTT broker password                          | <none>        |
| -mqtt-topic-prefix | MQTT_TOPIC_PRFEIX | MQTT topic prefix                             | `traindepart` |

3 digit station codes can be found here https://en.wikipedia.org/wiki/UK_railway_stations

## MQTT

The MQTT message payload is a JSON array containing the next few departures with scheduled and expected (accounting for delays) departure times. E.g.

```json
[
    {
      "scheduledTime": "2022-05-21T21:51:00Z",
      "expectedTime": "2022-05-21T21:51:00Z"
    },
    {
      "scheduledTime": "2022-05-21T21:55:00Z",
      "expectedTime": "2022-05-21T21:58:00Z"
    },
    ...
]
```

## Data

All data is provided by the Realtime Trains API. https://api.rtt.io/. This project is not affiliated with them.

## Docker Compose with Mosquitto

Only for Raspberry Pi right now.

```yaml
mosquitto:
    container_name: mosquitto
    image: eclipse-mosquitto:1.6.14
    restart: unless-stopped
    ports:
    - 1883:1883
    - 9001:9001

traindepart-mqtt:
    container_name: traindepartmqtt
    image: rskeenan/traindepart-mqtt:0.0.3-armv7
    restart: unless-stopped
    environment:
    - RTT_USERNAME=<username>
    - RTT_PASSWORD=<password>
    - RTT_FROM=BTN
    - RTT_TO=VIC
    - MQTT_BROKER_URL=mosquitto
```