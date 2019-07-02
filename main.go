package main

func main() {
	a := App{}
	a.Initialize("root", "")

	a.Run(":8080")
}
