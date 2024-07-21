// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library Pairing {
    uint256 constant Prime_Field_Mod =
        21888242871839275222246405745257275088696311157297823662689037894645226208583;

    /*
    G1Point represents a point in G1 field
    */
    struct G1Point {
        uint256 X;
        uint256 Y;
    }

    struct G2Point {
        uint256[2] X;
        uint256[2] Y;
    }

    /*
     * @return The negation of p, i.e. p.plus(p.negate()) should be zero.
     */
    function negate(G1Point memory p) internal pure returns (G1Point memory) {
        // The prime q in the base field F_q for G1
        if (p.X == 0 && p.Y == 0) {
            return G1Point(0, 0);
        } else {
            return G1Point(p.X, Prime_Field_Mod - (p.Y % Prime_Field_Mod));
        }
    }

    /*
     * @return r the sum of two points of G1
     */
    function plus(G1Point memory p1, G1Point memory p2)
        internal
        view
        returns (G1Point memory r)
    {
        uint256[4] memory input;
        input[0] = p1.X;
        input[1] = p1.Y;
        input[2] = p2.X;
        input[3] = p2.Y;
        bool success;

        // solium-disable-next-line security/no-inline-assembly
        assembly {
            success := staticcall(sub(gas(), 2000), 6, input, 0xc0, r, 0x60)
            // Use "invalid" to make gas estimation work
            switch success
            case 0 {
                invalid()
            }
        }

        require(success, "pairing-add-failed");
    }

    /*
     * @return r the product of a point on G1 and a scalar, i.e.
     *         p == p.scalar_mul(1) and p.plus(p) == p.scalar_mul(2) for all
     *         points p.
     */
    function scalar_mul(G1Point memory p, uint256 s)
        internal
        view
        returns (G1Point memory r)
    {
        uint256[3] memory input;
        input[0] = p.X;
        input[1] = p.Y;
        input[2] = s;
        bool success;
        // solium-disable-next-line security/no-inline-assembly
        assembly {
            success := staticcall(sub(gas(), 2000), 7, input, 0x80, r, 0x60)
            // Use "invalid" to make gas estimation work
            switch success
            case 0 {
                invalid()
            }
        }
        require(success, "pairing-mul-failed");
    }

    function pairing(
        G1Point memory a1,
        G2Point memory a2,
        G1Point memory b1,
        G2Point memory b2,
        G1Point memory c1,
        G2Point memory c2,
        G1Point memory d1,
        G2Point memory d2
    ) internal view returns (bool) {
        G1Point[4] memory p1 = [a1, b1, c1, d1];
        G2Point[4] memory p2 = [a2, b2, c2, d2];

        uint256 inputSize = 24;
        uint256[] memory input = new uint256[](inputSize);

        for (uint256 i = 0; i < 4; i++) {
            uint256 j = i * 6;
            input[j + 0] = p1[i].X;
            input[j + 1] = p1[i].Y;
            input[j + 2] = p2[i].X[0];
            input[j + 3] = p2[i].X[1];
            input[j + 4] = p2[i].Y[0];
            input[j + 5] = p2[i].Y[1];
        }

        uint256[1] memory out;
        bool success;

        // solium-disable-next-line security/no-inline-assembly
        assembly {
            success := staticcall(
                sub(gas(), 2000),
                8,
                add(input, 0x20),
                mul(inputSize, 0x20),
                out,
                0x20
            )
            // Use "invalid" to make gas estimation work
            switch success
            case 0 {
                invalid()
            }
        }

        require(success, "pairing-opcode-failed");

        return out[0] != 0;
    }
}

contract Homework5 {
    uint256 constant CurveOrder =
        21888242871839275222246405745257275088548364400416034343698204186575808495617;
    uint256 constant Prime_Field_Mod =
        21888242871839275222246405745257275088696311157297823662689037894645226208583;

    using Pairing for *;

    struct VerifyingKey {
        Pairing.G1Point alfa1;
        Pairing.G2Point beta2;
        Pairing.G2Point gamma2;
        Pairing.G2Point delta2;
        Pairing.G1Point[7] IC;
    }

    struct Proof {
        Pairing.G1Point A;
        Pairing.G2Point B;
        Pairing.G1Point C;
    }

      function verifyingKey() internal pure returns (VerifyingKey memory vk) {
        vk.alfa1 = Pairing.G1Point(
            uint256(1368015179489954701390400359078579693043519447331113978918064868415326638035), 
            uint256(9918110051302171585080402603319702774565515993150576347155970296011118125764)
        );
        vk.beta2 = Pairing.G2Point(
            [
                uint256(11559732032986387107991004021392285783925812861821192530917403151452391805634), 
                uint256(10857046999023057135944570762232829481370756359578518086990519993285655852781)
            ], [
                uint256(4082367875863433681332203403145435568316851327593401208105741076214120093531), 
                uint256(8495653923123431417604973247489272438418190587263600148770280649306958101930)
            ]
        );
        vk.gamma2 = Pairing.G2Point(
            [
                uint256(14583779054894525174450323658765874724019480979794335525732096752006891875705), 
                uint256(18029695676650738226693292988307914797657423701064905010927197838374790804409)
            ], [
                uint256(11474861747383700316476719153975578001603231366361248090558603872215261634898), 
                uint256(2140229616977736810657479771656733941598412651537078903776637920509952744750)
            ]
        );
        vk.delta2 = Pairing.G2Point(
            [
                uint256(14583779054894525174450323658765874724019480979794335525732096752006891875705), 
                uint256(18029695676650738226693292988307914797657423701064905010927197838374790804409)
            ], [
                uint256(11474861747383700316476719153975578001603231366361248090558603872215261634898), 
                uint256(2140229616977736810657479771656733941598412651537078903776637920509952744750)
            ]
        );
      }
}
