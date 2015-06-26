package patrun

import (
  "github.com/colmharte/patrun-golang/patrun"
  "testing"
  "regexp"
  "fmt"
  "strings"
  "sort"
)

func TestEmpty(t *testing.T) {

    r := patrun.Patrun{}

    if r.String() != "" {
      t.Error("pattern should be empty", r.String());
    }

    if r.Find(nil) != nil {
      t.Error("nil Find should return nil", r.Find(nil));
    }

    if r.Find(map[string]string{}) != nil {
      t.Error("empty Find should return nil", r.Find(map[string]string{}));
    }

    if r.Find(map[string]string{"a":"1"}) != nil {
      t.Error("{a:1} Find should return nil", r.Find(map[string]string{"a":"1"}));
    }

    r.Add(map[string]string{"a":"1"},"A")

    if r.Find(nil) != nil {
      t.Error("nil Find should return nil", r.Find(nil));
    }

    if r.Find(map[string]string{}) != nil {
      t.Error("empty Find should return nil", r.Find(map[string]string{}));
    }

    if r.Find(map[string]string{"a":"1"}) != "A" {
      t.Error("{a:1} Find should return A", r.Find(map[string]string{"a":"1"}).(string));
    }
}

func TestRoot(t *testing.T) {

    r := patrun.Patrun{}

    r.Add(map[string]string{}, "R")

    if r.String() != " -> <R>" {
      t.Error("pattern should be -> <R>", r.String());
    }

    if rs(r) != "<R>" {
      t.Error("pattern should be <R>", rs(r));
    }

    if string(r.ToJSON()[:]) != "[{\"Match\":{},\"Data\":\"R\",\"Modifier\":null}]" {
      t.Error("JSON pattern should be [{\"Match\":{},\"Data\":\"R\",\"Modifier\":null}]", string(r.ToJSON()[:]));
    }

    if r.Find(map[string]string{}) != "R" {
      t.Error("empty Find should return R", r.Find(map[string]string{}));
    }

    if r.Find(map[string]string{"x":"1"}) != "R" {
      t.Error("{x:1} Find should return R", r.Find(map[string]string{"x":"1"}));
    }

    if r.Find(map[string]string{"a":"1"}) != "R" {
      t.Error("{a:1} Find should return R", r.Find(map[string]string{"a":"1"}));
    }

    r.Add(map[string]string{"a":"1"},"r1")

    if r.String() != " -> <R>\na:1 -> <r1>" {
      t.Error("pattern should be  -> <R>\na:1 -> <r1>", r.String());
    }
    if rs(r) != "<R>a:1<r1>" {
      t.Error("pattern should be <R>a:1<r1>", rs(r));
    }

    if string(r.ToJSON()[:]) != "[{\"Match\":{},\"Data\":\"R\",\"Modifier\":null},{\"Match\":{\"a\":\"1\"},\"Data\":\"r1\",\"Modifier\":null}]" {
      t.Error("JSON pattern should be [{\"Match\":{},\"Data\":\"R\",\"Modifier\":null},{\"Match\":{\"a\":\"1\"},\"Data\":\"r1\",\"Modifier\":null}]", string(r.ToJSON()[:]));
    }

}

func TestAdd(t *testing.T) {

  r := patrun.Patrun{}

  r.Add(map[string]string{"a":"1"}, "r1")

  if r.String() != "a:1 -> <r1>" {
    t.Error("pattern should be a:1 -> <r1>", r.String());
  }

  pat := rs(r)
  if pat != "a:1<r1>" {
    t.Error("pattern should be a:1<r1>", pat);
  }

  if string(r.ToJSON()[:]) != "[{\"Match\":{\"a\":\"1\"},\"Data\":\"r1\",\"Modifier\":null}]" {
    t.Error("JSON pattern should be [{\"Match\":{\"a\":\"1\"},\"Data\":\"r1\",\"Modifier\":null}]", string(r.ToJSON()[:]));
  }

  r = patrun.Patrun{}
  r.Add(map[string]string{"a":"1", "b": "2"}, "r1")
  pat = rs(r)
  if pat != "a:1,b:2<r1>" {
    t.Error("pattern should be a:1,b:2<r1>", pat);
  }

  r = patrun.Patrun{}
  r.Add(map[string]string{"a":"1", "b": "2", "c":"3"}, "r1")
  pat = rs(r)
  if pat != "a:1,b:2,c:3<r1>" {
    t.Error("pattern should be a:1,b:2,c:3<r1>", pat);
  }

  r = patrun.Patrun{}
  r.Add(map[string]string{"a":"1", "b": "2"}, "r1")
  r.Add(map[string]string{"a":"1", "b": "3"}, "r2")
  pat = rs(r)
  if pat != "a:1,b:2<r1>a:1,b:3<r2>" {
    t.Error("pattern should be a:1,b:2<r1>a:1,b:3<r2>", pat);
  }

  r = patrun.Patrun{}
  r.Add(map[string]string{"a":"1", "b": "2"}, "r1")
  r.Add(map[string]string{"a":"1", "c": "3"}, "r2")
  pat = rs(r)
  if pat != "a:1,b:2<r1>a:1,c:3<r2>" {
    t.Error("pattern should be a:1,b:2<r1>a:1,c:3<r2>", pat);
  }
  r.Add(map[string]string{"a":"1", "d": "4"}, "r3")
  pat = rs(r)
  if pat != "a:1,b:2<r1>a:1,c:3<r2>a:1,d:4<r3>" {
    t.Error("pattern should be a:1,b:2<r1>a:1,c:3<r2>a:1,d:4<r3>", pat);
  }

  r = patrun.Patrun{}
  r.Add(map[string]string{"a":"1", "c": "2"}, "r1")
  r.Add(map[string]string{"a":"1", "b": "3"}, "r2")
  pat = rs(r)
  if pat != "a:1,b:3<r2>a:1,c:2<r1>" {
    t.Error("pattern should be a:1,b:3<r2>a:1,c:2<r1>", pat);
  }
  if string(r.ToJSON()[:]) != "[{\"Match\":{\"a\":\"1\",\"b\":\"3\"},\"Data\":\"r2\",\"Modifier\":null},{\"Match\":{\"a\":\"1\",\"c\":\"2\"},\"Data\":\"r1\",\"Modifier\":null}]" {
    t.Error("JSON pattern should be [{\"Match\":{\"a\":\"1\",\"b\":\"3\"},\"Data\":\"r2\",\"Modifier\":null},{\"Match\":{\"a\":\"1\",\"c\":\"2\"},\"Data\":\"r1\",\"Modifier\":null}]", string(r.ToJSON()[:]));
  }
}

