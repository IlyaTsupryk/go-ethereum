// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title StringStorage
 * @dev Store and retrieve string inputs for each user
 */
contract StringStorage {
    // Address of the Schnorr verification precompiled contract
    address constant SCHNORR_VERIFY = address(0x0101);

    // Mapping from user address to their array of strings
    mapping(address => string[]) private userStrings;

    // Event emitted when a new string is added
    event StringAdded(address indexed user, uint256 index, string value);

    /**
     * @dev Add a new string input for the caller
     * @param input The string to store
     */
    function addString(string memory input) public {
        userStrings[msg.sender].push(input);
        emit StringAdded(msg.sender, userStrings[msg.sender].length - 1, input);
    }

    /**
     * @dev Add a new string input with Schnorr signature verification
     * @param input The string to store
     * @param signature The 64-byte Schnorr signature (r || s)
     * @param pubkey The 64-byte uncompressed public key
     */
    function addStringSecure(
        string memory input,
        bytes memory signature,
        bytes memory pubkey
    ) public {
        require(signature.length == 64, "Invalid signature length");
        require(pubkey.length == 64, "Invalid pubkey length");

        // Check that pubkey really belongs to msg.sender
        // address = last 20 bytes of keccak256(pubkey)
        address derivedAddr = address(uint160(uint256(keccak256(pubkey))));
        require(derivedAddr == msg.sender, "Pubkey does not match sender");

        // Schorr pubKeyX
        bytes32 pubKeyX;
        assembly { pubKeyX := mload(add(pubkey, 32))}

        bytes32 messageHash = keccak256(abi.encodePacked(input));

        require(
            verifySchnorrSignature(messageHash, signature, abi.encodePacked(pubKeyX)),
            "Invalid Schnorr signature"
        );

        userStrings[msg.sender].push(input);
        emit StringAdded(msg.sender, userStrings[msg.sender].length - 1, input);
    }

    /**
     * @dev Internal function to verify Schnorr signature using precompiled contract
     * @param messageHash The 32-byte message hash
     * @param signature The 64-byte Schnorr signature (r || s)
     * @param pubkey The 32-byte compressed public key
     * @return bool True if signature is valid
     */
    function verifySchnorrSignature(
        bytes32 messageHash,
        bytes memory signature,
        bytes memory pubkey
    ) internal view returns (bool) {
        bytes memory input = abi.encodePacked(messageHash, signature, pubkey);

        (bool success, bytes memory result) = SCHNORR_VERIFY.staticcall(input);

        if (!success || result.length != 32) {
            return false;
        }
        return result[31] == 0x01;
    }

    /**
     * @dev Get a specific string by index for the user
     * @param user The address of the user
     * @param index The index of the string to retrieve
     * @return The string at the specified index for the given user
     */
    function getString(address user, uint256 index) public view returns (string memory) {
        require(index < userStrings[user].length, "Index out of bounds");
        return userStrings[user][index];
    }
}