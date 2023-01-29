package main

import "zinx/znet"

func main() {
	s := znet.NewServer("[zinxzzz]")

	s.Serve()
}
