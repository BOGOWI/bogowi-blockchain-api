// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/utils/Pausable.sol";

contract MockFailingPausable is Pausable {
    bool public shouldFailPause = true;
    bool public shouldFailUnpause = true;
    
    constructor() {
        // Start unpaused
    }
    
    function pause() external {
        if (shouldFailPause) {
            revert("Pause operation failed");
        }
        _pause();
    }
    
    function unpause() external {
        if (shouldFailUnpause) {
            revert("Unpause operation failed");
        }
        _unpause();
    }
    
    function setShouldFailPause(bool _shouldFail) external {
        shouldFailPause = _shouldFail;
    }
    
    function setShouldFailUnpause(bool _shouldFail) external {
        shouldFailUnpause = _shouldFail;
    }
    
    function forceSetPaused(bool _paused) external {
        if (_paused && !paused()) {
            _pause();
        } else if (!_paused && paused()) {
            _unpause();
        }
    }
}