package descend

type Descend struct {
	CaCert  string `flag:"ca"      description:"CA Cert"`
	Cert    string `flag:"cert"    description:"Client Cert"`
	Key     string `flag:"key"     description:"Client Key"`
	Host    string `flag:"host"    description:"Host to PUT to"`
	Port    int    `flag:"port"    description:"Port to PUT on"`
	Archive string `flag:"archive" description:"Archive to PUT to"`
}
