// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title StringStorage
 * @dev Store and retrieve string inputs for each user
 */
contract StringStorage {
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