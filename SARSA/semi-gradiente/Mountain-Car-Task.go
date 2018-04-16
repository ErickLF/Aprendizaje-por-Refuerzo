/*
Tile Coding Hashing

*/
package  main
import (
  "fmt"
  "math"
  "math/rand"
  "time"

)
type FuncionValor struct{
  W  []float64
	Mallas  int
	Caracteristicas    int
	Alfa float64
}
func (V *FuncionValor) inicializar(alfa float64, max_tam int, num_mallas int, caracteristicas int){
  V.Alfa = alfa
  V.Mallas =num_mallas
  V.W = make([]float64,max_tam)
  V.Caracteristicas = caracteristicas
}



type Estado struct{
  posicion    float64
  velocidad   float64
}

func (s *Estado) iniciar(){
  s.posicion = -0.5
  s.velocidad = 0
}
func (s *Estado) posicionRandom(){
  s.posicion = -0.5
  s.velocidad = 0
}

type MDP struct{
  estado *Estado
  acciones map[int]string
  pos_min float64
  pos_max float64 //meta
  vel_min float64
  vel_max float64
  V  *FuncionValor
  max_tam int
  escalaV float64
  escalaP float64
}
func (p *MDP) accionRandom() int{
  return rand.Intn(len(p.acciones))
}
func (p *MDP) inicializar(){
  p.pos_min  =  -1.2
  p.pos_max  =   0.5
  p.vel_min  =  -0.07
  p.vel_max =   0.07
  estado := Estado{}
  estado.iniciar()
  V := FuncionValor{}
  V.inicializar(0.4/8, 2048, 8, 1)
  //s.hash_table = make(map[string]int)
	p.max_tam = 2048
	p.escalaP = float64(8) / (p.pos_max - p.pos_min)
	p.escalaV = float64(8) / (p.vel_max - p.vel_min)
	p.acciones = map[int]string{-1:"reversa",0:"nada",1:"acelerar"}
}
func tomarAccion(p *MDP, estado *Estado, accion int) (int,Estado){
  velocidad := estado.velocidad + 0.001*float64(accion) - 0.0025*math.Cos(3*estado.posicion)
  //que respete las leyes del juego
  if velocidad > p.vel_max{
    velocidad = p.vel_max
  } else if velocidad < p.vel_min{
    velocidad = p.vel_min
  }
  posicion := estado.posicion + velocidad
  if posicion > p.pos_max{
    posicion = p.pos_max
  }else  if posicion < p.pos_min{
    posicion = p.pos_min
    velocidad = 0.0
  }
  s := Estado{velocidad:velocidad,posicion:posicion}
  return -1,s
}
func SarsaSemiGradiente(p *MDP){
  rand.Seed(time.Now().Unix())

	estado := Estado{}
  estado.posicionRandom()
  accion := p.accionRandom()

	episodios := 100
	for i := 0; i < episodios; i++ {
    paso := 0
    for estado.posicion < p.pos_max {
      paso++
      r,nuevoEstado := tomarAccion(p,estado,accion)
	   }
  }
}

func getActiveTiles(action string) []int {
	//fmt.Println("Hello: ", s.position, s.velocity, s.posScale)
	_pos := math.Floor(s.position * s.posScale * float64(s.v.Tilings))
	//fmt.Println("macho")
	_vel := math.Floor(s.velocity * s.velScale * float64(s.v.Tilings))
	tiles := make([]int, 0)

	for tile := 0; tile < s.v.Tilings; tile++ {

		div := math.Floor(((_pos + float64(tile)) / float64(s.v.Tilings))) //

		div2 := math.Floor(((_vel + 3*float64(tile)) / float64(s.v.Tilings)))

		key.WriteString(action)

		tiles = append(tiles, s.Idx(key.String()))
	}
	return tiles
}

func main(){
  problema := MDP{}
  problema.inicializar()

  fmt.Println("GO!")
}
