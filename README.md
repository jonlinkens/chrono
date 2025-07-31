# Chrono

A benchmarking TUI for measuring execution times and performance of shell commands.

![Demo of usage](https://github.com/user-attachments/assets/29d36268-f335-46df-9e7b-2afad5af5972)

```bash
brew install jonlinkens/tap/chrono
```

---

- Rich TUI experience and CLI for scripting usage
- Configurable number of runs with statistical analysis (mean, min, max, range)
- Optional warmup iterations before benchmarking
- Live output stream of stdout and stderr with scrollback buffer
- Shell startup calibration
- Phrase detection
- Timeout support
- Cross platform

## Usage

```bash
# TUI mode (default)
chrono echo "hello world"

# Simple CLI mode
chrono --cli echo "hello world"
```

### Options

```bash
chrono [OPTIONS] COMMAND [ARGS...]

  --runs N               Number of benchmark runs (default: 1)
  --warmups N            Number of warmup runs before benchmarking (default: 0)
  --phrase "text"        Stop timing when this phrase appears in output
  --timeout DURATION     Maximum time to wait (e.g., 5s, 1m30s)
  --calibration N        Number of shell overhead calibration runs (default: 5)
  --skip-calibration     Skip shell overhead calibration
  --cli                  Use CLI output instead of TUI
  --command "cmd args"   Command as quoted string (alternative to positional args)
  --version              Print version and exit
```
