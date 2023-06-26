package rcalc

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"github.com/shopspring/decimal"
	"strings"
	parser "troisdizaines.com/rcalc/rcalc/parser"
)

func tokenToPosition(ref []int, tokens []int) ([]int, error) {
	result := make([]int, len(tokens))
	var found bool
	for idx, token := range tokens {
		found = false
		for pos, candidate := range ref {
			if candidate == token {
				found = true
				result[idx] = pos
				continue
			}
		}
		if !found {
			//myLexer.getVocabulary.getSymbolicName(myTerminalNode.getSymbol.getType)

			return nil, fmt.Errorf("token %d at position %d is not expected", token, idx)
		}
	}
	return result, nil
}

func parserNumber(txt string) (Variable, error) {
	number, err := decimal.NewFromString(txt)
	if err != nil {
		return nil, err
	}
	return CreateNumericVariable(number), nil
}

func parseAction(txt string, registry *ActionRegistry) (Action, error) {
	lowerTxt := strings.ToLower(txt)
	if registry.ContainsOpCode(lowerTxt) {
		return registry.GetAction(txt), nil
	} else {
		return nil, fmt.Errorf("unknown action")
	}
}

type Location struct {
	start, stop antlr.Token
}

func toLocation(ruleContext antlr.ParserRuleContext) Location {
	return Location{
		start: ruleContext.GetStart(),
		stop:  ruleContext.GetStop(),
	}
}

type ValidationError struct {
	location Location
	err      error
}

func (ve *ValidationError) String() string {
	return ve.err.Error()
}

func toErrorMessage(validationErrors []ValidationError) string {
	var errorsAsString []string
	for _, validationError := range validationErrors {
		errorsAsString = append(errorsAsString, " - "+validationError.String())
	}
	return strings.Join(errorsAsString, "\n")
}

type ParseContext[T any] interface {
	GetParent() ParseContext[T]
	SetParent(ctx ParseContext[T])

	AddItem(item LocatedItem[T])
	AddIdentifier(id LocatedItem[string])
	ReportValidationError(location Location, err error)
	GetValidationErrors() []ValidationError

	// BackFromChild Add self as argument to get polymorphism
	BackFromChild(self ParseContext[T], child ParseContext[T])
	CreateFinalItem() ([]T, error)

	TokenVisited(token int)
}

type LocatedItem[T any] struct {
	Location
	item T
}

func newLocatedItem[T any](item T, start antlr.Token, stop antlr.Token) LocatedItem[T] {
	return LocatedItem[T]{
		Location: Location{
			start: start,
			stop:  stop,
		},
		item: item,
	}
}

func toNonLocated[T any](locatedItems []LocatedItem[T]) []T {
	if locatedItems == nil {
		return nil
	}
	var result []T = make([]T, len(locatedItems))
	for idx, locatedItem := range locatedItems {
		result[idx] = locatedItem.item
	}
	return result
}

type BaseParseContext[T any] struct {
	parent           ParseContext[T]
	location         Location
	items            []LocatedItem[T]
	idDeclarations   []LocatedItem[string]
	validationErrors []ValidationError
}

var _ ParseContext[string] = (*BaseParseContext[string])(nil)

func (g *BaseParseContext[T]) GetParent() ParseContext[T] {
	return g.parent
}

func (g *BaseParseContext[T]) SetParent(ctx ParseContext[T]) {
	g.parent = ctx
}

func (g *BaseParseContext[T]) AddItem(item LocatedItem[T]) {
	g.items = append(g.items, item)
}

func (g *BaseParseContext[T]) AddIdentifier(id LocatedItem[string]) {
	g.idDeclarations = append(g.idDeclarations, id)
}

func (g *BaseParseContext[T]) ReportValidationError(location Location, err error) {
	g.validationErrors = append(g.validationErrors, ValidationError{location: location, err: err})
}

func (g *BaseParseContext[T]) GetValidationErrors() []ValidationError {
	return g.validationErrors
}

func (g *BaseParseContext[T]) BackFromChild(self ParseContext[T], child ParseContext[T]) {
	childValidationErrors := child.GetValidationErrors()
	if len(childValidationErrors) > 0 {
		g.validationErrors = append(g.validationErrors, childValidationErrors...)
	} else {
		actions, err := child.CreateFinalItem()
		if err != nil {
			self.ReportValidationError(Location{}, err)
		} else {
			for _, action := range actions {
				self.AddItem(newLocatedItem(action, nil, nil))
			}
		}
	}
}

func (g *BaseParseContext[T]) CreateFinalItem() ([]T, error) {
	panic("CreateFinalItem must be implemented by sub structures")
}

func (g *BaseParseContext[T]) GetItems() []LocatedItem[T] {
	return g.items
}

