
TEXT ·accumAVX2(SB), NOSPLIT, $0-32
    MOVQ         seed+0(FP), AX
    MOVQ         paddr+8(FP), CX
    MOVQ         i+16(FP), BP
    VMOVDQU      (AX), Y1
    VMOVDQU      32(AX), Y2
	VPBROADCASTQ seed1<>+0(SB), Y0
    VPXOR
accum:
	CMPQ     BP, $0x40
	JLE      finalize

	VMOVDQU  (CX), Y0

	ADDQ     $0x00000040, CX
	SUBQ     $0x00000040, BP
	JMP      return

return:
	VMOVDQU Y1, (AX)
	VMOVDQU Y2, 32(AX)
	RET