func TestBasic(t *testing.T) {
  r := patrun.Patrun{}

  r.Add(map[string]string{"p1":"v1"}, "r1")

  if r.Find(map[string]string{"p1":"v1"}).(string) != "r1" {
    t.Error("p1:v1 Find should be r1", r.Find(map[string]string{"p1":"v1"}));
  }

  if r.Find(map[string]string{"p1":"v2"}) != nil {
    t.Error("p1:v2 Find should be nil", r.Find(map[string]string{"p1":"v2"}));
  }

  r.Add(map[string]string{"p1":"v1"}, "r1x")

  if r.Find(map[string]string{"p1":"v1"}).(string) != "r1x" {
    t.Error("p1:v1 Find should be r1x", r.Find(map[string]string{"p1":"v1"}));
  }
  if r.Find(map[string]string{"p2":"v1"}) != nil {
    t.Error("p2:v1 Find should be nil", r.Find(map[string]string{"p2":"v1"}));
  }
  r.Add(map[string]string{"p1":"v2"}, "r2")
  if r.Find(map[string]string{"p1":"v2"}).(string) != "r2" {
    t.Error("p1:v2 Find should be r2", r.Find(map[string]string{"p1":"v2"}));
  }
  if r.Find(map[string]string{"p2":"v2"}) != nil {
    t.Error("p2:v2 Find should be nil", r.Find(map[string]string{"p2":"v2"}));
  }

  r.Add(map[string]string{"p1":"v3"}, "r3")
  if r.Find(map[string]string{"p1":"v3"}).(string) != "r3" {
    t.Error("p1:v3 Find should be r3", r.Find(map[string]string{"p1":"v3"}));
  }
  if r.Find(map[string]string{"p2":"v2"}) != nil {
    t.Error("p2:v2 Find should be nil", r.Find(map[string]string{"p2":"v2"}));
  }
  if r.Find(map[string]string{"p2":"v1"}) != nil {
    t.Error("p2:v1 Find should be nil", r.Find(map[string]string{"p2":"v1"}));
  }

  r.Add(map[string]string{"p1":"v1", "p3":"v4"}, "r4")
  if r.Find(map[string]string{"p1":"v1", "p3":"v4"}).(string) != "r4" {
    t.Error("p1:v3,ps:v4 Find should be r4", r.Find(map[string]string{"p1":"v3", "p3":"v4"}));
  }
  if r.Find(map[string]string{"p1":"v1", "p3":"v5"}).(string) != "r1x" {
    t.Error("p1:v1,p3:v5 Find should be r1x", r.Find(map[string]string{"p1":"v1", "p3":"v5"}));
  }
  if r.Find(map[string]string{"p2":"v1"}) != nil {
    t.Error("p2:v1 Find should be nil", r.Find(map[string]string{"p2":"v1"}));
  }
}

func TestCulDesac(t *testing.T) {

  r := patrun.Patrun{}

  r.AddString("p1:v1", "r1")
  r.AddString("p1:v1, p2:v2", "r2")
  r.AddString("p1:v1, p3:v3", "r3")

  if r.FindString("p1:v1,p2:x").(string) != "r1" {
    t.Error("p1:v1,p2:x Find should be r1", r.FindString("p1:v1,p2:x"));
  }
  if r.FindString("p1:v1,p2:x,p3:v3").(string) != "r3" {
    t.Error("p1:v1,p2:x,p3:v3 Find should be r3", r.FindString("p1:v1,p2:x,p3:v3"));
  }
}