func (g *BaseParseContext[T]) TokenVisited(token int) {}

type IfThenElseContext struct {
	BaseParseContext[Action] // to avoid reimplementing the interface

	actions       [][]LocatedItem[Action]
	currentAction int
}

var _ ParseContext[Action] = (*IfThenElseContext)(nil)

func (i *IfThenElseContext) AddItem(item LocatedItem[Action]) {
	GetLogger().Debugf("AddItem called in IfThenElseContext")
	i.actions[i.currentAction] = append(i.actions[i.currentAction], item)
}

func (i *IfThenElseContext) TokenVisited(token int) {
	if i.actions == nil {
		i.actions = make([][]LocatedItem[Action], 3)
	}
	switch token {
	case parser.RcalcLexerKW_IF:
		GetLogger().Debugf("TokenVisited in IfThenElseContext : IF")
		i.currentAction = 0
	case parser.RcalcLexerKW_THEN:
		GetLogger().Debugf("TokenVisited in IfThenElseContext : THEN")
		i.currentAction = 1
	case parser.RcalcLexerKW_ELSE:
		GetLogger().Debugf("TokenVisited in IfThenElseContext : ELSE")
		i.currentAction = 2
	}
}

func (i *IfThenElseContext) CreateFinalItem() ([]Action, error) {
	return []Action{
		&IfThenElseActionDesc{
			ifActions:   toNonLocated(i.actions[0]),
			thenActions: toNonLocated(i.actions[1]),
			elseActions: toNonLocated(i.actions[2]),
		},
	}, nil
}

type StartEndLoopContext struct {
	BaseParseContext[Action]
}

var _ ParseContext[Action] = (*StartEndLoopContext)(nil)

func (pc *StartEndLoopContext) CreateFinalItem() ([]Action, error) {
	return []Action{
		&StartNextLoopActionDesc{actions: toNonLocated(pc.BaseParseContext.items)},
	}, nil
}

type ForNextLoopContext struct {
	BaseParseContext[Action]
}

var _ ParseContext[Action] = (*ForNextLoopContext)(nil)

func (pc *ForNextLoopContext) CreateFinalItem() ([]Action, error) {
	return []Action{
		&ForNextLoopActionDesc{
			varName: pc.BaseParseContext.idDeclarations[0].item,
			actions: toNonLocated(pc.BaseParseContext.items),
		},
	}, nil
}

type ProgramContext struct {
	BaseParseContext[Variable]
	parseContextManager *ParseContextManager
}

var _ ParseContext[Variable] = (*ProgramContext)(nil)

func (pc *ProgramContext) CreateFinalItem() ([]Variable, error) {
	programVariable := CreateProgramVariable(pc.parseContextManager.actionCtxStack.GetLastValues())
	return []Variable{
		programVariable,
	}, nil
}

type InstrLocalVarCreationContext struct {
	BaseParseContext[Action]
	parseContextManager *ParseContextManager
}

var _ ParseContext[Action] = (*InstrLocalVarCreationContext)(nil)

func (pc *InstrLocalVarCreationContext) CreateFinalItem() ([]Action, error) {
	// code is variable type agnostic, only the grammar ensures the variable to use
	// is either a program or algebraic expression
	return []Action{
		&VariableDeclarationActionDesc{
			varNames:           toNonLocated(pc.BaseParseContext.idDeclarations),
			variableToEvaluate: pc.parseContextManager.variableCtxStack.GetLastValues()[0],
		},
	}, nil
}

type ParseContextStack[T any] struct {
	currentActionPcIdx int
	rootActionPc       []ParseContext[T]
	currentActionPc    []ParseContext[T]
	lastActionValues   []T
}

func (pcs *ParseContextStack[T]) GetCurrent() ParseContext[T] {
	if pcs.currentActionPcIdx == -1 {
		return nil
	}
	return pcs.currentActionPc[pcs.currentActionPcIdx]
}

func (pcs *ParseContextStack[T]) GetCurrentRoot() ParseContext[T] {
	return pcs.rootActionPc[pcs.currentActionPcIdx]
}

func (pcs *ParseContextStack[T]) GetLastValues() []T {
	return pcs.lastActionValues
}

func (pcs *ParseContextStack[T]) startNewSubContext(ctx ParseContext[T]) {
	GetLogger().Debugf("CTX: startNewSubContext %T", ctx)
	ctx.SetParent(pcs.currentActionPc[pcs.currentActionPcIdx])
	pcs.currentActionPc[pcs.currentActionPcIdx] = ctx
}

