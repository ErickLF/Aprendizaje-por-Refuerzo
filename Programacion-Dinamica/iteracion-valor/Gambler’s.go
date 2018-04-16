/*
Un jugador tiene la oportunidad de hacer apuestas sobre los resultados
de una secuencia de lanzamientos de monedas.
Si la moneda sale cara, gana tantos dólares como ha apostado en ese lanzamiento;
si se trata de cruz, pierde su apuesta.
El juego termina cuando el jugador gana al alcanzar su meta de $ 100,
o pierde por quedarse sin dinero.

En cada lanzamiento, el jugador debe decidir qué parte de su capital apostar,
en números enteros de dólares.
Este problema se puede formular como un MDP finito no episódico, sin descontar.
El estado es el capital del jugador, s ∈ {1, 2,. . . , 99} y
las acciones son en juego, a ∈ {0, 1,. . . , min (s, 100 - s)}.
La recompensa es cero en todas las transiciones excepto aquellas en las que
el jugador alcanza su objetivo, cuando es +1.

La función de valor de estado da la probabilidad de ganar de cada estado.
Una política es un mapeo desde niveles de capital hasta apuestas.
La política óptima maximiza la probabilidad de alcanzar el objetivo.

ph la probabilidad de que la moneda salga cara.

Si se conoce ph, entonces se conoce todo el problema y se puede resolver,
por ejemplo, mediante iteración de valores.
*/
package main
import (
  "fmt"
  "math"
)
const (
  meta = 100
  gamma = 1
  ph = 0.25
  theta = 0.00000001
)
type MDP_Jugador struct{
}

var problema *MDP_Jugador


func (problema *MDP_Jugador) r(s int) int {
  //recompensa en ese estado
  if s == meta{
    return 1
  }
  return 0
}

func iteracion_valor(problema *MDP_Jugador) {
  V := make([]float64, meta+1)   //inicializamos los valores de V en 0
  pi := make([]float64, meta+1) //inicializamos los valores de la politica en
  delta := 1.0
  for delta > theta{
      delta = 0
      for s:= 1; s < meta; s++{
        v := V[s] //valor anterior
        bellmant(problema,s,V,pi)
        diff := math.Abs(v-V[s])
        delta = math.Max(delta,diff)
      }
  }
  fmt.Println(pi)
}

func bellmant(problema *MDP_Jugador, s int, V []float64, pi []float64){
  valor_optimo := 0.0
  for x := 0; x <= int(math.Min(float64(s),float64(meta-s))); x++{
    gan:= s + x
    per := s - x
    resultado := ph * (float64(problema.r(gan)) + gamma * V[gan]) + ((1 - ph) * (float64(problema.r(per)) + gamma * V[per]))
    //seleccionamos la accion que da mas recompensa
    if resultado > valor_optimo {
      valor_optimo = resultado
      V[s] = resultado //guarda la mejor recompensa
      pi[s] =float64(x) //guarda la mejor accion
    }
  }

}
func main() {
  iteracion_valor(problema)
  fmt.Println("Problema Gambler's")
}
