package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func Compare(op1, op2 IRExpression, ctx *IR_Context) ([]lib.Instruction, error) {

	result := []lib.Instruction{}
	returnType1, returnType2 := op1.ReturnType(ctx), op2.ReturnType(ctx)
	if returnType1 != returnType2 {
		return nil, fmt.Errorf("Unsupported types (%s, %s) in compare operation", returnType1, returnType2)
	}

	var reg1, reg2 encoding.Operand

	if op1.Type() == Variable {
		variable := op1.(*IR_Variable).Value
		reg1 = ctx.VariableMap[variable]
	} else {
		reg1 = ctx.AllocateRegister(returnType1)
		defer ctx.DeallocateRegister(reg1.(*encoding.Register))
		expr1, err := op1.Encode(ctx, reg1)
		if err != nil {
			return nil, err
		}
		result = lib.Instructions(result).Add(expr1)
	}

	if op2.Type() == Variable {
		variable := op2.(*IR_Variable).Value
		reg2 = ctx.VariableMap[variable]
	} else {
		reg2 = ctx.AllocateRegister(returnType1)
		defer ctx.DeallocateRegister(reg2.(*encoding.Register))
		expr2, err := op2.Encode(ctx, reg2)
		if err != nil {
			return nil, err
		}
		result = lib.Instructions(result).Add(expr2)
	}
	cmp := asm.CMP(reg2, reg1)
	result = append(result, cmp)
	ctx.AddInstruction(cmp)
	return result, nil
}
