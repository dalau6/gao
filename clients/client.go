// Add implements add.
func (s *clientsrvc) Add(ctx context.Context,
	p *client.AddPayload) (err error) {
	s.logger.Print("client.add started")
	newClient := client.ClientManagement{
				ClientID: p.ClientID,
				ClientName: p.ClientName,
	}
	err = CreateClient(&newClient)
	if err != nil {
				s.logger.Print("An error occurred...")
				s.logger.Print(err)
				return
	}
	s.logger.Print("client.add completed")
	return
}

// Get implements get.
func (s *clientsrvc) Get(ctx context.Context,
	p *client.GetPayload) (res *client.ClientManagement, err error) {
	s.logger.Print("client.get started")
	result, err := GetClient(p.ClientID)
	if err != nil {
				s.logger.Print("An error occurred...")
				s.logger.Print(err)
				return
	}
	s.logger.Print("client.get completed")
	return &result, err
}

// Show implements show.
func (s *clientsrvc) Show(ctx context.Context) (res client.ClientManagementCollection,
	err error) {
	s.logger.Print("client.show started")
	res, err = ListClients()
	if err != nil {
				s.logger.Print("An error occurred...")
				s.logger.Print(err)
				return
	}
	s.logger.Print("client.show completed")
	return
}