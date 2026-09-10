Monitors a mint plant's moisture level.

Two components:

- [firmware](firmware/): Code running on the microcontroller.
- [relay](relay/): Relay code from microcontroller to Telegram bot (needed since this specific microcontroller doesn't implement TLS yet).
