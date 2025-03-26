****************
Platform: PPC-64
Code:0x43 0x20 0x0c 0x07 0x41 0x56 0xff 0x17 0x80 0x20 0x00 0x00 0x80 0x3f 0x00 0x00 0x10 0x43 0x23 0x0e 0xd0 0x44 0x00 0x80 0x4c 0x43 0x22 0x02 0x2d 0x03 0x00 0x80 0x7c 0x43 0x20 0x14 0x7c 0x43 0x20 0x93 0x4f 0x20 0x00 0x21 0x4c 0xc8 0x00 0x21 0x40 0x82 0x00 0x14 
Disasm:
0x1000:	bdnzla+	0xc04
	op_count: 1
		operands[0].type: IMM = 0xc04
	Branch hint: 1

0x1004:	bdztla	4*cr5+eq, 0xffffffffffffff14
	op_count: 2
		operands[0].type: REG = cr5eq
		operands[1].type: IMM = 0xffffffffffffff14
	Branch code: 76
	Branch hint: 1

0x1008:	lwz	r1, 0(0)
	op_count: 2
		operands[0].type: REG = r1
		operands[1].type: MEM

0x100c:	lwz	r1, 0(r31)
	op_count: 2
		operands[0].type: REG = r1
		operands[1].type: MEM
			operands[1].mem.base: REG = r31

0x1010:	vpkpx	v2, v3, v4
	op_count: 3
		operands[0].type: REG = v2
		operands[1].type: REG = v3
		operands[2].type: REG = v4

0x1014:	stfs	f2, 0x80(r4)
	op_count: 2
		operands[0].type: REG = f2
		operands[1].type: MEM
			operands[1].mem.base: REG = r4
			operands[1].mem.disp: 0x80

0x1018:	crand	cr0eq, cr0un, cr1lt
	op_count: 3
		operands[0].type: REG = cr0eq
		operands[1].type: REG = cr0un
		operands[2].type: REG = cr1lt

0x101c:	cmpwi	cr2, r3, 0x80
	op_count: 3
		operands[0].type: REG = cr2
		operands[1].type: REG = r3
		operands[2].type: IMM = 0x80

0x1020:	addc	r2, r3, r4
	op_count: 3
		operands[0].type: REG = r2
		operands[1].type: REG = r3
		operands[2].type: REG = r4

0x1024:	mulhd.	r2, r3, r4
	op_count: 3
		operands[0].type: REG = r2
		operands[1].type: REG = r3
		operands[2].type: REG = r4
	Update-CR0: True

0x1028:	bdnzlrl+	
	Branch hint: 1

0x102c:	bgelrl-	cr2
	op_count: 1
		operands[0].type: REG = cr2
	Branch code: 4
	Branch hint: 2

0x1030:	bne	0x1044
	op_count: 1
		operands[0].type: IMM = 0x1044
	Branch code: 68

0x1034:

****************
Platform: PPC-64 + QPX
Code:0x10 0x60 0x2a 0x10 0x10 0x64 0x28 0x88 0x7c 0x4a 0x5d 0x0f 
Disasm:
0x1000:	qvfabs	q3, q5
	op_count: 2
		operands[0].type: REG = q3
		operands[1].type: REG = q5

0x1004:	qvfand	q3, q4, q5
	op_count: 3
		operands[0].type: REG = q3
		operands[1].type: REG = q4
		operands[2].type: REG = q5

0x1008:	qvstfsxa	q2, r10, r11
	op_count: 3
		operands[0].type: REG = q2
		operands[1].type: REG = r10
		operands[2].type: REG = r11

0x100c:

