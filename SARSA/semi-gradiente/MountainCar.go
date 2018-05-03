package main
import (
  "fmt"
  "math/rand"
  "time"
  "math"
)
/*
Para la realizacion de este programa se tomo como base la siguiente libreta
https://github.com/edjeffery/RL_Mountain_Car/blob/master/rlc3_mountain_car.ipynb
*/
type Mountain struct{
  posicion_min float64
  posicion_max float64
  velocidad_min float64
  velocidad_max float64
  acciones map[int]int
  gameOver bool
  posicion float64
  velocidad float64

}

func(this *Mountain) init(){
  this.posicion_min = -1.2
  this.posicion_max = 0.5
  this.velocidad_max = 0.7
  this.velocidad_min = -0.7
  this.acciones = map[int]int{0:-1,1:0, 2:1}
}

func (this *Mountain) reset(){
    this.gameOver = false
    rand.Seed(time.Now().UnixNano())
    //this.posicion = (rand.Float64() * this.posicion_min) + this.posicion_max
    //this.velocidad =(rand.Float64() * this.velocidad_min) + this.velocidad_max
    this.posicion = 0.85
    this.posicion = 0.2

}
func (this *Mountain) paso(accion int) float64{
  this.velocidad = this.velocidad + 0.001 * float64(accion) - 0.0025 * math.Cos(3*this.posicion)

  if this.velocidad > this.velocidad_max{
    this.velocidad = this.velocidad_max
  } else if this.velocidad < this.velocidad_min{
    this.velocidad = this.velocidad_min
  }
  this.posicion = this.posicion + this.velocidad
  if this.posicion > this.posicion_max{
    this.posicion = this.posicion_max
  }else  if this.posicion < this.posicion_min{
    this.posicion = this.posicion
    this.velocidad = 0.0
  }
  if this.posicion == this.posicion_max{
    this.gameOver = true
  }
  reward := -1.0
  return reward
}
type Matriz struct{
  ren [][]float64
}

func (m *Matriz) crear(tiles int){
  m.ren = make([][]float64, 2)
  for i := 0; i < 2; i++ {
		m.ren[i] = make([]float64, tiles)
	}
}

func (m *Matriz) modificar(x []float64, y []float64){
  for i := 0; i < 2; i++ {
    for j:=0; j< len(x); j++{
      if(i==0){
        m.ren[i][j] = x[j]
      }else{
        m.ren[i][j] = y[j]
      }
    }
	}
}

func (m *Matriz) getRen(i int)[]float64{
		return m.ren[i]

}


func linspace(min float64, max float64, n int, resta float64)[]float64{
  val := math.Abs(min-max)
  div := val/float64(n-1)
  var s []float64
  for i:=0; i <n; i++{
      if(i==0){
        s = append(s,min-resta)
      }else{
        s = append(s,(s[i-1]+div))
      }
  }
  return s
}
func CrearTiling()[]Matriz{

  num_tiles := 10
  min_x := -1.2
  max_x := 0.5
  x_tile_width := (max_x - min_x) / float64((num_tiles - 2))
  min_y := -0.07
  max_y := 0.07
  y_tile_width := (max_y - min_y) / float64((num_tiles - 2))
  //tilings = np.zeros((10, 2, 10))
  tilings := make([]Matriz, num_tiles)
  for i := 0; i < num_tiles; i++ {
		tilings[i].crear(num_tiles)
	}
  for i:=0; i<len(tilings); i++{
    x_offset := rand.Float64() + x_tile_width
    y_offset := rand.Float64() + y_tile_width
    xs := linspace(min_x,max_x + x_tile_width,num_tiles, x_offset)
    ys := linspace(min_y,max_y + y_tile_width,num_tiles, y_offset)
    tilings[i].modificar(xs,ys)
  }

  return tilings
}
func digitize(val float64, arr []float64)int{
  for i:=0; i< len(arr);i++{
    if val < arr[i]{
      return i
    }
  }
  return len(arr)-1
}

