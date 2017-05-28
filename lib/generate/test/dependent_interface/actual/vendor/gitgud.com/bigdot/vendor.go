package bigdot

// VendorDep is a vendor dependency
type VendorDep interface {
	Vending(VendorStruct)
}

// VendorStruct is a struct
type VendorStruct struct{}
