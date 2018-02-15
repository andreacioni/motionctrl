# motioncrtl

motionctrl is a RESTful API written in Golang that acts as a controller/proxy for [motion](https://github.com/Motion-Project/motion/). It also can help you to build an IP camera and control it from any other third-part application.

__Why motionctrl?__

motionctlr allows you to:
- start/stop motion through an easy REST api service
- provide only one point to access both stream and webcontrol
- improve motion security with HTTPS
- managing motion with JSON REST api that replace the old text webcontrol interface integrated in motion

__Download__

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
    }
}
```

__Launch__

Launch instruction

__Available API__


List of all available API and their description. Return value is always a JSON containing different values for each command (see table).


| Path | Description | Parameter | Return |
| ------------- | ------------- | ------------- | ------------- |
| /control/startup[?detection=true\|false] | Launch motion | **detection** parameter should be used to start motion with motion detection enabled at startup (default=__false__) | empty JSON object |
| /control/shutdown | Shutdown motion | no parameters | empty JSON object |
| /control/status | Report if motion is curretly running | no parameters  | **motionStarted**(bool): true if motion is running |

When call fail (HTTP status code != 200), returned JSON object contains only a **message** field with an additional description of the error.

__FAQ__

 - How can I obtain valid cert/key to enable HTTPS support?
   - BLA BLA BLA
