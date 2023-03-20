package shikakupuzzle

import (
  "fmt"
)

type Coordinate [2]int

type Rectangle struct {
  Location Coordinate
  Dimensions [2]int
}

type ShikakuPuzzle struct {
  Width int
  Height int
  NumRegions int
  RegionSize []int
  RegionLocation []Coordinate
}

type ShikakuState map[int]Rectangle

func (rect* Rectangle) Contains(c Coordinate) bool {
  if rect.Location[0] <= c[0] && c[0] < rect.Location[0] + rect.Dimensions[0] {
    if rect.Location[1] <= c[1] && c[1] < rect.Location[1] + rect.Dimensions[1] {
      return true
    }
  }
  return false
}

func (state* ShikakuState) FindAssignment(c Coordinate) int {
  for regionId, rect := range *state {
    if rect.Contains(c) {
      return regionId
    }    
  }
  return -1
}

func (p* ShikakuPuzzle) Print(state ShikakuState) {
  fmt.Printf("+")
  for y:= 0; y<p.Height; y++ {
    fmt.Printf("--+")
  }
  fmt.Printf("\n")
  
  for y:=p.Height-1; y>=0; y-- {
    fmt.Printf("|")
    for x:=0; x<p.Width; x++ {
      region := state.FindAssignment(Coordinate{x,y})
      if region > -1 {
        fmt.Printf("%2d|", region)
      } else {
        fmt.Printf("  |")
      }
    }
    fmt.Printf("\n")
    fmt.Printf("+")
    for y:= 0; y<p.Width; y++ {
      fmt.Printf("--+")
    }
    fmt.Printf("\n")
  }
  fmt.Printf("\nRegion info\n")
  for id:=0; id<p.NumRegions; id++ {
    fmt.Printf("%d\t%d\t%v\n", id, p.RegionSize[id], p.RegionLocation[id])
  }
}

func min(x, y int) int {
  if x<y {
    return x
  } else {
    return y
  }
}

func max(x, y int) int {
  if x>y {
    return x
  } else {
    return y
  }
}

func Overlap(r1, r2 Rectangle) bool {
  if max(r1.Location[0], r2.Location[0]) >= min(r1.Location[0]+r1.Dimensions[0], r2.Location[0]+r2.Dimensions[0]) {
    return false
  }
  if max(r1.Location[1], r2.Location[1]) >= min(r1.Location[1]+r1.Dimensions[1], r2.Location[1]+r2.Dimensions[1]) {
    return false
  }
  return true
}

func (p* ShikakuPuzzle) IsSolved(state ShikakuState) bool {
  // Is everything assigned
  if len(state) < p.NumRegions {
    return false
  }
  
  // Does each region contain its starting point
  for regionId, loc := range p.RegionLocation {
    r := state.FindAssignment(loc)
    if r != regionId {
      return false
    }
  }
  
  // Is each region the correct size
  for regionId, size := range p.RegionSize {
    rect := state[regionId]
    if rect.Dimensions[0] * rect.Dimensions[1] != size {
      return false
    }
  }
  
  // Regions do not overlap
  for regionId:=0; regionId<p.NumRegions; regionId++ {
    for regionId2:=regionId+1; regionId2<p.NumRegions; regionId2++ {
      if Overlap(state[regionId], state[regionId2]) {
        return false
      }
    }
  }
  
  return true
}
