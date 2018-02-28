# motionctrl <img src="https://travis-ci.org/andreacioni/motionctrl.svg?branch=master"> [![Go Report Card](https://goreportcard.com/badge/github.com/andreacioni/motionctrl)](https://goreportcard.com/report/github.com/andreacioni/motionctrl)

motionctrl is a RESTful API written in Golang that acts as a controller/proxy for [motion](https://github.com/Motion-Project/motion/) (with some sweet additional feature). It also can help you to build an IP camera and control it from any other third-part application.

__Why motionctrl?__

motionctlr allows you to:
- start/stop motion through an easy REST api service
- provide only one point to access both stream and webcontrol
- improve motion security with HTTPS
- managing motion with JSON REST api that replace the old text webcontrol interface integrated in motion
- backup old image/video in Google Drive* (archive & encryption support)
- notify event through Telegram* to every device you want (TODO)

*: more backup and notify services could be implemented easily, take a look inside backup/ notify/ folders!

__Download__

Download instruction

__Configuration__

In order to execute motionctrl you need a valid JSON configuration file, an example of it could be:

```json
{
    "address" : "127.0.0.1",
    "port" : 8888,
    "motionConfigFile" : "/home/andreacioni/motion/motion.conf",

    "username" : "user",
    "password" : "pass",

    "ssl" : {
        "key" : "/path/to/key.key",
        "cert" : "/path/to/cert.pem"
    },

    "backup" :  {
        "when" : "@every 1m",
        "method" : "google",
        "encryptionKey" : "NOT YET IMPLEMENTED",
        "archive":true,
        "filePerArchive" : 10,
        "keepLocalCopy":false
    },
}
```

__Launch__

Launch instruction

__Available API__


List of all available API and their description. Return value is always a JSON containing different values for each command (see table).


| Path | Description | Parameter | Return | Status Codes | 
| ------------- | ------------- | ------------- | ------------- | ------------- |
| /control/startup[?detection=(true\|false)] | Launch motion | **detection** parameter should be used to start motion with motion detection enabled at startup (default=__false__) | JSON object | **200**: motion started<br>**500**: there was an error on starting up motion |
| /control/shutdown | Shutdown motion | no parameters | JSON object |  **200**: motion stopped<br>**500**: there was an error on stopping motion  |
| /control/restart | Restart motion | no parameters | JSON object | **200**: motion restarted<br>**500**: there was an error on restarting motion
| /control/status | Report if motion is curretly running | no parameters  | **motionStarted**(bool): true if motion is running |  **200**: always |
| /detection/start | Enable motion detection | no parameters  |  JSON object |  **200**: motion detection enabled<br>**500**: there was an error on enabling motion detection<br>**409**: motion is not started |
| /detection/stop | Disable motion detection | no parameters  |  JSON object |  **200**: motion detection enabled<br>**500**: there was an error on enabling motion detection<br>**409**: motion is not started |
| /detection/status | Check if motion detection is enabled | no parameters |  **motionDetectionEnabled**(bool): true if motion detection is enabled|  **200**: if this checks succed<br>**500**: there was an error on checking motion detection enabled<br>**409**: motion is not started |
| /config/list | List current motion configuration | no parameters | JSON object which attributes contain all motion configuration parameters |  **200**: configuration obtained without errors<br>**500**: there was an error on retrieving motion configuraition<br>**409**: motion is not started |
| /config/get?query=\<param\> | Get a specific configuration | **query** parameter indicates requested configuration parameter | JSON object with only one attribute which name is \<param\> |  **200**: configuration obtained without errors<br>**400**: 'query' parameter is not specified<br>**500**: there was an error on retrieving motion configuraition<br>**409**: motion is not started |
| /config/set?\<name\>=\<value\>[&writeback=(true\|false)] | Set a specific configuration | **name** parameter indicates the configuration parameter to update<br>**writeback**(bool): 'true' indicates that configuration will be written to file (default=__false__) | JSON object with only one attribute which name is \<param\> |  **200**: configuration set without errors<br>**400**: 'writeback' parameter has an invalid value, allowed: 'true' or 'false'<br>**403**: attempting to write read-only configuration parameters<br>**500**: there was an error on setting motion configuraition<br>**409**: motion is not started |
| /config/write | Write current configuration to file | **name** parameter indicates the configuration parameter to update<br>**writeback**(bool): 'true' indicates that configuration will be written to file (default=__false__) | JSON object with only one attribute which name is \<param\> |  **200**: configuration set without errors<br>**500**: there was an error on writing motion configuration to file<br>**409**: motion is not started |

When call fail (HTTP status code != 200), returned JSON object has only a **message** field containing an additional description of the error.

__FAQ__

 - How can I obtain valid cert/key to enable HTTPS support?
   - BLA BLA BLA
   
 - How can I enable Google Drive backup?
   - BLA BLA BLA
