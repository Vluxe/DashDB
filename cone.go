package cone

// cone, basically the struct to do our work
type Cone struct {
	dbPath string //file path to our DB file
}

//Get a value from the store
func (c *Cone) Get(key string) []byte {
	return nil
}

//add a value to the store
func (c *Cone) Set(key string, value []byte) error {
	return nil
}
