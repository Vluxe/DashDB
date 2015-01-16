package cone

//This is where the logarithmic data magic happens. A node's basic structure on disk looks like so:
//
// |keyPart|\n|c|keyPart,seekInt,keyPart,seekInt...|\n|v,valueLength|valueData........|
// |keyPart|\n|c|keyPart,seekInt,keyPart,seekInt...|\n|v,valueLength|valueData........|
//
// It is possible that a node has no children and a "c" delimiter isn't present, thus making the structure:
//
// |keyPart|\n|v,valueLength|valueData........|

type childNode struct {
	keyPart string //the part of the key it represents
	start   int    //the start point in the file for this node
}

// The node struct is used to represent the key/value pairs on disk
type Node struct {
	parent   *Node       //the parent node.
	keyPart  string      //the part of the node this presents (e.g. if key is "user-1825") and this was the "5" part
	childern []childNode //the start location, seek path on disk
	value    []byte
}
