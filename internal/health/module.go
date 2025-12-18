package health

type Module struct {
    Controller *Controller
}

func Initialize() *Module {
    controller := NewController()
    return &Module{Controller: controller}
}

