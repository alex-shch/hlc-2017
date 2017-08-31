package prof

type ProfilingTool interface {
	Stop()
}

var Prof ProfilingTool = emptyProf{}

type emptyProf struct{}

func (emptyProf) Stop() {}
