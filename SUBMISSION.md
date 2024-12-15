# Implementation Details

## Backend service
Details for the backend service can be found in [GENOMIC_SERVICE_DOCUMENT.md](GENOMIC_SERVICE_DOCUMENT.md)

## Smart contract
- Contract was updated and pass all the tests
- Deployment contracts was saved in [deployment.json](genomicdao/scripts/deployments/deployments.json)

## Avalanche Subnet
For detailed instructions on setting up an Avalanche Subnet for Localnet, refer to the [official documentation](https://docs.avax.network/avalanche-l1s/build-first-avalanche-l1)
```
➜  ~ avalanche blockchain describe LIFENetwork 
+--------------------------------------------------------------------------------------------------------------------------------+
|                                                           LIFENETWORK                                                          |
+---------------+----------------------------------------------------------------------------------------------------------------+
| Name          | LIFENetwork                                                                                                    |
+---------------+----------------------------------------------------------------------------------------------------------------+
| VM ID         | abdW2t8SQ86UUf1LVL2WVhB4THZbQiEoSDBSVgTyuEe5hcpFH                                                              |
+---------------+----------------------------------------------------------------------------------------------------------------+
| VM Version    | v0.7.0                                                                                                         |
+---------------+----------------------------------------------------------------------------------------------------------------+
| Validation    | Proof Of Authority                                                                                             |
+---------------+--------------------------+-------------------------------------------------------------------------------------+
| Local Network | ChainID                  | 9999                                                                                |
|               +--------------------------+-------------------------------------------------------------------------------------+
|               | SubnetID                 | emGjKGWnF2YPU8ZvErRWewGQVY1RAncMV58Jvmtsy83nzurAB                                   |
|               +--------------------------+-------------------------------------------------------------------------------------+
|               | BlockchainID (CB58)      | 2bGAh54yzGQ3nj4txNDDSzoKZTNfghang9z92Cgz315ch1nAsA                                  |
|               +--------------------------+-------------------------------------------------------------------------------------+
|               | BlockchainID (HEX)       | 0xd17db3db9d086bf25e5c2e3fa53cb21bcdfa94560f37fae1315d7cf50ed3b98e                  |
|               +--------------------------+-------------------------------------------------------------------------------------+
|               | RPC Endpoint             | http://127.0.0.1:9650/ext/bc/2bGAh54yzGQ3nj4txNDDSzoKZTNfghang9z92Cgz315ch1nAsA/rpc |
+---------------+--------------------------+-------------------------------------------------------------------------------------+
```

## Integration Test
1. Deploy Avalanche Subnet or using any chain that you want then config the network in [hardhat.config.js](genomicdao/hardhat.config.js)
2. Deploy contract 
```bash
  npx hardhat run scripts/deploy.js --network lifenetwork
```
3. Setting deploy contract, chain configs in [local.yml](genomic-service/config/local.yml)
4. Set the APP_DIR to your directory in [integration_test.sh](integration_test.sh)
5. Run [integration_test.sh](integration_test.sh)

Test Result 

```json lines
Starting the application...
Application started with PID 32356.
Step 1 - Register user
Sending POST request to /v1/users...
HTTP Response Body:
{
  "meta": {
    "code": 200,
    "message": "Successful"
  },
  "data": {
    "public_key": "03fb1dee9df7b2f231a69b167ddecdb4223960aeee04fb3d80801c89372ff8a6a4"
  }
}
Step 2 - Upload Data
Sending POST request to /v1/tee/encrypt to encrypt raw data
HTTP Response Body:
{
  "meta": {
    "code": 200,
    "message": "Successful"
  },
  "data": {
    "encrypted_data": "BPHaUiEwn6Gcd/1Rd3uTW1l+KCerUGgAXs+zoYDIhOWZ4q7rw1kdRtauu0+jeMeOg+saa7sN8GbKnu6vE9/XBWUQmxDpDR0kvnIxds/2WvTL/XsIxuX9BcuuZz0wkc5zHy72VK96EfP+x/GfNAbjQ9BXWiUO+5ob6w=="
  }
}
Sending POST request to /v1/users/upload to upload data
HTTP Response Body:
{
  "meta": {
    "code": 200,
    "message": "Successful"
  },
  "data": {
    "file_id": "fab66c26376470f539b5dba0202c88c1",
    "user_address": "0x1234567890abcdef1234567890abcdef12345678",
    "data_hash": "+rZsJjdkcPU5tdugICyIwSHWz7/8ezc86fwMjYc4GoY=",
    "encrypted_data": "BPHaUiEwn6Gcd/1Rd3uTW1l+KCerUGgAXs+zoYDIhOWZ4q7rw1kdRtauu0+jeMeOg+saa7sN8GbKnu6vE9/XBWUQmxDpDR0kvnIxds/2WvTL/XsIxuX9BcuuZz0wkc5zHy72VK96EfP+x/GfNAbjQ9BXWiUO+5ob6w==",
    "is_confirmed": false
  }
}
Step 3 - Confirm and reward for user
Fetching uploaded data to confirm the data
HTTP Response Body:
{
  "meta": {
    "code": 200,
    "message": "Successful"
  },
  "data": {
    "file_id": "fab66c26376470f539b5dba0202c88c1",
    "user_address": "0x1234567890abcdef1234567890abcdef12345678",
    "data_hash": "+rZsJjdkcPU5tdugICyIwSHWz7/8ezc86fwMjYc4GoY=",
    "encrypted_data": "BPHaUiEwn6Gcd/1Rd3uTW1l+KCerUGgAXs+zoYDIhOWZ4q7rw1kdRtauu0+jeMeOg+saa7sN8GbKnu6vE9/XBWUQmxDpDR0kvnIxds/2WvTL/XsIxuX9BcuuZz0wkc5zHy72VK96EfP+x/GfNAbjQ9BXWiUO+5ob6w==",
    "is_confirmed": true
  }
}
Stopping the application...
Killing process on port 8080...
Application stopped.


```