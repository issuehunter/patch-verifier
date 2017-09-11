# Issuehunter patch verifier

Verify that a proposed solution is indeed a valid commit on the repo and that has been merged into the main
repo branch.


## Interactive logs


There are two utility functions in the `patchverifier` package that lets you
inspect live logs from the node. You can use them as in the following examples:

First of all start a node locally. In the following example we are connecting to the rinkeby.io test network:

    geth --networkid=4 --datadir=$HOME/.rinkeby --cache=512 --ethstats='yournode:Respect my authoritah!@stats.rinkeby.io' --bootnodes=enode://a24ac7c5484ef4ed0c5eb2d36620ba4e4aa13b8c84684e1b4aab0cebea2ae45cb4d375b77eab56516d34bfbd3c1a833fc51296ff084b770b94fb9028c4d25ccf@52.169.42.101:30303 --rpc

### Live logging of all transactions involving the issuehunter contract

    package main
      
    import (
        "github.com/issuehunter/patch-verifier"
    )
    
    func main() {
        patchverifier.InteractiveLogObserver("http://127.0.0.1:8545", "CONTRACT_ADDRESS_HERE")
    }


### Live logging of all ProposedResultion events on the issuehunter contract

    package main
      
    import (
        "github.com/issuehunter/patch-verifier"
    )
    
    func main() {
        patchverifier.InteractiveResolutionProposedLogObserver("http://127.0.0.1:8545", "CONTRACT_ADDRESS_HERE")
    }
