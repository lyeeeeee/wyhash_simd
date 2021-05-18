#include "textflag.h"

DATA mask1_3<>+0x00(SB)/8, $0xe7037ed1a0b428db
DATA mask1_3<>+0x08(SB)/8, $0
DATA mask1_3<>+0x10(SB)/8, $0x589965cc75374cc3
DATA mask1_3<>+0x18(SB)/8, $0
GLOBL mask1_3<>(SB), RODATA|NOPTR, $32

DATA mask2_4<>+0x00(SB)/8, $0
DATA mask2_4<>+0x08(SB)/8, $0xffffffffffffffff
DATA mask2_4<>+0x10(SB)/8, $0
DATA mask2_4<>+0x18(SB)/8, $0xffffffffffffffff
GLOBL mask2_4<>(SB), RODATA|NOPTR, $32

TEXT ·avx2(SB), $0
    MOVQ        result+0(FP), AX
    MOVQ        paddr+8(FP), BX
    MOVQ        seeds+16(FP), CX


    VMOVDQU     mask2_4<>(SB), Y0
	VPBROADCASTQ    (AX), Y2
	VPAND       Y0, Y2, Y0

    VMOVDQU     (BX), Y3

    VMOVDQU     mask1_3<>(SB), Y1
	VPAND       Y1, Y3, Y1

    MOVQ        $1, ret+24(FP)

    RET
