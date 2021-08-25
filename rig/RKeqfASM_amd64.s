#include "textflag.h"

DATA ulim+0x00(SB)/8, $(0x40A7840000000000)
DATA ulim+0x08(SB)/8, $(0x4079A00000000000)
DATA ulim+0x10(SB)/8, $(0x40A7840000000000)
DATA ulim+0x18(SB)/8, $(0x40A7840000000000)
GLOBL ulim(SB), RODATA, $32

TEXT ·eqfar(SB), NOSPLIT, $0
    VMOVUPD r+0(FP), Y0
    VMOVUPD v+24(FP), Y1
    MOVQ m+48(FP), AX
    VXORPD Y2, Y2, Y2
    //Y2 = (0,0,0,0) //default
    VPERMQ $(0x24), Y0, Y8
    //Y8 = (x, y, z, x)
    VMOVUPD ulim(SB), Y9
    //Y9 = (3000, 400, 3000, 3000)
    VPCMPEQQ Y3, Y3, Y3
    //Y3 all set
    VPSLLQ $(63), Y3, Y3
    //Y3=(-0, -0, -0, -0)
    VANDPD Y3, Y8, Y10
    //Y10 = (sign(x), sign(y), sign(z), sign(x))
    VXORPD Y10, Y8, Y10
    //Y10 = (|x|, |y|, |z|, |x|)
    VCMPPD $(0x05), Y9, Y10, Y11
    //Y10[i] >= Y9[i] ? Y11[i]=0xffffffff : Y11[i]=0x00000000
    VTESTPD Y11, Y11
    // sign(Y11[i]) == 0 ? ZF=ZF : ZF&=0
    // if any i : Y10[i] >= Y9[i] then ZF = 0
    JNZ end
    

    MOVQ $(0x4044800000000000), R10
    VXORPD Y10, Y10, Y10
    MOVQ R10, X10
    VPERMQ $(0x51), Y10, Y10
    //Y10 = (0, 41, 0, 0)
    MOVQ $(0x3fb999999999999a), R10
    MOVQ R10, X9
    VPERMQ $(0x0), Y9, Y9
    //Y9 = (0.1, 0.1, 0.1, 0.1)
    VFMADD231PD Y8, Y9, Y10
    //Y10 = (0.1x, 0.1y+41, 0.1z, 0.1x)
    VANDPD Y10, Y3, Y2
    //Y2 = (sign(0.1x)0, sign(0.1y+41)0, sign(0.1z)0, sign(0.1x)0)
    VXORPD Y10, Y2, Y11
    //Y11 = (|0.1x|, |0.1y+41|, |0.1z|, |0.1x|)
    VROUNDPD $(0x01), Y11, Y8
    //Y8 = floor(|0.1x|, |0.1y+41|, |0.1z|, |0.1x|)
    VSUBPD Y8, Y11, Y0
    //Y0 = (|0.1x|, |0.1y+41|, |0.1z|, |0.1x|) - floor(|0.1x|, |0.1y+41|, |0.1z|, |0.1x|)
    //Y0 = Y11 - Y8 // res

    VCVTPD2DQY Y8, X8
    //X8 4 int32s

    MOVQ X8, BX
    MOVQ BX, CX
    SHRQ $(32), CX
    VPERMQ $(0x01), Y8, Y8
    MOVQ X8, DX
    MOVQ $(0x00000000ffffffff), R8
    ANDQ R8, BX
    ANDQ R8, CX
    ANDQ R8, DX
    //BX = int(|x/cm|)
    //CX = int(|y/cm+40|)
    //DX = int(|z/cm|)
    LEAQ (BX)(BX*2), BX
    LEAQ (CX)(CX*2), CX
    LEAQ (DX)(DX*2), DX
    //IMUL3Q $(0x012d), CX, CX
    //IMUL3Q $(0x5f3d), DX, DX
    IMUL3Q $(0x012e), BX, BX
    IMUL3Q $(0x16444), CX, CX
    // BX *= 302
    // CX *= 302*302
    LEAQ (AX)(BX*8), AX
    LEAQ (AX)(CX*8), AX
    LEAQ (AX)(DX*8), AX

    //0, +z
    VPERMQ $(0xaa), Y0, Y3
    VMOVUPD 0x00000000(AX), Y8
    VMOVUPD 0x00000018(AX), Y9
    VFNMADD231PD Y8, Y3, Y8
    VFMADD231PD Y9, Y3, Y8

    //+x, +xz
    VMOVUPD 0x00001c50(AX), Y9
    VMOVUPD 0x00001c68(AX), Y10
    VFNMADD231PD Y9, Y3, Y9
    VFMADD231PD Y10, Y3, Y9

    //+y, +yz
    VMOVUPD 0x00216660(AX), Y10
    VMOVUPD 0x00216678(AX), Y11
    VFNMADD231PD Y10, Y3, Y10
    VFMADD231PD Y11, Y3, Y10

    //+xy, +xyz
    VMOVUPD 0x002182b0(AX), Y11
    VMOVUPD 0x002182c8(AX), Y12
    VFNMADD231PD Y11, Y3, Y11
    VFMADD231PD Y12, Y3, Y11

    VPERMQ $(0x00), Y0, Y3
    VFNMADD231PD Y8, Y3, Y8
    VFMADD231PD Y9, Y3, Y8
    VFNMADD231PD Y10, Y3, Y10
    VFMADD231PD Y11, Y3, Y10

    VPERMQ $(0x55), Y0, Y3
    VFNMADD231PD Y8, Y3, Y8
    VFMADD231PD Y10, Y3, Y8
    VXORPD Y8, Y2, Y8
    //Y8 = (Bx, By, Bz, ?)
    //Y1 = (vx, vy, vz, ?)
    VPERMQ $(0x09), Y1, Y2
    VPERMQ $(0x12), Y8, Y10
    VPERMQ $(0x12), Y1, Y11
    VPERMQ $(0x09), Y8, Y12
    VMULPD Y2, Y10, Y2
    VFMSUB231PD Y11, Y12, Y2
    //VMOVUPD Y8, Y2

