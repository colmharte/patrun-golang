package patrun

import "github.com/colmharte/patrun-golang/patrun"
import "testing"
import "regexp"
import "fmt"

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

    if string(r.ToJSON()[:]) != "[{\"Match\":{},\"Data\":\"R\"}]" {
      t.Error("JSON pattern should be [{\"Match\":{},\"Data\":\"R\"}]", string(r.ToJSON()[:]));
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

    if string(r.ToJSON()[:]) != "[{\"Match\":{},\"Data\":\"R\"},{\"Match\":{\"a\":\"1\"},\"Data\":\"r1\"}]" {
      t.Error("JSON pattern should be [{\"Match\":{},\"Data\":\"R\"},{\"Match\":{\"a\":\"1\"},\"Data\":\"r1\"}]", string(r.ToJSON()[:]));
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

  if string(r.ToJSON()[:]) != "[{\"Match\":{\"a\":\"1\"},\"Data\":\"r1\"}]" {
    t.Error("JSON pattern should be [{\"Match\":{\"a\":\"1\"},\"Data\":\"r1\"}]", string(r.ToJSON()[:]));
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
  if string(r.ToJSON()[:]) != "[{\"Match\":{\"a\":\"1\",\"b\":\"3\"},\"Data\":\"r2\"},{\"Match\":{\"a\":\"1\",\"c\":\"2\"},\"Data\":\"r1\"}]" {
    t.Error("JSON pattern should be [{\"Match\":{\"a\":\"1\",\"b\":\"3\"},\"Data\":\"r2\"},{\"Match\":{\"a\":\"1\",\"c\":\"2\"},\"Data\":\"r1\"}]", string(r.ToJSON()[:]));
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

  if fmt.Sprintf("%v", r.List(nil, false)) != "[{map[a:1] X} {map[b:2] Y}]" {
    t.Error("List all should be [{map[a:1] X} {map[b:2] Y}]", fmt.Sprintf("%v", r.List(nil, false)))
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
  if pat != "[{map[p1:v1] r0}]" {
    t.Error("List p1:v1 should be [{map[p1:v1] r0}]", pat);
  }

  pat = fmt.Sprintf("%v", r.ListString("p1:v1,p2:*", true))
  if pat != "[{map[p1:v1 p2:v2a] r1} {map[p1:v1 p2:v2b] r1}]" {
    t.Error("List p1:v1 should be [{map[p1:v1 p2:v2a] r1} {map[p1:v1 p2:v2b] r1}]", pat);
  }

  r.AddString("p1:v1,p2:v2c,p3:v3a", "r3a")
  r.AddString("p1:v1,p2:v2d,p3:v3b", "r3b")

  pat = fmt.Sprintf("%v", r.ListString("p1:v1,p2:*,p3:v3a", true))
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
  if pat != "[{map[p1:v1] r0}]" {
    t.Error("List p1:v1 should be [{map[p1:v1] r0}]", pat);
  }

  pat = fmt.Sprintf("%v", r.ListString("p1:v1,p2:*", true))
  if pat != "[{map[p1:v1 p2:v2a] r1} {map[p1:v1 p2:v2b] r1}]" {
    t.Error("List p1:v1 should be [{map[p1:v1 p2:v2a] r1} {map[p1:v1 p2:v2b] r1}]", pat);
  }

  r.AddString("p1:v1,p2:v2c,p3:v3a", "r3a")
  r.AddString("p1:v1,p2:v2d,p3:v3b", "r3b")

  pat = fmt.Sprintf("%v", r.ListString("p1:v1,p2:*,p3:v3a", true))
  if pat != "[{map[p1:v1 p2:v2c p3:v3a] r3a}]" {
    t.Error("List p1:v1 should be [{map[p1:v1 p2:v2c p3:v3a] r3a}]", pat);
  }

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