func (pcs *ParseContextStack[T]) backToParentContext() {
	//pcs.parserContextDepth--
	GetLogger().Debugf("CTX: backToParentContext %T", pcs.currentActionPc[pcs.currentActionPcIdx])

	oldCurrent := pcs.currentActionPc[pcs.currentActionPcIdx]
	pcs.currentActionPc[pcs.currentActionPcIdx] = pcs.currentActionPc[pcs.currentActionPcIdx].GetParent()
	pcs.currentActionPc[pcs.currentActionPcIdx].BackFromChild(pcs.currentActionPc[pcs.currentActionPcIdx], oldCurrent)
}

func (pcs *ParseContextStack[T]) switchToNewRootContext(ctx ParserProvider) {

	if len(pcs.rootActionPc)-1 == pcs.currentActionPcIdx {
		pcs.rootActionPc = append(pcs.rootActionPc, &RootContext[T]{
			BaseParseContext: BaseParseContext[T]{
				parent:   nil,
				location: toLocation(ctx),
			},
		})
		pcs.currentActionPcIdx = pcs.currentActionPcIdx + 1
		GetLogger().Debugf("switchToNewRootContext: depth = %d", pcs.currentActionPcIdx+1)
		pcs.currentActionPc = append(pcs.currentActionPc, pcs.rootActionPc[pcs.currentActionPcIdx])
	} else {
		pcs.currentActionPcIdx = pcs.currentActionPcIdx + 1
		GetLogger().Debugf("switchToNewRootContext: depth = %d", pcs.currentActionPcIdx+1)

		pcs.rootActionPc[pcs.currentActionPcIdx] = &RootContext[T]{
			BaseParseContext: BaseParseContext[T]{
				parent:   nil,
				location: toLocation(ctx),
			},
		}
		pcs.currentActionPc[pcs.currentActionPcIdx] = pcs.rootActionPc[pcs.currentActionPcIdx]
	}
}

func (pcs *ParseContextStack[T]) exitRootContext(ctx ParserProvider) {
	action, err := pcs.rootActionPc[pcs.currentActionPcIdx].CreateFinalItem()
	if err != nil {
		panic("Error in exitRootContext")
	}
	pcs.lastActionValues = action
	pcs.currentActionPc[pcs.currentActionPcIdx] = nil
	pcs.rootActionPc[pcs.currentActionPcIdx] = nil
	pcs.currentActionPcIdx = pcs.currentActionPcIdx - 1
	GetLogger().Debugf("exitRootContext: depth after = %d", pcs.currentActionPcIdx+1)
}

type ParseContextManager struct {
	registry *ActionRegistry

	// parserContextDepth int // for pretty logging

	actionCtxStack    *ParseContextStack[Action]
	variableCtxStack  *ParseContextStack[Variable]
	algebraicCtxStack *ParseContextStack[AlgebraicExpressionNode]
}

func CreateParseContextManager(registry *ActionRegistry) *ParseContextManager {
	rootPc := &RootContext[Action]{}

	return &ParseContextManager{
		registry: registry,
		actionCtxStack: &ParseContextStack[Action]{
			currentActionPcIdx: 0,
			rootActionPc:       []ParseContext[Action]{rootPc},
			currentActionPc:    []ParseContext[Action]{rootPc},
		},
		variableCtxStack: &ParseContextStack[Variable]{
			currentActionPcIdx: -1,
			rootActionPc:       nil,
			currentActionPc:    nil,
		},
		algebraicCtxStack: &ParseContextStack[AlgebraicExpressionNode]{
			currentActionPcIdx: -1,
			rootActionPc:       nil,
			currentActionPc:    nil,
		},
	}
}

func (l *ParseContextManager) AddAction(action LocatedItem[Action]) {
	l.actionCtxStack.GetCurrent().AddItem(action)
}

func (l *ParseContextManager) AddVariable(variable LocatedItem[Variable]) {
	l.variableCtxStack.GetCurrent().AddItem(variable)
}

func (l *ParseContextManager) AddVarName(varName LocatedItem[string]) {
	l.actionCtxStack.GetCurrent().AddIdentifier(varName)
}

func (l *ParseContextManager) TokenVisited(token int) {
	if token == parser.RcalcLexerWHITESPACE {
		// ignore whitespace
		// we cannot ask the grammar to skip it in order to make a diference between 2- 3 and 2 -3
		return
	}

	if l.actionCtxStack.GetCurrent() != nil {
		l.actionCtxStack.GetCurrent().TokenVisited(token)
	}
	l.actionCtxStack.GetCurrent().TokenVisited(token)
	if l.variableCtxStack.GetCurrent() != nil {
		l.variableCtxStack.GetCurrent().TokenVisited(token)
	}
	if l.algebraicCtxStack.GetCurrent() != nil {
		l.algebraicCtxStack.GetCurrent().TokenVisited(token)
	}
}

