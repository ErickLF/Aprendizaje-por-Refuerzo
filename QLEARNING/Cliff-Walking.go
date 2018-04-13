/*
Esta es una tarea episódica no descontada estándar, con inicio y
estados de objetivo y las acciones habituales que causan movimiento hacia arriba,
hacia abajo, hacia la derecha y hacia la izquierda. La recompensa es -1 en todos
las transiciones, excepto aquellas en la región marcada como "The Cliff".
Ingresar en esta región implica una recompensade -100 y envía al agente
instantáneamente de vuelta al inicio.
*/

package main
import (
  "fmt"
  "math/rand"
)

type QL struct{
  Q [][]float64
  renglones int
  columnas int
  acciones map[int]string
  epsilon float64
  gamma float64
  alfa float64

}

func (q *QL) Inicializar(problema *MDP, epsilon float64,gamma float64,alfa float64 ){
  q.epsilon = epsilon
  q.gamma = gamma
  q.alfa = alfa
  q.acciones = problema.acciones
  q.columnas = 12
  q.renglones = 4
//inicializamos Q en 0's
  q.Q = make([][]float64, q.renglones*q.columnas)
  for i := 0; i < q.renglones*q.columnas; i++ {
		q.Q[i] = make([]float64, len(problema.acciones))
	}
}
func (q *QL) getQ(estado[]int, accion int) float64{
  i := q.getIdx(estado)
  return q.Q[i][accion]
}

func (q *QL) setQ(estado[]int, accion int, valor float64){
  i := q.getIdx(estado)
  q.Q[i][accion] = valor
}


func (q *QL) tomarAccion(problema *MDP, s[]int, a int) (float64,[]int){
  estado := []int{s[0],s[1]}
  accion := a

  //checamos en que estado vamos a caer despues de tomar una accion
  switch accion {
  case 0://Arriba
  //Checando que no se salga del techo
    if s[1] > 0 {
      estado[1] = s[1] - 1 //restamos -1 porque nuestro 00 esta en la esquina izq superior
    }
  case 1://Abajo
  //Checando que no se pase del piso
    if s[1] != q.renglones-1 {
      estado[1] = s[1] + 1
    }
  case 2://Izquierda
    if s[0] > 0{
      estado[0] = s[0] - 1
    }
  case 3://Derecha
    //evitando que se salga del mapa a la derecha
    if s[0] != q.columnas-1{
      estado[0] = s[0] + 1
    }
  }

  if estado[0] == problema.estado_fin[0] && estado[1] == problema.estado_fin[1] {
		return 0.0, estado
	}
  r := problema.r(estado)
  if(r == -100.0){
    estado[0]=problema.estado_ini[0]
    estado[1]=problema.estado_ini[1]
    return r,estado
  }
  //-1 por cada accion y el estado siguiente
	return -1.0, estado

}

func Qlearning(problema *MDP, max_ep int){
  q := QL{}
  q.Inicializar(problema,0.1,1.0,0.5)
  for i := 0 ; i < max_ep; i++ {
    estado := problema.estado_ini
    fmt.Println(i)
    for estado[0] != problema.estado_fin[0] || estado[1] != problema.estado_fin[1]{ //mientras no lleguemos a la meta
      accion := q.epsilon_greedy(estado)
      r,estado_sig := q.tomarAccion(problema,estado,accion)
      QSA := q.getQ(estado,accion)
      maxAccion := q.maxAccion(estado)
      QSA_sig := q.getQ(estado_sig,maxAccion)//elegimos la maxima accion por si el egreedy toma una accion al azar
      QSA = QSA + q.alfa*(r+q.gamma*QSA_sig - QSA)
      fmt.Println(estado,"-->",q.acciones[accion],"---",estado_sig)
      q.setQ(estado,accion,QSA)
      estado = estado_sig
    }
  }
}

func (q *QL) getIdx(estado []int) int { //regresa un indice
    //posicion y*num_columnas + desplazamiento en x
  return estado[1]*q.columnas + estado[0]
}

func (q *QL) epsilon_greedy(estado []int) (int){//regresa una accion
  var accion int
  if rand.Float64() < q.epsilon{//seleccionamos una accion al azar si es menor a epsilon
     accion = rand.Intn(len(q.acciones)) //regresamos una accion de 0-4
  }else{ // si no buscamos la accion con mayor valor de ese estado
      i := q.getIdx(estado)
	    max := q.Q[i][0] //nos posicionamos en ese estado y agarramos su primer valor que da esa accion
	    accion = 0
    	for j := 0; j < len(q.acciones); j++ {//buscamos la accion con mayor valor
    		if max < q.Q[i][j] {
    			max = q.Q[i][j]
    			accion = j //seleccionamos la mejor accion
    		}
    	}
  }
  return accion
}
func(q *QL) maxAccion(estado []int) int{
  i := q.getIdx(estado)
  max := q.Q[i][0] //nos posicionamos en ese estado y agarramos su primer valor que da esa accion
  accion := 0
  for j := 0; j < len(q.acciones); j++ {//buscamos la accion con mayor valor
    if max < q.Q[i][j] {
      max = q.Q[i][j]
      accion = j //seleccionamos la mejor accion
    }
  }
  return accion
}

type MDP struct{
  acciones map[int]string
  estado_ini []int
  estado_fin []int
  //cliff []int //sona de perdición
}

func (problema *MDP) inicializar(){
  problema.acciones = map[int]string{0 : "Aarriba" ,1 : "Abajo", 2 : "Izquierda", 3 :"Derecha"}
  problema.estado_ini = []int{0,3}
  problema.estado_fin = []int{11,3}
  //cliff =[]int{1,2,3,4,5,6,7,8,9,10}//en x
}

func (problema *MDP) r(estado []int) float64{
  if estado[1]==problema.estado_ini[1]{
    for i:=problema.estado_ini[0]+1; i < problema.estado_fin[0]; i++{
      if estado[0]==i{
        return -100.0
      }
    }
  }
  return -1.0
}

func main() {
  fmt.Println("Problema Cliff Walking")
  problema := MDP{}
  problema.inicializar()
  Qlearning(&problema,100)
}
