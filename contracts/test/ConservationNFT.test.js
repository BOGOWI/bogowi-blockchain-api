const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("ConservationNFT", function () {
  let conservationNFT;
  let owner, minter, dao, pauser, user1, user2, user3;
  
  // Token type constants
  const ADOPTION_CERTIFICATE = 1;
  const CONSERVATION_BADGE = 2;
  const DONOR_RECOGNITION = 3;
  const IMPACT_CERTIFICATE = 4;
  
  // Token ID ranges
  const ADOPTION_ID_START = 1000;
  const ADOPTION_ID_END = 1999;
  const BADGE_ID_START = 2000;
  const BADGE_ID_END = 2999;
  const DONOR_ID_START = 3000;
  const DONOR_ID_END = 3999;
  const IMPACT_ID_START = 4000;
  const IMPACT_ID_END = 4999;

  beforeEach(async function () {
    [owner, minter, dao, pauser, user1, user2, user3] = await ethers.getSigners();
    
    const ConservationNFT = await ethers.getContractFactory("ConservationNFT");
    conservationNFT = await ConservationNFT.deploy();
    
    // Grant roles
    const MINTER_ROLE = await conservationNFT.MINTER_ROLE();
    const DAO_ROLE = await conservationNFT.DAO_ROLE();
    const PAUSER_ROLE = await conservationNFT.PAUSER_ROLE();
    
    await conservationNFT.grantRole(MINTER_ROLE, minter.address);
    await conservationNFT.grantRole(DAO_ROLE, dao.address);
    await conservationNFT.grantRole(PAUSER_ROLE, pauser.address);
  });

  describe("Deployment", function () {
    it("Should set correct initial values", async function () {
      expect(await conservationNFT.maxSupplyPerType(ADOPTION_CERTIFICATE)).to.equal(10000);
      expect(await conservationNFT.maxSupplyPerType(CONSERVATION_BADGE)).to.equal(5000);
      expect(await conservationNFT.maxSupplyPerType(DONOR_RECOGNITION)).to.equal(10000);
      expect(await conservationNFT.maxSupplyPerType(IMPACT_CERTIFICATE)).to.equal(5000);
    });
    
    it("Should grant correct roles to deployer", async function () {
      const DEFAULT_ADMIN_ROLE = await conservationNFT.DEFAULT_ADMIN_ROLE();
      const DAO_ROLE = await conservationNFT.DAO_ROLE();
      const MINTER_ROLE = await conservationNFT.MINTER_ROLE();
      
      expect(await conservationNFT.hasRole(DEFAULT_ADMIN_ROLE, owner.address)).to.be.true;
      expect(await conservationNFT.hasRole(DAO_ROLE, owner.address)).to.be.true;
      expect(await conservationNFT.hasRole(MINTER_ROLE, owner.address)).to.be.true;
    });
  });

  describe("Adoption Certificate Minting", function () {
    it("Should mint adoption certificate successfully", async function () {
      const tokenId = 1001;
      const species = "African Elephant";
      const location = "Kenya";
      const tokenUri = "https://example.com/metadata/1001";
      const impactScore = 750;
      
      await expect(
        conservationNFT.connect(minter).mintAdoptionCertificate(
          user1.address,
          tokenId,
          species,
          location,
          tokenUri,
          impactScore
        )
      ).to.emit(conservationNFT, "ConservationNFTMinted")
        .withArgs(user1.address, tokenId, ADOPTION_CERTIFICATE, species, location, impactScore);
      
      expect(await conservationNFT.balanceOf(user1.address, tokenId)).to.equal(1);
      expect(await conservationNFT.tokenExists(tokenId)).to.be.true;
      expect(await conservationNFT.isSoulbound(tokenId)).to.be.true;
      expect(await conservationNFT.uri(tokenId)).to.equal(tokenUri);
      
      const conservationData = await conservationNFT.conservationData(tokenId);
      expect(conservationData.species).to.equal(species);
      expect(conservationData.location).to.equal(location);
      expect(conservationData.impactScore).to.equal(impactScore);
      expect(conservationData.verified).to.be.true;
    });
    
    it("Should reject invalid parameters", async function () {
      await expect(
        conservationNFT.connect(minter).mintAdoptionCertificate(
          ethers.ZeroAddress,
          1001,
          "Species",
          "Location",
          "URI",
          500
        )
      ).to.be.revertedWith("Invalid recipient address");
      
      await expect(
        conservationNFT.connect(minter).mintAdoptionCertificate(
          user1.address,
          999, // Invalid ID range
          "Species",
          "Location",
          "URI",
          500
        )
      ).to.be.revertedWith("INVALID_PARAMETER");
      
      await expect(
        conservationNFT.connect(minter).mintAdoptionCertificate(
          user1.address,
          1001,
          "", // Empty species
          "Location",
          "URI",
          500
        )
      ).to.be.revertedWith("EMPTY_STRING");
      
      await expect(
        conservationNFT.connect(minter).mintAdoptionCertificate(
          user1.address,
          1001,
          "Species",
          "", // Empty location
          "URI",
          500
        )
      ).to.be.revertedWith("EMPTY_STRING");
      
      await expect(
        conservationNFT.connect(minter).mintAdoptionCertificate(
          user1.address,
          1001,
          "Species",
          "Location",
          "", // Empty URI
          500
        )
      ).to.be.revertedWith("URI cannot be empty");
      
      await expect(
        conservationNFT.connect(minter).mintAdoptionCertificate(
          user1.address,
          1001,
          "Species",
          "Location",
          "URI",
          0 // Invalid impact score
        )
      ).to.be.revertedWith("Invalid impact score");
      
      await expect(
        conservationNFT.connect(minter).mintAdoptionCertificate(
          user1.address,
          1001,
          "Species",
          "Location",
          "URI",
          1001 // Invalid impact score
        )
      ).to.be.revertedWith("Invalid impact score");
    });
    
    it("Should reject duplicate token IDs", async function () {
      const tokenId = 1001;
      
      await conservationNFT.connect(minter).mintAdoptionCertificate(
        user1.address,
        tokenId,
        "Species",
        "Location",
        "URI",
        500
      );
      
      await expect(
        conservationNFT.connect(minter).mintAdoptionCertificate(
          user2.address,
          tokenId,
          "Species2",
          "Location2",
          "URI2",
          600
        )
      ).to.be.revertedWith("Token ID already exists");
    });
    
    it("Should reject unauthorized minting", async function () {
      await expect(
        conservationNFT.connect(user1).mintAdoptionCertificate(
          user1.address,
          1001,
          "Species",
          "Location",
          "URI",
          500
        )
      ).to.be.reverted;
    });
  });

  describe("Conservation Badge Minting", function () {
    it("Should mint conservation badge successfully", async function () {
      const tokenId = 2001;
      const achievementType = "Forest Protector";
      const tokenUri = "https://example.com/metadata/2001";
      const impactScore = 850;
      
      await expect(
        conservationNFT.connect(minter).mintConservationBadge(
          user1.address,
          tokenId,
          impactScore,
          tokenUri,
          achievementType
        )
      ).to.emit(conservationNFT, "ConservationNFTMinted")
        .withArgs(user1.address, tokenId, CONSERVATION_BADGE, achievementType, "Global", impactScore);
      
      expect(await conservationNFT.balanceOf(user1.address, tokenId)).to.equal(1);
      expect(await conservationNFT.isSoulbound(tokenId)).to.be.true;
      
      const conservationData = await conservationNFT.conservationData(tokenId);
      expect(conservationData.species).to.equal(achievementType);
      expect(conservationData.location).to.equal("Global");
      expect(conservationData.impactScore).to.equal(impactScore);
    });
    
    it("Should reject invalid badge parameters", async function () {
      await expect(
        conservationNFT.connect(minter).mintConservationBadge(
          user1.address,
          1999, // Wrong ID range
          500,
          "URI",
          "Achievement"
        )
      ).to.be.revertedWith("INVALID_PARAMETER");
      
      await expect(
        conservationNFT.connect(minter).mintConservationBadge(
          user1.address,
          2001,
          500,
          "", // Empty URI
          "Achievement"
        )
      ).to.be.revertedWith("URI cannot be empty");
      
      await expect(
        conservationNFT.connect(minter).mintConservationBadge(
          user1.address,
          2001,
          500,
          "URI",
          "" // Empty achievement type
        )
      ).to.be.revertedWith("EMPTY_STRING");
    });
  });

  describe("Donor NFT Minting", function () {
    it("Should mint donor NFT successfully", async function () {
      const tokenId = 3001;
      const donationAmount = ethers.parseEther("2");
      const tokenUri = "https://example.com/metadata/3001";
      const campaign = "Save the Whales";
      
      await expect(
        conservationNFT.connect(minter).mintDonorNFT(
          user1.address,
          tokenId,
          donationAmount,
          tokenUri,
          campaign
        )
      ).to.emit(conservationNFT, "ConservationNFTMinted");
      
      expect(await conservationNFT.balanceOf(user1.address, tokenId)).to.equal(1);
      expect(await conservationNFT.donationAmounts(tokenId)).to.equal(donationAmount);
      expect(await conservationNFT.isSoulbound(tokenId)).to.be.true;
    });
    
    it("Should calculate correct impact scores for donations", async function () {
      const testCases = [
        { amount: ethers.parseEther("15"), expectedScore: 1000 },
        { amount: ethers.parseEther("7"), expectedScore: 750 },
        { amount: ethers.parseEther("2"), expectedScore: 500 },
        { amount: ethers.parseEther("0.7"), expectedScore: 250 },
        { amount: ethers.parseEther("0.2"), expectedScore: 100 },
        { amount: ethers.parseEther("0.05"), expectedScore: 50 }
      ];
      
      for (let i = 0; i < testCases.length; i++) {
        const tokenId = 3001 + i;
        await conservationNFT.connect(minter).mintDonorNFT(
          user1.address,
          tokenId,
          testCases[i].amount,
          `URI${i}`,
          `Campaign${i}`
        );
        
        const conservationData = await conservationNFT.conservationData(tokenId);
        expect(conservationData.impactScore).to.equal(testCases[i].expectedScore);
      }
    });
    
    it("Should reject zero donation amount", async function () {
      await expect(
        conservationNFT.connect(minter).mintDonorNFT(
          user1.address,
          3001,
          0,
          "URI",
          "Campaign"
        )
      ).to.be.revertedWith("ZERO_AMOUNT");
    });
  });

  describe("Impact Certificate Minting", function () {
    it("Should mint impact certificate with DAO role", async function () {
      const tokenId = 4001;
      const project = "Reforestation Project";
      const location = "Amazon";
      const impactScore = 900;
      const tokenUri = "https://example.com/metadata/4001";
      
      await expect(
        conservationNFT.connect(dao).mintImpactCertificate(
          user1.address,
          tokenId,
          project,
          location,
          impactScore,
          tokenUri
        )
      ).to.emit(conservationNFT, "ConservationNFTMinted")
        .withArgs(user1.address, tokenId, IMPACT_CERTIFICATE, project, location, impactScore);
      
      expect(await conservationNFT.balanceOf(user1.address, tokenId)).to.equal(1);
      expect(await conservationNFT.isSoulbound(tokenId)).to.be.true;
    });
    
    it("Should reject non-DAO minting", async function () {
      await expect(
        conservationNFT.connect(minter).mintImpactCertificate(
          user1.address,
          4001,
          "Project",
          "Location",
          500,
          "URI"
        )
      ).to.be.reverted;
    });
  });

  describe("Soulbound Token Logic", function () {
    beforeEach(async function () {
      await conservationNFT.connect(minter).mintAdoptionCertificate(
        user1.address,
        1001,
        "Species",
        "Location",
        "URI",
        500
      );
    });
    
    it("Should prevent transfer of soulbound tokens", async function () {
      await expect(
        conservationNFT.connect(user1).safeTransferFrom(
          user1.address,
          user2.address,
          1001,
          1,
          "0x"
        )
      ).to.be.revertedWith("Soulbound token cannot be transferred");
    });
    
    it("Should prevent batch transfer of soulbound tokens", async function () {
      await expect(
        conservationNFT.connect(user1).safeBatchTransferFrom(
          user1.address,
          user2.address,
          [1001],
          [1],
          "0x"
        )
      ).to.be.revertedWith("Soulbound token cannot be transferred");
    });
  });

  describe("Admin Functions", function () {
    it("Should set max supply for token types", async function () {
      const newMaxSupply = 15000;
      
      await expect(
        conservationNFT.connect(dao).setMaxSupply(ADOPTION_CERTIFICATE, newMaxSupply)
      ).to.emit(conservationNFT, "MaxSupplySet")
        .withArgs(ADOPTION_CERTIFICATE, newMaxSupply);
      
      expect(await conservationNFT.maxSupplyPerType(ADOPTION_CERTIFICATE)).to.equal(newMaxSupply);
    });
    
    it("Should reject invalid max supply updates", async function () {
      await expect(
        conservationNFT.connect(dao).setMaxSupply(5, 1000) // Invalid token type
      ).to.be.revertedWith("INVALID_PARAMETER");
      
      // Mint some tokens first
      await conservationNFT.connect(minter).mintAdoptionCertificate(
        user1.address,
        1001,
        "Species",
        "Location",
        "URI",
        500
      );
      
      await expect(
        conservationNFT.connect(dao).setMaxSupply(ADOPTION_CERTIFICATE, 0) // Less than minted
      ).to.be.revertedWith("INVALID_PARAMETER");
    });
    
    it("Should update token URI", async function () {
      const tokenId = 1001;
      const newUri = "https://newuri.com/metadata/1001";
      
      await conservationNFT.connect(minter).mintAdoptionCertificate(
        user1.address,
        tokenId,
        "Species",
        "Location",
        "URI",
        500
      );
      
      await expect(
        conservationNFT.connect(dao).updateTokenURI(tokenId, newUri)
      ).to.emit(conservationNFT, "TokenURIUpdated")
        .withArgs(tokenId, newUri);
      
      expect(await conservationNFT.uri(tokenId)).to.equal(newUri);
    });
    
    it("Should reject URI update for non-existent token", async function () {
      await expect(
        conservationNFT.connect(dao).updateTokenURI(9999, "newuri")
      ).to.be.revertedWith("DOES_NOT_EXIST");
    });
    
    it("Should reject empty URI", async function () {
      const tokenId = 1001;
      
      await conservationNFT.connect(minter).mintAdoptionCertificate(
        user1.address,
        tokenId,
        "Species",
        "Location",
        "URI",
        500
      );
      
      await expect(
        conservationNFT.connect(dao).updateTokenURI(tokenId, "")
      ).to.be.revertedWith("EMPTY_STRING");
    });
  });

  describe("Pause Functionality", function () {
    it("Should pause and unpause contract", async function () {
      await conservationNFT.connect(pauser).pause();
      expect(await conservationNFT.paused()).to.be.true;
      
      await expect(
        conservationNFT.connect(minter).mintAdoptionCertificate(
          user1.address,
          1001,
          "Species",
          "Location",
          "URI",
          500
        )
      ).to.be.revertedWithCustomError(conservationNFT, "EnforcedPause");
      
      await conservationNFT.connect(pauser).unpause();
      expect(await conservationNFT.paused()).to.be.false;
      
      // Should work after unpause
      await conservationNFT.connect(minter).mintAdoptionCertificate(
        user1.address,
        1001,
        "Species",
        "Location",
        "URI",
        500
      );
    });
    
    it("Should reject pause from unauthorized user", async function () {
      await expect(
        conservationNFT.connect(user1).pause()
      ).to.be.revertedWith("UNAUTHORIZED");
    });
  });

  describe("Utility Functions", function () {
    it("Should return correct token type", async function () {
      expect(await conservationNFT.getTokenType(1500)).to.equal(ADOPTION_CERTIFICATE);
      expect(await conservationNFT.getTokenType(2500)).to.equal(CONSERVATION_BADGE);
      expect(await conservationNFT.getTokenType(3500)).to.equal(DONOR_RECOGNITION);
      expect(await conservationNFT.getTokenType(4500)).to.equal(IMPACT_CERTIFICATE);
      
      await expect(
        conservationNFT.getTokenType(999)
      ).to.be.revertedWith("Invalid token ID");
    });
    
    it("Should support correct interfaces", async function () {
      // ERC1155 interface
      expect(await conservationNFT.supportsInterface("0xd9b67a26")).to.be.true;
      // AccessControl interface
      expect(await conservationNFT.supportsInterface("0x7965db0b")).to.be.true;
    });
  });

  describe("Supply Limits", function () {
    it("Should enforce max supply limits", async function () {
      // Set a low max supply for testing
      await conservationNFT.connect(dao).setMaxSupply(ADOPTION_CERTIFICATE, 1);
      
      // First mint should succeed
      await conservationNFT.connect(minter).mintAdoptionCertificate(
        user1.address,
        1001,
        "Species",
        "Location",
        "URI",
        500
      );
      
      // Second mint should fail
      await expect(
        conservationNFT.connect(minter).mintAdoptionCertificate(
          user2.address,
          1002,
          "Species2",
          "Location2",
          "URI2",
          500
        )
      ).to.be.revertedWith("MAX_REACHED");
    });
  });

  describe("Edge Cases", function () {
    it("Should handle boundary token IDs correctly", async function () {
      // Test boundary IDs for each range
      await conservationNFT.connect(minter).mintAdoptionCertificate(
        user1.address,
        ADOPTION_ID_START,
        "Species",
        "Location",
        "URI",
        500
      );
      
      await conservationNFT.connect(minter).mintAdoptionCertificate(
        user1.address,
        ADOPTION_ID_END,
        "Species",
        "Location",
        "URI",
        500
      );
      
      await conservationNFT.connect(minter).mintConservationBadge(
        user1.address,
        BADGE_ID_START,
        500,
        "URI",
        "Achievement"
      );
      
      await conservationNFT.connect(minter).mintConservationBadge(
        user1.address,
        BADGE_ID_END,
        500,
        "URI",
        "Achievement"
      );
    });
    
    it("Should handle maximum impact scores", async function () {
      await conservationNFT.connect(minter).mintAdoptionCertificate(
        user1.address,
        1001,
        "Species",
        "Location",
        "URI",
        1000 // Maximum score
      );
      
      const conservationData = await conservationNFT.conservationData(1001);
      expect(conservationData.impactScore).to.equal(1000);
    });
  });
});