type RootContext[T any] struct {
	BaseParseContext[T]
}

var _ ParseContext[struct{}] = (*RootContext[struct{}])(nil)

func (rac *RootContext[T]) CreateFinalItem() ([]T, error) {
	return toNonLocated(rac.items), nil
}

type RcalcParserListener struct {
	*parser.BaseRcalcListener

	registry       *ActionRegistry
	contextManager *ParseContextManager
}

var _ parser.RcalcListener = (*RcalcParserListener)(nil)

func CreateRcalcParserListener(registry *ActionRegistry) *RcalcParserListener {
	return &RcalcParserListener{
		registry:       registry,
		contextManager: CreateParseContextManager(registry),
	}
}

// ExitInstrActionOrVarCall is called when exiting the InstrActionOrVarCall.
func (l *RcalcParserListener) ExitInstrActionOrVarCall(ctx *parser.InstrActionOrVarCallContext) {
	//fmt.Println("ExitInstrActionOrVarCall")
	action, err := parseAction(ctx.GetText(), l.registry)
	if err != nil {
		//ctx.AddErrorNode(ctx.GetParser().GetCurrentToken())
		l.contextManager.AddAction(newLocatedItem[Action](&VariableEvaluationActionDesc{varName: ctx.GetText()}, ctx.GetStart(), ctx.GetStop()))
	} else {
		l.contextManager.AddAction(newLocatedItem[Action](action, ctx.GetStart(), ctx.GetStop()))
	}
}

// ExitDeclarationVariable is called when exiting the DeclarationVariable production.
func (l *RcalcParserListener) ExitDeclarationVariable(ctx *parser.DeclarationVariableContext) {
	action, err := parseAction(ctx.GetText(), l.registry)
	if err != nil {
		l.contextManager.AddVarName(newLocatedItem(ctx.GetText(), ctx.GetStart(), ctx.GetStop()))
	} else {
		//TODO we should raise error here
		l.contextManager.AddAction(newLocatedItem[Action](action, ctx.GetStart(), ctx.GetStop()))
	}
}

/*********************************************************************************/
/* Instructions */
/*********************************************************************************/

// ExitInstrOp is called when production InstrOp is exited.
func (l *RcalcParserListener) ExitInstrOp(ctx *parser.InstrOpContext) {
	//fmt.Println("ExitInstrOp")
	action, err := parseAction(ctx.GetText(), l.registry)
	if err != nil {
		ctx.AddErrorNode(ctx.GetParser().GetCurrentToken())
	} else {
		l.contextManager.AddAction(newLocatedItem[Action](action, ctx.GetStart(), ctx.GetStop()))
	}
}

func (l *RcalcParserListener) VisitTerminal(node antlr.TerminalNode) {
	l.contextManager.TokenVisited(node.GetSymbol().GetTokenType())
	//fmt.Printf("VisitTerminal : #%s# / #%d#\n", node.GetSymbol().GetText(), node.GetSymbol().GetTokenType())
}

type InstructionSequenceParseContext struct {
	BaseParseContext[Action]
}

var _ ParseContext[Action] = (*InstructionSequenceParseContext)(nil)

func (ispc *InstructionSequenceParseContext) CreateFinalItem() ([]Action, error) {
	return toNonLocated(ispc.BaseParseContext.items), nil
}

// EnterInstructionSequence is called when entering the InstructionSequence production.
func (l *RcalcParserListener) EnterInstructionSequence(c *parser.InstructionSequenceContext) {
	l.contextManager.actionCtxStack.startNewSubContext(&InstructionSequenceParseContext{})
}

// ExitInstructionSequence is called when exiting the InstructionSequence production.
func (l *RcalcParserListener) ExitInstructionSequence(c *parser.InstructionSequenceContext) {
	l.contextManager.actionCtxStack.backToParentContext()

}

// EnterInstrIfThenElse is called when entering the InstrIfThenElse production.
func (l *RcalcParserListener) EnterInstrIfThenElse(ctx *parser.InstrIfThenElseContext) {
	l.contextManager.actionCtxStack.startNewSubContext(&IfThenElseContext{})
}

// ExitInstrIfThenElse is called when entering the InstrIfThenElse production.
func (l *RcalcParserListener) ExitInstrIfThenElse(ctx *parser.InstrIfThenElseContext) {
	l.contextManager.actionCtxStack.backToParentContext()
}

// EnterInstrStartNextLoop is called when production InstrStartNextLoop is entered.
func (l *RcalcParserListener) EnterInstrStartNextLoop(ctx *parser.InstrStartNextLoopContext) {
	loopContext := &StartEndLoopContext{}
	l.contextManager.actionCtxStack.startNewSubContext(loopContext)
}

