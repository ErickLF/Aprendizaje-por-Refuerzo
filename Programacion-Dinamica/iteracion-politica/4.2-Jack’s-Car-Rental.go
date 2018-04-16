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

type MDP_Jack struct {

  car_loc_1 int
  car_loc_2 int

  //maximo numero de carros en cada locacion
  max_car int

  //maximo numero de carros a trasladar
  max_car_tran int

  //solicitudes de alquiler
  λ_sol_1 int
  λ_sol_2 int

  λ_dev_1 int
  λ_dev_2 int

  γ float64 //tasa de descuento
	θ float64

  r_car int
  r_no_car int

  //faltan los estados
  //acciones legales
}

var problemaJack *MDP_Jack
func NuevoModelo(){
  problemaJack = &MDP_Jack{
    max_car :20,
    λ_sol_1 : 3,
    λ_sol_2 : 4,
    λ_dev_1 : 3,
    λ_dev_2 : 2,
    r_car : 10,
    r_no_car : -2
  }
}
//calcula la probabilidad de autos solicitados y devueltos dependiendo la λ
func  (modelo *MDP_Jack) calcular_prob_car(λ int,int n) (float64){

}
func poisson(λ int ,n int) float64 {
  return ((math.Pow(float64(λ), float64(n)) / float64(factorial(n))) * math.Exp(float64(-λ)) )
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
func main() {
  NuevoModelo()
  fmt.Println("5!", factorial(5))
  fmt.Println("Problema Jack", problemaJack.max_car)
}
