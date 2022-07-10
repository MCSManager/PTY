# MCSManager pty application

[![--](https://img.shields.io/badge/Go_Version-1.18.3-green.svg)](https://github.com/MCSManager)
[![--](https://img.shields.io/badge/Support-Windows/Linux-yellow.svg)](https://github.com/MCSManager)
[![--](https://img.shields.io/badge/License-MIT-red.svg)](https://github.com/MCSManager)

<br />

## What is PTY?

tty = "teletype"，pty = "pseudo-teletype"

In UNIX, /dev/tty\* is any device that acts like a "teletype"

A pty is a pseudotty, a device entry that acts like a terminal to the process reading and writing there,
but is managed by something else.
They first appeared for X Window and screen and the like,
where you needed something that acted like a terminal but could be used from another program.

<br />

## Quickstart

1. Start a PTY and set window size.

-   Note: -cmd receives an array, and the parameters of the command are passed in the form of an array, such as：`["java","-jar","ser.jar","nogui"]`

```bash
go build main.go
./main -dir "." -cmd '["bash"]' -size 50,50
```

You can execute any command, just like the SSH terminal.

```
ping google.com
top
htop
```

<br />

2. Start a PTY and dynamically change the window size.

```bash
go build main.go
./main.exe -dir "." -cmd '["cmd.exe"]'
```

Ping google.com.

```bash
{"type":1,"data":"ping google.com\r\n"}\n
```

Resize pty window size.

```
{"type":2,"data":"20,20"}\n
```

<br />

## MCSManager

MCSManager is Distributed, out-of-the-box, supports docker,
supports Minecraft and other game server management panel for the Chinese market.

This application will provide PTY functionality for MCSManager,
it is specifically designed for MCSManager,
you can also try porting to your own application.

More info: [https://github.com/mcsmanager](https://github.com/mcsmanager)

<br />

## Contributing

Interested in getting involved?

-   If you want to add a new feature, please create an issue first to describe the new feature, as well as the implementation approach. Once a proposal is accepted, create an implementation of the new features and submit it as a pull request.
-   If you are just fixing bugs, you can simply submit PR.

<br />

## MIT license

Released under the [MIT License](https://opensource.org/licenses/MIT).

Copyright 2022 [zijiren233](https://github.com/zijiren233) and contributors.