// ExitInstrStartNextLoop is called when production InstrStartNextLoop is exited.
func (l *RcalcParserListener) ExitInstrStartNextLoop(ctx *parser.InstrStartNextLoopContext) {
	l.contextManager.actionCtxStack.backToParentContext()
}

// EnterInstrForNextLoop is called when exiting the InstrForNextLoop production.
func (l *RcalcParserListener) EnterInstrForNextLoop(ctx *parser.InstrForNextLoopContext) {
	loopContext := &ForNextLoopContext{
		BaseParseContext: BaseParseContext[Action]{
			location: toLocation(ctx),
		},
	}
	l.contextManager.actionCtxStack.startNewSubContext(loopContext)
}

// ExitInstrForNextLoop is called when exiting the InstrForNextLoop production.
func (l *RcalcParserListener) ExitInstrForNextLoop(c *parser.InstrForNextLoopContext) {
	l.contextManager.actionCtxStack.backToParentContext()
}

/*********************************************************************************/
/* Variables */
/*********************************************************************************/

type VariableActionContext struct {
	BaseParseContext[Action] // to avoid reimplementing the interface
	// TODO should be part of API or keep it here ?
	// Maybe interface to only be a variable/action/algebraicexpr provider
	parseContextManager *ParseContextManager
}

var _ ParseContext[Action] = (*VariableActionContext)(nil)

func (vac *VariableActionContext) CreateFinalItem() ([]Action, error) {
	if len(vac.parseContextManager.variableCtxStack.GetLastValues()) > 1 {
		return nil, fmt.Errorf("cannot have multiple variables at this place")
	}
	action := &VariablePutOnStackActionDesc{
		value: vac.parseContextManager.variableCtxStack.GetLastValues()[0],
	}
	return []Action{action}, nil
}

// EnterInstrVariable is called when entering the InstrVariable production.
func (l *RcalcParserListener) EnterInstrVariable(c *parser.InstrVariableContext) {
	l.contextManager.actionCtxStack.startNewSubContext(&VariableActionContext{
		parseContextManager: l.contextManager,
	})
	l.contextManager.variableCtxStack.switchToNewRootContext(c)
}

// ExitInstrVariable is called when exiting the InstrVariable production.
func (l *RcalcParserListener) ExitInstrVariable(c *parser.InstrVariableContext) {
	l.contextManager.variableCtxStack.exitRootContext(c)
	l.contextManager.actionCtxStack.backToParentContext()
}

// ExitVariableNumber is called when production InstrNumber is exited.
func (l *RcalcParserListener) ExitVariableNumber(ctx *parser.VariableNumberContext) {
	//fmt.Printf("ExitInstrNumber: %s\n", ctx.GetText())
	number, err := parserNumber(ctx.GetText())
	if err != nil {
		ctx.AddErrorNode(ctx.GetParser().GetCurrentToken())
	} else {
		//l.contextManager.AddAction(newLocatedItem[Action](&VariablePutOnStackActionDesc{number}, ctx.GetStart(), ctx.GetStop()))
		l.contextManager.AddVariable(newLocatedItem[Variable](number, ctx.GetStart(), ctx.GetStop()))
	}
}

type ParserProvider interface {
	antlr.InterpreterRuleContext
	GetParser() antlr.Parser
}

// EnterVariableAlgebraicExpression is called when production VariableAlgebraicExpression is entered.
func (l *RcalcParserListener) EnterVariableAlgebraicExpression(ctx *parser.VariableAlgebraicExpressionContext) {
	// root for Algebraic expressions
	text := ctx.GetText()
	lenText := len(text)
	exprText := text[1 : lenText-1]

	l.contextManager.variableCtxStack.startNewSubContext(&AlgebraicVariableContext{
		parseContextManager: l.contextManager,
		exprText:            exprText,
	})
	l.contextManager.algebraicCtxStack.switchToNewRootContext(ctx)
}

// ExitVariableAlgebraicExpression is called when production VariableAlgebraicExpression is exited.
func (l *RcalcParserListener) ExitVariableAlgebraicExpression(ctx *parser.VariableAlgebraicExpressionContext) {
	l.contextManager.algebraicCtxStack.exitRootContext(ctx)
	l.contextManager.variableCtxStack.backToParentContext()
}

// EnterAlgExprAddSub is called when entering the AlgExprAddSub production.
func (l *RcalcParserListener) EnterAlgExprAddSub(c *parser.AlgExprAddSubContext) {
	l.contextManager.algebraicCtxStack.startNewSubContext(&AlgebraicAddSubContext{
		AlgebraicExprContext{reg: l.registry},
	})
}

