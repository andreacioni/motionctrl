# motionctrl [![Travis CI](https://travis-ci.org/andreacioni/motionctrl.svg?branch=master)](https://travis-ci.org/andreacioni/motionctrl) [![Go Report Card](https://goreportcard.com/badge/github.com/andreacioni/motionctrl)](https://goreportcard.com/report/github.com/andreacioni/motionctrl)

motionctrl is a RESTful API written in Golang that acts as a controller/proxy for [motion](https://github.com/Motion-Project/motion/) (with some sweet additional feature). It also can help you to build an IP camera and control it from any other third-part application.

__Why motionctrl?__

motionctlr allows you to:
- start/stop motion through an easy REST api service
- provide only one point to access both stream and webcontrol
- improve motion stream/webcontrol security with HTTPS
- managing motion with JSON REST api that replace the old text webcontrol interface integrated in motion
- backup old image/video in Google Drive* (archive & encryption support)
- notify event through Telegram* to every device you want

*: more backup and notify services could be implemented easily, take a look inside backup/ notify/ folders!

__Download__

Download instruction

__Configuration__

In order to execute motionctrl you need a valid JSON configuration file, an example of it could be:

```json
{
    "address" : "127.0.0.1",
    "port" : 8888,
    "motionConfigFile" : "/etc/motion/motion.conf",

    "username" : "user",
    "password" : "pass",

    "appPath" : "/path/to/app",

    "ssl" : {
        "key" : "/path/to/key.key",
        "cert" : "/path/to/cert.pem"
    },

    "backup" :  {
        "when" : "@every 1m",
        "method" : "google",
        "encryptionKey" : "super_secret_key",
        "archive":true,
        "filePerArchive" : 10
    },

    "notify" : {
        "method" : "telegram",
        "token" : "YOUR TELEGRAM API KEY",
        "to": ["12345678"],
        "message": "Motion recognized",
        "photo": 2
    }
}
```

__Launch__

Simple usage: ```$> ./motionctrl```

Accepted command line arguments ( ```$> ./motionctrl -h```):
```
  -a    start motion right after motionctrl
  -c string
        configuration file path (default "config.json")
  -d    when -a is set, starts with motion detection enabled
  -l string
        set log level (default "WARN")
```

# Available APIs

All the following APIs are accessible from ```/api```

- [/control](#/control/startup)
  - [/startup](#/control/startup)
  - [/shutdown](#/control/shutdown)
  - [/restart](#/control/restart)
  - [/status](#/control/status)
- [/detection](#/detection/start)
  - [/start](#/detection/start)
  - [/stop](#/detection/stop)
  - [/status](#/detection/status)
- [/config](#/config/list)
  - [/list](#/config/list)
  - [/get](#/config/get/:config:)
  - [/set](#/config/set)
  - [/write](#/config/write)
- [/camera](#/camera/stream)
  - [/stream](#/camera/stream)
  - [/snapshot](#/camera/snapshot)
- [/targetdir](#/targetdir/list)
  - [/list](#/targetdir/list)
  - [/size](#/targetdir/size)
  - [/get](#/targetdir/get/:filename:)
  - [/remove](#/targetdir/remove/:filename:)
- [/backup](#/backup/status)
  - [/status](#/backup/status)
  - [/launch](#/backup/launch)

## /control/startup

- **Description**: launch motion
- **Method**: ``` GET ```
- **Parameters**:
  - *detection*: should be used to start motion with motion detection enabled at startup (default: ```false```)
- **Return**:
  - *Status Code + Body*:
    - 200: motion started succefully
    - JSON
    ```json
    {"message": <STRING>}
    ```
    - 400: detection parameter **must** be ```true``` or ```false```
    ```json
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```
## /control/shutdown

- **Description**: shutdown motion
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: motion shutdown succefully
    - JSON
    ```json
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```

## /control/restart

- **Description**: restart motion
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: motion restarted succefully
    - JSON
    ```json
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```   
## /control/status

- **Description**: restart motion
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: motion status retrieved succefully
    - JSON
    ```json
    {"motionStarted": true|false}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```   

## /detection/start

- **Description**: start motion detection
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: motion detection enabled
    - JSON
    ```json
    {"message": <STRING>}
    ```
    - 409: generic internal server error
    ```json
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```

## /detection/stop

- **Description**: stop motion detection
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: motion detection paused
    - JSON
    ```json
    {"message": <STRING>}
    ```
    - 409: generic internal server error
    ```json
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```

## /detection/status

- **Description**: return the current state of motion detection
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: motion detection status retrieved
    - JSON
    ```json
    {"motionDetectionEnabled": true|false}
    ```
    - 409: generic internal server error
    ```json
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```

## /config/list

- **Description**: list all motion configuration
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: motion detection status retrieved
    - JSON
    ```json
    {<CONFIG_KEY1>: <CONFIG_VALUE1>, <CONFIG_KEY2>: <CONFIG_VALUE2>, ...}
    ```
    - 409: generic internal server error
    ```json
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```

## /config/get/:config:

- **Description**: get the specified configuration parameter
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: motion detection status retrieved
    - JSON
    ```json
    { 
      <CONFIG_KEY1>: <CONFIG_VALUE1>,
      <CONFIG_KEY2>: <CONFIG_VALUE2>,
      ...
    }
    ```
    - 409: generic internal server error
    ```json
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```

## /config/set

- **Description**: set the specified configuration to a specified value
- **Method**: ``` GET ```
- **Parameters**: 
  - *\<key\>*:\<value\> set \<key\> configuration to \<value\>
  - *writeback* (optional): indicates if the configuration will be written to the motion configuration file (default: ```false```)
- **Return**:
  - *Status Code + Body*:
    - 200: configuration set correctly
    - JSON
    ```json
    {"message": <STRING>}
    ```
    - 409: generic internal server error
    ```json
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```

## /config/write

- **Description**: write current configuration to motion configuration file
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: configuration set correctly
    - JSON
    ```json
    {"message": <STRING>}
    ```
    - 409: generic internal server error
    ```json
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```
## /camera/stream

- **Description**: camera stream
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: streaming
    - MJPEG stream
    - 409: generic internal server error
    ```json
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```

## /camera/snapshot

- **Description**: capture and retrieve snapshot from camera
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: snapshot
    - Image
    - 409: generic internal server error
    ```json
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```

## /targetdir/list

- **Description**: list all files in *target_dir*
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: snapshot
    - JSON
    ```json
    [
      {
        "name": <STRING>,
        "creationDate": <DATE>
      },
      ...
    ]
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```

## /targetdir/size

- **Description**: evaluate the *target_dir* folder size
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: size evaluated succefully
    - JSON
    ```json
    { "size": <INTEGER> }
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```

## /targetdir/get/:filename:

- **Description**: retrieve *filename* from *target_dir*
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: file retrieved correctly
      - Requested file from *target_dir*
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```
## /targetdir/remove/:filename:

- **Description**: remove *filename* from *target_dir*
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: file removed correctly
    - JSON
    ```json
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```
## /backup/status

- **Description**: get the current state of backup service
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: status retrieved correctly
    - JSON
    ```json
    {"status": "ACTIVE_IDLE" | "ACTIVE_RUNNING" | "NOT_ACTIVE"}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```
## /backup/launch

- **Description**: run backup service now
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: backup service stared
    - JSON
    ```json
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```json
    {"message": <STRING>}
    ```

# FAQ

 - How can I obtain valid cert/key to enable HTTPS support?
   - You can obtain them by issuing: ```openssl genrsa -out key.pem 1024 && openssl req -new -x509 -sha256 -key key.pem -out cert.pem -days 365```. This will give you a self signed certificate valid for 365 days.
   
 - How can I open encrypted backup files?
   - In order to open *.aes file you need ```aescrypt``` installed on your system. AES Crypt is a cross-plattform AES file encryption/decryption tool that you can download [here](https://www.aescrypt.com/download/).
