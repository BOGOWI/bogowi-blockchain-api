const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("RoleManaged", function () {
  let roleManager;
  let testContract;
  let owner;
  let user1;
  let user2;

  // Role constants
  let DAO_ROLE;
  let BUSINESS_ROLE;

  beforeEach(async function () {
    [owner, user1, user2] = await ethers.getSigners();

    // Deploy RoleManager first
    const RoleManager = await ethers.getContractFactory("RoleManager");
    roleManager = await RoleManager.deploy();
    await roleManager.waitForDeployment();

    // Get role constants
    DAO_ROLE = await roleManager.DAO_ROLE();
    BUSINESS_ROLE = await roleManager.BUSINESS_ROLE();

    // Deploy test contract with RoleManager
    const TestRoleManaged = await ethers.getContractFactory("TestRoleManaged");
    testContract = await TestRoleManaged.deploy(await roleManager.getAddress());
    await testContract.waitForDeployment();

    // Register the test contract with RoleManager
    await roleManager.registerContract(await testContract.getAddress(), "TestRoleManaged");

    // Grant roles
    await roleManager.grantRole(DAO_ROLE, user1.address);
    await roleManager.grantRole(BUSINESS_ROLE, user2.address);
  });

  describe("Constructor", function () {
    it("Should set RoleManager correctly", async function () {
      expect(await testContract.roleManager()).to.equal(await roleManager.getAddress());
    });

    it("Should emit RoleManagerSet event", async function () {
      const TestRoleManaged = await ethers.getContractFactory("TestRoleManaged");
      const roleManagerAddress = await roleManager.getAddress();
      
      const tx = await TestRoleManaged.deploy(roleManagerAddress);
      const receipt = await tx.deploymentTransaction().wait();
      
      // Check event was emitted
      const event = receipt.logs.find(log => {
        try {
          const parsed = testContract.interface.parseLog(log);
          return parsed.name === "RoleManagerSet";
        } catch {
          return false;
        }
      });
      
      expect(event).to.not.be.undefined;
    });

    it("Should revert with zero address", async function () {
      const TestRoleManaged = await ethers.getContractFactory("TestRoleManaged");
      
      await expect(TestRoleManaged.deploy(ethers.ZeroAddress))
        .to.be.revertedWith("Invalid RoleManager address");
    });
  });

  describe("onlyRole Modifier", function () {
    it("Should allow DAO role to call DAO-restricted function", async function () {
      await testContract.connect(user1).setValueDAO(100);
      expect(await testContract.value()).to.equal(100);
    });

    it("Should allow BUSINESS role to call BUSINESS-restricted function", async function () {
      await testContract.connect(user2).setValueBusiness(200);
      expect(await testContract.value()).to.equal(200);
    });

    it("Should revert when unauthorized user calls DAO-restricted function", async function () {
      await expect(testContract.connect(user2).setValueDAO(100))
        .to.be.revertedWithCustomError(testContract, "UnauthorizedRole")
        .withArgs(DAO_ROLE, user2.address);
    });

    it("Should revert when unauthorized user calls BUSINESS-restricted function", async function () {
      await expect(testContract.connect(user1).setValueBusiness(200))
        .to.be.revertedWithCustomError(testContract, "UnauthorizedRole")
        .withArgs(BUSINESS_ROLE, user1.address);
    });
  });

  describe("hasRole Function", function () {
    it("Should return true for granted roles", async function () {
      expect(await testContract.checkUserRole(DAO_ROLE, user1.address)).to.be.true;
      expect(await testContract.checkUserRole(BUSINESS_ROLE, user2.address)).to.be.true;
    });

    it("Should return false for non-granted roles", async function () {
      expect(await testContract.checkUserRole(DAO_ROLE, user2.address)).to.be.false;
      expect(await testContract.checkUserRole(BUSINESS_ROLE, user1.address)).to.be.false;
    });
  });

  describe("getRoleManager Function", function () {
    it("Should return the correct RoleManager address", async function () {
      expect(await testContract.getRoleManagerAddress()).to.equal(await roleManager.getAddress());
    });
  });

  describe("Integration with RoleManager", function () {
    it("Should reflect role changes made in RoleManager", async function () {
      // Initially user1 doesn't have BUSINESS_ROLE
      expect(await testContract.checkUserRole(BUSINESS_ROLE, user1.address)).to.be.false;

      // Grant BUSINESS_ROLE to user1 via RoleManager
      await roleManager.grantRole(BUSINESS_ROLE, user1.address);

      // Now user1 should have BUSINESS_ROLE
      expect(await testContract.checkUserRole(BUSINESS_ROLE, user1.address)).to.be.true;

      // user1 should now be able to call BUSINESS-restricted function
      await testContract.connect(user1).setValueBusiness(300);
      expect(await testContract.value()).to.equal(300);
    });

    it("Should reflect role revocation made in RoleManager", async function () {
      // Initially user1 has DAO_ROLE
      expect(await testContract.checkUserRole(DAO_ROLE, user1.address)).to.be.true;

      // Revoke DAO_ROLE from user1 via RoleManager
      await roleManager.revokeRole(DAO_ROLE, user1.address);

      // Now user1 should not have DAO_ROLE
      expect(await testContract.checkUserRole(DAO_ROLE, user1.address)).to.be.false;

      // user1 should no longer be able to call DAO-restricted function
      await expect(testContract.connect(user1).setValueDAO(400))
        .to.be.revertedWithCustomError(testContract, "UnauthorizedRole")
        .withArgs(DAO_ROLE, user1.address);
    });
  });
});