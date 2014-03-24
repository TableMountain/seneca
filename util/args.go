/*
Copyright 2014 Gavin Bong.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific
language governing permissions and limitations under the
License.
*/

package util

import (
    "errors"
    "flag"
    "fmt"
    "io/ioutil"
    "strings"
    "time"
)



type Arguments struct {
    Help     bool
    Version  bool
    VideoIn  string
    Port     int
    FromTime Instant

    RequestedScaling bool
    ScaleFilter string
}

func New() *Arguments {
    args := new(Arguments)
    return args
}

func (a *Arguments) Parse(arguments []string) error {
    f := flag.NewFlagSet("seneca", flag.ContinueOnError)
    f.SetOutput(ioutil.Discard)

    f.BoolVar(&a.Help, "h", false, "")
    f.BoolVar(&a.Version, "version", false, "")
    f.StringVar(&a.VideoIn, "video-infile", a.VideoIn, "")
    f.IntVar(&a.Port, "port", 8080, "")
    f.Var(&a.FromTime, "from", "")

    var scalingArg string
    f.StringVar(&scalingArg, "scale", "_:_", "")

    if err := f.Parse(arguments); err != nil {
        return err
    }

    preprocessScaleArg(a, scalingArg)

    var _ = a.validate()
    return nil
}

func (a *Arguments) validate() error {
    // TODO
    fmt.Printf("Port %d\n", a.Port)
    return nil
}

func preprocessScaleArg(a *Arguments, scalingArg string) error {
    if scalingArg != "_:_" {
        a.RequestedScaling = true
    }
    return nil
}


/////////////////////////////////////////////////////////////////


type Instant []time.Duration

func (i *Instant) String() string {
    return fmt.Sprint(*i)
}

func (i *Instant) Set(value string) error {
    if len(*i) > 0 {
        return errors.New("interval flag already set")
    }
    for _, dt := range strings.Split(value, ",") {
        duration, err := time.ParseDuration(dt)
        if err != nil {
            return err
        }
        *i = append(*i, duration)
    }
    return nil
}




/////////////////////////////////////////////////////////////////

type predicate func(uint16) bool

type ScaleType uint16

const (
    // -scale width:_
    //   height scaled to keep aspect ratio
    WidthOnly ScaleType = 1 << iota

    // -scale _:height
    //   Width scaled to keep aspect ratio
    HeightOnly

    // -scale width:height
    WidthHeight
)

var scales = map[ScaleType]interface{}{
    WidthOnly:   nil,
    HeightOnly:  nil,
    WidthHeight: nil,
}

// Converts into a valid argument to the -vf option of ffmpeg
func (s ScaleType) interpolate(width, height uint16) string {
    switch s {
    case WidthHeight:
        return fmt.Sprintf("scale=%d:%d", width, height)
    case HeightOnly:
        return fmt.Sprintf("scale=trunc(oh*a/2)*2:%d", height)
    case WidthOnly:
        return fmt.Sprintf("scale=%d:trunc(ow/a/2)*2", width)
    default:
        return ""
    }
}

func mkFn(max uint16) func(uint16) bool {
    return func(arg uint16) bool {
        return arg < max
    }
}

var even predicate = func(x uint16) bool { return x%2 == 0 }

var ErrNoArgs = errors.New("At least 1 arg must be provided")
var ErrTwoArgs = errors.New("2 args must be provided")
var ErrScaleOutOfRange = errors.New("ScaleType not recognized")

// A priori - arguments have been checked to not exceed the
//            respective dimensions of the input video 
func (s ScaleType) Decode(args ...uint16) (string, error) {
    var c = len(args)

    if _, ok := scales[s]; !ok {
        return "", ErrScaleOutOfRange
    }

    switch {
    default:
        return "", ErrNoArgs

    case c == 1:
        if !even(args[0]) {
            return "", fmt.Errorf("%d must be even", args[0])
        }
        switch s {
        case WidthOnly:
            return s.interpolate(args[0], 0), nil
        case HeightOnly:
            return s.interpolate(0, args[0]), nil
        case WidthHeight:
            fallthrough
        default:
            return "", ErrTwoArgs
        }

    case c >= 2:
        switch {
        case !even(args[0]):
            return "", fmt.Errorf("%d must be even", args[0])
        case !even(args[1]):
            return "", fmt.Errorf("%d must be even", args[1])
        default:
            return s.interpolate(args[0], args[1]), nil
        }
    }
}