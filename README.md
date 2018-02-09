# motioncrtl

motionctrl is a REST API written in Golang that acts as a controller for motion.

**API**

| Path | Description |
| ------------- | ------------- |
| /startup[?detection=true\|false] | Startup motion instance, **detection** parameter should be used to start motion with motion detection enabled at startup|
| /shutdown | Shutdown motion |