func genIndices(estado []float64, acc int, tilings[]Matriz)[]int32{
    posicion := estado[0]
    velocidad := estado[1]

    accion := acc

    num_tilings := len(tilings)
    tiling_length := 10
    tiling_height := 10
    num_tiles := num_tilings * tiling_length * tiling_height

    tiles := make([]int32, 10)
    for i:=0; i < num_tilings; i++{
      xs := tilings[i].getRen(0)
      ys := tilings[i].getRen(1)
      xi := digitize(posicion, xs)
      yi := digitize(velocidad, ys)
      index := (accion * num_tiles) + (i * tiling_length * tiling_height) + xi + (yi * tiling_height)
      tiles[i] = int32(index)
    }
    return tiles
}

func calcQ(F []int32, theta []float64) float64{
    Qa := 0.0
    for i:=0; i < len(F); i++{
        Qa = Qa + theta[i]
    }
    return Qa
}

func max(Q []float64) float64{
  max := Q[0]
  for i:=1; i < len(Q); i++{
    if(Q[i] > max){
      max = Q[i]
    }
  }
  return max
}
func sum(Q []float64, Qmax float64)int{
 cont:=0
  for i:=0;i <len(Q);i++{
    if Q[i] == Qmax{
      cont++
    }
  }
  return cont
}
func argmax(Q []float64)int{
  ind := 0
  max:=Q[0]
  for i:=1; i<len(Q);i++{
    if(Q[i]> max){
      ind = i
    }
  }
  return ind
}
func mejor(Q []float64, maxQ float64, acciones map[int]int)[]int{
  var s []int
  for i:=0; i < len(acciones); i++{
    if(Q[i] == maxQ){
      s = append(s,i)
    }
  }
  return s
}
func e_greedy_accion(theta []float64, s []float64, acciones map[int]int, tilings[]Matriz)(int,float64){
  Q := make([]float64, 3)
  Qa:= 0.0
  for a:=0; a<len(acciones);a++{
    F := genIndices(s, a, tilings)
    Qa = calcQ(F, theta)
    Q[a] = Qa
  }
  maxQ := max(Q)
  i:=0
  if sum(Q,maxQ) > 1{
      best := mejor(Q,maxQ,acciones)
      i = rand.Intn(len(best))
  }else{
      i = argmax(Q)
  }
  accion := acciones[i]
  Qa = Q[i]
  return accion, Qa
}

func Play(car *Mountain,num_episodes int, alfa float64, gamma float64, epsilon float64, lambda float64)[]float64{
  theta := make([]float64, 3000)// 10 x 10 x 10 x 3
  recompensa_por_episodio := make([]float64, num_episodes)

  tilings := CrearTiling()
  for episodio :=0; episodio < num_episodes; episodio++{
    rewardAcomulado := 0.0
    paso := 0
    e := make([]float64, 3000)
    car.reset()

    estado := []float64{car.posicion, car.velocidad}
    accion := rand.Intn(len(car.acciones))
    for true{
      if(car.gameOver){
        break
      }
      F := genIndices(estado, accion, tilings)
      for j:=0;j<len(F); j++{
        e[j] = 1
      }
      reward := car.paso(accion)
      estado_ := []float64{car.posicion, car.velocidad}

      Qa := calcQ(F, theta)
      delta := reward - Qa

      //Epsilon greedy siguiente accion
      Qa_ := 0.0
      accion_ := 0
      if rand.Float64() < (1 -epsilon) {
        accion_, Qa_ = e_greedy_accion(theta, estado_, car.acciones, tilings)
      }else{
        accion_ = rand.Intn(len(car.acciones))
        F = genIndices(estado_, accion_, tilings)
        Qa_ = calcQ(F, theta)
      }
      delta = delta + gamma * Qa_
      for k:=0; k < len(theta); k++{
        theta[k] = theta[k]+alfa* e[k]*delta
      }
      for k:=0; k < len(e); k++{
        e[k] = gamma*lambda*e[k]
      }
      copy(estado, estado_)
      accion = accion_
      paso++
      rewardAcomulado += reward
      recompensa_por_episodio[episodio] = rewardAcomulado
    }
  }
  return recompensa_por_episodio
}

func main() {
  agentes := 20
  num_eps:= 150
  //recompensas := make([]float64,num_eps)
  for i :=0; i < agentes; i++{
    car := Mountain{}
    car.init()
    r := Play(&car, num_eps,0.05,1,0.05,.90)
    //recompensas[i] = r
    fmt.Println("Agente", i+1, "termino: ",r[num_eps-1])
  }
  fmt.Println("Mountain")
}
