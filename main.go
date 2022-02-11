package main

func main() {
	orch := new(Orchestrator)
	err := orch.Start()

	if err != nil {
		panic(err)
	}

	<-orch.terminated
}
