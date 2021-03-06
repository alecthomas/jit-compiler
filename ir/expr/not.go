package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_Not struct {
	*BaseIRExpression
	Op1 IRExpression
}

func NewIR_Not(op1 IRExpression) *IR_Not {
	return &IR_Not{
		BaseIRExpression: NewBaseIRExpression(Not),
		Op1:              op1,
	}
}

func (i *IR_Not) ReturnType(ctx *IR_Context) Type {
	return TBool
}
func (i *IR_Not) EncodeWithoutSETE(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	return i.encode(ctx, target, false)
}

func (i *IR_Not) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	return i.encode(ctx, target, true)
}

func (i *IR_Not) encode(ctx *IR_Context, target encoding.Operand, includeSETE bool) ([]lib.Instruction, error) {

	var reg1 encoding.Operand

	result := []lib.Instruction{}
	switch i.Op1.(type) {
	case *IR_Variable:
		variable := i.Op1.(*IR_Variable).Value
		reg1 = ctx.VariableMap[variable]
	case *IR_Uint64:
		value := i.Op1.(*IR_Uint64).Value
		reg1 = ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(reg1.(*encoding.Register))
		result = append(result, asm.MOV(encoding.Uint64(value), reg1))
	case *IR_Equals:
		var err error
		var eq []lib.Instruction
		if includeSETE {
			eq, err = i.Op1.Encode(ctx, target)
		} else {
			eq, err = i.Op1.(*IR_Equals).EncodeWithoutSETE(ctx, target)
		}
		if err != nil {
			return nil, err
		}
		for _, r := range eq {
			result = append(result, r)
		}
		reg1 = target
		// TODO: introduce logical operator interface
	case *IR_LT:
		var err error
		var eq []lib.Instruction
		if includeSETE {
			eq, err = i.Op1.Encode(ctx, target)
		} else {
			eq, err = i.Op1.(*IR_LT).EncodeWithoutSETE(ctx, target)
		}
		if err != nil {
			return nil, err
		}
		for _, r := range eq {
			result = append(result, r)
		}
		reg1 = target
		// TODO: introduce logical operator interface
	case *IR_LTE:
		var err error
		var eq []lib.Instruction
		if includeSETE {
			eq, err = i.Op1.Encode(ctx, target)
		} else {
			eq, err = i.Op1.(*IR_LTE).EncodeWithoutSETE(ctx, target)
		}
		if err != nil {
			return nil, err
		}
		for _, r := range eq {
			result = append(result, r)
		}
		reg1 = target
	case *IR_GT:
		var err error
		var eq []lib.Instruction
		if includeSETE {
			eq, err = i.Op1.Encode(ctx, target)
		} else {
			eq, err = i.Op1.(*IR_GT).EncodeWithoutSETE(ctx, target)
		}
		if err != nil {
			return nil, err
		}
		for _, r := range eq {
			result = append(result, r)
		}
		reg1 = target
		// TODO: introduce logical operator interface
	case *IR_GTE:
		var err error
		var eq []lib.Instruction
		if includeSETE {
			eq, err = i.Op1.Encode(ctx, target)
		} else {
			eq, err = i.Op1.(*IR_GTE).EncodeWithoutSETE(ctx, target)
		}
		if err != nil {
			return nil, err
		}
		for _, r := range eq {
			result = append(result, r)
		}
		reg1 = target
	case *IR_And:
		var err error
		var eq []lib.Instruction
		if includeSETE {
			eq, err = i.Op1.Encode(ctx, target)
		} else {
			eq, err = i.Op1.(*IR_And).EncodeWithoutSETE(ctx, target)
		}
		if err != nil {
			return nil, err
		}
		for _, r := range eq {
			result = append(result, r)
		}
		reg1 = target
	case *IR_Or:
		var err error
		var eq []lib.Instruction
		if includeSETE {
			eq, err = i.Op1.Encode(ctx, target)
		} else {
			eq, err = i.Op1.(*IR_Or).EncodeWithoutSETE(ctx, target)
		}
		if err != nil {
			return nil, err
		}
		for _, r := range eq {
			result = append(result, r)
		}
		reg1 = target
	default:
		return nil, fmt.Errorf("Unsupported ! operation: %s", i.Op1.Type())
	}

	tmpReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(tmpReg)
	instr := []lib.Instruction{}
	if includeSETE {
		// TODO: use test?
		instr = append(instr, asm.CMP_immediate(0, reg1))
		instr = append(instr, asm.XOR(tmpReg, tmpReg))
		instr = append(instr, asm.SETE(tmpReg.Get8BitRegister()))
		instr = append(instr, asm.MOV(tmpReg.ForOperandWidth(target.Width()), target))
	}
	for _, inst := range instr {
		result = append(result, inst)
		ctx.AddInstruction(inst)
	}
	return result, nil
}

func (i *IR_Not) String() string {
	return fmt.Sprintf("!(%s)", i.Op1.String())
}

func (b *IR_Not) AddToDataSection(ctx *IR_Context) error {
	return b.Op1.AddToDataSection(ctx)
}
func (b *IR_Not) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	if IsLiteralOrVariable(b.Op1) {
		return nil, b
	}
	rewrites, expr := b.Op1.SSA_Transform(ctx)
	v := ctx.GenerateVariable()
	rewrites = append(rewrites, NewSSA_Rewrite(v, expr))
	return rewrites, NewIR_Not(NewIR_Variable(v))
}
