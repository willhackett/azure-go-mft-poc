# azure-mft

Managed File Transfer for distributed systems using only Azure Blob Storage & Queue Storage.

This concept relies on a distributed-trust model — each transaction is verified by responding to the initiating agent to verify the integrity of the request. Requests are signed by the initiating agent and validated by requesting the public key. The agents don't directly communicate with eachother ensuring that they remain free from the restrictions of firewalls & network paths.

This model works well with a VWAN or private link, but can also work effectively via the public web.

## CLI

```
  Start the MFT service, this should be done by systemctl:
  $ mft init

  Copy a file to a specific agent:
  $ mft copy <file> --destAgent=<agentName> --destPath=<path> --overwrite=true|false

  Copy a file to multiple agents:
  $ mft copy <file> --config=<destinations.yaml>

  Request a file from a specific agent:
  $ mft req --sourceAgent=<agentName> --sourcePath=<path> <destinationPath>

  Check for updates & update the service as needed:
  $ mft update
```

### Configuration

The MFT configuration file is located by default in `/var/mft2/config.yaml` on unix systems and in `%PROGRAMDATA%\mft2\config.yaml` on Windows systems.

The user that executes the MFT service must have permission to write or read from a destination in order for a file transfer to be successful.

```yaml
version: 1
config:
  app_insights_key: ...
  agent:
    name: '...'
  mft:
    storage_account: '...'
  allow_files_from:
    - 'allowed_agent_name'
  allow_requests_from:
    - 'allowed_agent_name'
  exits:
    - source_agent: 'source_agent_name'
      file_match: '*\\.txt$'
      command: 'cat {fullFilePath}'
```

### Security

#### Service Account

The IAM rules lock down this solution and allow it to be operated in a distributed manner.
The agent name is used to determine access to resources.

- Queue Storage Restrictions
  - `{agent_name}` queue is read/write
  - Other agent names are write-only
- Blob Storage Restrictions
  - `{agent_name}` container is read/write/generate SAS token
  - `public_keys/{agent_name}` is write-only
  - `public_keys/*` is read-only
- AppInsights
  - Write to app insights

#### Allow Lists

Each agent must specify which agents it wishes to receive files and file requests from. Requests from an agent that is not specified in one of these lists will be rejected.

#### Handshake & Signed Messages

When a file transfer is initiated, to verify the identity of the sender a handshake is performed where a signed message of intent is verified by asking for the sending agent to send its public key. This ensures that the requesting agent and signature match before a file transfer is performed.

The requester agent name is checked against a KeyID file stored within the `public_keys/*` container of the blob storage account. This allows for highly available agents to store their keys and for key rotation to perform effectively.

#### SAS Tokens

No agent shoule be able to access another agent's blob storage container. They can only use the SAS token to perform a download.

#### Transfer ID

Each transfer is accompanied with a transfer ID which is consistent across the file transfer. This is used to ensure that all the events can be correlated in Application Insights.

### Exits

Exits are a mechanism for executing arbitrary comands after a successful file transfer. When a transfer is complete and an exit exists that matches the source agent and file regex, the contents of command will be executed by the `azmft` process.

You can have as many exits as you want.
