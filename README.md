# CheckPower

A simple scheduled power outage checker hastily created for give you friendly remainder about power outages

## Usage

```sh
./checkPower <accountNumber> <timeAhead?>
```

- `accountNumber` - is 8 digits of your account
- `timeAhead` - is optional time to add to current time for period check (default `15m`) [`1h15m40s` - format]

## Cron

This cron job will start every hour at 46 minutes and check for power outages 15 minutes before. When it is scheduled some - it will trigger `say` and `notification` standard MacOS utilities to warn you.

```cron
46 * * * * /full/path/to/checkPower 00000000
```
