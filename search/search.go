package search

import (
  "fmt"
  
  "sluai/shikaku/shikakupuzzle"
)

type RectangleOptions map[int][]shikakupuzzle.Rectangle

func (ro *RectangleOptions) Print() {
  for regionId, options := range *ro {
    fmt.Printf("%d: ", regionId)
    for _, r := range options {
      fmt.Printf("%d by %d rectangle at (%d,%d), ", r.Dimensions[0], r.Dimensions[1], r.Location[0], r.Location[1])
    }
    fmt.Printf("\n")
  }
}

func Search(p *shikakupuzzle.ShikakuPuzzle) shikakupuzzle.ShikakuState {
  rootState, rootOptions := root(p)
  
  return backtrack(p, rootState, rootOptions)
}

func backtrack(p* shikakupuzzle.ShikakuPuzzle, state shikakupuzzle.ShikakuState, options RectangleOptions) shikakupuzzle.ShikakuState {  
  if len(options) == 0 {
    if p.IsSolved(state) {
      return state
    } else {
      return nil
    }
  }
  
  fmt.Printf("\nCalling backtrack\n")
  fmt.Printf("State:\n")
  p.Print(state)
  fmt.Printf("\nOptions:\n")
  options.Print()
  fmt.Printf("\n")
  
  regionId := selectRegion(p, state, options)
  fmt.Printf("Selected region %d\n", regionId)
  for _, rect := range options[regionId] {
    fmt.Printf("For region %d selecting %d by %d rectangle at (%d,%d)\n", regionId, rect.Dimensions[0], rect.Dimensions[1], rect.Location[0], rect.Location[1])
    newState := make(shikakupuzzle.ShikakuState)
    newOptions := make(RectangleOptions)
    
    for k,v := range state {
      newState[k] = v
    }
    for index, item := range options {
      if index != regionId {
        newOptions[index] = item
      }
    }    
    newState[regionId] = rect
    consistent := infer(p, newState, newOptions)
    if consistent {
      result := backtrack(p, newState, newOptions)
      if result != nil {
        return result
      }
    }
  }
  return nil
}

func root(p *shikakupuzzle.ShikakuPuzzle) (shikakupuzzle.ShikakuState, RectangleOptions) {
  state := make(shikakupuzzle.ShikakuState)
  options := make(RectangleOptions)
  
  // This can be significantly improved
  for regionId:=0; regionId<p.NumRegions; regionId++ {
    options[regionId] = make([]shikakupuzzle.Rectangle, 0)
    size := p.RegionSize[regionId]
    for width:=1; width<=size; width++ {
      if size % width == 0 {
        height := size / width
        for x:=0; x<p.Width - width + 1; x++ {
          for y:=0; y<p.Height - height + 1; y++ {
            rect := shikakupuzzle.Rectangle{shikakupuzzle.Coordinate{x,y}, shikakupuzzle.Coordinate{width, height}}
            if rect.Contains(p.RegionLocation[regionId]) {
              options[regionId] = append(options[regionId], rect)
            }
          }
        }
      }
    }
  }
  
  return state, options
}

func infer(p *shikakupuzzle.ShikakuPuzzle, s shikakupuzzle.ShikakuState, o RectangleOptions) bool {
  for id, r := range s {
    for i := 0 ; i < p.NumRegions ; i ++ {
      if o[i] != nil && i != id  {
        no := []shikakupuzzle.Rectangle{}
        for _, op := range o[i]{
          if shikakupuzzle.Overlap(r, op) == false {
            no = append(no, op)
          }
        }
        o[i] = no
        if len(no) == 0{
          return false
        }  
      }
    }
  }
return true
}


func selectRegion(p* shikakupuzzle.ShikakuPuzzle, state shikakupuzzle.ShikakuState, options RectangleOptions) int {
  // Return the first region that is unassigned
  for regionId:=0; regionId<p.NumRegions; regionId++ {
    _, assigned := options[regionId]
    if assigned {
      return regionId
    }
  }
  return -1
}
