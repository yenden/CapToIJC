/* Impose an ordering on the State. */
open util/ordering[State]

/*virtual machine modelisation
* a virtual machine is composed of  a heap for memory allocation  and stack of frame
* each Frame represent the current executing method
*/

one sig vm{
	heap : Heap,
	call_stack : CallStack
}
/*Heap representation*/
one sig Heap {
	map : seq Mappings		--mappings maps each object to its location
}

fact {
	all hp : Heap |  #hp.map >= 1 
	all hp : Heap | not hp.map.hasDups
}

sig Mappings  {
	location: Int,
	value : JcvmRef
}

/*StackFrame representation*/

one sig CallStack {
	cs :  seq Frame
}

fact {
	all cst : CallStack |  #cst.cs = 1  -- to simplify we suppose that there is only one Frame
	all cst : CallStack | not cst.cs.hasDups
}

/* A frame is composed of a PC, a local varible array and an operand stack for the operations of 
* the method associated with the frame
*/

sig Frame {
	pc : PCHead,
	local_var_array : one LocalVarArr,
	op_stack : OpHead
}

sig LocalVarArr {
	var :  seq JcvmOpVarType
}

/* the state of our system is associated with a PC and  an OpStack*/
sig State { 
	rel1 : PC,
	rel2 : OpStack	
}

/*Program Counter modelisation
* PC is modelized as a linked list
* We suppose there is no jump instruction
* and each PC is associated with one instruction
*/
sig PC{
    nxt : lone PC,
	ins : ByteCode
}
abstract sig ByteCode {}
one sig Dup, Bspush extends ByteCode{} -- we only implement 2 bytecodes : Dup and Bspush

one sig PCHead in PC {}                     -- head PC is the root and it is still a PC

fact{
    all n : PC | n not in n.^nxt         -- no cycles (acyclic property)
    no nxt.PCHead                            -- no PC points to Head
    all n : PC - PCHead | some nxt.n       -- for all other PCs, there has to be someone pointing to them
	all disj  n,m : PC |  disj [n.ins, m.ins]  -- each PC is associated with one instruction, no PC share the same instruction
}

/*OpStack is also a linked list
* each node of opStack show the state of the stack
*/

sig OpStack {
	operands :  seq JcvmOpType,
	nxt : lone OpStack
}

one sig OpHead in OpStack {}                    
fact{
    all n : OpStack | n not in n.^nxt        
    no nxt.OpHead                           
    all n : OpStack - OpHead | some nxt.n       
}

/* move form one instruction to another*/
pred moveIns [p: PC, p' : PC,  op : OpStack, op' : OpStack] {
	Dup in p.ins =>
		{	
			op.nxt.operands = op.operands.insert[0, op.operands.first]
			op' = op.nxt	
			p' = p.nxt		
		}
	else  --Bspush
		{
			one b: JcvmByte | op.nxt.operands = op.operands.insert[0, b]
			op' = op.nxt
			p' = p.nxt
	    }

}

/*Transition between states*/
fact {
  all v_m : vm {first.rel2 = (v_m.call_stack.cs.first).op_stack} -- the initial state of the system always point on the operand of the VM first frame
  all s: State, s': s.next {
	moveIns[s.rel1,s'.rel1,s.rel2,s'.rel2]
  }
}

/*Differents types used in the JCVM
*								JCVM_Op_Type
*								/					\
*				JCVM_Var_type				JCVM_Ref_Type
*					/		\			\						/				\
*				Byte		Short	Int				Array		ClassOrInterface
*
*/

abstract sig  JcvmOpType {}
abstract sig JcvmOpVarType, JcvmRef extends JcvmOpType{}
sig JcvmByte, JcvmInt, JcvmShort  extends  JcvmOpVarType{}

sig JcvmArray extends JcvmRef {
	array : seq JcvmOpType
}
sig JcvmClassOrInterface extends JcvmRef {}
run {} for 2  but 1 Frame, 2 OpStack, 2 State
