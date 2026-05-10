import pytest
from ast_firewall import validate_equation
from fastapi import HTTPException
from main import to_wgsl_module
import sympy as sp

def test_firewall_valid():
    # Standard variables and funcs
    validate_equation("sin(x) + cos(y) - z")
    # State variables
    validate_equation("state.x + state.y - state.z")
    # Sqrt and power
    validate_equation("sqrt(x**2 + y**2)")
    # Min/Max
    validate_equation("Min(x, Max(y, z))")

def test_firewall_invalid_vars():
    with pytest.raises(HTTPException) as excinfo:
        validate_equation("x + a")
    assert "Unauthorized variable: 'a'" in str(excinfo.value.detail)

def test_firewall_invalid_funcs():
    with pytest.raises(HTTPException) as excinfo:
        validate_equation("tan(x)")
    assert "Unauthorized function call: 'tan'" in str(excinfo.value.detail)

def test_firewall_forbidden_nodes():
    with pytest.raises(HTTPException) as excinfo:
        validate_equation("x if x > 0 else -x")
    assert "Forbidden AST node type" in str(excinfo.value.detail)

def test_wgsl_conversion():
    x, y, z = sp.symbols('x y z')
    expr = x**2 + y**2 + z**2 - 100
    wgsl = to_wgsl_module(expr)
    assert "let dist = x*x + y*y + z*z - 100.0;" in wgsl or "let dist = pow(x, 2.0) + pow(y, 2.0) + pow(z, 2.0) - 100.0;" in wgsl
    assert "return vec2<f32>(dist, 1.0);" in wgsl
    assert "struct State" in wgsl
    assert "fn map(p: vec3<f32>) -> vec2<f32>" in wgsl

def test_wgsl_state_vars():
    state_x = sp.symbols('state_x')
    expr = state_x - 10
    wgsl = to_wgsl_module(expr)
    assert "state.x - 10.0" in wgsl
