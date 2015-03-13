package main

import (
    "github.com/daaku/go.httpgzip"
    "net/http"
    "image"
    "image/png"
    "image/color"
    "strings"
    "strconv"
    "fmt"
    "math"
)


func handler(w http.ResponseWriter, r *http.Request) {

  // URLs here look like http://localhost:8080/map/13/4100/2724.png
  //                                               z  x    y

  tile_size   := 256
  tile_size_f := float64(256)

  spliturl := strings.Split(r.URL.Path, "/")
  tile_zi, _ := strconv.Atoi(spliturl[2])
  tile_z := float64(tile_zi)
  tile_xi, _ := strconv.Atoi(spliturl[3])
  tile_x := float64(tile_xi)-1
  tile_yi, _ := strconv.Atoi(strings.Split(spliturl[4],".")[0])
  tile_y := float64(tile_yi)-1
  fmt.Printf("Input: %f %f %f\n", tile_z,tile_x,tile_y)

  w.Header().Set("Content-Type","image/png")

  myimage := image.NewRGBA(image.Rectangle{image.Point{0,0},image.Point{tile_size,tile_size}})


  // This loop just fills the image tile with fractal data
  for cx := 0; cx < tile_size; cx++ {
    for cy := 0; cy < tile_size; cy++ {

      cx_f := float64(cx)
      cy_f := float64(cy)

      i := complex128(complex(0,1))

      zoom := float64(math.Pow(2,float64(tile_z-2)))

      tile_range   := 1/zoom
      tile_start_x := 1/zoom + (tile_range*tile_x)
      tile_start_y := 1/zoom + (tile_range*tile_y)

      x := -2 + tile_start_x + (cx_f/tile_size_f)*tile_range
      y := -2 + tile_start_y + (cy_f/tile_size_f)*tile_range

      // x and y are now in the range ~-2 -> +2

      z := complex128(complex(x,0)) + complex128(complex(y,0))*complex128(i)

      c := complex(0.274,0.008)
      for n := 0; n < 100; n++ {
        z = z*z + complex128(c)
      }

      z = z *10
      ratio := float64(2 * (real(z)/2))
      r     := math.Max(0, float64(255*(ratio - 1)))
      b     := math.Max(0, float64(255*(1 - ratio)))
      g     := float64(255 - b - r)
      col := color.RGBA{uint8(r),uint8(g),uint8(b),255}
      myimage.Set(cx,cy,col)
    }
  }

  png.Encode(w, myimage)
}

func main() {
  http.HandleFunc("/map/", handler)
  http.Handle("/", httpgzip.NewHandler(http.FileServer(http.Dir("."))))

  http.ListenAndServe(":8080", nil)
}
