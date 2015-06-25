package patrun

import (
  "github.com/colmharte/patrun-golang/patrun"
  "testing"
  "fmt"
  "math/rand"
  "time"
)

func setup() (patrun.Patrun, [100]string) {

  r := patrun.Patrun{}

  var k0 [100]string

  for i :=0; i< 100; i++ {
    k0[i] = string(rand.Intn(2000))

    p := map[string]string{}
    p[k0[i]] = k0[i]

    r.Add(p, k0[i])
  }


  // validate data
  for i := 0; i < 100; i++ {


    p := map[string]string{}
    p[k0[i]] = k0[i]
    if r.Find(p).(string) != k0[i] {
      fmt.Println("Data does't validate")
    }

  }


  var I = 300
  var J = 1000

  var bs = time.Now()

  var p1 = patrun.Patrun{}

  var k1x [300]string
  var k1y [1000]string
  for i:=0; i < I; i++  {
    k1x[i] = fmt.Sprintf("x%v", i)
  }

  for j := 0; j < J; j++  {
    k1y[j] = fmt.Sprintf("y%v", j)
  }

  for i := 0; i < I; i++  {
    for j := 0; j < J; j++  {
      var p = map[string]string{}
      p[k1x[i]] = k1x[i]
      p[k1y[j]] = k1y[j]

      p1.Add(p,fmt.Sprintf("%v~%v", k1x[i], k1y[j]))
    }
  }

  var be = time.Now()

  fmt.Println("BUILT: ", be.Sub(bs))

  return p1, k0
}


func BenchmarkSimple(b *testing.B) {
  r, k0 := setup()

  var bs = time.Now()
  // run the function b.N times
  for n := 0; n < b.N; n++ {

    for i := 0; i < 100; i++  {
        for j := 0; j < i; j++ {

          p := map[string]string{}
          p[k0[i]] = k0[i]
          r.Find(p)

        }
    }

  }

  var be = time.Now()

  fmt.Println("EXECUTED: ", be.Sub(bs))

}