// ExitAlgExprAddSub is called when production AlgExprRoot is exited.
func (l *RcalcParserListener) ExitAlgExprAddSub(ctx *parser.AlgExprAddSubContext) {
	l.contextManager.algebraicCtxStack.backToParentContext()
}

// EnterAlgExprMulDiv is called when production AlgExprMulDiv is entered.
func (l *RcalcParserListener) EnterAlgExprMulDiv(ctx *parser.AlgExprMulDivContext) {
	l.contextManager.algebraicCtxStack.startNewSubContext(&AlgebraicMulDivContext{
		AlgebraicExprContext{reg: l.registry},
	})
}

// ExitAlgExprMulDiv is called when production AlgExprMulDiv is exited.
func (l *RcalcParserListener) ExitAlgExprMulDiv(ctx *parser.AlgExprMulDivContext) {
	//fmt.Printf("ExitAlgExprMulDiv %s\n", ctx.GetText())
	l.contextManager.algebraicCtxStack.backToParentContext()
}

// EnterAlgExprPow is called when entering the AlgExprPow production.
func (l *RcalcParserListener) EnterAlgExprPow(c *parser.AlgExprPowContext) {
	l.contextManager.algebraicCtxStack.startNewSubContext(&AlgebraicPowerContext{
		AlgebraicExprContext{reg: l.registry},
	})
}

// ExitAlgExprPow is called when exiting the AlgExprPow production.
func (l *RcalcParserListener) ExitAlgExprPow(c *parser.AlgExprPowContext) {
	l.contextManager.algebraicCtxStack.backToParentContext()
}

// EnterAlgExprFuncCall is called when production AlgExprFuncAtom is entered.
func (l *RcalcParserListener) EnterAlgExprFuncCall(ctx *parser.AlgExprFuncCallContext) {
	functionName := ctx.GetFunction_name().GetText()
	l.contextManager.algebraicCtxStack.startNewSubContext(&AlgebraicFunctionContext{
		AlgebraicExprContext: AlgebraicExprContext{reg: l.registry},
		functionName:         functionName,
	})
}

// ExitAlgExprFuncCall is called when production AlgExprFuncAtom is exited.
func (l *RcalcParserListener) ExitAlgExprFuncCall(ctx *parser.AlgExprFuncCallContext) {
	l.contextManager.algebraicCtxStack.backToParentContext()
}

// EnterAlgExprNumber is called when entering the AlgExprNumber production.
func (l *RcalcParserListener) EnterAlgExprNumber(ctx *parser.AlgExprNumberContext) {
	value, err := decimal.NewFromString(ctx.GetText())
	if err != nil {
		panic(fmt.Sprintf("Cannot parse number %s", ctx.GetText()))
	}
	l.contextManager.algebraicCtxStack.startNewSubContext(&AlgebraicNumberContext{
		AlgebraicExprContext: AlgebraicExprContext{reg: l.registry},
		value:                value,
	})
}

// ExitAlgExprNumber is called when exiting the AlgExprNumber production.
func (l *RcalcParserListener) ExitAlgExprNumber(ctx *parser.AlgExprNumberContext) {
	l.contextManager.algebraicCtxStack.backToParentContext()
}

// EnterAlgExprVariable is called when production AlgExprVariable is entered.
func (l *RcalcParserListener) EnterAlgExprVariable(ctx *parser.AlgExprVariableContext) {
	l.contextManager.algebraicCtxStack.startNewSubContext(&AlgebraicVariableNameContext{
		AlgebraicExprContext: AlgebraicExprContext{reg: l.registry},
		value:                ctx.GetText()})
}

// ExitAlgExprVariable is called when production AlgExprVariable is exited.
func (l *RcalcParserListener) ExitAlgExprVariable(ctx *parser.AlgExprVariableContext) {
	l.contextManager.algebraicCtxStack.backToParentContext()
}

// EnterAlgExprAddSignedAtom is called when production AlgExprAddSignedAtom is entered.
func (l *RcalcParserListener) EnterAlgExprAddSignedAtom(ctx *parser.AlgExprAddSignedAtomContext) {
	l.contextManager.algebraicCtxStack.startNewSubContext(&AlgebraicAtomContext{
		AlgebraicExprContext{
			BaseParseContext: BaseParseContext[AlgebraicExpressionNode]{
				location: Location{
					start: ctx.GetStart(),
					stop:  ctx.GetStop(),
				},
			},
			reg: l.registry},
	})
}

// ExitAlgExprAddSignedAtom is called when production AlgExprAddSignedAtom is exited.
func (l *RcalcParserListener) ExitAlgExprAddSignedAtom(ctx *parser.AlgExprAddSignedAtomContext) {
	l.contextManager.algebraicCtxStack.backToParentContext()
}

