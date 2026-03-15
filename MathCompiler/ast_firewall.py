import ast
from fastapi import HTTPException

ALLOWED_FUNCS = {"sin", "cos", "Min", "Max"}
ALLOWED_VARS  = {"x", "y", "z"}
ALLOWED_OPS   = {
    ast.Add, ast.Sub, ast.Mul, ast.Div, ast.Pow,
    ast.USub, ast.UAdd,
    ast.BinOp, ast.UnaryOp, ast.Call, ast.Expression,
    ast.Constant, ast.Name, ast.Load,
}

class _FirewallVisitor(ast.NodeVisitor):
    def visit_Name(self, node: ast.Name):
        if node.id not in ALLOWED_VARS:
            raise ValueError(f"Unauthorized variable: '{node.id}'")
        self.generic_visit(node)

    def visit_Call(self, node: ast.Call):
        if not isinstance(node.func, ast.Name) or node.func.id not in ALLOWED_FUNCS:
            raise ValueError(f"Unauthorized function call: '{ast.unparse(node.func)}'")
        self.generic_visit(node)

    def generic_visit(self, node: ast.AST):
        if type(node) not in ALLOWED_OPS:
            raise ValueError(f"Forbidden AST node type: {type(node).__name__}")
        super().generic_visit(node)


def validate_equation(equation: str) -> None:

    try:
        tree = ast.parse(equation, mode="eval")
    except SyntaxError as exc:
        raise HTTPException(status_code=400, detail=f"Syntax error: {exc}")

    try:
        _FirewallVisitor().visit(tree)
    except ValueError as exc:
        raise HTTPException(status_code=400, detail=f"Firewall rejection: {exc}")
