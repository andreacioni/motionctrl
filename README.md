# motionctrl [![Travis CI](https://travis-ci.org/andreacioni/motionctrl.svg?branch=master)](https://travis-ci.org/andreacioni/motionctrl) [![Release](https://img.shields.io/badge/release-passing-brightgreen.svg)](https://github.com/andreacioni/motionctrl/releases) [![Go Report Card](https://goreportcard.com/badge/github.com/andreacioni/motionctrl)](https://goreportcard.com/report/github.com/andreacioni/motionctrl)

motionctrl is a RESTful API written in Golang that acts as a controller/proxy for [motion](https://github.com/Motion-Project/motion/) (with some sweet additional feature). It also can help you to build an IP camera and control it from any other third-part application.

__Why motionctrl?__

motionctlr allows you to:
- start/stop motion through an easy REST api service
- provide only one point to access both stream and webcontrol
- improve motion stream/webcontrol security with HTTPS
- managing motion with JSON REST api that replace the old text webcontrol interface integrated in motion
- backup old image/video in Google Drive* (archive & encryption support) (more [here](#backup))
- notify event through Telegram* to every device you want (more [here](#notification))
- host simple frontend application (more [here](#application-path))

*: more backup and notify services could be implemented easily, take a look inside backup/ notify/ folders!

__Download__

Download of precompiled version is available [here](https://github.com/andreacioni/motionctrl/releases)

__Configuration__

In order to execute *motionctrl* you need a valid JSON configuration file, an example of it could be:

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

To allow *motionctrl* to interact with motion correctly you MUST set some *motion* parameter (defined inside ```motionConfigFile```) to following values:
|Name|Value|
|----|-----|
|webcontrol_port| TCP/IP port |
|stream_port| TCP/IP port |
|stream_auth_method| 0 |
|stream_authentication| comment this |
|webcontrol_html_output| off |
|webcontrol_parms| 2 |
|webcontrol_authentication| comment this |
|process_id_file| pid file path |
|target_dir| target directory file path |

An example of a valid configuration should be:


```
target_dir /home/pi/motion/output

process_id_file /home/pi/motion/run/motion.pid

webcontrol_port 8080

webcontrol_parms 2

#webcontrol_authentication

webcontrol_html_output off

stream_port 8081

stream_auth_method 0

#stream_authentication

[...]

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

- [/control](#controlstartup)
  - [/startup](#controlstartup)
  - [/shutdown](#controlshutdown)
  - [/restart](#controlrestart)
  - [/status](#controlstatus)
- [/detection](#detectionstart)
  - [/start](#detectionstart)
  - [/stop](#detectionstop)
  - [/status](#detectionstatus)
- [/config](#configlist)
  - [/list](#configlist)
  - [/get](#configgetconfig)
  - [/set](#configset)
  - [/write](#configwrite)
- [/camera](#camerastream)
  - [/stream](#camerastream)
  - [/snapshot](#camerasnapshot)
- [/targetdir](#targetdirlist)
  - [/list](#targetdirlist)
  - [/size](#targetdirsize)
  - [/get](#targetdirgetfilename)
  - [/remove](#targetdirremovefilename)
- [/backup](#backupstatus)
  - [/status](#backupstatus)
  - [/launch](#backuplaunch)

### /control/startup

- **Description**: launch motion
- **Method**: ``` GET ```
- **Parameters**:
  - *detection*: should be used to start motion with motion detection enabled at startup (default: ```false```)
- **Return**:
  - *Status Code + Body*:
    - 200: motion started succefully
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 400: detection parameter **must** be ```true``` or ```false```
    ```
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```
    {"message": <STRING>}
    ```
 - Example:
 ```
$> curl http://10.8.0.1:8888/api/control/startup?detection=true

Output: {"message":"motion started"}
 ```
### /control/shutdown

- **Description**: shutdown motion
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: motion shutdown succefully
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
 - Example:
 ```
$> curl http://10.8.0.1:8888/api/control/shutdown

Output: {"message":"motion stopped"}
 ```

### /control/restart

- **Description**: restart motion
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: motion restarted succefully
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 409: motion not started yet
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```   
 - Example:
 ```
$> curl http://10.8.0.1:8888/api/control/restart

Output: {"message":"motion restarted"}
 ```

### /control/status

- **Description**: restart motion
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: motion status retrieved succefully
    - Response type: JSON
    ```
    {"motionStarted": true|false}
    ```
    - 500: generic internal server error
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```   
- Example:
 ```
$> curl http://10.8.0.1:8888/api/control/status

Output: {"motionStarted":false}
 ```

### /detection/start

- **Description**: start motion detection
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: motion detection enabled
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 409: motion not started yet
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
- Example:
 ```
$> curl http://10.8.0.1:8888/api/detection/start

{"message":"motion detection started"}
 ```

### /detection/stop

- **Description**: stop motion detection
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: motion detection paused
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 409: motion not started yet
    ```
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```
    {"message": <STRING>}
    ```
- Example:
 ```
$> curl http://10.8.0.1:8888/api/detection/stop

{"message":"motion detection paused"}
 ```

### /detection/status

- **Description**: return the current state of motion detection
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: motion detection status retrieved
    - Response type: JSON
    ```
    {"motionDetectionEnabled": true|false}
    ```
    - 409: motion not started yet
    ```
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```
    {"message": <STRING>}
    ```
- Example:
 ```
$> curl http://10.8.0.1:8888/api/detection/status

{"motionDetectionEnabled":false}
 ```

### /config/list

- **Description**: list all motion configuration
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: configuration list retrieved correctly
    - Response type: JSON
    ```
    { 
      <CONFIG_KEY1>: <CONFIG_VALUE1>,
      <CONFIG_KEY2>: <CONFIG_VALUE2>,
      ...
    }
    ```
    - 409: motion not started yet
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
- Example:
 ```
$> curl http://10.8.0.1:8888/api/config/list

{"area_detect":null,"auto_brightness":0,"camera":null,"camera_dir":null,"camera_id":0,"camera_name":null,"daemon":true," ... }
 ```

### /config/get/:config:

- **Description**: get the specified configuration parameter
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: configuration parameter retrieved correctly
    - Response type: JSON
    ```
    { <CONFIG_KEY>: <CONFIG_VALUE> }
    ```
    - 409: motion not started yet
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
- Example:
 ```
$> curl http://10.8.0.1:8888/api/config/get/event_gap

{"event_gap":20}
 ```

### /config/set

- **Description**: set the specified configuration to a specified value
- **Method**: ``` GET ```
- **Parameters**: 
  - *\<key\>*:\<value\> set \<key\> configuration to \<value\>
  - *writeback* (optional): indicates if the configuration will be written to the motion configuration file (default: ```false```)
- **Return**:
  - *Status Code + Body*:
    - 200: configuration set correctly
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 409: motion not started yet
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
- Example:
 ```
$> curl http://10.8.0.1:8888/api/config/set?event_gap=60&writeback=true

{"event_gap":60}
 ```

### /config/write

- **Description**: write current configuration to motion configuration file
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: configuration wrote correctly to file
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 409: motion not started yet
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
- Example:
 ```
$> curl http://10.8.0.1:8888/api/config/write

{"message":"configuration written to file"}
 ```

### /camera/stream

- **Description**: camera stream
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: streaming
    - Response type: MJPEG stream
    - 409: motion not started yet
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
- Example:
 ```
$> curl http://10.8.0.1:8888/api/control/startup

Open your browser and go to: http://localhost:8888/api/camera/stream
 ```

### /camera/snapshot

- **Description**: capture and retrieve snapshot from camera
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: snapshot
    - Response type: image
    - 409: motion not started yet
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    ```
    {"message": <STRING>}
    ```
- Example:
 ```
$> curl http://10.8.0.1:8888/api/control/startup

Open your browser and go to: http://localhost:8888/api/camera/snapshot
 ```

### /targetdir/list

- **Description**: list all files in *target_dir*
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: snapshot
    - Response type: JSON
    ```
    [
      {
        "name": <STRING>,
        "creationDate": <DATE>
      },
      ...
    ]
    ```
    - 500: generic internal server error
    ```
    {"message": <STRING>}
    ```
- Example:
 ```
$> curl http://10.8.0.1:8888/api/targetdir/list

Output: [{"name":"01-20180314152202-01.jpg","creationDate":"2018-03-14T15:22:02.88866395+01:00"},{"name":"01-20180314152202.mkv","creationDate":"2018-03-14T15:22:16.728497457+01:00"}, ... ]
 ```

### /targetdir/size

- **Description**: evaluate the *target_dir* folder size
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: size evaluated succefully
    - Response type: JSON
    ```
    { "size": <INTEGER> }
    ```
    - 500: generic internal server error
    ```
    {"message": <STRING>}
    ```
- Example:
 ```
$> curl http://10.8.0.1:8888/api/targetdir/size

Output: {"size":3753260} ~ 3.75 MB
 ```

### /targetdir/get/:filename:

- **Description**: retrieve *filename* from *target_dir*
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: file retrieved correctly
    - Response type: file from *target_dir*
    - 500: generic internal server error
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
- Example:
 ```
Open your browser and go to: http://10.8.0.1:8888/api/targetdir/get/01-20180314152211-01.jpg
 ```

### /targetdir/remove/:filename:

- **Description**: remove *filename* from *target_dir*
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: file removed correctly
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
- Example:
 ```
$> curl http://10.8.0.1:8888/api/targetdir/remove/06-20180314114422-01.jpg

Output: {"message":"06-20180314114422-01.jpg successfully removed"}
 ```

### /backup/status

- **Description**: get the current state of backup service
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: status retrieved correctly
    - Response type: JSON
    ```
    {"status": "ACTIVE_IDLE" | "ACTIVE_RUNNING" | "NOT_ACTIVE"}
    ```
    - 500: generic internal server error
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
- Example:
 ```
$> curl http://10.8.0.1:8888/api/backup/status

Output: {"status":"ACTIVE_IDLE"}
 ```

### /backup/launch

- **Description**: run backup service now
- **Method**: ``` GET ```
- **Parameters**: N.D.
- **Return**:
  - *Status Code + Body*:
    - 200: backup service stared
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
    - 500: generic internal server error
    - Response type: JSON
    ```
    {"message": <STRING>}
    ```
- Example:
 ```
$> curl http://10.8.0.1:8888/api/backup/launch

Output: {"message":"backup service is running now"}
 ```

### Internal APIs

There are some APIs that are not accessible directly by the user. These APIs (accessible from ```/internal```) are necessary to let *motion* communicate events to *motionctrl*.

This APIs are required by built-in [notification service](#notification) of *motionctrl*

# Backup

Following steps are needed only if you want to enable backup service available in *motionctrl*

*motionctrl* allows you to backup files produced by *motion* (images, videos) to your Google Drive account.

Here below some example showing, only, backup section of motionctrl config file:

1. Manual backup trigger it with (/api/backup/launch)[#backuplaunch]
```json
"backup" :  {
        "when" : "manual",
        "method" : "google",
    }
```

2. Automatic backup when ```target_dir``` size is greater than 10 Mbyte
```json
"backup" :  {
        "when" : "10MB",
        "method" : "google",
    }
```

3. Automatic backup every day at 22:00
```json
"backup" :  {
        "when" : "0 22 * * *",
        "method" : "google",
    }
```

4. Automatic backup every day at 22:00, encrypt every single file before upload
```json
"backup" :  {
        "when" : "0 22 * * *",
        "method" : "google",
        "encryptionKey": "secret_password",
    }
```

5. Automatic backup every day at 22:00, create archives with max 10 files and encrypt them before upload
```json
"backup" :  {
        "when" : "0 22 * * *",
        "method" : "google",
        "encryptionKey": "secret_password",
        "archive":true,
        "filePerArchive" : 10
    }
```

In order to correctly login to your account you must simply run *motionctrl* and follow the istructions on the command line.

# Notification

Following steps are needed only if you want to enable notification service available in *motionctrl*

- Install ```curl```
- Open your motion configuration file (e.g. /etc/motion/motion.conf)
- Set ```on_event_start``` and ```on_event_end``` ```on_picture_save``` to:

```
# Command to be executed when an event starts. (default: none)
# An event starts at first motion detected after a period of no motion defined by event_gap
on_event_start curl http://localhost:8888/internal/event/start

# Command to be executed when an event ends after a period of no motion
# (default: none). The period of no motion is defined by option event_gap.
on_event_end curl http://localhost:8888/internal/event/end

# Command to be executed when a picture (.ppm|.jpg) is saved (default: none)
# To give the filename as an argument to a command append it with %f
on_picture_save curl http://localhost:8888/internal/event/picture/saved?picturepath=%f
```

**NOTE**: curl command syntax could differ in case you have enabled HTTPS (replace ```http``` with ```https```).

Now you can add *notify* section to your *motionctrl* configuration file.

```json
"notify" : {
        "method" : "telegram",
        "token" : "324565775:JHBFEIFEIBFedae-2neuifbEDEEGEFEAF",
        "to": ["12345678", "87654321"],
        "message": "Motion recognized",
        "photo": 2
    }
```

```photo``` parameter indicates how many photos are sent to configured chats after an event starts.

# Application Path

In *motionctrl* configuration file you could specify the ```appPath``` parameter to point to the directory that contains the frontend application files.
Those files are accessible from: ```http://<IP>:<PORT>/app/```

# FAQ

 - How can I obtain valid cert/key to enable HTTPS support?
   - You can obtain them by issuing: ```openssl genrsa -out key.pem 1024 && openssl req -new -x509 -sha256 -key key.pem -out cert.pem -days 365```. This will give you a self signed certificate valid for 365 days.
   
 - How can I open encrypted backup files?
   - In order to open *.aes file you need ```aescrypt``` installed on your system. AES Crypt is a cross-plattform AES file encryption/decryption tool that you can download [here](https://www.aescrypt.com/download/).
