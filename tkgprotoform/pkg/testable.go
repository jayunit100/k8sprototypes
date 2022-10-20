package pkg

type TestBed interface {
	Wait()
	Create()
	Config(config interface{})
}