func TestRemove(t *testing.T) {
  r := patrun.Patrun{}
  r.RemoveString("p1:v1")

  r.AddString("", "s0")
  if r.FindString("").(string) != "s0" {
    t.Error("empty Find should be s0", r.FindString(""));
  }
  r.RemoveString("")
  if r.FindString("") != nil {
    t.Error("empty Find should be nil", r.FindString(""));
  }

  r.AddString("p1:v1", "r0" )
  if r.Find(map[string]string{"p1":"v1"}).(string) != "r0" {
    t.Error("p1:v1 Find should be r0", r.Find(map[string]string{"p1":"v1"}));
  }


  r.RemoveString("p1:v1")
  if r.FindString("p1:v1") != nil {
    t.Error("p1:v1 Find should be nil", r.FindString("p1:v2"));
  }

  r.AddString("p2:v2,p3:v3", "r1" )
  r.AddString("p2:v2,p4:v4", "r2" )

  if r.FindString("p2:v2,p3:v3").(string) != "r1" {
    t.Error("p2:v2,p3:v3 Find should be r1", r.FindString("p2:v2,p3:v3"));
  }
  if r.FindString("p2:v2,p4:v4").(string) != "r2" {
    t.Error("p2:v2,p4:v4 Find should be r1", r.FindString("p2:v2,p4:v4"));
  }

  r.RemoveString("p2:v2,p3:v3")
  if r.FindString("p2:v2,p3:v3") != nil {
    t.Error("p2:v2 ,p3:v3Find should be nil", r.FindString("p2:v2,p3:v3"));
  }
  if r.FindString("p2:v2,p4:v4").(string) != "r2" {
    t.Error("p2:v2,p4:v4 Find should be r1", r.FindString("p2:v2,p4:v4"));
  }
}

func TestRemoveIntermediate(t *testing.T) {

  r := patrun.Patrun{}

  r.AddString("a:1, b:2,d:4", "XX" )
  r.AddString("c:3,d:4", "YY" )
  r.AddString("a:1, b:2", "X" )
  r.AddString("c:3", "Y" )

  r.RemoveString("c:3")

  if r.FindString("c:3") != nil {
    t.Error("c3 Find should be nil", r.FindString("c:3"));
  }
  if r.FindString("a:1,c:3, d:4").(string) != "YY" {
    t.Error("a:1,c:3,d:4 Find should be YY", r.FindString("a:1,c:3, d:4"));
  }
  if r.FindString("a:1,b:2, d:4").(string) != "XX" {
    t.Error("a:1,b:2,d:4 Find should be XX", r.FindString("a:1,b:2, d:4"));
  }
  if r.FindString("a:1,b:2").(string) != "X" {
    t.Error("a:1,b:2 Find should be X", r.FindString("a:1,b:2"));
  }

  r.RemoveString("a:1,b:2")
  if r.FindString("c:3") != nil {
    t.Error("c3 Find should be nil", r.FindString("c:3"));
  }
  if r.FindString("a:1,c:3, d:4").(string) != "YY" {
    t.Error("a:1,c:3,d:4 Find should be YY", r.FindString("a:1,c:3, d:4"));
  }
  if r.FindString("a:1,b:2, d:4").(string) != "XX" {
    t.Error("a:1,b:2,d:4 Find should be XX", r.FindString("a:1,b:2, d:4"));
  }
  if r.FindString("a:1,b:2") != nil {
    t.Error("a:1,b:2 Find should be nil", r.FindString("a:1,b:2"));
  }
}

func TestExact(t *testing.T) {
  r := patrun.Patrun{}

  r.AddString("a:1", "X" )

  if r.FindExactString("a:1").(string) != "X" {
    t.Error("a:1 FindExact should be X", r.FindExactString("a:1"));
  }
  if r.FindExactString("a:1,b:2") != nil {
    t.Error("a:1,b:2 FindExact should be nil", r.FindExactString("a:1,b:2"));
  }

}

func TestAll(t *testing.T) {
  r := patrun.Patrun{}

  r.AddString("a:1", "X" )
  r.AddString("b:2", "Y" )

  if fmt.Sprintf("%v", r.List(nil, false)) != "[{map[a:1] X <nil>} {map[b:2] Y <nil>}]" {
    t.Error("List all should be [{map[a:1] X <nil>} {map[b:2] Y <nil>}]", fmt.Sprintf("%v", r.List(nil, false)))
  }

}

func TestMultiStar(t *testing.T) {
  r := patrun.Patrun{}


  r.AddString("a:1", "A" )
  r.AddString("a:1,b:2", "B" )
  r.AddString("c:3", "C" )
  r.AddString("b:1,c:4", "D" )


  pat := rs(r)
  if pat != "a:1<A>a:1,b:2<B>b:1,c:4<D>c:3<C>" {
    t.Error("pattern should be a:1<A>a:1,b:2<B>b:1,c:4<D>c:3<C>", pat);
  }

  if r.FindString("c:3").(string) != "C" {
    t.Error("c:3 Find should be C", r.FindString("c:3"));
  }
  if r.FindString("c:3,a:0").(string) != "C" {
    t.Error("c:3,a:0 Find should be C", r.FindString("c:3,a:0"));
  }
  if r.FindString("c:3,a:0,b:0").(string) != "C" {
    t.Error("c:3,a:0,b:0 Find should be C", r.FindString("c:3,a:0,b:0"));
  }

}

