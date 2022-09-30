// TODO: remove when we update go.mod to go1.19
package atomicx

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
