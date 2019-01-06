package PacketData

type NOAAProduct struct {
	ID          int
	Name        string
	SubProducts map[int]NOAASubProduct
}

func MakeNOAAProduct(id int) NOAAProduct {
	return MakeNOAAProductWithName(id, "")
}

func MakeNOAAProductWithName(id int, name string) NOAAProduct {
	return MakeNOAAProductWithSubProductsAndName(id, name, map[int]NOAASubProduct{})
}

func MakeNOAAProductWithSubProductsAndName(id int, name string, subProducts map[int]NOAASubProduct) NOAAProduct {
	return NOAAProduct{
		ID:          id,
		Name:        "",
		SubProducts: subProducts,
	}
}

func (np *NOAAProduct) GetSubProduct(id int) NOAASubProduct {
	val, ok := np.SubProducts[id]

	if !ok {
		return MakeSubProduct(id, "Unknown")
	}

	return val
}