// EnterAlgExprSubSignedAtom is called when production AlgExprSubSignedAtom is entered.
func (l *RcalcParserListener) EnterAlgExprSubSignedAtom(ctx *parser.AlgExprSubSignedAtomContext) {
	l.contextManager.algebraicCtxStack.startNewSubContext(&AlgebraicAtomContext{
		AlgebraicExprContext{reg: l.registry},
	})
}

// ExitAlgExprSubSignedAtom is called when production AlgExprSubSignedAtom is exited.
func (l *RcalcParserListener) ExitAlgExprSubSignedAtom(ctx *parser.AlgExprSubSignedAtomContext) {
	l.contextManager.algebraicCtxStack.backToParentContext()
}

// EnterAlgExprAtom is called when production AlgExprAtom is entered.
func (l *RcalcParserListener) EnterAlgExprAtom(ctx *parser.AlgExprAtomContext) {
	l.contextManager.algebraicCtxStack.startNewSubContext(&AlgebraicAtomContext{
		AlgebraicExprContext{reg: l.registry},
	})
}

// ExitAlgExprAtom is called when production AlgExprAtom is exited.
func (l *RcalcParserListener) ExitAlgExprAtom(ctx *parser.AlgExprAtomContext) {
	l.contextManager.algebraicCtxStack.backToParentContext()
}

// EnterVariableProgramDeclaration is called when entering the VariableProgramDeclaration production.
func (l *RcalcParserListener) EnterVariableProgramDeclaration(c *parser.VariableProgramDeclarationContext) {
	l.contextManager.variableCtxStack.startNewSubContext(&ProgramContext{
		parseContextManager: l.contextManager,
	})
	l.contextManager.actionCtxStack.switchToNewRootContext(c)
}

// ExitVariableProgramDeclaration is called when exiting the VariableProgramDeclaration production.
func (l *RcalcParserListener) ExitVariableProgramDeclaration(c *parser.VariableProgramDeclarationContext) {
	l.contextManager.actionCtxStack.exitRootContext(c)
	l.contextManager.variableCtxStack.backToParentContext()
}

type VariableListContext struct {
	BaseParseContext[Action]
}

var _ ParseContext[Action] = (*VariableListContext)(nil)

func (pc *VariableListContext) CreateFinalItem() ([]Action, error) {

	actions := toNonLocated(pc.BaseParseContext.items)
	var itemVars []Variable
	for _, item := range actions {
		itemAction := item.(*VariablePutOnStackActionDesc)
		itemVars = append(itemVars, itemAction.value)
	}
	listVar := CreateListVariable(itemVars)
	return []Action{
		&VariablePutOnStackActionDesc{value: listVar},
	}, nil
}

type ListItemContext struct {
	BaseParseContext[Variable] // to avoid reimplementing the interface
}

func (rlc *ListItemContext) CreateFinalItem() ([]Variable, error) {
	listVariable := CreateListVariable(toNonLocated(rlc.BaseParseContext.items))
	return []Variable{listVariable}, nil
}

var _ ParseContext[Variable] = (*ListItemContext)(nil)

// EnterVariableList is called when entering the VariableList production.
func (l *RcalcParserListener) EnterVariableList(c *parser.VariableListContext) {
	l.contextManager.variableCtxStack.startNewSubContext(&ListItemContext{})
}

// ExitVariableList is called when exiting the VariableList production.
func (l *RcalcParserListener) ExitVariableList(c *parser.VariableListContext) {
	l.contextManager.variableCtxStack.backToParentContext()
}

func (l *RcalcParserListener) EnterListItem(c *parser.ListItemContext) {
	l.contextManager.variableCtxStack.startNewSubContext(&RootContext[Variable]{})
}

func (l *RcalcParserListener) ExitListItem(c *parser.ListItemContext) {
	l.contextManager.variableCtxStack.backToParentContext()
}

/*********************************************************************************/
/* Local var creation */
/*********************************************************************************/

// EnterLocalVarCreation is called when entering the LocalVarCreation production.
func (l *RcalcParserListener) EnterLocalVarCreation(c *parser.LocalVarCreationContext) {
	GetLogger().Debugf("EnterLocalVarCreation: %s", c.GetText())
	l.contextManager.actionCtxStack.startNewSubContext(&InstrLocalVarCreationContext{
		parseContextManager: l.contextManager,
	})
}

// ExitLocalVarCreation is called when exiting the LocalVarCreation production.
func (l *RcalcParserListener) ExitLocalVarCreation(c *parser.LocalVarCreationContext) {
	GetLogger().Debugf("ExitLocalVarCreation: %s", c.GetText())
	l.contextManager.actionCtxStack.backToParentContext()
}

