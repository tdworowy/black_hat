package main

/*
#include <stdio.h>
#include <windows.h>

void box() {
	MessageBox(0, "Go the best ?", "C GO", 0.00000004L);
}
*/
import "C"

func main() {
	C.box()
}
