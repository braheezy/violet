# Violet
Give [Vagrant](https://developer.hashicorp.com/vagrant) a splash of color :art:

Violet is a colorful TUI frontend to manage Vagrant virtual machines. Quickly view the state of all VMs and issue commands against them!

![Violet Gif](./assets/demo.gif)

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

    violet

See the following table for how to interact with Violet:
| Action                  | Key        | Description                                               |
|-------------------------|------------|-----------------------------------------------------------|
| Switch Environment Tab  | Tab/Shift+Tab | Cycle through found Vagrant environments       |
| Select Command | Left/Right | Cycle through the supported Vagrant commands |
| Run command | Enter | Run the highlighted command on the selected entity |
| Toggle Environments/VM control | Space bar | Operate on the environment as a whole or individual machines |


Note that Violet does not aim to support all Vagrant commands and will provide a poor interface for troubleshooting issues with Vagrant, VMs, hypervisors, etc.

## Development

The `Makefile` contains the most common developer actions to perform. See `make help` for everything, or build and run for your machine:

    make run

## Acknowledgements

* [bubbletea](https://github.com/charmbracelet/bubbletea) - Main TUI framework
* [lipgloss](https://github.com/charmbracelet/lipgloss) - Styling and colors
* [bubbletint](https://github.com/lrstanley/bubbletint) - Pre-made lipgloss colors

## Contributing

## Roadmap
In somewhat priority order:
- [x] Reduce magic numbers in sizing
- [x] Better screen resize support and smarter app sizing in general
- [x] Bulk operations on VMs, including the entire Environment (e.g. `vagrant up` instead of `vagrant up <machine>`)
- [x] Error handling:
  - Violet errors: Show in message area.
- [x] File logging
- [x] Icons
---
- [ ] Support destroy
- [ ] Remember user selections between
- [ ] Load VMs and Envs deterministically
- [ ] Launch SSH sessions in external apps
- [ ] Mouse support
- [ ] Add VHS to CI
- [ ] Pagination to handle many environments and/or VMs
- [ ] Configuration w/ cobra for Env Vars, CLI, and config file:
  - Theme
  - Log control

## Inspiration
My interest in TUI applications was growing and I wanted to build something complicated and useful (more than a [game](https://github.com/braheezy/hangman)).
