package device

func (s *Service) UpdateName(deviceTypeID uint, name string) error {
	updateData := map[string]interface{}{
		"name": name,
	}
	err := s.dao.UpdateDeviceType(s.dao.DB(), deviceTypeID, updateData)
	if err != nil {
		return err
	}
	return nil
}
