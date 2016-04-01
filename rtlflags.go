// Implements command-line flag control over settings for rtl-sdr's.
package rtlflags

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/bemasher/rtltcp/si"
)

// Consists of methods common to rtl-sdr packages.
type Radio interface {
	SetAgcMode(bool) error
	SetCenterFreq(int) error
	SetDirectSampling(int) error
	SetFreqCorrection(int) error
	SetOffsetTuning(bool) error
	SetSampleRate(int) error
	SetTestMode(bool) error
	SetTunerBw(int) error
	SetTunerGain(int) error
	SetTunerGainMode(bool) error
}

// Embeds a Radio and applies flag values by calling appropriate methods.
type Context struct {
	Radio

	agcMode        bool
	centerFreq     si.ScientificNotation
	directSampling SamplingMode
	freqCorrection int
	offsetTuning   bool
	sampleRate     si.ScientificNotation
	testMode       bool
	tunerBandwidth si.ScientificNotation
	tunerGain      float64
	tunerGainMode  bool
}

// Register flags using the default flagset.
func (c *Context) RegisterFlags() {
	flag.BoolVar(&c.agcMode, "agcmode", false, "enable rtl2832u agc")
	flag.Var(&c.centerFreq, "centerfreq", "center frequency to receive on")
	flag.Lookup("centerfreq").DefValue = "100M"
	flag.Var(&c.directSampling, "directsampling", "set sampling mode: none, inphase, quadrature")
	flag.Lookup("directsampling").DefValue = "none"
	flag.IntVar(&c.freqCorrection, "freqcorrection", 0, "frequency correction in ppm")
	flag.BoolVar(&c.offsetTuning, "offsettuning", false, "enable offset tuning")
	flag.Var(&c.sampleRate, "samplerate", "sample rate")
	flag.Lookup("samplerate").DefValue = "2.4M"
	flag.BoolVar(&c.testMode, "testmode", false, "enable test mode")
	flag.Var(&c.tunerBandwidth, "tunerbandwidth", "tuner bandwidth")
	flag.Lookup("tunerbandwidth").DefValue = "2.4M"
	flag.Float64Var(&c.tunerGain, "tunergain", 0.0, "set tuner gain in dB")
	flag.BoolVar(&c.tunerGainMode, "tunergainmode", false, "eanble manual gain")
}

// Applies settings given by flags. Must be called after flag.Parse()
func (c *Context) HandleFlags() {
	flag.Visit(func(f *flag.Flag) {
		var err error

		switch f.Name {
		case "agcmode":
			err = c.SetAgcMode(c.agcMode)
		case "centerfreq":
			err = c.SetCenterFreq(int(c.centerFreq))
		case "directsampling":
			err = c.SetDirectSampling(int(c.directSampling))
		case "freqcorrection":
			err = c.SetFreqCorrection(c.freqCorrection)
		case "offsettuning":
			err = c.SetOffsetTuning(c.offsetTuning)
		case "samplerate":
			err = c.SetSampleRate(int(c.sampleRate))
		case "testmode":
			err = c.SetTestMode(c.testMode)
		case "tunerbandwidth":
			err = c.SetTunerBw(int(c.tunerBandwidth))
		case "tunergain":
			err = c.SetTunerGain(int(c.tunerGain * 10))
		case "tunergainmode":
			err = c.SetTunerGainMode(c.tunerGainMode)
		}

		if err != nil {
			log.Fatal(err)
		}
	})
}

// Implements flag.Value for sampling mode.
type SamplingMode int

const (
	SamplingNone SamplingMode = iota
	SamplingIADC
	SamplingQADC
	SamplingUnknown
)

func (s SamplingMode) String() string {
	switch s {
	case SamplingNone:
		return "None"
	case SamplingIADC:
		return "In-Phase ADC"
	case SamplingQADC:
		return "Quadrature ADC"
	}
	return "Unknown"
}

func (s *SamplingMode) Set(val string) (err error) {
	switch strings.ToLower(val) {
	case "none":
		*s = SamplingNone
	case "inphase":
		*s = SamplingIADC
	case "quadrature":
		*s = SamplingQADC
	}
	return fmt.Errorf("invalid sampling mode: %q", val)
}
