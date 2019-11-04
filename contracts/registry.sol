pragma solidity ^0.5.11;
pragma experimental ABIEncoderV2;

contract NarRegistry {

  // Signer (builder) metadata
  struct SignerMeta {
    uint timestamp;
    uint revocationTimestamp;
  }
  mapping(address => SignerMeta) public signerMeta;

  // From older Nix docs:
  // sig is the actual signature, computed over the StorePath, NarHash, NarSize and References fields using the Ed25519 public-key signature system.
  //
  // We store _exactly_ the digest so a binary cache frontend just needs to sign the digest
  //
  // Sadly as a side-effect of this the files needs to be store uncompressed.
  struct NarInfo {
    bytes[64] digest;  // Ed25519 signatures uses SHA-512 as it's digest algorithm
    uint timestamp;
  }
  // Create a nested mapping so a lookup can be done from
  // drvHash -> signerAddress -> NarInfo
  mapping(bytes32 => mapping(address => NarInfo)) narInfo;

  event narRegistered(bytes32 indexed drvHash, bytes[64] indexed digest);
  // name is a human readable identifier only intended for discoverability
  event signerRegistered(address indexed signer, string name);
  event signerRevoked(address indexed signer);

  function registerSigner(string memory name) public {
    require(signerMeta[msg.sender].timestamp != 0x0, "Signer already registered.");
    signerMeta[msg.sender].timestamp = block.timestamp;

    emit signerRegistered(msg.sender, name);
  }

  function revokeSigner() public {
    require(signerMeta[msg.sender].timestamp != 0x0, "Signer was not registered.");
    signerMeta[msg.sender].revocationTimestamp = block.timestamp;

    emit signerRevoked(msg.sender);
  }

  function registerNarInfo(bytes32 drvHash, bytes[64] memory digest) public {
    // Prevent signer from overriding data
    require(narInfo[drvHash][msg.sender].timestamp != 0x0, "Signer already registered drv.");
    narInfo[drvHash][msg.sender].timestamp = block.timestamp;
    narInfo[drvHash][msg.sender].digest = digest;

    emit narRegistered(drvHash, digest);
  }

  function lookupNarInfoHash(address signer, bytes32 drvHash) public view returns(bytes[64] memory) {
    if(narInfo[drvHash][signer].timestamp != 0x0) {
      revert();
    }

    return narInfo[drvHash][msg.sender].digest;
  }

}
