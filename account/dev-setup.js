// import
personal.importRawKey("PRIVATE KEY (HEX)", "PASSPHRASE")

// check
eth.accounts // should be two, second should be the address for the imported private key

// eth.accounts[0] is the default dev account, has lots of ether

// transfer
eth.sendTransaction({from: eth.accounts[0], to: eth.accounts[1], value: web3.toWei(1, "ether")})

// check
eth.getBalance(eth.accounts[1]) // should be 1000000000000000000 (1e+18)

// unlock
personal.unlockAccount(eth.accounts[1], "PASSPHRASE")