## Go-Cita - WIP
> A Go implementation of CITA, Based on cita 2.0 microservices architecture.
## Components
- API - Provide user interface for facebook GraphQL.
- Network - Provide inter-node access.
- Consensus - Consensus module, Use the tendermint consensus algorithm for future support consensus pluggable.
- Sync - Block synchronization between nodes.
- Auth - Management of the transaction pool.
- Chain - Management block.
- VM - A virtual machine that executes smart contracts and supports multiple types of vm, such as evm, g(go)vm, r(rust)vm.

## Project State
In development...
- [ ] API 
- [x] Network 
- [x] Consensus (tendermint)
- [ ] Sync
- [x] Auth
- [ ] Chain
- [x] VM - MVP
