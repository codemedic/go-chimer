package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codemedic/chimer"

	"github.com/ardanlabs/conf/v3"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

type ChimeSoundGetter interface {
	Get(w chimer.Chime) (*beep.Buffer, error)
}

// build is the version of this program; set via Makefile.
var build = "develop"

func logErrorf(format string, a ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf(format, a...))
}

func chimeNow(now time.Time, tolerance time.Duration, sounds ChimeSoundGetter, cfg chimer.Config) error {
	hour, chime := chimer.GetChime(now, tolerance)
	if chime == chimer.None {
		return nil
	}

	vol := GetVolumeState()

	bufferedChime, err := sounds.Get(chime)
	if err != nil {
		return err
	}

	if err = speaker.Init(bufferedChime.Format().SampleRate, bufferedChime.Format().SampleRate.N(time.Second/10)); err != nil {
		return err
	}

	var s beep.Streamer
	s = bufferedChime.Streamer(0, bufferedChime.Len())
	if chime == chimer.Hour && cfg.RepeatHourlySound {
		s = beep.Loop(hour, s.(beep.StreamSeeker))
	}

	vol.SetVolume(cfg.MinimumVolumeLevel)

	done := make(chan struct{})
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(done)
	})))

	<-done

	vol.Restore()

	return nil
}

func main() {
	cfg := struct {
		conf.Version
		chimer.Config
		TestTime        time.Time `conf:"help:Specify a time to test chimer in cron-mode. This option is ignored when cron-mode is not enabled."`
		CronModeEnabled bool      `conf:"help:Enable cron-mode where it acts on the current time to decide whether to chime or not, and quits."`
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Chimer Clock",
		},
	}

	const prefix = "CHIMER"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			note := `
NOTE: All sound paths can be either a relative or full path to an MP3 file or
a directory containing multiple MP3s. In the latter case, one of the files will
be chosen at random for each chime.`
			fmt.Println(help, note)
			os.Exit(0)
		}

		logErrorf("Error parsing config: %v\n", err)
		os.Exit(1)
	}

	cache := chimer.NewSoundCache(&cfg.Config)
	scheduler := chimer.NewScheduler(20 * time.Second)

	signals := make(chan os.Signal, 1)
	defer close(signals)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM) // Quit on SIGINT or SIGTERM
	go func() {
		select {
		case <-signals:
			scheduler.Stop()
		}
	}()

	if cfg.CronModeEnabled {
		var now time.Time
		if cfg.TestTime.IsZero() {
			now = time.Now()
		} else {
			now = cfg.TestTime
		}
		err = chimeNow(now, 2*time.Second, cache, cfg.Config)
		if err != nil {
			logErrorf("failed to chime: %v", err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	scheduler.EveryQuarterHour(func(t time.Time, tolerance time.Duration) {
		err := chimeNow(t, tolerance, cache, cfg.Config)
		if err != nil {
			logErrorf("failed to chime: %v", err)
		}
	})

	speaker.Close()
}
