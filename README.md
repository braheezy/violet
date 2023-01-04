# Violet
Give [Vagrant](https://developer.hashicorp.com/vagrant) a splash of color :art:

Violet is a colorful TUI frontend to manage Vagrant virtual machines. Quickly view the state of all VMs and issue commands against them!

![Violet Gif](./assets/demo.gif)

## Project Status
Violet is in early stages of development and is not recommended for production use cases. It probably handles most error cases poorly. It hardly cares about terminal types and sizes.

## Getting Started

Violet is delivered as a single binary for various platforms. See the [Releases](https://github.com/braheezy/violet/releases) page for the latest builds.

### Prerequisites

Violet does absolutely nothing without Vagrant installed. See the [Vagrant docs](https://developer.hashicorp.com/vagrant/downloads) to install it for your platform.

Vagrant itself does absolutely nothing unless you have a Hypervisor installed and configured. Here's a few popular ones:
- [VirtualBox](https://www.virtualbox.org/)
- [Libvirt/QEMU](https://libvirt.org/)

For best results, it helps to have existing Vagrant VMs.

### Usage
Open a terminal and run the program:

    $ violet

See the following table for how to interact with Violet:
| Action                  | Key        | Description                                               |
|-------------------------|------------|-----------------------------------------------------------|
| Switch Environment Tab  | Tab        | Cycle through different Vagrant environments.             |
| Select Virtual Machines | Left/Right | Cycle through the different VMs in a Vagrant environment. |
| Select Command          | 1,2,3,4    | Choose the command by number                              |
| Scroll Vagrant Output   | Up/Down    | Scroll the Vagrant output area to see more text           |

Note that Violet does not aim to support all Vagrant commands and will provide a poor interface for troubleshooting issues with Vagrant, VMs, hypervisors, etc.

## Development

The `Makefile` contains the most common developer actions to perform. See `make help`.

## Built with :heart: using other people's tools

* [bubbletea](https://github.com/charmbracelet/bubbletea) - Main TUI framework
* [lipgloss](https://github.com/charmbracelet/lipgloss) - Styling and colors
* [bubbletint](https://github.com/lrstanley/bubbletint) - Pre-made lipgloss colors

## Contributing

## Inspiration
My interest in TUI applications is growing and I wanted something more complicated and useful (than a [game](https://github.com/braheezy/hangman)) to build. And I got to a learn lots of Go!