****************
Platform: PPC + PS
Code:0x10 0x00 0x1f 0xec 0xe0 0x6d 0x80 0x04 0xe4 0x6d 0x80 0x04 0x10 0x60 0x1c 0x4c 0x10 0x60 0x1c 0x0c 0xf0 0x6d 0x80 0x04 0xf4 0x6d 0x80 0x04 0x10 0x60 0x1c 0x4e 0x10 0x60 0x1c 0x0e 0x10 0x60 0x1a 0x10 0x10 0x60 0x1a 0x11 0x10 0x63 0x20 0x2a 0x10 0x63 0x20 0x2b 0x10 0x83 0x20 0x40 0x10 0x83 0x20 0xc0 0x10 0x83 0x20 0x00 0x10 0x83 0x20 0x80 0x10 0x63 0x20 0x24 0x10 0x63 0x20 0x25 0x10 0x63 0x29 0x3a 0x10 0x63 0x29 0x3b 0x10 0x63 0x29 0x1c 0x10 0x63 0x29 0x1d 0x10 0x63 0x29 0x1e 0x10 0x63 0x29 0x1f 0x10 0x63 0x24 0x20 0x10 0x63 0x24 0x21 0x10 0x63 0x24 0x60 0x10 0x63 0x24 0x61 0x10 0x63 0x24 0xa0 0x10 0x63 0x24 0xa1 0x10 0x63 0x24 0xe0 0x10 0x63 0x24 0xe1 0x10 0x60 0x20 0x90 0x10 0x60 0x20 0x91 0x10 0x63 0x29 0x38 0x10 0x63 0x29 0x39 0x10 0x63 0x01 0x32 0x10 0x63 0x01 0x33 0x10 0x63 0x01 0x18 0x10 0x63 0x01 0x19 0x10 0x63 0x01 0x1a 0x10 0x63 0x01 0x1b 0x10 0x60 0x19 0x10 0x10 0x60 0x19 0x11 0x10 0x60 0x18 0x50 0x10 0x60 0x18 0x51 0x10 0x63 0x29 0x3e 0x10 0x63 0x29 0x3f 0x10 0x63 0x29 0x3c 0x10 0x63 0x29 0x3d 0x10 0x60 0x18 0x30 0x10 0x60 0x18 0x31 0x10 0x60 0x18 0x34 0x10 0x60 0x18 0x35 0x10 0x63 0x29 0x2e 0x10 0x63 0x29 0x2f 0x10 0x63 0x20 0x28 0x10 0x63 0x20 0x29 0x10 0x63 0x29 0x14 0x10 0x63 0x29 0x15 0x10 0x63 0x29 0x16 0x10 0x63 0x29 0x17 
Disasm:
0x1000:	dcbz_l	r0, r3
	op_count: 2
		operands[0].type: REG = r0
		operands[1].type: REG = r3

0x1004:	psq_l	f3, 4(r13), 1, 0
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: MEM
			operands[1].mem.base: REG = r13
			operands[1].mem.disp: 0x4
		operands[2].type: IMM = 0x1
		operands[3].type: IMM = 0x0

0x1008:	psq_lu	f3, 4(r13), 1, 0
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: MEM
			operands[1].mem.base: REG = r13
			operands[1].mem.disp: 0x4
		operands[2].type: IMM = 0x1
		operands[3].type: IMM = 0x0

0x100c:	psq_lux	f3, r0, r3, 1, 0
	op_count: 5
		operands[0].type: REG = f3
		operands[1].type: REG = r0
		operands[2].type: REG = r3
		operands[3].type: IMM = 0x1
		operands[4].type: IMM = 0x0

0x1010:	psq_lx	f3, r0, r3, 1, 0
	op_count: 5
		operands[0].type: REG = f3
		operands[1].type: REG = r0
		operands[2].type: REG = r3
		operands[3].type: IMM = 0x1
		operands[4].type: IMM = 0x0

0x1014:	psq_st	f3, 4(r13), 1, 0
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: MEM
			operands[1].mem.base: REG = r13
			operands[1].mem.disp: 0x4
		operands[2].type: IMM = 0x1
		operands[3].type: IMM = 0x0

0x1018:	psq_stu	f3, 4(r13), 1, 0
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: MEM
			operands[1].mem.base: REG = r13
			operands[1].mem.disp: 0x4
		operands[2].type: IMM = 0x1
		operands[3].type: IMM = 0x0

0x101c:	psq_stux	f3, r0, r3, 1, 0
	op_count: 5
		operands[0].type: REG = f3
		operands[1].type: REG = r0
		operands[2].type: REG = r3
		operands[3].type: IMM = 0x1
		operands[4].type: IMM = 0x0

0x1020:	psq_stx	f3, r0, r3, 1, 0
	op_count: 5
		operands[0].type: REG = f3
		operands[1].type: REG = r0
		operands[2].type: REG = r3
		operands[3].type: IMM = 0x1
		operands[4].type: IMM = 0x0

0x1024:	ps_abs	f3, f3
	op_count: 2
		operands[0].type: REG = f3
		operands[1].type: REG = f3

0x1028:	ps_abs.	f3, f3
	op_count: 2
		operands[0].type: REG = f3
		operands[1].type: REG = f3
	Update-CR0: True

0x102c:	ps_add	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4

0x1030:	ps_add.	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
	Update-CR0: True

