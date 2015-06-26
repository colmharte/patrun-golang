
// Package patrun is a fast pattern matcher on Go map properties
//
// For a full guide visit https://github.com/colmharte/patrun-golang
//
//Need to pick out an object based on a subset of its properties? Say you've got:
//
//{ x: 1       } -> A
//{ x: 1, y: 1 } -> B
//{ x: 1, y: 2 } -> C
//
//Then patrun can give you the following results:
//
//{ x: 1 }      -> A
//{ x: 2 }      -> no match
//{ x: 1, y: 1 } -> B
//{ x: 1, y: 2 } -> C
//{ x: 2, y: 2 } -> no match
//{ y: 1 }      -> no match
//
//It's basically _query-by-example_ for property sets.
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

//Returned by the List method to idenfity the Match and Data stored for each pattern
type Pattern struct {
    Match map[string]string
    Data interface{}
    Modifier Modifiers
}

//Modifiers allow you to customise the results for the Find and Remove methods
type Modifiers interface {
  Find(pm *Patrun, pat map[string]string, data interface{}) interface{}
  Remove(pm *Patrun, pat map[string]string, data interface{}) bool
}

//Customisers allow custom logic to be added when processing patterns
type Customiser interface {
  Add(pm *Patrun, pat map[string]string, data interface{}) Modifiers
}

//Patrun is the main object, specify Custom when creating to allow custom logic to be applied when manipulating patterns
type Patrun struct {
  tree node
  Custom Customiser
}


//Register a pattern, and the object that will be returned if an input matches.
func (p *Patrun) Add(pat map[string]string, data interface{}) *Patrun {

    var custom Modifiers

    if p.Custom != nil {
      custom = p.Custom.Add(p, pat, data)
    }

    if p.tree.key == "" {
      p.tree = node{"root", map[string]node{}, nil, nil}
    }

    var keys = sortKeys(pat)


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
      p.tree.modifier = custom
    }

    return p
}

//Same as Add but using simple string notation instead of a map type. eg "a:1,b:2" is the equivalent to map[string]string{"a":"1","b":"2"}
func (p *Patrun) AddString(pat string, data interface{}) *Patrun {
  mapData := createMap(pat)

  return p.Add(mapData, data)
}

//Return the unique match for this subject, or nil if not found. The
//properties of the subject are matched against the patterns previously
//added, and the most specifc pattern wins. Unknown properties in the
//subject are ignored.
func (p *Patrun) Find(pat map[string]string) interface{} {
  return p.findItem(pat, false)
}

//Same as Find but using simple string notation
func (p *Patrun) FindString(pat string) interface{} {
  mapData := createMap(pat)

  return p.Find(mapData)
}

//Same as Find but only matches where all properties match will be returned.
func (p *Patrun)FindExact(pat map[string]string) interface{} {
  return p.findItem(pat, true)
}

//Same as FindExact but using simple string notation
func (p *Patrun) FindExactString(pat string) interface{} {
  mapData := createMap(pat)

  return p.FindExact(mapData)
}

func (p *Patrun) findItem(pat map[string]string, exact bool) interface{} {
  var keys = sortKeys(pat)

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

//Remove this pattern, and it's object, from the matcher.
func (p *Patrun) Remove(pat map[string]string) {
  var keys = sortKeys(pat)

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
    var item node

    if len(pat) == 0 {
      fmt.Println("afsdf")
      item = p.tree
    } else {
      item = lastParent.value[val]
    }

    var okToDel = true

    if lastGoodNode.modifier != nil {
      okToDel = lastGoodNode.modifier.Remove(p, pat, item.data)
    }
    if okToDel {
      item.data = nil
      item.modifier = nil
      if len(pat) == 0 {
        p.tree = node{p.tree.key, p.tree.value, nil, nil}
      } else {
        lastParent.value[val] = item
      }
    }
  }

}

//Same as Remove but using simple string notation
func (p *Patrun) RemoveString(pat string)  {
  mapData := createMap(pat)

  p.Remove(mapData)
}

//Return the list of registered patterns that contain this partial
//pattern. You can use wildcards for property values.  Omitted values
//are *not* equivalent to a wildcard of _"*"_, you must specify each
//property explicitly. You can provide a second boolean
//parameter, _exact_. If true, then only those patterns matching the
//pattern-partial exactly are returned.
func (p *Patrun)List(pat map[string]string, exact bool) []Pattern {
  var items []Pattern
  var keyMap []string

  if pat == nil {
    pat = map[string]string{}
  }

  if p.tree.data != nil {
    items = append(items, createMatchList(keyMap, p.tree.data, p.tree.modifier))
  }

  if p.tree.key != "" {

    descendTree(&items, pat, exact, true, p.tree.value, keyMap)
  }
  return items
}

//Same as List but using simepl, string notation
func (p *Patrun) ListString(pat string, exact bool) []Pattern {

  mapData := createMap(pat)

  return p.List(mapData, exact)
}

//Generate JSON representation of the tree.
func (p *Patrun)ToJSON() []byte {
  b, _ := json.Marshal(p.List(nil, false))


  return b
}

//Generate a string representation of the decision tree for debugging.
func (p Patrun)String() string {

  items := p.List(nil, false)

  var data []string

  for k := range items {
    v := items[k]
    data = append(data, fmt.Sprintf("%v -> <%v>", formatMatch(v.Match), formatData(v.Data)))
  }

  return strings.Join(data, "\n")
}

//Generate a string representation of the decision tree for debugging and alows you to specicy a custom formatting function.
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

  var keys  = sortKeys(items)

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
        descendTree(items, pat, exact, false, val.value, append(keyMap, key))

      } else if val.data != nil {
        localKeyMap = append(keyMap, val.key)
        if validatePatternMatch(pat, exact, localKeyMap) {
          *items = append(*items, createMatchList(localKeyMap, val.data, val.modifier))
        }

        if len(val.value) > 0 {
          descendTree(items, pat, exact, false, val.value, append(keyMap, val.key))
        }
      }
  }
}

func validatePatternMatch(pat map[string]string, exact bool, matchedKeys []string) bool {
  var keys = sortKeys(pat)

  if len(keys) == 0 {
    return true
  }

  pathMap := convertListToMap(matchedKeys)

  var matched = true
  for k, v := range pat {

    if pathMap[k] == "" || !gexval(v, pathMap[k]) {
      matched = false
      break
    }
  }

  if exact && len(pat) != len(pathMap) {
    matched = false
  }

  return matched
}


func convertListToMap(listItems []string) map[string]string {
  var mapData = map[string]string{}

  for k := 0; k < len(listItems); k+=2 {
    mapData[listItems[k]] = listItems[k+1]
  }

  return mapData
}

func createMatchList(keyMap []string, dataItem interface{}, modifier Modifiers) Pattern {

  var keys map[string]string = map[string]string{}
  var item Pattern = Pattern{}

  for i := 0; i < len(keyMap); i+=2 {
      if i + 1 < len(keyMap) {
        keys[keyMap[i]] = keyMap[i+1]
      }
  }

  item.Match = keys
  item.Data = dataItem
  item.Modifier = modifier

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

func sortKeys(pat map[string]string) []string {
  var keys []string
  for k := range pat {
    keys = append(keys, k)
  }
  sort.Strings(keys)

  return keys
}
