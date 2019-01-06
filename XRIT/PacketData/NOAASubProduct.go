package PacketData

type NOAASubProduct struct {
	ID   int
	Name string
}

func MakeSubProduct(id int, name string) NOAASubProduct {
	return NOAASubProduct{
		ID:   id,
		Name: name,
	}
}
