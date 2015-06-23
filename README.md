# patrun-golang

### A fast pattern matcher on Go map properties

Need to pick out an object based on a subset of its properties? Say you've got:

```Go
{ x: 1          } -> A
{ x: 1, y: 1 } -> B
{ x: 1, y: 2 } -> C
```

Then patrun can give you the following results:

```Go
{ x: 1 }      -> A
{ x: 2 }      -> no match
{ x: 1, y: 1 } -> B
{ x: 1, y: 2 } -> C
{ x: 2, y: 2 } -> no match
{ y: 1 }      -> no match
```

It's basically _query-by-example_ for property sets.


### Support

If you're using this library, feel free to contact me on twitter if you have any questions! :) [@colmharte](http://twitter.com/colmharte)


Current Version: 0.1.0

Tested on: Go v1.4.2


### Quick example

Here's how you register some patterns, and then search for matches:

```Go
package main

import "github.com/colmharte/patrun-golang/patrun"
import "fmt"

func main() {
  pm := patrun.Patrun{}

  pm.Add(map[string]string{"a":"1"},"A")
  pm.Add(map[string]string{"b":"2"},"B")

  // You can also use the AddString method when adding patterns, this allows you to
  specify a map using simple notation
  // format is "key:value,key:value" etc
  // spaces around the : divider and , divider will be stripped out when processed
  pm.AddString("c:3","C")
  pm.AddString("c:3, d:4","D")

  // prints A
  fmt.Println(pm.Find(map[string]string{"a":"1"}))

  // prints nil
  fmt.Println(pm.Find(map[string]string{"a":"2"}))

  // prints A, b:1 is ignored, it was never registered
  fmt.Println(pm.Find(map[string]string{"a":"1", "b":"1"}))

  // prints B, :c:3 is ignored, it was never registered
  fmt.Println(pm.Find(map[string]string{"b":"2", "c":"3"}))

  // can also use the FindString method using the simple string notation

  //prints C
  fmt.Println(pm.FindString("c:3"))

	//prints C, z:4 is ignore
  fmt.Println(pm.FindString("c:3,z:4"))
}
```

You're matching a subset, so your input can contain any number of other properties.


# The Why

This module lets you build a simple decision tree so you can avoid
writing _if_ statements. It tries to make the minimum number of
comparisons necessary to pick out the most specific match.

This is very useful for handling situations where you have lots of
"cases", some of which have "sub-cases", and even "sub-sub-sub-cases".