// EnterStatementLocalVarProgram is called when entering the StatementLocalVarProgram production.
func (l *RcalcParserListener) EnterStatementLocalVarProgram(c *parser.StatementLocalVarProgramContext) {
	l.contextManager.variableCtxStack.switchToNewRootContext(c)
	l.contextManager.variableCtxStack.startNewSubContext(&ProgramContext{
		parseContextManager: l.contextManager,
	})
	l.contextManager.actionCtxStack.switchToNewRootContext(c)
}

// ExitStatementLocalVarProgram is called when exiting the StatementLocalVarProgram production.
func (l *RcalcParserListener) ExitStatementLocalVarProgram(c *parser.StatementLocalVarProgramContext) {
	l.contextManager.actionCtxStack.exitRootContext(c)
	l.contextManager.variableCtxStack.backToParentContext()
	l.contextManager.variableCtxStack.exitRootContext(c)
}

// EnterStatementLocalVarAlgebraicExpression is called when entering the StatementLocalVarAlgebraicExpression production.
func (l *RcalcParserListener) EnterStatementLocalVarAlgebraicExpression(c *parser.StatementLocalVarAlgebraicExpressionContext) {
	l.contextManager.variableCtxStack.switchToNewRootContext(c)
	l.contextManager.variableCtxStack.startNewSubContext(&AlgebraicVariableContext{
		parseContextManager: l.contextManager,
		// TODO
		exprText: "TODO",
	})
	l.contextManager.algebraicCtxStack.switchToNewRootContext(c)
}

// ExitStatementLocalVarAlgebraicExpression is called when exiting the StatementLocalVarAlgebraicExpression production.
func (l *RcalcParserListener) ExitStatementLocalVarAlgebraicExpression(c *parser.StatementLocalVarAlgebraicExpressionContext) {
	l.contextManager.algebraicCtxStack.exitRootContext(c)
	l.contextManager.variableCtxStack.backToParentContext()
	l.contextManager.variableCtxStack.exitRootContext(c)
}

/* Error Reporting */

type RcalcParserErrorListener struct {
	messages []string
}

var _ antlr.ErrorListener = (*RcalcParserErrorListener)(nil)

func (el *RcalcParserErrorListener) HasErrors() bool {
	return len(el.messages) > 0

}

func (el *RcalcParserErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	antlrParser := recognizer.(antlr.Parser)
	stack := antlrParser.GetRuleInvocationStack(antlrParser.GetParserRuleContext())
	message := fmt.Sprintf("SyntaxError (%d, %d) : %s with stack %v", line, column, msg, stack)
	el.messages = append(el.messages, message)
}

func (el *RcalcParserErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
	/*	message := "ReportAmbiguity"
		el.messages = append(el.messages, message)*/
}

func (el *RcalcParserErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
	/*	message := "ReportAttemptingFullContext"
		el.messages = append(el.messages, message)*/
}

func (el *RcalcParserErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
	message := "ReportContextSensitivity"
	el.messages = append(el.messages, message)
}

func ParseToActions(cmds string, lexerName string, registry *ActionRegistry) ([]Action, error) {
	return parseToActionsImpl(cmds, lexerName, registry, func(listener parser.RcalcListener) parser.RcalcListener {
		return listener
	})
}

func parseToActionsImpl(cmds string, lexerName string, registry *ActionRegistry, listenerTransformer func(listener parser.RcalcListener) parser.RcalcListener) ([]Action, error) {

	is := antlr.NewInputStream(cmds)

	// Create the Lexer
	lexer := parser.NewRcalcLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := parser.NewRcalcParser(stream)

	// Error Listener
	el := &RcalcParserErrorListener{}

	// Finally parse the expression (by walking the tree)
	var listener *RcalcParserListener = CreateRcalcParserListener(registry)
	//p.RemoveErrorListeners()
	p.AddErrorListener(el)
	parseResult := p.Start_()
	if el.HasErrors() {
		return nil, fmt.Errorf("There are %d error(s):\n - %s", len(el.messages), strings.Join(el.messages, "\n - "))
	}

	var pluggedListener parser.RcalcListener = listenerTransformer(listener)
	antlr.ParseTreeWalkerDefault.Walk(pluggedListener, parseResult)
	if len(listener.contextManager.actionCtxStack.GetCurrentRoot().GetValidationErrors()) > 0 {
		errorsAsString := toErrorMessage(listener.contextManager.actionCtxStack.GetCurrentRoot().GetValidationErrors())
		return nil, fmt.Errorf("There are %d validations error(s):\n%s", len(listener.contextManager.actionCtxStack.GetCurrentRoot().GetValidationErrors()), errorsAsString)
	}

	return listener.contextManager.actionCtxStack.GetCurrentRoot().CreateFinalItem()
}