func TestStarBacktrack(t *testing.T) {
  r := patrun.Patrun{}

  r.AddString("a:1,b:2", "X" )
  r.AddString("c:3", "Y" )

  if r.FindString("a:1,b:2").(string) != "X" {
    t.Error("a:1,b:2 Find should be X", r.FindString("a:1,b:2"));
  }
  if r.FindString("a:1,b:0,c:3").(string) != "Y" {
    t.Error("a:1,b:0,c:3 Find should be Y", r.FindString("a:1,b:0,c:3"));
  }

  r.AddString("a:1,b:2,d:4", "XX" )
  r.AddString("c:3,d:4", "YY" )

  if r.FindString("a:1,b:2,d:4").(string) != "XX" {
    t.Error("a:1,b:2,d:4 Find should be XX", r.FindString("a:1,b:2,d:4"));
  }
  if r.FindString("a:1,c:3,d:4").(string) != "YY" {
    t.Error("a:1,c:3,d:4 Find should be YY", r.FindString("a:1,c:3,d:4"));
  }
  if r.FindString("a:1,b:2").(string) != "X" {
    t.Error("a:1,b:2 Find should be X", r.FindString("a:1,b:2"));
  }
  if r.FindString("a:1,b:0,c:3").(string) != "Y" {
    t.Error("a:1,b:0,c:3 Find should be Y", r.FindString("a:1,b:0,c:3"));
  }

  if r.ListString("a:1,b:*", false)[0].Data.(string) != "X" {
    t.Error("List a:1,b:* item 0 should be X", r.ListString("a:1,b:*", false)[0].Data);
  }
  if r.ListString("c:3", false)[0].Data.(string) != "Y" {
    t.Error("List c:3 item 0 should be Y", r.ListString("c:3", false)[0].Data);
  }
  if r.ListString("c:3,d:*", false)[0].Data.(string) != "YY" {
    t.Error("List c:3,d:* item 0 should be YY", r.ListString("c:3,d:*", false)[0].Data);
  }
  if r.ListString("a:1,b:*,d:*", false)[0].Data.(string) != "XX" {
    t.Error("List a:1,b:*,d:* item 0 should be XX", r.ListString("a:1,b:*,d:*", false)[0].Data);
  }

  pat := rs(r)
  if pat != "a:1,b:2<X>a:1,b:2,d:4<XX>c:3<Y>c:3,d:4<YY>" {
    t.Error("pattern should be a:1,b:2<X>a:1,b:2,d:4<XX>c:3<Y>c:3,d:4<YY>", pat);
  }

}


func TestListTopVals(t *testing.T) {
  r := patrun.Patrun{}

  //r.AddString("a:1", "x")

  //'subvals==mode' && rt1.add( {a:'1'}, 'x' )

  r.AddString("p1:v1", "r0")
  r.AddString("p1:v1,p2:v2a", "r1")
  r.AddString("p1:v1,p2:v2b", "r1")

  var pat = fmt.Sprintf("%v", r.ListString("p1:v1", true))
  if pat != "[{map[p1:v1] r0 <nil>}]" {
    t.Error("List p1:v1 should be [{map[p1:v1] r0 <nil>}]", pat);
  }

  pat = convertListToString(r.ListString("p1:v1,p2:*", true))
  if pat != "[{map[p1:v1 p2:v2a] r1} {map[p1:v1 p2:v2b] r1}]" {
    t.Error("List p1:v1 should be [{map[p1:v1 p2:v2a] r1} {map[p1:v1 p2:v2b] r1}]", pat);
  }

  r.AddString("p1:v1,p2:v2c,p3:v3a", "r3a")
  r.AddString("p1:v1,p2:v2d,p3:v3b", "r3b")

  pat = convertListToString(r.ListString("p1:v1,p2:*,p3:v3a", true))
  if pat != "[{map[p1:v1 p2:v2c p3:v3a] r3a}]" {
    t.Error("List p1:v1 should be [{map[p1:v1 p2:v2c p3:v3a] r3a}]", pat);
  }

}

func TestListSubVals(t *testing.T) {
  r := patrun.Patrun{}

  r.AddString("a:1", "x")

  r.AddString("p1:v1", "r0")
  r.AddString("p1:v1,p2:v2a", "r1")
  r.AddString("p1:v1,p2:v2b", "r1")

  var pat = fmt.Sprintf("%v", r.ListString("p1:v1", true))
  if pat != "[{map[p1:v1] r0 <nil>}]" {
    t.Error("List p1:v1 should be [{map[p1:v1] r0 <nil>}]", pat);
  }

  pat = convertListToString(r.ListString("p1:v1,p2:*", true))
  if pat != "[{map[p1:v1 p2:v2a] r1} {map[p1:v1 p2:v2b] r1}]" {
    t.Error("List p1:v1 should be [{map[p1:v1 p2:v2a] r1} {map[p1:v1 p2:v2b] r1}]", pat);
  }

  r.AddString("p1:v1,p2:v2c,p3:v3a", "r3a")
  r.AddString("p1:v1,p2:v2d,p3:v3b", "r3b")

  pat = convertListToString(r.ListString("p1:v1,p2:*,p3:v3a", true))
  if pat != "[{map[p1:v1 p2:v2c p3:v3a] r3a}]" {
    t.Error("List p1:v1 should be [{map[p1:v1 p2:v2c p3:v3a] r3a}]", pat);
  }

}

type customHappy struct{}

func (a customHappy) Add(pm *patrun.Patrun, pat map[string]string, data interface{}) patrun.Modifiers {
	pat["q"] = "9"

  return nil
}
func TestCustomHappy(t *testing.T) {
  r := patrun.Patrun{Custom: new(customHappy)}

  r.AddString("a:1", "Q")

  if r.FindString("a:1") != nil {
    t.Error("a:1 Find should be nil", r.FindString("a:1"));
  }

  if r.FindString("a:1,q:9").(string) != "Q" {
    t.Error("a:1,q:9 Find should be Q", r.FindString("a:1,q:9"));
  }

}

