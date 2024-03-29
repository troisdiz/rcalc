// Rcalc.g4
grammar Rcalc;

// Tokens
fragment DOT: '.' ;
fragment INT_NUMBER: [0-9]+ ;
fragment DECIMAL_NUMBER: [0-9]*DOT[0-9]+ | [0-9]+(DOT[0-9]*)? ;
fragment SCIENTIFIC_NUMBER: [0-9]*(DOT[0-9])?[eE][+-]?[0-9]+ ;

// Numbers are not signed, signs are added at the parser level to handle various different
// cases which are different between the RPN and arithmetics expressions
NUMBER
    : INT_NUMBER
    | DECIMAL_NUMBER
    | SCIENTIFIC_NUMBER
    ;

OP_ADD: '+';
OP_SUB: '-';
OP_MUL: '*';
OP_DIV: '/';
OP_POW: '^' ;

OP_TEST_EQUAL: '==';
OP_TEST_LT: '<';
OP_TEST_GT: '>';
OP_TEST_LET: '<=';
OP_TEST_GET: '>=';

DQUOTE: '"';
QUOTE: '\'';
COMMA: ',';

PAREN_OPEN: '(' ;
PAREN_CLOSE: ')' ;

CURLY_OPEN: '{';
CURLY_CLOSE: '}';

BRACKET_OPEN: '[';
BRACKET_CLOSE: ']';

PROG_OPEN: '<<';
PROG_CLOSE: '>>';

KW_START: 'start';
KW_FOR: 'for';
KW_NEXT: 'next';

KW_IF: 'if';
KW_THEN: 'then';
KW_ELSE: 'else';
KW_END: 'end';

NAME: [a-zA-Z_][a-zA-Z0-9_]*;

// We define whitespaces but we cannot skip them since in RPN mode
// 2-3 must not parse and 2 - 3 and 2 -3 are not the same thing
// This is still useful to specify them at various places in the grammar
WHITESPACE: [ \r\n\t]+;

// Rules
start : instr_seq EOF;

instr_seq: WHITESPACE* instr (WHITESPACE+ instr)* WHITESPACE* # InstructionSequence;

instr
    : action_or_var_call         # InstrActionOrVarCall
    | op                         # InstrOp
    | variable                   # InstrVariable
    | if_then_else               # InstrIfThenElse
    | start_next_loop            # InstrStartNextLoop
    | for_next_loop              # InstrForNextLoop
    | local_var_creation         # InstrLocalVarCreation
    ;

op
    : OP_ADD | OP_SUB | OP_MUL | OP_DIV | OP_POW
    | OP_TEST_EQUAL | OP_TEST_GT | OP_TEST_GET | OP_TEST_LT | OP_TEST_LET
    ;

if_then_else
    : KW_IF instr_seq KW_THEN instr_seq (KW_ELSE instr_seq)* KW_END ;

start_next_loop: KW_START instr_seq KW_NEXT ;
for_next_loop: KW_FOR WHITESPACE* variableDeclaration instr_seq KW_NEXT ;

program_declaration: PROG_OPEN instr_seq PROG_CLOSE ;

local_var_creation
    : '->' (WHITESPACE* variableDeclaration)+ WHITESPACE* statement_for_local_var_creation # LocalVarCreation
    ;

variableDeclaration: NAME #DeclarationVariable;

statement_for_local_var_creation
    : program_declaration         # StatementLocalVarProgram
    | quoted_algebraic_expression # StatementLocalVarAlgebraicExpression
    ;

variable
    : number                      # VariableNumber
    | quoted_algebraic_expression # VariableAlgebraicExpression
    | program_declaration         # VariableProgramDeclaration
    | list                        # VariableList
    | vector                      # VariableVector
    ;

number: (OP_ADD|OP_SUB)?NUMBER ;

quoted_algebraic_expression: QUOTE WHITESPACE* alg_expression WHITESPACE* QUOTE ;

alg_expression
   : alg_mulExpression WHITESPACE* ((OP_ADD | OP_SUB) WHITESPACE* alg_mulExpression)* # AlgExprAddSub
   ;

alg_mulExpression
   : alg_powExpression WHITESPACE* ((OP_MUL | OP_DIV) WHITESPACE* alg_powExpression)* # AlgExprMulDiv
   ;

alg_powExpression
   : alg_signedAtom (OP_POW alg_signedAtom)* #AlgExprPow
   ;

alg_signedAtom
   : OP_ADD WHITESPACE* alg_signedAtom # AlgExprAddSignedAtom
   | OP_SUB WHITESPACE* alg_signedAtom # AlgExprSubSignedAtom
   | alg_func_call         # AlgExprFuncAtom
   | alg_atom              # AlgExprAtom
   ;

alg_atom
   : NUMBER                                # AlgExprNumber
   | alg_variable                          # AlgExprVariable
   | PAREN_OPEN alg_expression PAREN_CLOSE # AlgExprParen
   ;

alg_variable
   : NAME
   ;

alg_func_call
   : function_name=NAME PAREN_OPEN alg_expression (COMMA alg_expression)* PAREN_CLOSE # AlgExprFuncCall
   ;

list : CURLY_OPEN WHITESPACE* (list_item WHITESPACE*)* CURLY_CLOSE;

list_item : variable # ListItem;

vector : BRACKET_OPEN (vector+|number+) BRACKET_CLOSE ;

action_or_var_call: NAME;
