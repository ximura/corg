package node

type Role string

const (
	Worker  Role = "worker"
	Manager Role = "manager"
)

// Node - represent any machnie in cluster
// Example of node types:
// * Manager
// * Workers
type Node struct {
	Name            string
	Ip              string
	Cores           int
	Memory          int
	MemoryAllocated int
	Disk            int
	DiskAllocated   int
	Role            Role
	TaskCount       int
}