For example, here are some sales tax rules:

   * default: no sales tax
   * here's a list of countries with known rates: Ireland: 23%, UK: 20%, Germany: 19%, ...
   * but wait, that's only standard rates, here's [the other rates](http://www.vatlive.com/vat-rates/european-vat-rates/eu-vat-rates/)
   * Oh, and we also have the USA, where we need to worry about each state...

Do this:

```Go
package main

import "github.com/colmharte/patrun-golang/patrun"
import "fmt"

// queries return a func, in case there is some
// really custom logic (and there is, see US, NY below)
// in the normal case, just pass the rate back out with
// an identity function
// also record the rate for custom printing later

func main() {

  I := func(val float64) func(float64) float64 {

  	var rate = func(amt float64) float64 {
  		return val
  	}

  	return rate
  }

  salestax := patrun.Patrun{}
  salestax.AddString("", I(0.0) )
  salestax.Add(map[string]string{"country":"IE" }, I(0.25) )
  salestax.AddString("country:UK", I(0.20) )
  salestax.AddString("country:DE" , I(0.19) )
  salestax.AddString("country:IE, type:reduced" , I(0.135) )
  salestax.Add(map[string]string{"country":"IE", "type":"food" },    I(0.048) )
  salestax.AddString("country:UK,type:food",    I(0.0) )
  salestax.Add(map[string]string{"country":"DE", "type":"reduced" }, I(0.07) )
  salestax.Add(map[string]string{"country":"US" }, I(0.0) ) // no federeal rate (yet!)
  salestax.Add(map[string]string{"country":"US", "state":"AL" }, I(0.04) )
  salestax.AddString("country:US , state:AL, city : Montgomery ", I(0.10) )
  salestax.AddString("country:US, state:NY" , I(0.07) )

  under110 := func(net float64) float64 {
    if net < 110 {
      return 0.0
    } else {

  		f := salestax.Find(map[string]string{"country":"US", "state":"NY"})
  		fmt.Println(reflect.ValueOf(f))
  		return 1.0//f.(func())(net)
    }
  }

  salestax.Add(map[string]string{"country":"US", "state":"NY", "type":"reduced" }, under110)


  fmt.Println("Default Rate: ", salestax.Find(map[string]string{}).(func(float64) float64)(99))
  fmt.Println("Standard rate in Ireland on E99: ",salestax.FindString("country:IE").(func(float64) float64)(99))
  fmt.Println("Food rate in Ireland on E99: ",salestax.FindString("country:IE,type:food").(func(float64) float64)(99))
  fmt.Println("Reduced rate in Germany on E99: ",salestax.Find(map[string]string{"country": "DE", "type": "reduced"}).(func(float64) float64)(99))
  fmt.Println("Standard rate in Alabama on $99: ",salestax.FindString("country: US, state: AL").(func(float64) float64)(99))
  fmt.Println("Standard rate in Montgomery, Alabama on $99: ",salestax.Find(map[string]string{"country": "US", "state": "AL", "city": "Montgomery"}).(func(float64) float64)(99))
  fmt.Println("Reduced rate in New York for clothes $99: ",salestax.FindString("country:US, state: NY, type: reduced").(func(float64) float64)(99))


  fmt.Println(salestax.ToString(func(data interface{}) string {
      return fmt.Sprintf(":%v", data.(func(amt float64) float64)(99))
    }))


  // prints:
  // Default rate: 0
  // Standard rate in Ireland on E99: 0.25
  // Food rate in Ireland on E99:     0.048
  // Reduced rate in Germany on E99:  0.07
  // Standard rate in Alabama on $99: 0.04
  // Standard rate in Montgomery, Alabama on $99: 0.1
  // Reduced rate in New York for clothes on $99: 0.0
}
```

You can take a look a the decision tree at any time:

```Go

// print out patterns, using a custom format function
fmt.Println(salestax.ToString(func(data interface{}) string {
    return fmt.Sprintf("%v", data.(func(amt float64) float64)(99))
  }))


// prints:
-> :0.0
city=Montgomery, country=US, state=AL -> :0.1
country=IE -> :0.25
country=IE, type=reduced -> :0.135
country=IE, type=food -> :0.048
country=UK -> :0.2
country=UK, type=food -> :0.0
country=DE -> :0.19
country=DE, type=reduced -> :0.07
country=US -> :0.0
country=US, state=AL -> :0.04
country=US, state=NY -> :0.07
country=US, state=NY, type=reduced -> :0.0
```


# The Rules

   * 1: More specific matches beat less specific matches. That is, more property values beat fewer.
   * 2: Property names are checked in alphabetical order.

And that's it.


# Customization

You can customize the way that data is stored. For example, you might want to add a constant property to each pattern.

To do this, you implement the Customiser interface and pass this in when you create the _patrun_ object:

```Go
package main

import "github.com/colmharte/patrun-golang/patrun"
import "fmt"

type exCustomiser struct{}

func (a exCustomiser) Add(pm *patrun.Patrun, pat map[string]string, data interface{}) patrun.Modifiers {
	pat["foo"] = "true"

	return nil
}


func customiseExample() {
	alwaysAddFoo := patrun.Patrun{Custom: new (exCustomiser)}

	alwaysAddFoo.Add( map[string]string{"a":"1"}, "foorbar" )

	fmt.Println(alwaysAddFoo.FindExactString("a:1" )) // nothing!
	fmt.Println(alwaysAddFoo.FindExact( map[string]string{"a":"1", "foo":"true"} )) // == "foobar"

}

func main() {
  customiseExample()
}
```

The custom interface can also be used to modify found
data, and a modifier method can be used when removing data.

Here's an example that modifies found data:

```Go
package main

import (
  "github.com/colmharte/patrun-golang/patrun"
  "strings"
  "fmt"
)

type exCustomiser1 struct{}

type exModifier1 struct{}

func (a exCustomiser1) Add(pm *patrun.Patrun, pat map[string]string, data interface{}) patrun.Modifiers {
	return new(exModifier1)
}

func (a exModifier1) Find(pm *patrun.Patrun, pat map[string]string, data interface{}) interface{} {
	return strings.ToUpper(data.(string))
}
func (a exModifier1) Remove(pm *patrun.Patrun, pat map[string]string, data interface{}) bool {
	return true
}

func findModifierExample1() {
  upperify := patrun.Patrun{Custom: new (exCustomiser1)}

  upperify.Add( map[string]string{"a":"1"}, "bar" )

  fmt.Println(upperify.Find( map[string]string{"a":"1"} )) // BAR
}

func main() {
  findModifierExample()
}
```

Finally, here's an example that allows you to add multiple matches for a given pattern:

```Go
package main

import "fmt"
import "github.com/colmharte/patrun-golang/patrun"

ype exCustomiser2 struct {
}
type exModifier2 struct {
	items []string
}

func (a *exCustomiser2) Add(pm *patrun.Patrun, pat map[string]string, data interface{}) patrun.Modifiers {
	var items []string

	i := pm.FindExact(pat)
	if i == nil {
		items = []string{}
	} else {
		items = i.([]string)
	}
	items = append(items, data.(string))

	mod := new(exModifier2)
	mod.items = items

	return mod

}
func (a *exModifier2) Find(pm *patrun.Patrun, pat map[string]string, data interface{}) interface{} {


	if 0 < len(a.items) {
		return a.items
	} else {
		return nil
	}

}
func (a *exModifier2) Remove(pm *patrun.Patrun, pat map[string]string, data interface{}) bool {

	if len(a.items) > 0 {
		a.items = a.items[:len(a.items)-1]
	}

	return 0 == len(a.items)
}


func findModifierExample2() {

	many := patrun.Patrun{Custom: new (exCustomiser2)}

	many.AddString( "a:1", "A" )
	many.AddString( "a:1", "B" )
	many.AddString( "b:1", "C" )

	fmt.Println(many.Find( map[string]string{"a":"1"} ))  // [ 'A', 'B' ]
	fmt.Println(many.Find( map[string]string{"b":"1"} )) // [ 'C' ]

	many.Remove( map[string]string{"a":"1"} )
	fmt.Println(many.Find( map[string]string{"a":"1"} )) // [ 'A' ]

	many.Remove( map[string]string{"b":"1"} )
	fmt.Println(many.Find( map[string]string{"b":"1"} )) // nil

}



func main() {
  findModifierExample2()
}
```


# API

## patrun.Patrun{ [Customiser] }

Generates a new pattern matcher instance. Optionally provide a customisation implementation.


## .Add( map[string]string{...pattern...}, object )

Register a pattern, and the object that will be returned if an input
matches.

## .AddString( string{...pattern...}, object )

Same as Add but allows for a pattern to be specified using the simple pattern notation rather then having to create a map object.

Format of notation is as below
"key:value, key:value"
eg: "a:1, b:2, c:3"

White space is optional. This notation will be turned into a map object when the method is called


## .Find( map[string]string{...subject...} )

Return the unique match for this subject, or nil if not found. The
properties of the subject are matched against the patterns previously
added, and the most specifc pattern wins. Unknown properties in the
subject are ignored.

## .FindSring( string{...subject...} )

Same as Find but with simple string notation

## .FindExact( map[string]string{...subject...} )

Same as Find but only matches where all properties match will be returned.

## .FindExactString( string{...subject...} )

Same as FindExact but with simple string notation

## .List( map[string]string{...pattern-partial...}, exact bool )

Return the list of registered patterns that contain this partial
pattern. You can use wildcards for property values.  Omitted values
are *not* equivalent to a wildcard of _"*"_, you must specify each
property explicitly. You can provide a second boolean
parameter, _exact_. If true, then only those patterns matching the
pattern-partial exactly are returned.

```
pm = patrun.Patrun{}
pm.AddString("a:1, b:1","B1").AddString("a:1, b:2","B2")

// finds:
// [ { match: { a:1, b:1 }, :data:B1 },
//   { match: { a:1, b:2 }, :data:B2 } ]
fmt.Println(pm.ListString("a:1", false))

// finds:
// [ { match: { a:1, b:1 }, data:B1 },
//   { match: { a:1, b:2 }, data:B2 } ]
fmt.Println(pm.ListString("a:1, b:*", false))

// finds:
// [ { match: { a:1, b:1 }, data:B1 }]
fmt.Println(pm.ListString("a:1, b:1", false))

// finds nothing: []
fmt.Println(pm.ListString("c:1", false))
```

If you provide no pattern argument at all, _list_ will list all patterns that have been added.
```Ruby
# finds everything
fmt.Println(pm.List(nil, false))
```

## .ListString( string{...pattern-partial...}, exact bool )

Same as List but with simple string notation

## .Remove( map[string]string{...pattern...} )

Remove this pattern, and it's object, from the matcher.

## .RemoveString( string{...pattern...} )

Same as Remove but with simple string notation.


## .ToString( proc )

Generate a string representation of the decision tree for debugging. Provide a formatting function for objects.

   * proc: format proc for data

## .String( )

Generate a string representation of the decision tree for debugging.


## .ToJSON()

Generate JSON representation of the tree.


# Development

From the Irish patr&uacute;n: [pattern](http://www.focloir.ie/en/dictionary/ei/pattern). Pronounced _pah-troon_.
