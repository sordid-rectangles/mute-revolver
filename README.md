# discord-roulette
A discord russian roulette bot.

## Commands

### /load
Loads the revolver

### /spin 
Randomly spins the 6 chambers

### /shoot
Attempts to fire the revolver. If it is loaded and chambered it will ban the user who fired it for 3 days.

## Test

## Notes
This is still janky, and as such revolver state is managed in memory in a single instance across alllllll command use. Will be shifting this to 1 instance per guild in an embedded db or something. Research required.
