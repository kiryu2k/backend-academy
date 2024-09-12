package device

type Device struct {
	SerialNum string `json:"serial_num"`
	Model     string `json:"model"`
	IP        string `json:"ip"`
}

type Usecase interface {
	CreateDevice(*Device) error
	GetDevice(string) (*Device, error)
	UpdateDevice(*Device) error
	DeleteDevice(string) error
}

type Repository interface {
	GetDevice(string) (*Device, error)
	CreateDevice(*Device) error
	DeleteDevice(string) error
	UpdateDevice(*Device) error
}
