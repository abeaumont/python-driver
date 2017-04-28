package normalizer

import (
	"github.com/bblfsh/python-driver/driver/normalizer/pyast"
	. "github.com/bblfsh/sdk/uast"
	. "github.com/bblfsh/sdk/uast/ann"
)

/*
   Tip: to quickly see the native AST generated by Python you can do:
   from ast import *
   code = 'print("test code")'
   dump(parse(code))
*/

/*
For a description of Python AST nodes:

https://greentreesnakes.readthedocs.io/en/latest/nodes.html

 */

 // TODO: add the "stride", third element of slices in some way once we've the index/slice roles
var AnnotationRules = On(Any).Self(
	On(Not(HasInternalType(pyast.Module))).Error("root must be Module"),
	On(HasInternalType(pyast.Module)).Roles(File).Descendants(
		// Binary Expressions
		On(HasInternalType(pyast.BinOp)).Roles(BinaryExpression).Children(
			On(HasInternalRole("op")).Roles(BinaryExpressionOp),
			On(HasInternalRole("left")).Roles(BinaryExpressionLeft),
			On(HasInternalRole("right")).Roles(BinaryExpressionRight),
		),
		// Comparison operators
		On(HasInternalType(pyast.Eq)).Roles(OpEqual),
		On(HasInternalType(pyast.NotEq)).Roles(OpNotEqual),
		On(HasInternalType(pyast.Lt)).Roles(OpLessThan),
		On(HasInternalType(pyast.LtE)).Roles(OpLessThanEqual),
		On(HasInternalType(pyast.Gt)).Roles(OpGreaterThan),
		On(HasInternalType(pyast.GtE)).Roles(OpGreaterThanEqual),
		On(HasInternalType(pyast.Is)).Roles(OpSame),
		On(HasInternalType(pyast.IsNot)).Roles(OpNotSame),
		On(HasInternalType(pyast.In)).Roles(OpContains),
		On(HasInternalType(pyast.NotIn)).Roles(OpNotContains),

		// Aritmetic operators
		On(HasInternalType(pyast.Add)).Roles(OpAdd),
		On(HasInternalType(pyast.Sub)).Roles(OpSubstract),
		On(HasInternalType(pyast.Mult)).Roles(OpMultiply),
		On(HasInternalType(pyast.Div)).Roles(OpDivide),
		On(HasInternalType(pyast.Mod)).Roles(OpMod),
		// TODO: currently without mapping in the UAST
		//On(HasInternalType(pyast.FloorDiv)).Roles(OpDivide),
		//On(HasInternalType(pyast.Pow)).Roles(???),
		//On(HasInternalType(pyast.MatMult)).Roles(???),

		// Bitwise operators
		On(HasInternalType(pyast.LShift)).Roles(OpBitwiseLeftShift),
		On(HasInternalType(pyast.RShift)).Roles(OpBitwiseRightShift),
		On(HasInternalType(pyast.BitOr)).Roles(OpBitwiseOr),
		On(HasInternalType(pyast.BitXor)).Roles(OpBitwiseXor),
		On(HasInternalType(pyast.BitAnd)).Roles(OpBitwiseAnd),

		// Boolean operators
		On(HasInternalType(pyast.And)).Roles(OpBooleanAnd),
		On(HasInternalType(pyast.Or)).Roles(OpBooleanOr),
		On(HasInternalType(pyast.Not)).Roles(OpBooleanNot),

		// UnaryExpression TODO: change it to an specific UAST role if added
		On(HasInternalType(pyast.UnaryOp)).Roles(Expression),

		// Unary operators
		On(HasInternalType(pyast.Invert)).Roles(OpBitwiseComplement),
		On(HasInternalType(pyast.UAdd)).Roles(OpPositive),
		On(HasInternalType(pyast.USub)).Roles(OpNegative),

		On(HasInternalType(pyast.StringLiteral)).Roles(StringLiteral),
		On(HasInternalType(pyast.ByteLiteral)).Roles(ByteStringLiteral),
		On(HasInternalType(pyast.NumLiteral)).Roles(NumberLiteral),
		On(HasInternalType(pyast.Str)).Roles(StringLiteral),
		On(HasInternalType(pyast.BoolLiteral)).Roles(BooleanLiteral),
		// FIXME: JoinedStr are the fstrings (f"my name is {name}"), they have a composite AST
		// with a body that is a list of StringLiteral + FormattedValue(value, conversion, format_spec)
		On(HasInternalType(pyast.JoinedStr)).Roles(StringLiteral),
		On(HasInternalType(pyast.NoneLiteral)).Roles(NullLiteral),
		On(HasInternalType(pyast.Set)).Roles(SetLiteral),
		On(HasInternalType(pyast.List)).Roles(ListLiteral),
		On(HasInternalType(pyast.Dict)).Roles(MapLiteral).Children(
			On(HasInternalRole("keys")).Roles(MapKey),
			On(HasInternalRole("values")).Roles(MapValue),
		),
		On(HasInternalType(pyast.Tuple)).Roles(TupleLiteral),

		// FIXME: decorators
		// FIXME: the FunctionDeclarationReceiver is not set for methods; it should be taken from the parent
		// Type node Token (2 levels up) but the SDK doesn't allow this
		// TODO: create an issue for the SDK
		On(HasInternalType(pyast.FunctionDef)).Roles(FunctionDeclaration, FunctionDeclarationName).Children(
			On(HasInternalType("FunctionDef.body")).Roles(FunctionDeclarationBody),
			// FIXME: change to FunctionDeclarationArgumentS once the PR has been merged
			On(HasInternalType("arguments")).Roles(FunctionDeclarationArgument).Children(
				On(HasInternalRole("args")).Roles(FunctionDeclarationArgument, FunctionDeclarationArgumentName),
				On(HasInternalRole("vararg")).Roles(FunctionDeclarationArgument, FunctionDeclarationVarArgsList,
													FunctionDeclarationArgumentName),
				// FIXME: this is really different from vararg, change it when we have FunctionDeclarationMap
				// or something similar in the UAST
				On(HasInternalRole("kwarg")).Roles(FunctionDeclarationArgument, FunctionDeclarationVarArgsList,
												   FunctionDeclarationArgumentName),
				// Default arguments: Python's AST puts default arguments on a sibling list to the one of
				// arguments that must be mapped to the arguments right-aligned like:
				// a, b=2, c=3 ->
				//		args    [a,b,c],
				//		defaults  [2,3]
				// TODO: create an issue for the SDK
				On(HasInternalType("arguments.defaults")).Roles(FunctionDeclarationArgumentDefaultValue),
			),
		),

		On(HasInternalType(pyast.Call)).Roles(Call).Children(
			On(HasInternalRole("args")).Roles(CallPositionalArgument),
			On(HasInternalRole("keywords")).Roles(CallNamedArgument).Children(
				On(HasInternalRole("value")).Roles(CallNamedArgumentValue),
			),
			On(HasInternalRole("func")).Self(On(HasInternalRole("id"))).Roles(CallCallee),
			On(HasInternalRole("func")).Self(On(HasInternalRole("attr"))).Roles(CallCallee),
			On(HasInternalRole("func")).Self(On(HasInternalType(pyast.Attribute))).Children(
				On(HasInternalRole("id")).Roles(CallReceiver),
			),
		),

		//
		//	Assign => Assigment:
		//		targets[] => AssignmentVariable
		//		value	  => AssignmentValue
		//
		On(HasInternalType(pyast.Assign)).Roles(Assignment).Children(
			On(HasInternalRole("targets")).Roles(AssignmentVariable),
			On(HasInternalRole("value")).Roles(AssignmentValue),
		),

		On(HasInternalType(pyast.AugAssign)).Roles(AugmentedAssignment).Children(
			On(HasInternalRole("op")).Roles(AugmentedAssignmentOperator),
			On(HasInternalRole("target")).Roles(AugmentedAssignmentVariable),
			On(HasInternalRole("value")).Roles(AugmentedAssignmentValue),
		),

		On(HasInternalType(pyast.Expression)).Roles(Expression),
		On(HasInternalType(pyast.Expr)).Roles(Expression),
		On(HasInternalType(pyast.Name)).Roles(SimpleIdentifier),
		On(HasInternalType(pyast.Attribute)).Roles(QualifiedIdentifier),

		// Comments and non significative whitespace
		On(HasInternalType(pyast.SameLineNoops)).Roles(Comment),
		On(HasInternalType(pyast.PreviousNoops)).Roles(Whitespace).Children(
			On(HasInternalRole("lines")).Roles(Comment),
		),
		On(HasInternalType(pyast.RemainderNoops)).Roles(Whitespace).Children(
			On(HasInternalRole("lines")).Roles(Comment),
		),

		// TODO: check what Constant nodes are generated in the python AST and improve this
		On(HasInternalType(pyast.Constant)).Roles(SimpleIdentifier),
		On(HasInternalType(pyast.Try)).Roles(Try).Children(
			On(HasInternalRole("body")).Roles(TryBody),
			On(HasInternalRole("finalbody")).Roles(TryFinally),
			// TODO: this is really a list, use descendents and search for ExceptHandlers?
			On(HasInternalRole("handlers")).Roles(TryCatch),
			// TODO: check that the IfElse really goes here
			On(HasInternalRole("orelse")).Roles(IfElse),
		),
		// FIXME: add OnPath Try.body (uast_type=ExceptHandler) => TryBody
		On(HasInternalType(pyast.TryExcept)).Roles(TryCatch),
		On(HasInternalType(pyast.TryFinally)).Roles(TryFinally),
		On(HasInternalType(pyast.Raise)).Roles(Throw),
		// FIXME: review, add path for the body and items childs
		// FIXME: withitem on Python to RAII on a resource and can aditionally create and alias on it,
		// both of which currently doesn't have representation in the UAST
		On(HasInternalType(pyast.With)).Roles(BlockScope),
		On(HasInternalType(pyast.Return)).Roles(Return),
		On(HasInternalType(pyast.Break)).Roles(Break),
		On(HasInternalType(pyast.Continue)).Roles(Continue),
		// FIXME: extract the test, orelse and the body to test-> IfCondition, orelse -> IfElse, body -> IfBody
		// UAST are first level members
		On(HasInternalType(pyast.If)).Roles(If).Children(
			On(HasInternalType("If.body")).Roles(IfBody),
			On(HasInternalType(pyast.Compare)).Roles(IfCondition),
			On(HasInternalType("If.orelse")).Roles(IfElse),
		),
		On(HasInternalType(pyast.IfExp)).Roles(If, Expression).Children(
			// These are used on ifexpressions (a = 1 if x else 2)
			On(HasInternalRole("body")).Roles(IfBody),
			On(HasInternalRole("test")).Roles(IfCondition),
			On(HasInternalRole("orelse")).Roles(IfElse),
		),
		// One liner if, like a normal If but it will be inside an Assign (like the ternary if in C)
		// also applies the comment about the If
		On(HasInternalType(pyast.IfExp)).Roles(If),
		// FIXME: Import and ImportFrom can make an alias (name -> asname), extract it and put it as
		// uast.ImportAlias
		On(HasInternalType(pyast.Import)).Roles(ImportDeclaration),
		On(HasInternalType(pyast.ImportFrom)).Roles(ImportDeclaration),
		On(HasInternalType(pyast.Alias)).Roles(ImportPath),
		On(HasInternalType(pyast.ClassDef)).Roles(TypeDeclaration).Children(
			On(HasInternalType("ClassDef.body")).Roles(TypeDeclarationBody),
			On(HasInternalType("ClassDef.bases")).Roles(TypeDeclarationBases),
		),

		// FIXME: Internal keys for the ForEach: iter -> ?, target -> ?, body -> ForBody,
		//
		//	For => Foreach:
		//		body => ForBody
		//		iter => ForIter
		//		target => ForTarget
		//		else => IfElse
		//
		On(HasInternalType(pyast.For)).Roles(ForEach).Children(
			On(HasInternalType("For.body")).Roles(ForBody),
			On(HasInternalRole("iter")).Roles(ForExpression),
			On(HasInternalRole("target")).Roles(ForUpdate),
			On(HasInternalType("For.orelse")).Roles(IfElse),
		),
		On(HasInternalType(pyast.While)).Roles(While).Children(
			On(HasInternalType("While.body")).Roles(WhileBody),
			On(HasInternalRole("test")).Roles(WhileCondition),
			On(HasInternalType("While.orelse")).Roles(IfElse),
		),
		On(HasInternalType(pyast.Pass)).Roles(Noop),
		On(HasInternalType(pyast.Num)).Roles(NumberLiteral),
		// FIXME: this is the annotated assignment (a: annotation = 3) not exactly Assignment
		// it also lacks AssignmentValue and AssignmentVariable (see how to add them)
		On(HasInternalType(pyast.AnnAssign)).Roles(Assignment),
		On(HasInternalType(pyast.AugAssign)).Roles(Assignment),
		On(HasInternalType(pyast.Assert)).Roles(Assert),

		// These are AST nodes in Python2 but we convert them to functions in the UAST like
		// they are in Python3
		On(HasInternalType(pyast.Exec)).Roles(Call).Children(
			On(HasInternalRole("body")).Roles(CallPositionalArgument),
			On(HasInternalRole("globals")).Roles(CallPositionalArgument),
			On(HasInternalRole("locals")).Roles(CallPositionalArgument),
		),
		// Repr already comes as a Call \o/
		// Print as a function too.
		// FIXME: the UAST generated is missing the name node
		On(HasInternalType(pyast.Print)).Roles(Call, CallCallee).Children(
			On(HasInternalRole("dest")).Roles(CallPositionalArgument),
			On(HasInternalRole("nl")).Roles(CallPositionalArgument),
			On(HasInternalRole("values")).Roles(CallPositionalArgument).Children(
				On(Any).Roles(CallPositionalArgument),
			),
		),

		// Python annotations for variables, function argument or return values doesn't have any semantic
		// information by themselves and this we consider it comments (some preprocessors or linters can use
		// them, the runtimes ignore them). The TOKEN will take the annotation in the UAST node so
		// the information is keept in any case.
		On(HasInternalRole("annotation")).Roles(Comment),
		On(HasInternalRole("returns")).Roles(Comment),

		// Python very odd ellipsis operator. Has a special rule in tonoder synthetic tokens
		// map to load it with the token "PythonEllipsisOperator" and gets the role SimpleIdentifier
		On(HasInternalType(pyast.Ellipsis)).Roles(SimpleIdentifier),
	),
)

