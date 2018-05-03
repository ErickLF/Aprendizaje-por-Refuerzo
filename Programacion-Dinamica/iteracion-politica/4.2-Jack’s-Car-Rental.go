/*
DESCRIPCION DEL PROBLEMA:
Jack administra dos ubicaciones para una empresa de alquiler de automóviles
en todo el país.
Cada día, algunos clientes llegan a cada locación para alquilar automóviles.
Si Jack tiene un auto disponible, éllo alquila y la empresa nacional
le acredita $ 10.
Si él no tiene autos en ese lugar, entonces el negocio se pierde.
Los automóviles están disponibles para alquilar el día posterior a su devolución.
Para ayudar a garantizar que los automóviles estén disponibles donde se necesiten,
Jack puede trasladarlos entre las dos ubicaciones de la noche a la mañana,
a un costo de $ 2 por automóvil trasladado.

Suponemos que la cantidad de autos solicitados y devueltos en cada ubicación
son variables aleatorias de Poisson, lo que significa que la probabilidad de que
el número sea n = (λ^n)/n!*(e^-λ), donde λ es el número esperado.
Supongamos que λ es 3 y 4 para las solicitudes de alquiler en la primera y segunda
ubicación y 3 y 2 para las devoluciones.

Para simplificar un poco el problema, suponemos que no puede haber más de 20
automóviles en cada ubicación (todos los automóviles adicionales se devuelven
a la empresa nacional, y por lo tanto desaparecen del problema) y se puede trasladar
un máximo de 5 automóviles de una ubicación al otro en una noche.

Consideramos que la tasa de descuento es γ = 0.9 y formulamos esto como un MDP
finito continuo, donde los pasos de tiempo son días, el estado es el número de
automóviles en cada ubicación al final del día y las acciones son los números
netos de automóviles se movieron entre los dos lugares durante la noche.
*/

package main
import (
  "fmt"
  "math"
)

type Poisson struct{
  tasaTransferencia[]float64
}

func (po *Poisson) Poisson(lambda float64, max int){
  po.tasaTransferencia = make([]float64, max)

		for i := 0; i < len(po.tasaTransferencia); i++ {
      tasa:= ((math.Pow(float64(lambda), float64(i)) / float64(factorial(i))) * math.Exp(float64(-lambda)) )
			po.tasaTransferencia[i] = tasa;
		}
}
func (po *Poisson) getTransf(maxCar int, NumDisponibles int)float64{
  if maxCar >= NumDisponibles{
    suma := 0.0
    for i := 0; i < maxCar; i++{
        suma +=po.tasaTransferencia[i]
    }
    return 1.0 -suma
  }else{
    return po.tasaTransferencia[maxCar]
  }
}

type Estado struct{
  carA int
  carB int
  V float64
  accion int
}
func (estado *Estado)init(A int, B int){
  estado.carA = A
  estado.carB = B
  estado.V = 0.0
  estado.accion = 0
}

func (estado *Estado) maxAccion(problema *MDP_Jack) int{
    carDisponiblesDeA := estado.carA;
		carDisponiblesParaB := problema.maxCar - estado.carB;
		moviendo := math.Min(float64(carDisponiblesDeA), float64(carDisponiblesParaB));
		return int(math.Min(moviendo, float64(problema.maxMov)));
}

func (estado *Estado)minAccion(problema *MDP_Jack) int{
    carDisponiblesDeB := estado.carB;
		carDisponiblesParaA := problema.maxCar - estado.carA;
		moviendo := math.Min(float64(carDisponiblesDeB), float64(carDisponiblesParaA));
		return -int(math.Min(moviendo, float64(problema.maxMov)));
}


type MDP_Jack struct {

  maxCar int //maximo numero de carros en cada locacion
  maxMov int //maximo numero de carros a trasladar

  //solicitudes de alquiler
  lambdaS_A float64
  lambdaS_B float64
  //Devoluciones
  lambdaD_A float64
  lambdaD_B float64
  gamma float64 //tasa de descuento
  delta float64

  recompensa int //precio de rentar un auto
  costo_transferencia int

  estado [][]Estado
  poisson_S_A Poisson
  poisson_D_A Poisson
  poisson_S_B Poisson
  poisson_D_B Poisson


}


func (problema *MDP_Jack) Iniciar(){
    problema.maxCar = 20
    problema.maxMov = 5
    problema.recompensa = 10
    problema.lambdaS_A = 3.0
    problema.lambdaS_B = 4.0
    problema.lambdaD_A = 3.0
    problema.lambdaD_B = 2.0
    problema.costo_transferencia = 2
    problema.gamma = 0.9
    problema.delta = 1.0
    problema.estado = make([][]Estado, problema.maxCar+1)
    for i := 0; i < problema.maxCar+1; i++ {
  		problema.estado[i] = make([]Estado, problema.maxCar+1)
  	}

    for j:= 0; j < problema.maxCar + 1; j++ {
			for k:= 0; k < problema.maxCar + 1; k++ {
				problema.estado[j][k].init(j,k)
			}
		}
    problema.poisson_S_A.Poisson(problema.lambdaS_A,problema.maxCar+1)
    problema.poisson_D_A.Poisson(problema.lambdaD_A,problema.maxCar+1)
    problema.poisson_S_B.Poisson(problema.lambdaS_B,problema.maxCar+1)
    problema.poisson_D_B.Poisson(problema.lambdaD_B,problema.maxCar+1)
    //Sfmt.Println(problema.estado)

}

