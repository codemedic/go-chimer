package main

import "github.com/itchyny/volume-go"

type VolumeState struct {
	level int
	muted bool
	err   error
}

func GetVolumeState() *VolumeState {
	ret := &VolumeState{}
	ret.level, ret.err = volume.GetVolume()
	if ret.err == nil {
		ret.muted, ret.err = volume.GetMuted()
	}

	return ret
}

func (v *VolumeState) Restore() {
	if !v.Ok() {
		return
	}

	_ = volume.SetVolume(v.level)
	if v.muted {
		_ = volume.Mute()
	}
}

func (v *VolumeState) Ok() bool {
	return v.err == nil
}

func (v *VolumeState) SetVolume(l int) {
	if !v.Ok() {
		return
	}

	if l < 0 {
		return
	}

	if l > 100 {
		l = 100
	}

	if v.muted {
		if v.err = volume.Unmute(); v.err != nil {
			return
		}
	}

	if v.level < l {
		v.err = volume.SetVolume(l)
	}
}
