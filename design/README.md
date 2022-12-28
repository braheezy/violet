## Requirements
1. The program will need to initialize the third party framework for rendering. This will allow the program to display a graphical view of the VMs to the user.
2. The program will need to connect to Vagrant and retrieve a list of the VMs that are currently available on the workstation.
3. The program will need to display the list of VMs to the user in the graphical view, allowing the user to select one or more VMs to interact with.
4. The program will need to allow the user to enter Vagrant commands and execute them on the selected VMs.
5. The program will need to monitor the status of the VMs and update the graphical view accordingly.

## Project Structure
- `cmd/violet/`: Code for `violet` CLI. Entry point for the TUI. Imports and uses code from `internal`.
- `internal/app/`: Core logic of Violet
    - `app`: Application state and logic
    - `view`: User interface
    - `update`: How to update the application state based on user events
- `pkg/vagrant/`: Code for Vagrant client and its API. Used by the `app` package to interact with Vagrant
- `test/`: Unit tests for packages

`cmd/violet/main.go` contains the entry point for the program. This file is responsible for any configuration required before launching the app.

`internal/app/` contains the core logic for the program.
  - `app.go` initializes the rendering framework and starts the main application loop. Manages internal state.
  - `models.go` contains all app model definitions and logic
  - `update.go` handles everything from user input to the result of that input:
    - execute Vagrant commands with `pkg/vagrant` client
    - handle Vagrant output
    - update model according to the current state of VMs, driving changes in View.
  - `view.go` contains the code for rendering the graphical view of the app

`pkg/vagrant/vagrant.go` contains the logic for connecting to Vagrant, executing Vagrant commands, and returning sensible data to the app.

## Wireframe
+------------------------------------------------------------+
|
|                     Violet
|      A splash of color for `vagrant`  :art:
|
|  Environments:
|        [/path/to/env1] (1 VM)
|    --> [/path/to/env2] (2 VM)
|        [/path/to/env3] (1 VM)
|
|  VMs in [env2]:
|        [ ] vm1 (provider: virtualbox, state: running)
|        [x] vm2 (provider: vmware,     state: not created)
|        [ ] vm3 (provider: virtualbox, state: running)
|
|  Commands:
|        up       <description>
|        halt     <description>
|        ssh      <description>
|
|  Vagrant Output:
|  [Output from Vagrant commands goes here]
|
|
|
|
|
|
|
|  Input:
|  [Input for Vagrant prompts goes here]
|
|
|  (q to quit)
|
+------------------------------------------------------------+