package jump

import (
  "fmt"
	"log"
	"os"
)

var similarFile *os.File

func init() {
	similarFile, _ = os.OpenFile(basePath+"/similar.ai", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
}

func NewSimilar(ratio float64) *Similar {
	similar := &Similar{
		distances:    []float64{},
		ratios:       map[float64]float64{},
		defaultRatio: ratio,
	}
//  scanner := bufio.NewScanner(similarFile)
//	for scanner.Scan() {
//		line := strings.Split(scanner.Text(), ",")
//		if len(line) == 2 {
//			distance, err1 := strconv.ParseFloat(line[0], 64)
//			ratio, err2 := strconv.ParseFloat(line[1], 64)
//			if err1 == nil && err2 == nil {
//				similar.Add(distance, ratio)
//			}
//		}
//	}

	return similar
}

type Similar struct {
	distances    []float64
	ratios       map[float64]float64
	defaultRatio float64
}

func (s *Similar) Add(distance, ratio float64) {
	similarFile.Write([]byte(fmt.Sprintf("%v,%v\n", distance, ratio)))

	s.distances = append(s.distances, distance)
	s.ratios[distance] = ratio
}

func (s *Similar) Find(nowDistance float64) (es, di float64) {
  hDistance := 0.0
  hRatio := 0.0
  lDistance := 0.0
  lRatio := 0.0
  hDelta := 10000.0
  lDelta := -10000.0
  delta := 0.0
	for _, distance := range s.distances {
    delta = distance - nowDistance
    if delta >= 0 && delta < hDelta {
      hDelta = delta
      hDistance = distance
      hRatio = s.ratios[distance]
    }
    if delta <= 0 && delta > lDelta {
      lDelta = delta
      lDistance = distance
      lRatio = s.ratios[distance]
    }
	}

  log.Printf("hDistance %.2f hRatio %.2f hDelta %.2f", hDistance, hRatio, hDelta)
  log.Printf("lDistance %.2f lRatio %.2f lDelta %.2f", lDistance, lRatio, lDelta)

  if hDistance == 0 && lDistance != 0 {
    return lRatio, -lDelta
  }

  if lDistance == 0 && hDistance != 0 {
    return hRatio, hDelta
  }
  if hDistance == 0 && lDistance == 0 {
    return s.defaultRatio, 10000
  }
  est := ((-lDelta * hRatio)  + (hDelta * lRatio)) / (hDelta - lDelta)
  div := (hDelta - lDelta)/2
  return est, div
}
