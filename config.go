package chimer

import (
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
)

type Sounds []string

func (c Sounds) GetOne() string {
	if c == nil {
		panic("nil sound collection")
	}

	switch len(c) {
	case 0:
		panic("empty sound collection")
	case 1:
		return c[0]
	default:
		return c[rand.Intn(len(c))]
	}
}

type SoundGetter interface {
	Get(path string) (*beep.Buffer, error)
}

type SoundCache struct {
	cfg         *Config
	chimeSounds *Cache[Chime, Sounds]
	sounds      SoundGetter
}

func NewSoundCache(cfg *Config) *SoundCache {
	rand.Seed(time.Now().UnixNano())

	ret := &SoundCache{
		cfg:         cfg,
		chimeSounds: NewCache[Chime, Sounds](cfg.GetSounds),
	}

	if cfg.CacheSounds {
		ret.sounds = NewCache[string, *beep.Buffer](LoadSound)
	} else {
		ret.sounds = SoundLoader{}
	}

	return ret
}

func (c *SoundCache) Get(w Chime) (*beep.Buffer, error) {
	cc, err := c.chimeSounds.Get(w)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup collection: %w", err)
	}

	if cc == nil {
		cc, err = c.chimeSounds.Get(None) // get default sound(s)
		if err != nil {
			return nil, fmt.Errorf("failed to lookup falback collection: %w", err)
		}

		if cc == nil {
			return nil, fmt.Errorf("%v: %w", w, ErrNoSoundConfigured)
		}
	}

	return c.sounds.Get(cc.GetOne())
}

type Config struct {
	Sound struct {
		Default     string `conf:"flag:sound,help:Default path for sound(s) if none are not specified."`
		Hour        string `conf:"help:Path for sound(s) at zero past an hour."`
		QuarterPast string `conf:"help:Path for sound(s) at quarter past an hour."`
		HalfPast    string `conf:"help:Path for sound(s) at half past an hour."`
		QuarterTo   string `conf:"help:Path for sound(s) at quarter to an hour."`
	}
	CacheSounds        bool `conf:"default:true,help:Cache sounds in memory"`
	RepeatHourlySound  bool `conf:"short:r"`
	MinimumVolumeLevel int  `conf:"default:50,short:l"`
}

var ErrNoSoundConfigured = errors.New("no sound configured")
var ErrUnusableSound = errors.New("unusable sound")

func (c *Config) GetSoundPath(w Chime) string {
	switch {
	case w == Hour:
		return c.Sound.Hour
	case w == QuarterPast:
		return c.Sound.QuarterPast
	case w == HalfPast:
		return c.Sound.HalfPast
	case w == QuarterTo:
		return c.Sound.QuarterTo
	default:
		return c.Sound.Default
	}
}

func errUnusableSound(w Chime, path string, msg string) error {
	return fmt.Errorf("%w; %s : chime:%v path:%v", ErrUnusableSound, msg, w, path)
}

func (c *Config) GetSounds(w Chime) (Sounds, error) {
	path := c.GetSoundPath(w)
	if path == "" {
		return nil, nil
	}

	fi, err := os.Stat(path)
	if err != nil {
		return nil, errUnusableSound(w, path, "failed to stat path")
	}

	if fi.Mode().IsRegular() {
		return Sounds{path}, nil
	}

	if !fi.Mode().IsDir() {
		return nil, ErrUnusableSound
	}

	fsys := os.DirFS(path)

	mp3s, err := fs.Glob(fsys, "*.[mM][pP]3")
	if err != nil {
		return nil, errUnusableSound(w, path, "failed to find sounds")
	}

	ret := mp3s[:0]
	for _, mp3 := range mp3s {
		if stat, err := fs.Stat(fsys, mp3); err == nil && stat.Mode().IsRegular() {
			ret = append(ret, filepath.Join(path, mp3))
		}
	}

	if len(ret) == 0 {
		return nil, errUnusableSound(w, path, "no sounds in path")
	}

	return ret, nil
}
