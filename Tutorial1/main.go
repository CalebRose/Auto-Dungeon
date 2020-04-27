package main

import (
	"fmt"
	"github.com/calebrose/Auto-Dungeon/Tutorial1/strutil"
	"strconv"
	"math"
)

func greeting(name string) string {
	return "Hello, " + name + "!"
}

func getSum(a, b int) int {
	return a + b
}

func arrays() {
	var fruitArr [2] string

	// Assign values
	fruitArr[0] = "apple"
	fruitArr[1] = "orange"

	fmt.Println(fruitArr)
	fmt.Println(fruitArr[1])

	// Declare & assign
	fruitArr2 := [2]string {"banana", "peach"}
	fmt.Println(fruitArr2)

	fruitSlice := []string {"grape", "pear", "mango", "cherry"}
	fmt.Println(fruitSlice)
	fmt.Println(len(fruitSlice))
	fmt.Println(fruitSlice[1:3])
}

func Conditionals () {
	x := 10
	y := 10

	if x <= y {
		fmt.Printf("%d is less than / equal to %d\n", x, y)
	} else {
		fmt.Printf("%d is less than %d\n", y, x)
	}

	color := "green"
	switch color {
	case "red":
		fmt.Println("Color is red")
	case "blue":
		fmt.Println("Color is blue")
	default:
		fmt.Println("I DON'T KNOW THAT")
	}

	if color == "red" {
		fmt.Println("Color is red")
	} else if color == "blue" {
		fmt.Println("Color is blue")
	} else {
		fmt.Println("I don't recognize that color")
	}
}

func loops(){
	// Long Method
	// i := 1
	// for i <= 10 {
	// 	fmt.Println(i)
	// 	i++
	// }

	// for j := 0; j < 10; j++ {
	// 	fmt.Println("THIS IS NOT OVER 9000")
	// }

	// fizz buzz
	for k := 0; k <= 100; k++ {
		if k % 15 ==0 {
			fmt.Println("Fizz Buzz")
		} else if k % 3 == 0{
			fmt.Println("Fizz")
		} else if k % 5 == 0 {
			fmt.Println("Buzz")
		} else {
			fmt.Println(k)
		}
	}
}

func maps() {
	// emails := make(map[string]string)
	// // emails["Bob"] = "bob@gmail.com"
	// // emails["Pat"] = "pat@gmail.com"
	// // emails["Mike"] = "mike@gmail.com"
	// // fmt.Println(emails)

	// // // Delete
	// // delete(emails, "Bob")
	// Declare map
	emails := map[string]string{"Bob":"bob@gmail.com", "Mike":"mike@gmail.com"}
	fmt.Println(emails)
}

func ranger() {
	// To Loop Throughs Arrays & Maps
	ids := []int {33,76,54,23,11,2}

	for i, id := range ids {
		fmt.Printf("%d - ID: %d\n", i, id)
	}

	for _, id := range ids {
		fmt.Printf("%d \n", id)
	}

	sum := 0
	for _, id := range ids {
		sum += id
	}
	fmt.Println("Sum", sum)

	// Maps
	emails := map[string]string{"Bob":"bob@gmail.com", "Mike":"mike@gmail.com"}

	for k, v := range emails {
		fmt.Printf("%s: %s\n", k, v)
	}

	for k := range emails {
		fmt.Println("Name: " + k)
	}
}

func pointers(){
	a := 5
	b := &a

	fmt.Println(a, b)
	fmt.Printf("%T\n", b)

	// Use * to read val from address
	fmt.Println(*b)
	fmt.Println(*&a)

	// Change Val with Pointer
	*b = 10
	fmt.Println(a) // LE PLOT TWIST
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func closures() {
	sum := adder()
	for i := 0; i < 10; i++ {
		fmt.Println(sum(i))
	}
}

type Person struct {
	firstName string
	lastName string
	age int
	city string
	gender string
}

func (p Person) greet() string {
	return "Hello, my name is " + p.firstName + " " + p.lastName + " and I am " + strconv.Itoa(p.age)
}

func (p *Person) hasBirthday() {
	p.age++
}

func(p * Person) gotMarried(spouseLastName string){
	if p.gender == "M" {
		return
	} else {
		p.lastName = spouseLastName
	}
}

func structures() {
	// The classes of go
	Caleb := Person{"Caleb", "Rose", 27, "Austin", "M"}
	fmt.Println(Caleb)
	fmt.Println(Caleb.age)
	fmt.Println(Caleb.greet())
	Caleb.hasBirthday()
	fmt.Println(Caleb.age)
	Kate := Person{"Kait", "K", 28, "Anchorage", "F"}
	Kate.gotMarried("Schroeder")
	fmt.Println(Kate)
}

type Shape interface {
	area() float64
}

type Circle struct {
	x, y, radius float64
}

type Rectangle struct {
	width, height float64
}

func (c Circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (r Rectangle) area() float64 {
	return r.width * r.height
}

func getArea(s Shape) float64 {
	return s.area()
}

func interfaces() {
	// A template for structs -- datatypes that represent a set of method signatures for structs
	circle := Circle {0, 0, 5}
	rectangle := Rectangle {10, 5}

	fmt.Printf("Circle Area: %f\n", getArea(circle))
	fmt.Printf("Rectangle Area: %f\n", getArea(rectangle))
}

func main() {
	fmt.Println("Hello, World!")
	fmt.Println(strutil.Reverse("This is an algorithm"))
	fmt.Println(greeting("Caleb"))
	// fmt.Println(getSum(3, 4))

	interfaces()
}

/* 
	Variables

*/

/* Packages 


*/

/* Functions
Very straightforward


*/

/* Arrays & Slices 
Arrays have to be a fixed length, and types need to be named.

Slices are arrays that don't have a fixed type

Need to learn of appending slices

*/

/* Conditionals */

/* Loops */

/* Map */