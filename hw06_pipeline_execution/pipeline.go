package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	doneStage := func(in In) Out {
		out := make(Bi)
		go func() {
			defer close(out)
			for {
				select {
				case v, ok := <-in:
					if !ok {
						return
					}
					out <- v
				case <-done:
					return
				}
			}
		}()
		return out
	}

	out := doneStage(in)

	if len(stages) == 1 && stages[0] == nil {
		return out
	}

	for _, stage := range stages {
		out = doneStage(stage(out))
	}

	return out
}