end:
    VPERMQ $(0x90), Y2, Y2
    VMOVUPD Y2, r2m+88(FP)
    VMOVUPD Y1, r1+72(FP)
    MOVQ X2, r2+96(FP)
    RET




TEXT ·Fetcheqfar(SB), NOSPLIT, $0
    VMOVUPD r+0(FP), Y0
    VMOVUPD v+24(FP), Y1
    MOVQ m+48(FP), AX
    VXORPD Y2, Y2, Y2
    //Y2 = (0,0,0,0) //default
    VPERMQ $(0x24), Y0, Y8
    //Y8 = (x, y, z, x)
    VMOVUPD ulim(SB), Y9
    //Y9 = (3000, 400, 3000, 3000)
    VPCMPEQQ Y3, Y3, Y3
    //Y3 all set
    VPSLLQ $(63), Y3, Y3
    //Y3=(-0, -0, -0, -0)
    VANDPD Y3, Y8, Y10
    //Y10 = (sign(x), sign(y), sign(z), sign(x))
    VXORPD Y10, Y8, Y10
    //Y10 = (|x|, |y|, |z|, |x|)
    VCMPPD $(0x05), Y9, Y10, Y11
    //Y10[i] >= Y9[i] ? Y11[i]=0xffffffff : Y11[i]=0x00000000
    VTESTPD Y11, Y11
    // sign(Y11[i]) == 0 ? ZF=ZF : ZF&=0
    // if any i : Y10[i] >= Y9[i] then ZF = 0
    JNZ end
    

    MOVQ $(0x4044800000000000), R10
    VXORPD Y10, Y10, Y10
    MOVQ R10, X10
    VPERMQ $(0x51), Y10, Y10
    //Y10 = (0, 41, 0, 0)
    MOVQ $(0x3fb999999999999a), R10
    MOVQ R10, X9
    VPERMQ $(0x0), Y9, Y9
    //Y9 = (0.1, 0.1, 0.1, 0.1)
    VFMADD231PD Y8, Y9, Y10
    //Y10 = (0.1x, 0.1y+41, 0.1z, 0.1x)
    VANDPD Y10, Y3, Y2
    //Y2 = (sign(0.1x)0, sign(0.1y+41)0, sign(0.1z)0, sign(0.1x)0)
    VXORPD Y10, Y2, Y11
    //Y11 = (|0.1x|, |0.1y+41|, |0.1z|, |0.1x|)
    VROUNDPD $(0x01), Y11, Y8
    //Y8 = floor(|0.1x|, |0.1y+41|, |0.1z|, |0.1x|)
    VSUBPD Y8, Y11, Y0
    //Y0 = (|0.1x|, |0.1y+41|, |0.1z|, |0.1x|) - floor(|0.1x|, |0.1y+41|, |0.1z|, |0.1x|)
    //Y0 = Y11 - Y8 // res

    VCVTPD2DQY Y8, X8
    //X8 4 int32s

    MOVQ X8, BX
    MOVQ BX, CX
    SHRQ $(32), CX
    VPERMQ $(0x01), Y8, Y8
    MOVQ X8, DX
    MOVQ $(0x00000000ffffffff), R8
    ANDQ R8, BX
    ANDQ R8, CX
    ANDQ R8, DX
    //BX = int(|x/cm|)
    //CX = int(|y/cm+40|)
    //DX = int(|z/cm|)
    LEAQ (BX)(BX*2), BX
    LEAQ (CX)(CX*2), CX
    LEAQ (DX)(DX*2), DX
    IMUL3Q $(0x012e), BX, BX
    IMUL3Q $(0x16444), CX, CX
    // BX *= 302
    // CX *= 302*302
    LEAQ (AX)(BX*8), AX
    LEAQ (AX)(CX*8), AX
    LEAQ (AX)(DX*8), AX

    //0, +z
    VPERMQ $(0xaa), Y0, Y3
    VMOVUPD 0x00000000(AX), Y8
    VMOVUPD 0x00000018(AX), Y9
    VFNMADD231PD Y8, Y3, Y8
    VFMADD231PD Y9, Y3, Y8

    //+x, +xz
    VMOVUPD 0x00001c50(AX), Y9
    VMOVUPD 0x00001c68(AX), Y10
    VFNMADD231PD Y9, Y3, Y9
    VFMADD231PD Y10, Y3, Y9

    //+y, +yz
    VMOVUPD 0x00216660(AX), Y10
    VMOVUPD 0x00216678(AX), Y11
    VFNMADD231PD Y10, Y3, Y10
    VFMADD231PD Y11, Y3, Y10

    //+xy, +xyz
    VMOVUPD 0x002182b0(AX), Y11
    VMOVUPD 0x002182c8(AX), Y12
    VFNMADD231PD Y11, Y3, Y11
    VFMADD231PD Y12, Y3, Y11

    VPERMQ $(0x00), Y0, Y3
    VFNMADD231PD Y8, Y3, Y8
    VFMADD231PD Y9, Y3, Y8
    VFNMADD231PD Y10, Y3, Y10
    VFMADD231PD Y11, Y3, Y10

    VPERMQ $(0x55), Y0, Y3
    VFNMADD231PD Y8, Y3, Y8
    VFMADD231PD Y10, Y3, Y8
    VXORPD Y8, Y2, Y8
    //Y8 = (Bx, By, Bz, ?)
    //Y1 = (vx, vy, vz, ?)
    VPERMQ $(0x09), Y1, Y2
    VPERMQ $(0x12), Y8, Y10
    VPERMQ $(0x12), Y1, Y11
    VPERMQ $(0x09), Y8, Y12
    VMULPD Y2, Y10, Y2
    VFMSUB231PD Y11, Y12, Y2
    

end:
    VMOVUPD Y2, r1+72(FP)
    RET


