package patrun

import (
  "sort"
  "fmt"
  "strings"
  "encoding/json"
  "regexp"
  "reflect"
)

type node struct {
  key string
  value map[string]node
  data interface{}
  modifier Modifiers
}

type Pattern struct {
    Match map[string]string
    Data interface{}
}

type Modifiers interface {
  Find(pm *Patrun, pat map[string]string, data interface{}) interface{}
  Remove(pm *Patrun, pat map[string]string, data interface{}) bool
}

type Customiser interface {
  Add(pm *Patrun, pat map[string]string, data interface{}) Modifiers
}

type Patrun struct {
  tree node
  Custom Customiser
}

func (p *Patrun) AddString(pat string, data interface{}) *Patrun {
  mapData := createMap(pat)

  return p.Add(mapData, data)
}

func (p *Patrun) Add(pat map[string]string, data interface{}) *Patrun {

    var custom Modifiers

    if p.Custom != nil {
      custom = p.Custom.Add(p, pat, data)
    }

    if p.tree.key == "" {
      p.tree = node{"root", map[string]node{}, nil, nil}
    }

    var keys []string
    for k := range pat {
      keys = append(keys, k)
    }
    sort.Strings(keys)


    var currentNode node = p.tree
    var lastNode node
    var key, val string
    var justCreated = false

    for k := range keys {
      key = keys[k]
      val = pat[key]

      lastNode = currentNode
      currentNode = currentNode.value[key]

      if currentNode.key == "" {
        lastNode.value[key] = node{key, map[string]node{}, nil, nil}
        currentNode = lastNode.value[key]
      }

      lastNode = currentNode
      currentNode = currentNode.value[val]
      if currentNode.key == "" {
        justCreated = true
        if k == len(keys) - 1 {
          lastNode.value[val] = node{val, map[string]node{}, data, custom}
        } else {
          lastNode.value[val] = node{val, map[string]node{}, nil, custom}
        }
        currentNode = lastNode.value[val]
      } else {
        justCreated = false
      }

      if k == len(keys) - 1 && !justCreated {
        var item = lastNode.value[val]
        item.data = data
        item.modifier = custom
        lastNode.value[val] = item
      }
    }

    if len(keys) == 0 {
      p.tree.data = data
    }

    return p
}

func (p *Patrun) FindExactString(pat string) interface{} {
  mapData := createMap(pat)

  return p.FindExact(mapData)
}

func (p *Patrun)FindExact(pat map[string]string) interface{} {
  return p.findItem(pat, true)
}

func (p *Patrun) FindString(pat string) interface{} {
  mapData := createMap(pat)

  return p.Find(mapData)
}

func (p *Patrun) Find(pat map[string]string) interface{} {
  return p.findItem(pat, false)
}

func (p *Patrun) findItem(pat map[string]string, exact bool) interface{} {
  var keys []string
  for k := range pat {
    keys = append(keys, k)
  }
  sort.Strings(keys)

  var currentNode = p.tree
  var lastGoodNode = currentNode
  var foundKeys []string
  var lastData interface{} = p.tree.data
  var lastModifier Modifiers = p.tree.modifier
  var stars []node
  var keyPointer = 0

  for keyPointer < len(keys) {
    var key = keys[keyPointer]
    var val = pat[key]


    currentNode = currentNode.value[key].value[val]

    if currentNode.key != "" {
      if len(lastGoodNode.value) > 0 {
        stars = append(stars, lastGoodNode)
      }

      lastGoodNode = currentNode
      foundKeys = append(foundKeys, key)
      if lastGoodNode.data != nil {
        lastData = lastGoodNode.data
      }
      lastModifier = lastGoodNode.modifier

      keyPointer++

    } else if lastData == nil && len(stars) > 0 {
        currentNode = stars[len(stars) - 1]
    
        stars = stars[:len(stars)-1]
        lastGoodNode = currentNode

    } else {

      currentNode = lastGoodNode
      keyPointer++
    }

  }

  if exact && len(foundKeys) != len(keys) {
    lastData = nil
  }

  if lastModifier != nil {
    lastData = lastModifier.Find(p, pat, lastData)
  }


  return lastData


}

func (p *Patrun) RemoveString(pat string)  {
  mapData := createMap(pat)

  p.Remove(mapData)
}

func (p *Patrun) Remove(pat map[string]string) {
  var keys []string
  for k := range pat {
    keys = append(keys, k)
  }
  sort.Strings(keys)

  var currentNode = p.tree
  var lastGoodNode = currentNode
  var foundKeys []string
  var key, val string
  var lastParent = currentNode

  for k := range keys {
    key = keys[k]
    val = pat[key]

    currentNode = currentNode.value[key]

    if currentNode.key != "" {
      lastParent = currentNode
      lastGoodNode = currentNode
    }

    currentNode = currentNode.value[val]

    if currentNode.key != "" {
      lastGoodNode = currentNode
      foundKeys = append(foundKeys, key)
    }
  }

  //found a match so delete the data element
  if len(foundKeys) == len(keys) {
    var item = lastParent.value[val]
    var okToDel = true

    if lastGoodNode.modifier != nil {
      okToDel = lastGoodNode.modifier.Remove(p, pat, item.data)
    }
    if okToDel {
      item.data = nil
      lastParent.value[val] = item


    }
  }

}

func (p *Patrun) ListString(pat string, exact bool) []Pattern {

  mapData := createMap(pat)

  return p.List(mapData, exact)
}