type customMany struct{}
type customModifierMany struct {
	items []string
}
func (a *customMany) Add(pm *patrun.Patrun, pat map[string]string, data interface{}) patrun.Modifiers {
	var items []string

	i := pm.FindExact(pat)
	if i == nil {
		items = []string{}
	} else {
		items = i.([]string)
	}
	items = append(items, data.(string))

	mod := new(customModifierMany)
	mod.items = items

	return mod

}
func (a *customModifierMany) Find(pm *patrun.Patrun, pat map[string]string, data interface{}) interface{} {


	if 0 < len(a.items) {
		return a.items
	} else {
		return nil
	}

}
func (a *customModifierMany) Remove(pm *patrun.Patrun, pat map[string]string, data interface{}) bool {

	if len(a.items) > 0 {
		a.items = a.items[:len(a.items)-1]
	}

	return 0 == len(a.items)
}


func TestCustomMany(t *testing.T) {
  r := patrun.Patrun{Custom: new(customMany)}

  r.AddString("a:1", "A")
  r.AddString("a:1", "B")
  r.AddString("b:1", "C")

  if fmt.Sprintf("%v",r.FindString("a:1").([]string)) != "[A B]" {
    t.Error("a:1 Find should be [A B]", r.FindString("a:1"));
  }
  if fmt.Sprintf("%v", r.FindString("b:1").([]string)) != "[C]" {
    t.Error("b:1 Find should be [C]", r.FindString("b:1"));
  }

  if len(r.List(nil, false)) != 2 {
    t.Error("List should have length of 2", len(r.List(nil, false)));
  }

  r.RemoveString("b:1")
  if len(r.List(nil, false)) != 1 {
    t.Error("List should have length of 1", len(r.List(nil, false)));
  }
  if r.FindString("b:1")!= nil {
    t.Error("b:1 Find should be nil", r.FindString("b:1"));
  }
  if fmt.Sprintf("%v",r.FindString("a:1").([]string)) != "[A B]" {
    t.Error("a:1 Find should be [A B]", r.FindString("a:1"));
  }

  r.RemoveString("a:1")
  if len(r.List(nil, false)) != 1 {
    t.Error("List should have length of 1", len(r.List(nil, false)));
  }
  if r.FindString("b:1")!= nil {
    t.Error("b:1 Find should be nil", r.FindString("b:1"));
  }
  if fmt.Sprintf("%v",r.FindString("a:1").([]string)) != "[A]" {
    t.Error("a:1 Find should be [A]", r.FindString("a:1"));
  }

  r.RemoveString("a:1")
  if len(r.List(nil, false)) != 0 {
    t.Error("List should have length of 0", len(r.List(nil, false)));
  }
  if r.FindString("b:1")!= nil {
    t.Error("b:1 Find should be nil", r.FindString("b:1"));
  }
  if r.FindString("a:1")!= nil {
    t.Error("a:1 Find should be nil", r.FindString("a:1"));
  }


}


func TestFindExact(t *testing.T) {
  r := patrun.Patrun{}

  r.AddString("a:1", "A")
  r.AddString("a:1,b:2", "B")
  r.AddString("a:1,b:2,c:3", "C")

  if r.FindString("a:1").(string) != "A" {
    t.Error("a:1 Find should be A", r.FindString("a:1"));
  }
  if r.FindExactString("a:1").(string) != "A" {
    t.Error("a:1 Find should be A", r.FindString("a:1"));
  }
  if r.FindString("a:1,b:8").(string) != "A" {
    t.Error("a:1,b:8 Find should be A", r.FindString("a:1,b:8"));
  }
  if r.FindExactString("a:1,b:8") != nil {
    t.Error("a:1,b:8 Find should be nil", r.FindString("a:1,b:8"));
  }
  if r.FindString("a:1,b:8,c:3").(string) != "A" {
    t.Error("a:1,b:8c:3 Find should be A", r.FindString("a:1,b:8,c:3"));
  }
  if r.FindExactString("a:1,b:8,c:3") != nil {
    t.Error("a:1,b:8,c:3 Find should be nil", r.FindString("a:1,b:8,c:3"));
  }

  if r.FindString("a:1,b:2").(string) != "B" {
    t.Error("a:1,b:2 Find should be B", r.FindString("a:1,b:2"));
  }
  if r.FindExactString("a:1,b:2").(string) != "B" {
    t.Error("a:1,b:2 Find should be B", r.FindString("a:1,b:2"));
  }
  if r.FindString("a:1,b:2,c:9").(string) != "B" {
    t.Error("a:1,b:2c:9 Find should be B", r.FindString("a:1,b:2,c:9"));
  }
  if r.FindExactString("a:1,b:2,c:9") != nil {
    t.Error("a:1,b:2,c:9 Find should be nil", r.FindString("a:1,b:2,c:9"));
  }

  if r.FindString("a:1,b:2,c:3").(string) != "C" {
    t.Error("a:1,b:2,c:3 Find should be C", r.FindString("a:1,b:2,c:3"));
  }
  if r.FindExactString("a:1,b:2,c:3").(string) != "C" {
    t.Error("a:1,b:2,c:3 Find should be C", r.FindString("a:1,b:2,c:3"));
  }
  if r.FindString("a:1,b:2,c:3,d:7").(string) != "C" {
    t.Error("a:1,b:2,c:3,d:7 Find should be C", r.FindString("a:1,b:2,c:3,d:7"));
  }
  if r.FindExactString("a:1,b:2,c:3,d:7") != nil {
    t.Error("a:1,b:2,c:3,d:7 Find should be nil", r.FindString("a:1,b:2,c:3,d:7"));
  }

}

