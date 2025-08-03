const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("RoleManager", function () {
  let roleManager;
  let owner;
  let admin;
  let user1;
  let user2;
  let contract1;
  let contract2;

  // Role constants
  let DEFAULT_ADMIN_ROLE;
  let DAO_ROLE;
  let BUSINESS_ROLE;
  let MINTER_ROLE;
  let PAUSER_ROLE;
  let TREASURY_ROLE;
  let DISTRIBUTOR_BACKEND_ROLE;

  beforeEach(async function () {
    [owner, admin, user1, user2, contract1, contract2] = await ethers.getSigners();

    // Deploy RoleManager
    const RoleManager = await ethers.getContractFactory("RoleManager");
    roleManager = await RoleManager.deploy();
    await roleManager.waitForDeployment();

    // Get role constants
    DEFAULT_ADMIN_ROLE = await roleManager.DEFAULT_ADMIN_ROLE();
    DAO_ROLE = await roleManager.DAO_ROLE();
    BUSINESS_ROLE = await roleManager.BUSINESS_ROLE();
    MINTER_ROLE = await roleManager.MINTER_ROLE();
    PAUSER_ROLE = await roleManager.PAUSER_ROLE();
    TREASURY_ROLE = await roleManager.TREASURY_ROLE();
    DISTRIBUTOR_BACKEND_ROLE = await roleManager.DISTRIBUTOR_BACKEND_ROLE();
  });

  describe("Deployment", function () {
    it("Should set deployer as DEFAULT_ADMIN", async function () {
      expect(await roleManager.hasRole(DEFAULT_ADMIN_ROLE, owner.address)).to.be.true;
    });

    it("Should set correct role admins", async function () {
      expect(await roleManager.getRoleAdmin(DAO_ROLE)).to.equal(DEFAULT_ADMIN_ROLE);
      expect(await roleManager.getRoleAdmin(BUSINESS_ROLE)).to.equal(DEFAULT_ADMIN_ROLE);
      expect(await roleManager.getRoleAdmin(MINTER_ROLE)).to.equal(DEFAULT_ADMIN_ROLE);
      expect(await roleManager.getRoleAdmin(PAUSER_ROLE)).to.equal(DEFAULT_ADMIN_ROLE);
      expect(await roleManager.getRoleAdmin(TREASURY_ROLE)).to.equal(DEFAULT_ADMIN_ROLE);
      expect(await roleManager.getRoleAdmin(DISTRIBUTOR_BACKEND_ROLE)).to.equal(DEFAULT_ADMIN_ROLE);
    });
  });

  describe("Contract Registration", function () {
    it("Should register a contract", async function () {
      await expect(roleManager.registerContract(contract1.address, "TestContract1"))
        .to.emit(roleManager, "ContractRegistered")
        .withArgs(contract1.address, "TestContract1");

      expect(await roleManager.registeredContracts(contract1.address)).to.be.true;
      expect(await roleManager.contractNames(contract1.address)).to.equal("TestContract1");
    });

    it("Should reject zero address registration", async function () {
      await expect(
        roleManager.registerContract(ethers.ZeroAddress, "TestContract")
      ).to.be.revertedWith("Invalid contract address");
    });

    it("Should reject duplicate registration", async function () {
      await roleManager.registerContract(contract1.address, "TestContract1");
      
      await expect(
        roleManager.registerContract(contract1.address, "TestContract1")
      ).to.be.revertedWith("Contract already registered");
    });

    it("Should only allow admin to register contracts", async function () {
      await expect(
        roleManager.connect(user1).registerContract(contract1.address, "TestContract1")
      ).to.be.revertedWithCustomError(roleManager, "AccessControlUnauthorizedAccount");
    });
  });

  describe("Contract Deregistration", function () {
    beforeEach(async function () {
      await roleManager.registerContract(contract1.address, "TestContract1");
    });

    it("Should deregister a contract", async function () {
      await expect(roleManager.deregisterContract(contract1.address))
        .to.emit(roleManager, "ContractDeregistered")
        .withArgs(contract1.address);

      expect(await roleManager.registeredContracts(contract1.address)).to.be.false;
      expect(await roleManager.contractNames(contract1.address)).to.equal("");
    });

    it("Should reject deregistering non-registered contract", async function () {
      await expect(
        roleManager.deregisterContract(contract2.address)
      ).to.be.revertedWith("Contract not registered");
    });

    it("Should only allow admin to deregister contracts", async function () {
      await expect(
        roleManager.connect(user1).deregisterContract(contract1.address)
      ).to.be.revertedWithCustomError(roleManager, "AccessControlUnauthorizedAccount");
    });
  });

  describe("Role Checking", function () {
    beforeEach(async function () {
      await roleManager.registerContract(contract1.address, "TestContract1");
      await roleManager.grantRole(DAO_ROLE, user1.address);
    });

    it("Should allow registered contracts to check roles", async function () {
      // Impersonate the registered contract
      await ethers.provider.send("hardhat_impersonateAccount", [contract1.address]);
      await owner.sendTransaction({ to: contract1.address, value: ethers.parseEther("1") });
      const contractSigner = await ethers.provider.getSigner(contract1.address);
      
      const hasRole = await roleManager.connect(contractSigner).checkRole(DAO_ROLE, user1.address);
      expect(hasRole).to.be.true;

      const noRole = await roleManager.connect(contractSigner).checkRole(DAO_ROLE, user2.address);
      expect(noRole).to.be.false;

      await ethers.provider.send("hardhat_stopImpersonatingAccount", [contract1.address]);
    });

    it("Should reject role checking from non-registered contracts", async function () {
      await expect(
        roleManager.connect(user1).checkRole(DAO_ROLE, user1.address)
      ).to.be.revertedWith("Not a registered contract");
    });
  });

  describe("Role Management with Events", function () {
    it("Should emit global event when granting role", async function () {
      await expect(roleManager.grantRole(DAO_ROLE, user1.address))
        .to.emit(roleManager, "RoleGrantedGlobally")
        .withArgs(DAO_ROLE, user1.address, owner.address);

      expect(await roleManager.hasRole(DAO_ROLE, user1.address)).to.be.true;
    });

    it("Should emit global event when revoking role", async function () {
      await roleManager.grantRole(DAO_ROLE, user1.address);

      await expect(roleManager.revokeRole(DAO_ROLE, user1.address))
        .to.emit(roleManager, "RoleRevokedGlobally")
        .withArgs(DAO_ROLE, user1.address, owner.address);

      expect(await roleManager.hasRole(DAO_ROLE, user1.address)).to.be.false;
    });

    it("Should only allow role admin to grant roles", async function () {
      await expect(
        roleManager.connect(user1).grantRole(DAO_ROLE, user2.address)
      ).to.be.revertedWithCustomError(roleManager, "AccessControlUnauthorizedAccount");
    });

    it("Should only allow role admin to revoke roles", async function () {
      await roleManager.grantRole(DAO_ROLE, user1.address);

      await expect(
        roleManager.connect(user1).revokeRole(DAO_ROLE, user1.address)
      ).to.be.revertedWithCustomError(roleManager, "AccessControlUnauthorizedAccount");
    });
  });

  describe("Batch Operations", function () {
    it("Should batch grant roles", async function () {
      const accounts = [user1.address, user2.address];

      await roleManager.batchGrantRole(DAO_ROLE, accounts);

      expect(await roleManager.hasRole(DAO_ROLE, user1.address)).to.be.true;
      expect(await roleManager.hasRole(DAO_ROLE, user2.address)).to.be.true;
    });

    it("Should batch revoke roles", async function () {
      const accounts = [user1.address, user2.address];
      await roleManager.batchGrantRole(DAO_ROLE, accounts);

      await roleManager.batchRevokeRole(DAO_ROLE, accounts);

      expect(await roleManager.hasRole(DAO_ROLE, user1.address)).to.be.false;
      expect(await roleManager.hasRole(DAO_ROLE, user2.address)).to.be.false;
    });

    it("Should only allow role admin for batch operations", async function () {
      const accounts = [user1.address, user2.address];

      await expect(
        roleManager.connect(user1).batchGrantRole(DAO_ROLE, accounts)
      ).to.be.revertedWithCustomError(roleManager, "AccessControlUnauthorizedAccount");
    });
  });

  describe("Admin Transfer", function () {
    it("Should transfer admin role", async function () {
      await roleManager.transferAdmin(admin.address);

      expect(await roleManager.hasRole(DEFAULT_ADMIN_ROLE, admin.address)).to.be.true;
      expect(await roleManager.hasRole(DEFAULT_ADMIN_ROLE, owner.address)).to.be.false;
    });

    it("Should reject zero address as new admin", async function () {
      await expect(
        roleManager.transferAdmin(ethers.ZeroAddress)
      ).to.be.revertedWith("Invalid new admin");
    });

    it("Should only allow current admin to transfer", async function () {
      await expect(
        roleManager.connect(user1).transferAdmin(user2.address)
      ).to.be.revertedWithCustomError(roleManager, "AccessControlUnauthorizedAccount");
    });
  });

  describe("Pausable", function () {
    it("Should pause and unpause", async function () {
      await roleManager.pause();
      expect(await roleManager.paused()).to.be.true;

      await roleManager.unpause();
      expect(await roleManager.paused()).to.be.false;
    });

    it("Should only allow admin to pause/unpause", async function () {
      await expect(
        roleManager.connect(user1).pause()
      ).to.be.revertedWithCustomError(roleManager, "AccessControlUnauthorizedAccount");

      await roleManager.pause();

      await expect(
        roleManager.connect(user1).unpause()
      ).to.be.revertedWithCustomError(roleManager, "AccessControlUnauthorizedAccount");
    });
  });

  describe("Interface Support", function () {
    it("Should support AccessControl interface", async function () {
      const accessControlInterface = "0x7965db0b"; // IAccessControl interface ID
      expect(await roleManager.supportsInterface(accessControlInterface)).to.be.true;
    });

    it("Should support ERC165 interface", async function () {
      const erc165Interface = "0x01ffc9a7"; // IERC165 interface ID
      expect(await roleManager.supportsInterface(erc165Interface)).to.be.true;
    });
  });

  describe("Get Registered Contracts", function () {
    it("Should return empty arrays for registered contracts", async function () {
      const result = await roleManager.getRegisteredContracts();
      expect(result.addresses).to.be.an('array').that.is.empty;
      expect(result.names).to.be.an('array').that.is.empty;
    });
  });
});