func (p *Patrun)List(pat map[string]string, exact bool) []Pattern {
  var items []Pattern
  var keyMap []string

  if pat == nil {
    pat = map[string]string{}
  }

  if p.tree.data != nil {
    items = append(items, createMatchList(keyMap, p.tree.data))
  }

  if p.tree.key != "" {
    descendTree(&items, pat, exact, true, p.tree.value, keyMap)
  }
  return items
}


func (p *Patrun)ToJSON() []byte {
  b, _ := json.Marshal(p.List(nil, false))


  return b
}



func (p Patrun)String() string {

  items := p.List(nil, false)

  var data []string

  for k := range items {
    v := items[k]
    data = append(data, fmt.Sprintf("%v -> <%v>", formatMatch(v.Match), formatData(v.Data)))
  }

  return strings.Join(data, "\n")
}

func (p Patrun)ToString(custom func(data interface{}) string) string {

  items := p.List(nil, false)

  var data []string

  for k := range items {
    v := items[k]
    data = append(data, fmt.Sprintf("%v -> %v", formatMatch(v.Match), custom(v.Data)))
  }

  return strings.Join(data, "\n")
}

func formatMatch(items map[string]string) string {
  var points []string

  var keys []string
  for k := range items {
    keys = append(keys, k)
  }
  sort.Strings(keys)

  for k := range keys {
    var key = keys[k]
    var val = items[key]

    points = append(points, fmt.Sprintf("%v:%v", key, val))
  }

  return strings.Join(points, ", ")
}

func descendTree(items *[]Pattern, pat map[string]string, exact bool, rootLevel bool, values map[string]node, keyMap []string) {

  var localKeyMap []string

  copy(localKeyMap, keyMap)

  var keys []string
  for k := range values {
    keys = append(keys, k)
  }
  sort.Strings(keys)

  for k := range keys {
    var key = keys[k]
    var val = values[key]

      if rootLevel {
        keyMap = []string{}
      }

      if val.data == nil && len(val.value) > 0 {
      //  keyMap = append(keyMap, key)

        descendTree(items, pat, exact, false, val.value, append(keyMap, key))

      } else if val.data != nil {
        localKeyMap = append(keyMap, val.key)
        if validatePatterMatch(pat, exact, localKeyMap) {
          *items = append(*items, createMatchList(localKeyMap, val.data))
        }

        if len(val.value) > 0 {
          descendTree(items, pat, exact, false, val.value, append(keyMap, val.key))
        }
      }
  }
}

func validatePatterMatch(pat map[string]string, exact bool, matchedKeys []string) bool {
  var keys []string
  for k := range pat {
    keys = append(keys, k)
  }
  sort.Strings(keys)


  if len(keys) == 0 {
    return true
  }

  var foundKeys []string
  var key, val string

  for k := range keys {
    key = keys[k]
    val = pat[key]

    foundKeys = append(foundKeys, key)
    foundKeys = append(foundKeys, val)
  }


  var matched = true
  for i := 0; i < len(foundKeys); i++ {
    //regex here


    if i >= len(matchedKeys) {
      matched = false
    }
    if matched && i % 2 == 0 && foundKeys[i] != matchedKeys[i] {
      matched = false
    }
    if matched && i % 2 != 0 && !gexval(foundKeys[i], matchedKeys[i]) {
    //if foundKeys[i] != matchedKeys[i] && founKeys[i] != "*"{
      matched = false
    }

    if !matched {
      break
    }
  }

  if exact && len(foundKeys) != len(matchedKeys) {
    matched = false
  }

  return matched
}

func createMatchList(keyMap []string, dataItem interface{}) Pattern {

  var keys map[string]string = map[string]string{}
  var item Pattern = Pattern{}

  for i := 0; i < len(keyMap); i+=2 {
      if i + 1 < len(keyMap) {
        keys[keyMap[i]] = keyMap[i+1]
      }
  }

  item.Match = keys
  item.Data = dataItem

  return item
}


func formatData(data interface{}) string {

  if reflect.ValueOf(data).Kind().String() == "func" {
    return "function"
  } else {
    return fmt.Sprintf("%v", data)
  }

}

func createMap(pat string) map[string]string {
    mapData := map[string]string{}

    items := strings.Split(pat, ",")

    for k := range items {
      item := strings.TrimSpace(items[k])

      if len(item) > 0 {
        parts := strings.Split(item, ":")
        if len(parts) == 2 {
          mapData[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
        }
      }

    }

    return mapData
}

func gexval(pattern string, value string) bool {

  pattern = escregexp(pattern)

  // use [\s\S] instead of . to match newlines

  r := regexp.MustCompile(`\\\*`)

  pattern = r.ReplaceAllString(pattern, "[\\s\\S]*")

  r = regexp.MustCompile(`\\\?`)
  pattern = r.ReplaceAllString(pattern, "[\\s\\S]")

  // escapes ** and *?
  r = regexp.MustCompile(`\[\\s\\S\]\*\[\\s\\S\]\*`)
  pattern = r.ReplaceAllString(pattern, `\\\*`)

  r = regexp.MustCompile(`\[\\s\\S\]\*\[\\s\\S\]`)
  pattern = r.ReplaceAllString(pattern, `\\\?`)

  pattern = fmt.Sprintf("^%v$", pattern)

  r = regexp.MustCompile(pattern)

  return r.MatchString(value)

}

func escregexp(restr string) string {

  r := regexp.MustCompile(`([-\[\]{}()*+?.,\\^$|#\s])`)

  return r.ReplaceAllString(restr, "\\$1")

}