type customTop struct{}
type customModifierTop struct {}

func (a *customTop) Add(pm *patrun.Patrun, pat map[string]string, data interface{}) patrun.Modifiers {
	return new(customModifierTop)
}
func (a *customModifierTop) Find(pm *patrun.Patrun, pat map[string]string, data interface{}) interface{} {

  return fmt.Sprintf("%v!", data)
}
func (a *customModifierTop) Remove(pm *patrun.Patrun, pat map[string]string, data interface{}) bool {
  return true
}

func TestCustomTop(t *testing.T) {
  r := patrun.Patrun{Custom: new(customTop)}

  r.AddString("", "Q")
  r.AddString("a:1", "A")
  r.AddString("a:1,b:2", "B")
  r.AddString("a:1,b:2,c:3", "C")

  if r.FindString("").(string) != "Q!" {
    t.Error("{} Find should be Q!", r.FindString(""));
  }
  if r.FindString("a:1").(string) != "A!" {
    t.Error("a:1 Find should be A!", r.FindString("a:1"));
  }
  if r.FindString("a:1,b:2").(string) != "B!" {
    t.Error("a:1,b:2 Find should be B!", r.FindString("a:1,b:2"));
  }
  if r.FindString("a:1,b:2,c:3").(string) != "C!" {
    t.Error("a:1,b:2,c:3 Find should be C!", r.FindString("a:1,b:2,c:3"));
  }

}