//para calcular la funccion de probabilidad de poisson
func factorial(n int) (uint64) {
  var fact uint64 = 1; //para especificar que es uint64 lo hacemos explicito
  if n < 0{
  fmt.Println("Factorial negativo")
  }else{
    for i := 1; i <= n; i++ {
      fact *= uint64(i) //para que sean del mismo tipo
    }
  }
  return fact
}

func (problema *MDP_Jack) estimarV(CarA int, CarB int) float64{
  val := 0.0
  for clientesA:=0; clientesA <= CarA; clientesA++{
    tasaA := problema.poisson_S_A.getTransf(clientesA, CarA)
    espacioA := problema.maxCar - (CarA - clientesA)

    for regresoA:=0; regresoA <= espacioA; regresoA++{
      tasaRegresoA := problema.poisson_D_A.getTransf(regresoA, espacioA)

      for clientesB:=0; clientesB <= CarB; clientesB++{
        tasaB := problema.poisson_S_B.getTransf(clientesB, CarB)
        espacioB := problema.maxCar - (CarB - clientesB)

        for regresoB:=0; regresoB <= espacioB; regresoB++{
          tasaRegresoB := problema.poisson_D_B.getTransf(regresoB, espacioB)

          tasa := tasaA * tasaB * tasaRegresoA * tasaRegresoB

          r := (clientesA + clientesB) * problema.recompensa

          nuevoCarA := CarA - clientesA +regresoA
          nuevoCarB := CarB - clientesB + regresoB
          estado := problema.estado[nuevoCarA][nuevoCarB]
          dval := tasa *(float64(r) + problema.gamma*estado.V)
          val += dval
        }
      }
    }
  }
  return val
}
func (problema *MDP_Jack) estimarAccionValor(estado Estado, accion int) float64{
  costoMov := -math.Abs(float64(accion))*float64(problema.costo_transferencia)
  carA := (estado.carA - accion)
  carB := estado.carB + accion
  v := problema.estimarV(carA,carB) + costoMov
  return v
}

func (problema *MDP_Jack) iteracionPolitica(maxE int){
  V := make([][]float64, problema.maxCar+1)
  for i := 0; i < problema.maxCar+1; i++ {
    V[i] = make([]float64, problema.maxCar+1)
  }
  max := maxE
  for episodios:=0; episodios < max; episodios++{
    //Actualizamos valores
    for true{
      for i:= 0; i < problema.maxCar + 1; i++ {
        for j:= 0; j < problema.maxCar + 1; j++ {
          s:= problema.estado[i][j]
          v := problema.estimarAccionValor(s,s.accion)
          V[i][j] = v
        }
      }
      maxDeltaV := -99999999999999999999.0

      for i:= 0; i < problema.maxCar + 1; i++ {
        for j:= 0; j < problema.maxCar + 1; j++ {
          s:= problema.estado[i][j]
          deltaV := math.Abs(V[i][j] - s.V)
          if deltaV > maxDeltaV {
            maxDeltaV = deltaV;
          }
          problema.estado[i][j].V = V[i][j]
        }
      }
      //fmt.Println(maxDeltaV)
      if maxDeltaV < problema.delta {
        break;
      }
    }
    //Actualizamos Politica
    for i:= 0; i < problema.maxCar + 1; i++ {
      for j:= 0; j < problema.maxCar + 1; j++ {
        s:= problema.estado[i][j]
        movCarMin := s.minAccion(problema)
        movCarMax := s.maxAccion(problema)
        mejorAccion := -1
        mejorValor := -9999999999999.0

        for accion:=movCarMin; accion <= movCarMax; accion++{
          v:=problema.estimarAccionValor(s, accion)
          if v > mejorValor{
            mejorValor = v
            mejorAccion = accion
          }
        }
        problema.estado[i][j].accion = mejorAccion
      }
    }
  }
}
func(problema *MDP_Jack) ImprimirA(){
  for i:= 0; i < problema.maxCar + 1; i++ {
    for j:= 0; j < problema.maxCar + 1; j++ {
      fmt.Print(math.Abs(float64(problema.estado[i][j].accion)), " ")
    }
    fmt.Println("")
  }
}

func(problema *MDP_Jack) ImprimirV(){
  for i:= 0; i < problema.maxCar + 1; i++ {
    for j:= 0; j < problema.maxCar + 1; j++ {
      fmt.Print(problema.estado[i][j].V, " ")
    }
    fmt.Println("")
  }
}

func main() {
  problema := MDP_Jack{}
  problema.Iniciar()
  fmt.Println("Problema Jack")
  //prueba para 4 iteraciones
  problema.iteracionPolitica(4)


  fmt.Println("====Acciones====")
  problema.ImprimirA()
  fmt.Println("====V====")
  problema.ImprimirV()

}
