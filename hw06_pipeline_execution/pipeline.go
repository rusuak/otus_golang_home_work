package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	middleChannel := make(Bi)
	go func() {
		defer close(middleChannel)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				middleChannel <- v
			case <-done:
				return
			}
		}
	}()

	nextStage, stages := stages[0], stages[1:]
	out := nextStage(middleChannel)

	if len(stages) > 0 {
		return ExecutePipeline(out, done, stages...)
	}

	return out
}