func TestListAny(t *testing.T) {
  r := patrun.Patrun{}

  r.AddString("a:1", "A")
  r.AddString("a:1,b:2", "B")
  r.AddString("a:1,b:2,c:3", "C")


  var mA = "{map[a:1] A}"
  var mB = "{map[a:1 b:2] B}"
  var mC = "{map[a:1 b:2 c:3] C}"

  if convertListToString(r.List(nil, false)) != fmt.Sprintf("[%v %v %v]", mA, mB, mC) {
    t.Error("List should be ", fmt.Sprintf("[%v %v %v]", mA, mB, mC) , convertListToString(r.List(nil, false)));
  }

  if convertListToString(r.ListString("a:1", false)) != fmt.Sprintf("[%v %v %v]", mA, mB, mC) {
    t.Error("a:1 List should be ", fmt.Sprintf("[%v %v %v]", mA, mB, mC) , convertListToString(r.ListString("a:1", false)));
  }
  if convertListToString(r.ListString("b:2", false)) != fmt.Sprintf("[%v %v]", mB, mC) {
    t.Error("b:2 List should be ", fmt.Sprintf("[%v %v]", mB, mC) , convertListToString(r.ListString("b:2", false)));
  }
  if convertListToString(r.ListString("c:3", false)) != fmt.Sprintf("[%v]", mC) {
    t.Error("c:3 List should be ", fmt.Sprintf("[%v]", mC) , convertListToString(r.ListString("c:3", false)));
  }

  if convertListToString(r.ListString("a:*", false)) != fmt.Sprintf("[%v %v %v]", mA, mB, mC) {
    t.Error("a:* List should be ", fmt.Sprintf("[%v %v %v]", mA, mB, mC) , convertListToString(r.ListString("a:*", false)));
  }
  if convertListToString(r.ListString("b:*", false)) != fmt.Sprintf("[%v %v]", mB, mC) {
    t.Error("b:* List should be ", fmt.Sprintf("[%v %v]", mB, mC) , convertListToString(r.ListString("b:*", false)));
  }
  if convertListToString(r.ListString("c:*", false)) != fmt.Sprintf("[%v]", mC) {
    t.Error("c:* List should be ", fmt.Sprintf("[%v]", mC) , convertListToString(r.ListString("c:*", false)));
  }

  if convertListToString(r.ListString("a:1,b:2", false)) != fmt.Sprintf("[%v %v]", mB, mC) {
    t.Error("a:1,b:2 List should be ", fmt.Sprintf("[%v %v]", mB, mC) , convertListToString(r.ListString("a:1,b:2", false)));
  }
  if convertListToString(r.ListString("a:1,b:*", false)) != fmt.Sprintf("[%v %v]", mB, mC) {
    t.Error("a:1,b:* List should be ", fmt.Sprintf("[%v %v]", mB, mC) , convertListToString(r.ListString("a:1, b:*", false)));
  }
  if convertListToString(r.ListString("a:1,b:*,c:3", false)) != fmt.Sprintf("[%v]", mC) {
    t.Error("a:1,b*,c:3 List should be ", fmt.Sprintf("[%v]", mC) , convertListToString(r.ListString("a:1,b:*,c:3", false)));
  }
  if convertListToString(r.ListString("a:1,b:*,c:*", false)) != fmt.Sprintf("[%v]", mC) {
    t.Error("a:1,b*,c:* List should be ", fmt.Sprintf("[%v]", mC) , convertListToString(r.ListString("a:1,b:*,c:*", false)));
  }
  if convertListToString(r.ListString("a:1,c:*", false)) != fmt.Sprintf("[%v]", mC) {
    t.Error("a:1,c:* List should be ", fmt.Sprintf("[%v]", mC) , convertListToString(r.ListString("a:1,c:*", false)));
  }

  r.AddString("a:1,d:4", "D")
  var mD = "{map[a:1 d:4] D}"

  if convertListToString(r.ListString("", false)) != fmt.Sprintf("[%v %v %v %v]", mA, mB, mC, mD) {
    t.Error("List all should be ", fmt.Sprintf("[%v %v %v %v]", mA,mB, mC, mD) , convertListToString(r.ListString("", false)));
  }
  if convertListToString(r.ListString("a:1", false)) != fmt.Sprintf("[%v %v %v %v]", mA, mB, mC, mD) {
    t.Error("a:1 List should be ", fmt.Sprintf("[%v %v %v %v]", mA,mB, mC, mD) , convertListToString(r.ListString("a:1", false)));
  }
  if convertListToString(r.ListString("d:4", false)) != fmt.Sprintf("[%v]",  mD) {
    t.Error("d:4 List should be ", fmt.Sprintf("[%v]",  mD) , convertListToString(r.ListString("d:4", false)));
  }
  if convertListToString(r.ListString("a:1,d:4", false)) != fmt.Sprintf("[%v]",  mD) {
    t.Error("a:1,d:4 List should be ", fmt.Sprintf("[%v]",  mD) , convertListToString(r.ListString("a:1,d:4", false)));
  }
  if convertListToString(r.ListString("a:1,d:*", false)) != fmt.Sprintf("[%v]",  mD) {
    t.Error("a:1,d:* List should be ", fmt.Sprintf("[%v]",  mD) , convertListToString(r.ListString("a:1,d:*", false)));
  }
  if convertListToString(r.ListString("d:*", false)) != fmt.Sprintf("[%v]",  mD) {
    t.Error("d:* List should be ", fmt.Sprintf("[%v]",  mD) , convertListToString(r.ListString("d:*", false)));
  }

  r.AddString("a:1,c:33", "CC")
  var mCC = "{map[a:1 c:33] CC}"

  if convertListToString(r.ListString("", false)) != fmt.Sprintf("[%v %v %v %v %v]", mA, mB, mC, mCC, mD) {
    t.Error("List all should be ", fmt.Sprintf("[%v %v %v %v %v]", mA,mB, mC, mCC, mD) , convertListToString(r.ListString("", false)));
  }
  if convertListToString(r.ListString("a:1", false)) != fmt.Sprintf("[%v %v %v %v %v]", mA, mB, mC, mCC, mD) {
    t.Error("a:1 List should be ", fmt.Sprintf("[%v %v %v %v %v]", mA,mB, mC, mCC, mD) , convertListToString(r.ListString("a:1", false)));
  }
  if convertListToString(r.ListString("d:4", false)) != fmt.Sprintf("[%v]",  mD) {
    t.Error("d:4 List should be ", fmt.Sprintf("[%v]",  mD) , convertListToString(r.ListString("d:4", false)));
  }
  if convertListToString(r.ListString("a:1,d:4", false)) != fmt.Sprintf("[%v]",  mD) {
    t.Error("a:1,d:4 List should be ", fmt.Sprintf("[%v]",  mD) , convertListToString(r.ListString("a:1,d:4", false)));
  }
  if convertListToString(r.ListString("a:1,d:*", false)) != fmt.Sprintf("[%v]",  mD) {
    t.Error("a:1,d:* List should be ", fmt.Sprintf("[%v]",  mD) , convertListToString(r.ListString("a:1,d:*", false)));
  }
  if convertListToString(r.ListString("d:*", false)) != fmt.Sprintf("[%v]",  mD) {
    t.Error("d:* List should be ", fmt.Sprintf("[%v]",  mD) , convertListToString(r.ListString("d:*", false)));
  }
  if convertListToString(r.ListString("c:33", false)) != fmt.Sprintf("[%v]",  mCC) {
    t.Error("c:33 List should be ", fmt.Sprintf("[%v]",  mCC) , convertListToString(r.ListString("c:33", false)));
  }
  if convertListToString(r.ListString("a:1,c:33", false)) != fmt.Sprintf("[%v]",  mCC) {
    t.Error("a:1,c:33 List should be ", fmt.Sprintf("[%v]",  mCC) , convertListToString(r.ListString("a:1,c:33", false)));
  }
  if convertListToString(r.ListString("a:1,c:*", false)) != fmt.Sprintf("[%v %v]",  mC, mCC) {
    t.Error("a:1,c:* List should be ", fmt.Sprintf("[%v %v]",  mC, mCC) , convertListToString(r.ListString("a:1,c:*", false)));
  }
  if convertListToString(r.ListString("c:*", false)) != fmt.Sprintf("[%v %v]", mC, mCC) {
    t.Error("c:* List should be ", fmt.Sprintf("[%v %v]", mC, mCC) , convertListToString(r.ListString("c:*", false)));
  }

  if convertListToString(r.ListString("a:1", true)) != fmt.Sprintf("[%v]", mA) {
    t.Error("a:1 List Exact should be ", fmt.Sprintf("[%v]", mA) , convertListToString(r.ListString("a:1", true)));
  }
  if convertListToString(r.ListString("a:*", true)) != fmt.Sprintf("[%v]", mA) {
    t.Error("a:* List Exact should be ", fmt.Sprintf("[%v]", mA) , convertListToString(r.ListString("a:*", true)));
  }
  if convertListToString(r.ListString("a:1,b:2", true)) != fmt.Sprintf("[%v]", mB) {
    t.Error("a:1,b:2 List Exact should be ", fmt.Sprintf("[%v]", mB) , convertListToString(r.ListString("a:1,b:2", true)));
  }
  if convertListToString(r.ListString("a:1,b:*", true)) != fmt.Sprintf("[%v]", mB) {
    t.Error("a:1,b:* List Exact should be ", fmt.Sprintf("[%v]", mB) , convertListToString(r.ListString("a:1,b:*", true)));
  }
  if convertListToString(r.ListString("a:1,c:3", true)) != "[]" {
    t.Error("a:1,c:3 List Exact should be []", convertListToString(r.ListString("a:1,c:3", true)));
  }
  if convertListToString(r.ListString("a:1,c:33", true)) != fmt.Sprintf("[%v]", mCC) {
    t.Error("a:1,c:33 List Exact should be ", fmt.Sprintf("[%v]", mCC) , convertListToString(r.ListString("a:1,c:33", true)));
  }
  if convertListToString(r.ListString("a:1,c:*", true)) != fmt.Sprintf("[%v]", mCC) {
    t.Error("a:1,c:* List Exact should be ", fmt.Sprintf("[%v]", mCC) , convertListToString(r.ListString("a:1,c:*", true)));
  }

}