0x1034:	ps_cmpo0	cr1, f3, f4
	op_count: 3
		operands[0].type: REG = cr1
		operands[1].type: REG = f3
		operands[2].type: REG = f4

0x1038:	ps_cmpo1	cr1, f3, f4
	op_count: 3
		operands[0].type: REG = cr1
		operands[1].type: REG = f3
		operands[2].type: REG = f4

0x103c:	ps_cmpu0	cr1, f3, f4
	op_count: 3
		operands[0].type: REG = cr1
		operands[1].type: REG = f3
		operands[2].type: REG = f4

0x1040:	ps_cmpu1	cr1, f3, f4
	op_count: 3
		operands[0].type: REG = cr1
		operands[1].type: REG = f3
		operands[2].type: REG = f4

0x1044:	ps_div	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4

0x1048:	ps_div.	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
	Update-CR0: True

0x104c:	ps_madd	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5

0x1050:	ps_madd.	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5
	Update-CR0: True

0x1054:	ps_madds0	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5

0x1058:	ps_madds0.	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5
	Update-CR0: True

0x105c:	ps_madds1	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5

0x1060:	ps_madds1.	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5
	Update-CR0: True

0x1064:	ps_merge00	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4

0x1068:	ps_merge00.	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
	Update-CR0: True

0x106c:	ps_merge01	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4

0x1070:	ps_merge01.	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
	Update-CR0: True

0x1074:	ps_merge10	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4

0x1078:	ps_merge10.	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
	Update-CR0: True

0x107c:	ps_merge11	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4

0x1080:	ps_merge11.	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
	Update-CR0: True

0x1084:	ps_mr	f3, f4
	op_count: 2
		operands[0].type: REG = f3
		operands[1].type: REG = f4

0x1088:	ps_mr.	f3, f4
	op_count: 2
		operands[0].type: REG = f3
		operands[1].type: REG = f4
	Update-CR0: True

0x108c:	ps_msub	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5

0x1090:	ps_msub.	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5
	Update-CR0: True

0x1094:	ps_mul	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4

0x1098:	ps_mul.	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
	Update-CR0: True

0x109c:	ps_muls0	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4

0x10a0:	ps_muls0.	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
	Update-CR0: True

0x10a4:	ps_muls1	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4

0x10a8:	ps_muls1.	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
	Update-CR0: True

0x10ac:	ps_nabs	f3, f3
	op_count: 2
		operands[0].type: REG = f3
		operands[1].type: REG = f3

0x10b0:	ps_nabs.	f3, f3
	op_count: 2
		operands[0].type: REG = f3
		operands[1].type: REG = f3
	Update-CR0: True

0x10b4:	ps_neg	f3, f3
	op_count: 2
		operands[0].type: REG = f3
		operands[1].type: REG = f3

0x10b8:	ps_neg.	f3, f3
	op_count: 2
		operands[0].type: REG = f3
		operands[1].type: REG = f3
	Update-CR0: True

0x10bc:	ps_nmadd	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5

0x10c0:	ps_nmadd.	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5
	Update-CR0: True

0x10c4:	ps_nmsub	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5

0x10c8:	ps_nmsub.	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5
	Update-CR0: True

0x10cc:	ps_res	f3, f3
	op_count: 2
		operands[0].type: REG = f3
		operands[1].type: REG = f3

0x10d0:	ps_res.	f3, f3
	op_count: 2
		operands[0].type: REG = f3
		operands[1].type: REG = f3
	Update-CR0: True

0x10d4:	ps_rsqrte	f3, f3
	op_count: 2
		operands[0].type: REG = f3
		operands[1].type: REG = f3

0x10d8:	ps_rsqrte.	f3, f3
	op_count: 2
		operands[0].type: REG = f3
		operands[1].type: REG = f3
	Update-CR0: True

0x10dc:	ps_sel	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5

0x10e0:	ps_sel.	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5
	Update-CR0: True

0x10e4:	ps_sub	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4

0x10e8:	ps_sub.	f3, f3, f4
	op_count: 3
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
	Update-CR0: True

0x10ec:	ps_sum0	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5

0x10f0:	ps_sum0.	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5
	Update-CR0: True

0x10f4:	ps_sum1	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5

0x10f8:	ps_sum1.	f3, f3, f4, f5
	op_count: 4
		operands[0].type: REG = f3
		operands[1].type: REG = f3
		operands[2].type: REG = f4
		operands[3].type: REG = f5
	Update-CR0: True

0x10fc:

