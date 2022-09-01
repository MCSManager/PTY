# Pseudo-teletype App

[![--](https://img.shields.io/badge/Go_Version-1.19-green.svg)](https://github.com/MCSManager)
[![--](https://img.shields.io/badge/Support-Windows/Linux-yellow.svg)](https://github.com/MCSManager)
[![--](https://img.shields.io/badge/License-MIT-red.svg)](https://github.com/MCSManager)


[English](READNE.md) | [简体中文](README_CN.md)

<br />

## What is PTY?


<div align=center>

![term](https://user-images.githubusercontent.com/18360009/180396380-b2ec74c4-dcab-4405-a72a-2c66c4b3eac4.png)

</div>



tty = "teletype"，pty = "pseudo-teletype"

In UNIX, /dev/tty\* is any device that acts like a "teletype"

A pty is a pseudotty, a device entry that acts like a terminal to the process reading and writing there,
but is managed by something else.
They first appeared for X Window and screen and the like,
where you needed something that acted like a terminal but could be used from another program.

<br />

## Quickstart

Start a PTY and set window size.

- Note: -cmd receives an array, and the parameters of the command are passed in the form of an array and needs to be serialized, such as：`[\"java\",\"-jar\",\"ser.jar\",\"nogui\"]`

```bash
go build
./pty -dir "." -cmd [\"bash\"] -size 50,50
```

You can execute any command, just like the SSH terminal.

```
ping google.com
top
htop
```

<br />

## Flags:

```
  -cmd string
        command
  -coder string
        Coder (default "UTF-8")
  -color
        colorable (default false)
  -dir string
        command work path (default ".")
  -size string
        Initialize pty size, stdin will be forwarded directly (default "50,50")
  -test
        Test whether the system environment is pty compatible
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

- If you want to add a new feature, please create an issue first to describe the new feature, as well as the implementation approach. Once a proposal is accepted, create an implementation of the new features and submit it as a pull request.
- If you are just fixing bugs, you can simply submit PR.

<br />

## MIT license

Released under the [MIT License](https://opensource.org/licenses/MIT).

Copyright 2022 [zijiren233](https://github.com/zijiren233) and contributors.
