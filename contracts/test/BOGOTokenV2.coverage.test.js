const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("BOGOTokenV2 - Coverage Tests", function () {
    let bogoToken;
    let owner, daoWallet, businessWallet, user1, user2;
    let mockContract, nonContract;

    beforeEach(async function () {
        [owner, daoWallet, businessWallet, user1, user2] = await ethers.getSigners();

        // Deploy BOGOTokenV2
        const BOGOTokenV2 = await ethers.getContractFactory("BOGOTokenV2");
        bogoToken = await BOGOTokenV2.deploy();
        await bogoToken.waitForDeployment();

        // Deploy a mock contract for testing _isContract
        const MockContract = await ethers.getContractFactory("BOGOTokenV2"); // Any contract will do
        mockContract = await MockContract.deploy();
        await mockContract.waitForDeployment();

        // EOA address (not a contract)
        nonContract = user1.address;
    });

    describe("Contract Detection (_isContract)", function () {
        it("Should detect contract addresses correctly via queueRegisterFlavoredToken", async function () {
            // This should succeed - mockContract is a contract
            await expect(
                bogoToken.queueRegisterFlavoredToken("Ocean", mockContract.address)
            ).to.not.be.reverted;

            // Verify the operation was queued
            const operationId = ethers.keccak256(
                ethers.solidityPack(
                    ["string", "string", "address"],
                    ["registerFlavoredToken", "Ocean", mockContract.address]
                )
            );
            const executeTime = await bogoToken.timelockOperations(operationId);
            expect(executeTime).to.be.gt(0);
        });

        it("Should reject non-contract addresses", async function () {
            // This should fail - nonContract is an EOA, not a contract
            await expect(
                bogoToken.queueRegisterFlavoredToken("Ocean", nonContract)
            ).to.be.revertedWith("Address must be a contract");
        });

        it("Should reject zero address", async function () {
            await expect(
                bogoToken.queueRegisterFlavoredToken("Ocean", ethers.ZeroAddress)
            ).to.be.revertedWith("Invalid token address");
        });

        it("Should properly identify newly deployed contracts", async function () {
            // Deploy a new contract during the test
            const NewContract = await ethers.getContractFactory("BOGOTokenV2");
            const newContract = await NewContract.deploy();
            await newContract.waitForDeployment();

            // Should recognize it as a contract
            await expect(
                bogoToken.queueRegisterFlavoredToken("Wildlife", newContract.address)
            ).to.not.be.reverted;
        });

        it("Should handle contracts with code", async function () {
            // Deploy another contract to test _isContract functionality
            const AnotherContract = await ethers.getContractFactory("BOGOTokenV2");
            const anotherContract = await AnotherContract.deploy();
            await anotherContract.waitForDeployment();

            // Should be recognized as a contract
            await expect(
                bogoToken.queueRegisterFlavoredToken("Another", anotherContract.address)
            ).to.not.be.reverted;
        });
    });

    describe("ERC165 Interface Support", function () {
        it("Should support AccessControl interface", async function () {
            // AccessControl interface ID
            const accessControlInterfaceId = "0x7965db0b";
            expect(await bogoToken.supportsInterface(accessControlInterfaceId)).to.be.true;
        });

        it("Should support ERC165 interface itself", async function () {
            // ERC165 interface ID
            const erc165InterfaceId = "0x01ffc9a7";
            expect(await bogoToken.supportsInterface(erc165InterfaceId)).to.be.true;
        });

        it("Should not support random interface", async function () {
            // Random interface ID
            const randomInterfaceId = "0x12345678";
            expect(await bogoToken.supportsInterface(randomInterfaceId)).to.be.false;
        });

        it("Should not support invalid interface ID", async function () {
            // Invalid interface ID (0xffffffff is specifically invalid in ERC165)
            const invalidInterfaceId = "0xffffffff";
            expect(await bogoToken.supportsInterface(invalidInterfaceId)).to.be.false;
        });

        it("Should handle interface queries for multiple standards", async function () {
            // Test multiple interface IDs in sequence
            const interfaces = [
                { id: "0x7965db0b", name: "AccessControl", expected: true },
                { id: "0x01ffc9a7", name: "ERC165", expected: true },
                { id: "0x36372b07", name: "ERC20", expected: false }, // Not explicitly supported
                { id: "0x80ac58cd", name: "ERC721", expected: false }, // Not supported
                { id: "0xd9b67a26", name: "ERC1155", expected: false }, // Not supported
            ];

            for (const iface of interfaces) {
                const result = await bogoToken.supportsInterface(iface.id);
                expect(result).to.equal(iface.expected, `Interface ${iface.name} (${iface.id})`);
            }
        });
    });

    describe("Edge Cases for Assembly Code Coverage", function () {
        it("Should handle contract creation during same transaction", async function () {
            // Test that demonstrates the assembly code works correctly
            // by attempting to register a flavored token with various addresses
            
            const addresses = [
                { addr: mockContract.address, shouldPass: true, name: "Deployed contract" },
                { addr: owner.address, shouldPass: false, name: "EOA owner" },
                { addr: user1.address, shouldPass: false, name: "EOA user" },
                { addr: "0x0000000000000000000000000000000000000001", shouldPass: false, name: "Precompile" },
            ];

            for (const test of addresses) {
                if (test.shouldPass) {
                    await expect(
                        bogoToken.queueRegisterFlavoredToken(`Test-${test.name}`, test.addr)
                    ).to.not.be.reverted;
                } else {
                    await expect(
                        bogoToken.queueRegisterFlavoredToken(`Test-${test.name}`, test.addr)
                    ).to.be.revertedWith("Address must be a contract");
                }
            }
        });

        it("Should correctly identify contract with no code but existing storage", async function () {
            // This tests the edge case where a contract might have storage but no code
            // The _isContract function should return false for addresses with no code
            
            // Use a deterministic address that we know has no code
            const emptyAddress = "0x1234567890123456789012345678901234567890";
            
            await expect(
                bogoToken.queueRegisterFlavoredToken("Empty", emptyAddress)
            ).to.be.revertedWith("Address must be a contract");
        });
    });

    describe("Combined Functionality Tests", function () {
        it("Should queue and execute flavored token registration with interface support", async function () {
            // Queue registration
            await bogoToken.queueRegisterFlavoredToken("Ocean", mockContract.address);

            // Get operation ID
            const operationId = ethers.keccak256(
                ethers.solidityPack(
                    ["string", "string", "address"],
                    ["registerFlavoredToken", "Ocean", mockContract.address]
                )
            );

            // Fast forward time
            await ethers.provider.send("evm_increaseTime", [2 * 24 * 60 * 60 + 1]); // 2 days + 1 second
            await ethers.provider.send("evm_mine");

            // Execute registration
            await bogoToken.executeRegisterFlavoredToken("Ocean", mockContract.address);

            // Verify registration
            expect(await bogoToken.flavoredTokens("Ocean")).to.equal(mockContract.address);

            // Verify contract still supports its interfaces after operations
            expect(await bogoToken.supportsInterface("0x7965db0b")).to.be.true;
        });
    });
});