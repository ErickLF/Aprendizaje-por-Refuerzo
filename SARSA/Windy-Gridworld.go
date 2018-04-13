//DESCRIPCION DEL PROBLEMA
//EJEMPLO 6.5 Y EJERCICIOS 6.9 Y 6.10
/*
Utiliza un es un gridworld estándar Pag. 106 libro.
con estados de inicio y meta, pero con una diferencia:
Hay un viento cruzado hacia arriba a través del medio de la cuadricula.

Las acciones son el estándar cuatro arriba, abajo, derecha e izquierda,
pero en la región media los siguientes estados siguientes se desplazan
hacia arriba por un "viento", cuya intensidad varía de columna a columna.

La fuerza del viento se da debajo de cada columna, en el número de celdas
que te desplaza hacia arriba.
Por ejemplo, si está a una celda a la derecha de la meta, entonces la acción
que queda lo lleva a la celda justo arriba de la meta.

Tratemos esto como una tarea episódica no descontada, con recompensas
constantes de -1 hasta que se alcance el estado objetivo.
*/
package main
import (
  "fmt"
  "math/rand"
)

type SARSA struct {
  Q [][]float64

  epsilon float64
  gamma float64
  acciones map[int]string
  alfa float64
  renglones int
  columnas int
}
func (s *SARSA) Inicializar(windy *MDP_WindyGrid, epsilon float64,gamma float64,alfa float64 ){
  s.epsilon = epsilon
  s.gamma = gamma
  s.alfa = alfa
  s.acciones = windy.acciones
  s.columnas = 10
  s.renglones = 7
//inicializamos Q en 0's
  s.Q = make([][]float64, s.renglones*s.columnas)
  for i := 0; i < s.renglones*s.columnas; i++ {
		s.Q[i] = make([]float64, len(windy.acciones))
	}
}

type MDP_WindyGrid struct {
  acciones map[int]string
  viento []int
  estado_ini []int
  estado_fin []int
  //dimensiones del tablero

}
func (problema *MDP_WindyGrid) inicializar() {
    problema.acciones =  map[int]string{0 : "Aarriba" ,1 : "Abajo", 2 : "Izquierda", 3 :"Derecha"}
    problema.viento = []int{0, 0, 0, 1, 1, 1, 2, 2, 1, 0}
    //iniciamos las coordenadas de el estado inicial y final
    problema.estado_ini =  []int{0,3}
    problema.estado_fin = []int{7,3}
    //dimensiones del tablero

}

func (s *SARSA) getIdx(estado []int) int {
	return estado[1]*s.columnas + estado[0]
}

func (s *SARSA) getQ(estado[]int, accion int) float64{
  i := s.getIdx(estado)
  return s.Q[i][accion]
}

func (s *SARSA) setQ(estado[]int, accion int, valor float64){
  i := s.getIdx(estado)
  s.Q[i][accion] = valor
}

func (s *SARSA) epsilon_greedy(estado []int) int{
  var accion int
  if rand.Float64() < s.epsilon{//seleccionamos una accion al azar si es menor al epsilon
     accion = rand.Intn(len(s.acciones)) //una accion de 0-4
  }else{ // si no buscamos la accion con mayor valor de ese estado
      i := s.getIdx(estado)
	    max := s.Q[i][0] //nos posicionamos en ese estado
	    accion = 0
    	for j := 0; j < len(s.acciones); j++ {//buscamos la accion con mayor valor
    		if max < s.Q[i][j] {
    			max = s.Q[i][j]
    			accion = j
    		}
    	}
  }
  return accion
}

func (s *SARSA) tomarAccion(st []int, a int, viento []int, sf []int) (float64,[]int){
  estado := []int{st[0],st[1]}
  accion := a

  //checamos en que estado vamos a caer despues de tomar una accion
  switch accion {
  case 0://Arriba
  //Checando que no se salga del techo
    if st[1] > 0 {
      estado[1] = st[1] - 1 //restamos -1 porque nuestro 00 esta en la esquina izq superior
    }
  case 1://Abajo
  //Checando que no se pase del piso
    if st[1] != s.renglones-1 {
      estado[1] = st[1] + 1
    }
  case 2://Izquierda
    if st[0] > 0{
      estado[0] = st[0] - 1
    }
  case 3://Derecha
    //evitando que se salga del mapa a la derecha
    if st[0] != s.columnas-1{
      estado[0] = st[0] + 1
    }
  }
  //Aqui funciona si la suma del viento cae exactamente en el estado final, mas no si pasa
  estado[1] -= viento[estado[0]]
	if estado[1] < 0 {
    //regresamos al borde del mundo
		estado[1] = 0
	}
  //fmt.Println(st,s.acciones[accion],estado)
	if estado[0] == sf[0] && estado[1] == sf[1] {
		return 0.0, estado
	}
  //-1 por cada accion y el estado siguiente
	return -1.0, estado

}
func Sarsa(problema *MDP_WindyGrid, max_ep int){
  S := SARSA{}
  //Inicializamos el estado inicial
  S.Inicializar(problema,0.1,1,0.5)//epsilon, gamma, alfa
  for i:=0; i < max_ep; i++{
    fmt.Println(i)
    estado := problema.estado_ini
    accion := S.epsilon_greedy(estado)//seleccionamos una Accion mediante epsilon greedy
    for estado[0] != problema.estado_fin[0] || estado[1] != problema.estado_fin[1]{
      //Tomar una accion viendo la recompensa y el siguiente estado
      r,estado_sig := S.tomarAccion(estado,accion,problema.viento,problema.estado_fin)
      accion_sig := S.epsilon_greedy(estado_sig)
      QSA := S.getQ(estado,accion)
      QSA_sig := S.getQ(estado_sig,accion_sig)
      fmt.Println(estado,"-->",S.acciones[accion],"---",estado_sig,"-->",S.acciones[accion_sig])
      //Actualizamos
      QSA = QSA + S.alfa*(r+S.gamma*QSA_sig - QSA)
      S.setQ(estado,accion,QSA)
      estado = estado_sig
      accion = accion_sig
    }
  }
  fmt.Println(S.Q)
}

func main() {
  fmt.Println("Problema Windy Grid World")
  problema := MDP_WindyGrid{}
  problema.inicializar()
  Sarsa(&problema,1000)
}
