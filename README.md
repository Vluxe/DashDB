# DashDB
consistent, scalable, and safe key-value DB in Go.

Working theory on how to build an ACID compliant, yet fast key value store in Go.

Current theory:
- Build a schemaless file that stores data in a btree format. 
- Seek through the file from node to node of the btree. 
- Divide keys into nodes based on letters/words depending on uniqueness. 
- (e.g. users-18 and users-17, 4 nodes create first node of "users-", second node of "1" and 2 childern of "7" and "8"). This should scale logarithmically in theory.
- Store most recently used keys in memory for speed.

Going to test this out and see how it works, hoping it will vote for me and make all my wildest dreams come true.
