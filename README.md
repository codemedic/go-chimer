# Chimer

A chiming clock for the command line, written in Go.


## Installation

You can install `chimer` either from the latest release archive or be compiled locally.

### From Release Archive

1. Download the latest archive, specific to your platform, from https://github.com/codemedic/go-chimer/releases/latest
2. Verify SHA256
3. Extract the archive; e.g. into `~/apps/chimer`
4. That's it!

### Compile

1. Clone the repo
   ```
   git clone https://github.com/codemedic/go-chimer.git
   ```
2. Build the binary
   ```shell
   # On Linux, dependencies for oto library need to be installed
   sudo apt install -y libasound2-dev build-essential

   cd go-chimer
   CGO_ENABLED=1 go build ./cmd/chimer/ 
   ```
3. Copy binary and sounds
   ```shell
   mkdir -p ~/apps/chimer
   cp -r chimer chimes ~/apps/chimer
   ```

## Usage

The `chimer` can be used either as a one-shot command run from a cron, every 15 minutes
or as a continuously running service, from your service manager of choice. You may also
run the binary directly from a terminal window.

### Cron Usage

You need to add the below line to your `crontab`.

```
*/15 * * * * /home/username/apps/chimer/chimer --sound /home/username/apps/chimer/chimes` --cron
```

### From a Terminal Window

Run the command below from a terminal window; make sure you leave the window open afterwards. You may press `Ctrl+C` to stop `chimer`.

```shell
~/apps/chimer/chimer --sound ~/apps/chimer/chimes
```

Alternatively, you may add the below line to your `.bashrc` or your shell's rc file to simplify the invocation.

```shell
export CHIMER_SOUND_DEFAULT="$HOME/apps/chimer/chimes"
```

Then `chimer` can be invoked with no commandline option.

```shell
~/apps/chimer/chimer
```

## Usage Help

You can see the below usage help by using the command line option `--help`.

    $> ~/apps/chimer/chimer --help
    Usage: chimer [options] [arguments]

    OPTIONS
    --sound/$CHIMER_SOUND_DEFAULT  <string>
    Default path for sound(s) if none are not specified.
    --sound-hour/$CHIMER_SOUND_HOUR  <string>
    Path for sound(s) at zero past an hour.
    --sound-quarter-past/$CHIMER_SOUND_QUARTER_PAST  <string>
    Path for sound(s) at quarter past an hour.
    --sound-half-past/$CHIMER_SOUND_HALF_PAST  <string>
    Path for sound(s) at half past an hour.
    --sound-quarter-to/$CHIMER_SOUND_QUARTER_TO  <string>
    Path for sound(s) at quarter to an hour.
    --cache-sounds/$CHIMER_CACHE_SOUNDS  <bool>  (default: true)
    Cache sounds in memory
    --repeat-hourly-sound/-r/$CHIMER_REPEAT_HOURLY_SOUND  <bool>  (default: false)
    Repeat the hourly sound T time at T'O clock.
    --minimum-volume-level/-l/$CHIMER_MINIMUM_VOLUME_LEVEL  <int>    (default: 50)
    --test-time/$CHIMER_TEST_TIME                           <value>
    Specify a time to test chimer in cron-mode. This option is ignored when cron-mode is not enabled.
    --cron/-c/$CHIMER_CRON_MODE_ENABLED  <bool>  (default: false)
    Enable cron-mode where it acts on the current time to decide whether to chime or not
    --help/-h
    display this help message
    --version/-v
    display version information

    NOTE: All sound paths can be either a relative or full path to an MP3 file or
    a directory containing multiple MP3s. In the latter case, one of the files will
    be chosen at random for each chime.
