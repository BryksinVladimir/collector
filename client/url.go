package client

func (client *MobildaClient) makeUrl(accIndex int, path string) string {
	return client.accounts[accIndex].Url + path
}
