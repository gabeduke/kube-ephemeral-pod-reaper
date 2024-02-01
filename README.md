
# kube-ephemeral-pod-reaper

`kube-ephemeral-pod-reaper` is a Kubernetes tool designed to manage the lifecycle of ephemeral containers. It comprises two primary components: `Scout` and `Reaper`. `Scout` is responsible for monitoring ephemeral containers and marking them for deletion, whereas `Reaper` (still in development) will handle the actual deletion process.

## Features

- **Scout**: Watches for ephemeral containers and annotates them with an expiration time.
- **Reaper**: (TODO) Responsible for deleting the marked ephemeral containers.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Installing

_Building from Source_
```bash
git clone https://github.com/gabeduke/kube-ephemeral-pod-reaper.git
cd kube-ephemeral-pod-reaper
go build
```

_Go Install_
```bash
go install github.com/gabeduke/kube-ephemeral-pod-reaper
```

### Usage

#### Scout

To start the Scout controller:

```bash
./kube-ephemeral-pod-reaper scout --selector="your-label-selector"
```

#### Reaper (TODO)

The Reaper controller is still under development. Once completed, it will be used as follows:

```bash
./kube-ephemeral-pod-reaper reaper [options]
```

#### Using Controllers as CLI Commands in Your Own Tools

To integrate `Scout` or `Reaper` into your own CLI tools, import the respective package and add it as a subcommand:

```go
import (
    "github.com/gabeduke/kube-ephemeral-pod-reaper/pkg/scout"
    "github.com/gabeduke/kube-ephemeral-pod-reaper/pkg/reaper"
)

// In your Cobra command setup
var rootCmd = &cobra.Command{ /* ... */ }

func init() {
    rootCmd.AddCommand(scout.NewScoutCmd())
    rootCmd.AddCommand(reaper.NewReaperCmd())
}
```

#### Using the Controller as a Package

To use the controller directly in your code:

```go
import "github.com/gabeduke/kube-ephemeral-pod-reaper/pkg/scout"

cfg := scout.Config{
    Annotation: annotations,
    Duration:   duration,
    Name:       name,
    Selector:   labelSelector,
}

controller := scout.Controller{Cfg: cfg}
controller.Run()
```
