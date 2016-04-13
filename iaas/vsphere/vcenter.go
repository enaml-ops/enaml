package vsphere

func NewVCenter(address, user, password string) *VCenter {
	return &VCenter{
		Address:  address,
		User:     user,
		Password: password,
	}
}

func (s *VCenter) AddDataCenter(d DataCenter) (err error) {
	s.DataCenters = append(s.DataCenters, d)
	return
}