type customGex struct{}
type customModifierGex struct {
	gexers map[string]string
  prevfind patrun.Modifiers
  prevdata interface{}
}

func (a *customGex) Add(pm *patrun.Patrun, pat map[string]string, data interface{}) patrun.Modifiers {
   gexers := map[string]string{}
  for k, v := range pat {
    if strings.Index(v, "*") > -1 {
      gexers[k] = v
      delete(pat, k)
    }
  }

  // handle previous patterns that match this pattern
  var prev = pm.List(pat, false)
  var prevfind patrun.Modifiers
  var prevdata interface{}

  if len(prev) > 0 {
    prevfind = prev[0].Modifier

    prevdata = pm.FindExact(prev[0].Match)
  }

  mod := new(customModifierGex)
	mod.gexers = gexers
  mod.prevfind = prevfind
  mod.prevdata = prevdata

	return mod

}
func (a *customModifierGex) Find(pm *patrun.Patrun, pat map[string]string, data interface{}) interface{} {
  var out = data

  for k, _ := range a.gexers {
    val := pat[k]

    if !gexval(a.gexers[k], val) {
      out = nil
    }
  }

  if a.prevfind != nil && out != nil {
    out = a.prevfind.Find(pm, pat , a.prevdata)
  }

  return out

}
func (a *customModifierGex) Remove(pm *patrun.Patrun, pat map[string]string, data interface{}) bool {
  return true
}

func aaTestCustomGex(t *testing.T) {

  r := patrun.Patrun{Custom: new (customGex)}

  r.AddString( "a:1,b:*", "X")

  if r.FindString("a:1").(string) != "X" {
    t.Error("a:1 Find should be X", r.FindString("a:1"));
  }
  if r.FindString("a:1,b:x").(string) != "X" {
    t.Error("a:1,b:x Find should be X", r.FindString("a:1,b:x"));
  }

  r.AddString( "a:1,b:*,c:q*z", "Y")
  if r.FindString("a:1").(string) != "X" {
    t.Error("a:1 Find should be X", r.FindString("a:1"));
  }
  if r.FindString("a:1,b:x").(string) != "X" {
    t.Error("a:1,b:x Find should be X", r.FindString("a:1,b:x"));
  }
  if r.FindString("a:1,b:x,c:qza").(string) != "Y" {
    t.Error("a:1,b:x,c:qaz Find should be Y", r.FindString("a:1,b:x,c:qaz"));
  }

  r.AddString( "w:1", "W")
  if r.FindString("w:1").(string) != "W" {
    t.Error("w:1 Find should be W", r.FindString("w:1"));
  }
  if r.FindString("w:1,q:x").(string) != "W" {
    t.Error("w:1,q:x Find should be W", r.FindString("w:1,q:x"));
  }


  r.AddString("w:1,q:*", "Q")
  if r.FindString("w:1").(string) != "W" {
    t.Error("w:1 Find should be W", r.FindString("w:1"));
  }
  if r.FindString("w:1,q:x").(string) != "Q" {
    t.Error("w:1,q:x Find should be Q", r.FindString("w:1,q:x"));
  }
  if r.FindString("w:1,q:y").(string) != "W" {
    t.Error("w:1,q:y Find should be Y", r.FindString("w:1,q:y"));
  }

}



func convertListToString(items []patrun.Pattern) string {
  var data []string

  for k := range items {
    v := items[k]
    data = append(data, fmt.Sprintf("{map[%v] %v}", formatMatch(v.Match), fmt.Sprintf("%v", v.Data)))
  }

  return fmt.Sprintf("[%v]", strings.Join(data, " "))
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

  return strings.Join(points, " ")
}


func rs(x patrun.Patrun) string {
  value := x.String()

  r := regexp.MustCompile(`\s+`)

  value = r.ReplaceAllString(value, "")

  r = regexp.MustCompile(`\n+`)
  value = r.ReplaceAllString(value, "")

  r = regexp.MustCompile(`->`)
  value = r.ReplaceAllString(value, "")


  return value
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
