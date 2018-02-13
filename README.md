# motioncrtl

motionctrl is a REST API written in Golang that acts as a controller for motion.

__Why motionctrl?__

motionctlr allows you

__Launch__

Launch instruction

__Configuration__

Configuration istruction

__Available API__


List of all available API and their description. Return value is always a JSON containing different values for each command (see table).


| Path | Description | Parameter | Return |
| ------------- | ------------- | ------------- | ------------- |
| /control/startup[?detection=true\|false] | Launch motion | **detection** parameter should be used to start motion with motion detection enabled at startup (default=__false__) | empty JSON object |
| /control/shutdown | Shutdown motion | no parameters | empty JSON object |
| /control/status | Report if motion is curretly running | no parameters  | **motionStarted**(bool): true if motion is running |

When call fail (HTTP status code != 200), returned JSON object contains only a **message** field with an additional description of the error.
