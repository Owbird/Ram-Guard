# Ram Guard

**Ram Guard** is a lightweight, utility tool for managing memory usage on your system. It helps prevent running out of memory by listening for user-defined intervals and thresholds; if memory usage exceeds the threshold, Ram Guard finds the most resource-hungry process and stops it in its tracks.

## 🛠️ Installation

1. **Clone the repository**:
    ```bash
    git clone https://github.com/owbird/ram-guard.git
    cd ram-guard
    ```

2. **Install dependencies**:
    ```bash
    go mod tidy
    ```

3. **Run Ram Guard**:
    ```bash
    go run .
    ```

## ⚙️ Usage

Configure `interval` and `threshold` settings in the configuration file or pass them as command-line arguments:

```bash
ram-guard --interval <time_in_seconds> --threshold <RAM_limit_in_%>
```

### Example

To set Ram Guard to check every 10 seconds and act if RAM usage exceeds 80%:
```bash
ram-guard --interval 10 --threshold 80
```

## 🔍 How It Works

1. **Interval Listener**: Ram Guard listens for the specified interval.
2. **Threshold Check**: Compares current RAM usage with the defined threshold.
3. **Process Termination**: If usage exceeds the threshold, Ram Guard identifies the largest process by memory usage and terminates it while notifying you.

## 🛡️ Safety Notice

Use Ram Guard cautiously, as it automagically terminates high-memory processes. Avoid setting the threshold too low, which could inadvertently stop essential services.

## 🤝 Contributing

Contributions are welcome! If you’d like to contribute, please open an issue or submit a pull request with suggested changes.
